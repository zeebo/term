package term

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"io/ioutil"
	"os"
	"strings"

	"github.com/zeebo/errs"
)

const (
	t_enter_ca     = 28
	t_exit_ca      = 40
	t_show_cursor  = 16
	t_hide_cursor  = 13
	t_clear_screen = 5
	t_enter_keypad = 89
	t_exit_keypad  = 88
)

var funcs = map[int16]string{
	t_enter_ca:     "",
	t_exit_ca:      "",
	t_show_cursor:  "",
	t_hide_cursor:  "",
	t_clear_screen: "",
	t_enter_keypad: "",
	t_exit_keypad:  "",
}

const (
	ti_magic         = 0432
	ti_header_length = 12
)

func load_terminfo() (data []byte, err error) {
	term := os.Getenv("TERM")
	if term == "" {
		return nil, errs.New("termbox: TERM not set")
	}

	terminfo := os.Getenv("TERMINFO")
	if terminfo != "" {
		return ti_try_path(terminfo)
	}

	home := os.Getenv("HOME")
	if home != "" {
		data, err = ti_try_path(home + "/.terminfo")
		if err == nil {
			return data, nil
		}
	}

	dirs := os.Getenv("TERMINFO_DIRS")
	if dirs != "" {
		for _, dir := range strings.Split(dirs, ":") {
			if dir == "" {
				dir = "/usr/share/terminfo"
			}
			data, err = ti_try_path(dir)
			if err == nil {
				return data, nil
			}
		}
	}

	data, err = ti_try_path("/lib/terminfo")
	if err == nil {
		return data, nil
	}

	return ti_try_path("/usr/share/terminfo")
}

func ti_try_path(path string) (data []byte, err error) {
	// load_terminfo already made sure it is set
	term := os.Getenv("TERM")

	// first try, the typical *nix path
	terminfo := path + "/" + term[0:1] + "/" + term
	data, err = ioutil.ReadFile(terminfo)
	if err == nil {
		return data, nil
	}

	// fallback to darwin specific dirs structure
	terminfo = path + "/" + hex.EncodeToString([]byte(term[:1])) + "/" + term
	data, err = ioutil.ReadFile(terminfo)
	if err != nil {
		return nil, errs.Wrap(err)
	}

	return data, nil
}

func setup_term_builtin() error {
	name := os.Getenv("TERM")
	if name == "" {
		return errs.New("termbox: TERM environment variable not set")
	}

	for _, t := range terms {
		if t.name == name {
			funcs = t.funcs
			return nil
		}
	}

	compat_table := []struct {
		partial string
		funcs   map[int16]string
	}{
		{"xterm", xterm_funcs},
		{"rxvt", rxvt_unicode_funcs},
		{"linux", linux_funcs},
		{"Eterm", eterm_funcs},
		{"screen", screen_funcs},
		// let's assume that 'cygwin' is xterm compatible
		{"cygwin", xterm_funcs},
		{"st", xterm_funcs},
	}

	// try compatibility variants
	for _, it := range compat_table {
		if strings.Contains(name, it.partial) {
			funcs = it.funcs
			return nil
		}
	}

	return errs.New("termbox: unsupported terminal")
}

func setup_term() (err error) {
	var data []byte
	var header [6]int16
	var str_offset, table_offset int16

	data, err = load_terminfo()
	if err != nil {
		return setup_term_builtin()
	}

	rd := bytes.NewReader(data)
	if err := binary.Read(rd, binary.LittleEndian, header[:]); err != nil {
		return errs.Wrap(err)
	}

	number_sec_len := int16(2)
	if header[0] == 542 {
		number_sec_len = 4
	}

	if (header[1]+header[2])%2 != 0 {
		header[2] += 1
	}

	str_offset = ti_header_length + header[1] + header[2] + number_sec_len*header[3]
	table_offset = str_offset + 2*header[4]

	for i, _ := range funcs {
		funcs[i], err = ti_read_string(rd, str_offset+2*i, table_offset)
		if err != nil {
			return errs.Wrap(err)
		}
	}

	return nil
}

func ti_read_string(rd *bytes.Reader, str_off, table int16) (string, error) {
	var off int16

	_, err := rd.Seek(int64(str_off), 0)
	if err != nil {
		return "", errs.Wrap(err)
	}
	err = binary.Read(rd, binary.LittleEndian, &off)
	if err != nil {
		return "", errs.Wrap(err)
	}
	_, err = rd.Seek(int64(table+off), 0)
	if err != nil {
		return "", errs.Wrap(err)
	}
	var bs []byte
	for {
		b, err := rd.ReadByte()
		if err != nil {
			return "", errs.Wrap(err)
		}
		if b == byte(0x00) {
			break
		}
		bs = append(bs, b)
	}
	return string(bs), nil
}
