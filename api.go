package term

import (
	"os"

	"github.com/mattn/go-runewidth"
	"github.com/zeebo/errs"
)

func Init() error {
	if err := setup_term(); err != nil {
		return errs.New("termbox: error while reading terminfo data: %v", err)
	}

	os.Stdout.WriteString(funcs[t_enter_ca])
	os.Stdout.WriteString(funcs[t_enter_keypad])
	os.Stdout.WriteString(funcs[t_hide_cursor])
	os.Stdout.WriteString(funcs[t_clear_screen])

	termw, termh = get_term_size()
	back_buffer.init(termw, termh)
	front_buffer.init(termw, termh)
	back_buffer.clear()
	front_buffer.clear()

	return nil
}

func Close() {
	os.Stdout.WriteString(funcs[t_show_cursor])
	os.Stdout.WriteString(funcs[t_clear_screen])
	os.Stdout.WriteString(funcs[t_exit_ca])
	os.Stdout.WriteString(funcs[t_exit_keypad])
}

func Flush() error {
	update_size_maybe()

	for y := 0; y < front_buffer.height; y++ {
		line_offset := y * front_buffer.width
		for x := 0; x < front_buffer.width; {
			offset := line_offset + x
			back := &back_buffer.data[offset]
			front := &front_buffer.data[offset]
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

			if w == 2 && x == front_buffer.width-1 {
				send_char(x, y, ' ')
			} else {
				send_char(x, y, *back)
				if w == 2 {
					next := offset + 1
					front_buffer.data[next] = 0
				}
			}
			x += w
		}
	}

	return flush()
}

func Set(x, y int, ch rune) {
	if x < 0 || x >= back_buffer.width {
		return
	}
	if y < 0 || y >= back_buffer.height {
		return
	}

	back_buffer.data[y*back_buffer.width+x] = ch
}

func Size() (width int, height int) {
	return termw, termh
}

func Clear() error {
	err := update_size_maybe()
	back_buffer.clear()
	return err
}
