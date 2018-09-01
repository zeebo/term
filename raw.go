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
	back_buffer  cellbuf
	front_buffer cellbuf
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

func get_term_size() (int, int) {
	var sz winsize
	_, _, _ = syscall.Syscall(syscall.SYS_IOCTL,
		os.Stdout.Fd(), uintptr(syscall.TIOCGWINSZ), uintptr(unsafe.Pointer(&sz)))
	return int(sz.cols), int(sz.rows)
}

func write_cursor(x, y int) {
	outbuf.WriteString("\033[")
	outbuf.Write(strconv.AppendUint(intbuf, uint64(y+1), 10))
	outbuf.WriteString(";")
	outbuf.Write(strconv.AppendUint(intbuf, uint64(x+1), 10))
	outbuf.WriteString("H")
}

func send_char(x, y int, ch rune) {
	var buf [8]byte
	n := utf8.EncodeRune(buf[:], ch)
	if x-1 != lastx || y != lasty {
		write_cursor(x, y)
	}
	lastx, lasty = x, y
	outbuf.Write(buf[:n])
}

func flush() error {
	_, err := outbuf.WriteTo(os.Stdout)
	outbuf.Reset()
	return errs.Wrap(err)
}

func send_clear() error {
	outbuf.WriteString(funcs[t_clear_screen])
	lastx, lasty = -2, -2
	return flush()
}

func update_size_maybe() error {
	w, h := get_term_size()
	if w != termw || h != termh {
		termw, termh = w, h
		back_buffer.resize(termw, termh)
		front_buffer.resize(termw, termh)
		front_buffer.clear()
		return send_clear()
	}
	return nil
}
