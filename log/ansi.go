package log

const (
	ANSI_ESC = "\x1b"
	ANSI_CSI = ANSI_ESC + "["

	ANSI_RESET = ANSI_CSI + "0m"
	ANSI_BOLD  = ANSI_CSI + "1m"

	ANSI_FG_BLUE   = ANSI_CSI + "34m"
	ANSI_FG_YELLOW = ANSI_CSI + "33m"
	ANSI_FG_RED    = ANSI_CSI + "31m"
	ANSI_FG_BRIGHT = ANSI_CSI + "97m"
)
