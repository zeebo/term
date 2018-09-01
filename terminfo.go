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
	fnEnterCa     = 28
	fnExitCa      = 40
	fnShowCursor  = 16
	fnHideCursor  = 13
	fnClearScreen = 5
)

var funcs = map[int16]string{
	fnEnterCa:     "",
	fnExitCa:      "",
	fnShowCursor:  "",
	fnHideCursor:  "",
	fnClearScreen: "",
}

const (
	tiMagic        = 0432
	tiHeaderLength = 12
)

func loadTerminfo() (data []byte, err error) {
	term := os.Getenv("TERM")
	if term == "" {
		return nil, errs.New("termbox: TERM not set")
	}

	terminfo := os.Getenv("TERMINFO")
	if terminfo != "" {
		return tryTerminfoPath(terminfo)
	}

	home := os.Getenv("HOME")
	if home != "" {
		data, err = tryTerminfoPath(home + "/.terminfo")
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
			data, err = tryTerminfoPath(dir)
			if err == nil {
				return data, nil
			}
		}
	}

	data, err = tryTerminfoPath("/lib/terminfo")
	if err == nil {
		return data, nil
	}

	return tryTerminfoPath("/usr/share/terminfo")
}

func tryTerminfoPath(path string) (data []byte, err error) {
	// loadTerminfo already made sure it is set
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

func setupTermBuiltin() error {
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

	compatTable := []struct {
		partial string
		funcs   map[int16]string
	}{
		{"xterm", xtermFuncs},
		{"rxvt", rxvtUnicodeFuncs},
		{"linux", linuxFuncs},
		{"Eterm", etermFuncs},
		{"screen", screenFuncs},
		{"cygwin", xtermFuncs},
		{"st", xtermFuncs},
	}

	// try compatibility variants
	for _, it := range compatTable {
		if strings.Contains(name, it.partial) {
			funcs = it.funcs
			return nil
		}
	}

	return errs.New("termbox: unsupported terminal")
}

func setupTerm() (err error) {
	var data []byte
	var header [6]int16
	var strOffset, tableOffset int16

	data, err = loadTerminfo()
	if err != nil {
		return setupTermBuiltin()
	}

	rd := bytes.NewReader(data)
	if err := binary.Read(rd, binary.LittleEndian, header[:]); err != nil {
		return errs.Wrap(err)
	}

	numberSecLen := int16(2)
	if header[0] == 542 {
		numberSecLen = 4
	}

	if (header[1]+header[2])%2 != 0 {
		header[2] += 1
	}

	strOffset = tiHeaderLength + header[1] + header[2] + numberSecLen*header[3]
	tableOffset = strOffset + 2*header[4]

	for i := range funcs {
		funcs[i], err = readTerminfoString(rd, strOffset+2*i, tableOffset)
		if err != nil {
			return errs.Wrap(err)
		}
	}

	return nil
}

func readTerminfoString(rd *bytes.Reader, strOffset, table int16) (string, error) {
	var off int16

	_, err := rd.Seek(int64(strOffset), 0)
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
