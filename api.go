package term

import (
	"os"

	"github.com/mattn/go-runewidth"
	"github.com/zeebo/errs"
)

// Init starts up the library.
func Init() error {
	if err := setupTerm(); err != nil {
		return errs.New("termbox: error while reading terminfo data: %v", err)
	}

	os.Stdout.WriteString(funcs[fnEnterCa])
	os.Stdout.WriteString(funcs[fnHideCursor])
	os.Stdout.WriteString(funcs[fnClearScreen])

	termw, termh = getTermSize()
	backBuffer.init(termw, termh)
	frontBuffer.init(termw, termh)
	backBuffer.clear()
	frontBuffer.clear()

	return nil
}

// Close should be called before the process exits.
func Close() {
	os.Stdout.WriteString(funcs[fnShowCursor])
	os.Stdout.WriteString(funcs[fnClearScreen])
	os.Stdout.WriteString(funcs[fnExitCa])
}

// Flush writes any changed data to the screen.
func Flush() error {
	updateSize()

	for y := 0; y < frontBuffer.height; y++ {
		lineOffset := y * frontBuffer.width
		for x := 0; x < frontBuffer.width; {
			offset := lineOffset + x
			back := &backBuffer.data[offset]
			front := &frontBuffer.data[offset]
			if *back < ' ' {
				*back = ' '
			}
			w := runewidth.RuneWidth(*back)
			if w == 0 || w == 2 && runewidth.IsAmbiguousWidth(*back) {
				w = 1
			}
			if *back == *front {
				x += w
				continue
			}
			*front = *back

			if w == 2 && x == frontBuffer.width-1 {
				sendChar(x, y, ' ')
			} else {
				sendChar(x, y, *back)
				if w == 2 {
					next := offset + 1
					frontBuffer.data[next] = 0
				}
			}
			x += w
		}
	}

	return flush()
}

// Set puts the character at position x,y to be the given rune.
func Set(x, y int, ch rune) {
	if x < 0 || x >= backBuffer.width {
		return
	}
	if y < 0 || y >= backBuffer.height {
		return
	}

	backBuffer.data[y*backBuffer.width+x] = ch
}

// Size returns the size of the terminal.
func Size() (width int, height int) {
	return termw, termh
}

// Clear clears the terminal.
func Clear() error {
	err := updateSize()
	backBuffer.clear()
	return err
}
