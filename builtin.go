package term

var eterm_funcs = map[int16]string{
	t_enter_ca:     "\x1b7\x1b[?47h",
	t_exit_ca:      "\x1b[2J\x1b[?47l\x1b8",
	t_show_cursor:  "\x1b[?25h",
	t_hide_cursor:  "\x1b[?25l",
	t_clear_screen: "\x1b[H\x1b[2J",
	t_enter_keypad: "",
	t_exit_keypad:  "",
}

var screen_funcs = map[int16]string{
	t_enter_ca:     "\x1b[?1049h",
	t_exit_ca:      "\x1b[?1049l",
	t_show_cursor:  "\x1b[34h\x1b[?25h",
	t_hide_cursor:  "\x1b[?25l",
	t_clear_screen: "\x1b[H\x1b[J",
	t_enter_keypad: "\x1b[?1h\x1b=",
	t_exit_keypad:  "\x1b[?1l\x1b>",
}

var xterm_funcs = map[int16]string{
	t_enter_ca:     "\x1b[?1049h",
	t_exit_ca:      "\x1b[?1049l",
	t_show_cursor:  "\x1b[?12l\x1b[?25h",
	t_hide_cursor:  "\x1b[?25l",
	t_clear_screen: "\x1b[H\x1b[2J",
	t_enter_keypad: "\x1b[?1h\x1b=",
	t_exit_keypad:  "\x1b[?1l\x1b>",
}

var rxvt_unicode_funcs = map[int16]string{
	t_enter_ca:     "\x1b[?1049h",
	t_exit_ca:      "\x1b[r\x1b[?1049l",
	t_show_cursor:  "\x1b[?25h",
	t_hide_cursor:  "\x1b[?25l",
	t_clear_screen: "\x1b[H\x1b[2J",
	t_enter_keypad: "\x1b=",
	t_exit_keypad:  "\x1b>",
}

var linux_funcs = map[int16]string{
	t_enter_ca:     "",
	t_exit_ca:      "",
	t_show_cursor:  "\x1b[?25h\x1b[?0c",
	t_hide_cursor:  "\x1b[?25l\x1b[?1c",
	t_clear_screen: "\x1b[H\x1b[J",
	t_enter_keypad: "",
	t_exit_keypad:  "",
}

var rxvt_256color_funcs = map[int16]string{
	t_enter_ca:     "\x1b7\x1b[?47h",
	t_exit_ca:      "\x1b[2J\x1b[?47l\x1b8",
	t_show_cursor:  "\x1b[?25h",
	t_hide_cursor:  "\x1b[?25l",
	t_clear_screen: "\x1b[H\x1b[2J",
	t_enter_keypad: "\x1b=",
	t_exit_keypad:  "\x1b>",
}

var terms = []struct {
	name  string
	funcs map[int16]string
}{
	{"Eterm", eterm_funcs},
	{"screen", screen_funcs},
	{"xterm", xterm_funcs},
	{"rxvt-unicode", rxvt_unicode_funcs},
	{"linux", linux_funcs},
	{"rxvt-256color", rxvt_256color_funcs},
}
