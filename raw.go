package term

import (
	"bytes"
	"os"
	"strconv"
	"syscall"
	"unicode/utf8"
	"unsafe"

	"github.com/zeebo/errs"
)

var (
	backBuffer   cellbuf
	frontBuffer  cellbuf
	termw        int
	termh        int
	outbuf       bytes.Buffer
	intbuf       = make([]byte, 0, 8)
	lastx, lasty int
)

type winsize struct {
	rows    uint16
	cols    uint16
	xpixels uint16
	ypixels uint16
}

func getTermSize() (int, int) {
	var sz winsize
	_, _, _ = syscall.Syscall(syscall.SYS_IOCTL,
		os.Stdout.Fd(), uintptr(syscall.TIOCGWINSZ), uintptr(unsafe.Pointer(&sz)))
	return int(sz.cols), int(sz.rows)
}

func writeCursor(x, y int) {
	outbuf.WriteString("\033[")
	outbuf.Write(strconv.AppendUint(intbuf, uint64(y+1), 10))
	outbuf.WriteString(";")
	outbuf.Write(strconv.AppendUint(intbuf, uint64(x+1), 10))
	outbuf.WriteString("H")
}

func sendChar(x, y int, ch rune) {
	var buf [8]byte
	n := utf8.EncodeRune(buf[:], ch)
	if x-1 != lastx || y != lasty {
		writeCursor(x, y)
	}
	lastx, lasty = x, y
	outbuf.Write(buf[:n])
}

func flush() error {
	_, err := outbuf.WriteTo(os.Stdout)
	outbuf.Reset()
	return errs.Wrap(err)
}

func sendClear() error {
	outbuf.WriteString(funcs[fnClearScreen])
	lastx, lasty = -2, -2
	return flush()
}

func updateSize() error {
	w, h := getTermSize()
	if w != termw || h != termh {
		termw, termh = w, h
		backBuffer.resize(termw, termh)
		frontBuffer.resize(termw, termh)
		frontBuffer.clear()
		return sendClear()
	}
	return nil
}
