package log

import (
	"log"
	"os"
)

var (
	Info = log.New(
		os.Stdout,
		ANSI_BOLD+ANSI_FG_BLUE+"[info]"+ANSI_RESET+" ",
		log.Ltime|log.Lshortfile,
	).Printf

	Warn = log.New(
		os.Stderr,
		ANSI_BOLD+ANSI_FG_YELLOW+"[warn]"+ANSI_RESET+" ",
		log.Ltime|log.Lshortfile,
	).Printf

	Fail = log.New(
		os.Stderr,
		ANSI_BOLD+ANSI_FG_RED+"[warn]"+ANSI_RESET+" ",
		log.Ltime|log.Lshortfile,
	).Printf

	Debg = log.New(
		os.Stdout,
		ANSI_BOLD+ANSI_FG_BRIGHT+"[debg]"+ANSI_RESET+" ",
		log.Ltime|log.Lshortfile,
	).Printf
)
