package term

var etermFuncs = map[int16]string{
	fnEnterCa:     "\x1b7\x1b[?47h",
	fnExitCa:      "\x1b[2J\x1b[?47l\x1b8",
	fnShowCursor:  "\x1b[?25h",
	fnHideCursor:  "\x1b[?25l",
	fnClearScreen: "\x1b[H\x1b[2J",
}

var screenFuncs = map[int16]string{
	fnEnterCa:     "\x1b[?1049h",
	fnExitCa:      "\x1b[?1049l",
	fnShowCursor:  "\x1b[34h\x1b[?25h",
	fnHideCursor:  "\x1b[?25l",
	fnClearScreen: "\x1b[H\x1b[J",
}

var xtermFuncs = map[int16]string{
	fnEnterCa:     "\x1b[?1049h",
	fnExitCa:      "\x1b[?1049l",
	fnShowCursor:  "\x1b[?12l\x1b[?25h",
	fnHideCursor:  "\x1b[?25l",
	fnClearScreen: "\x1b[H\x1b[2J",
}

var rxvtUnicodeFuncs = map[int16]string{
	fnEnterCa:     "\x1b[?1049h",
	fnExitCa:      "\x1b[r\x1b[?1049l",
	fnShowCursor:  "\x1b[?25h",
	fnHideCursor:  "\x1b[?25l",
	fnClearScreen: "\x1b[H\x1b[2J",
}

var linuxFuncs = map[int16]string{
	fnEnterCa:     "",
	fnExitCa:      "",
	fnShowCursor:  "\x1b[?25h\x1b[?0c",
	fnHideCursor:  "\x1b[?25l\x1b[?1c",
	fnClearScreen: "\x1b[H\x1b[J",
}

var rxvt256colorFuncs = map[int16]string{
	fnEnterCa:     "\x1b7\x1b[?47h",
	fnExitCa:      "\x1b[2J\x1b[?47l\x1b8",
	fnShowCursor:  "\x1b[?25h",
	fnHideCursor:  "\x1b[?25l",
	fnClearScreen: "\x1b[H\x1b[2J",
}

var terms = []struct {
	name  string
	funcs map[int16]string
}{
	{"Eterm", etermFuncs},
	{"screen", screenFuncs},
	{"xterm", xtermFuncs},
	{"rxvt-unicode", rxvtUnicodeFuncs},
	{"linux", linuxFuncs},
	{"rxvt-256color", rxvt256colorFuncs},
}
