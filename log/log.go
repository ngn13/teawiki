package log

import (
	"fmt"
	"log"
	"os"
)

const (
	COLOR_BLUE   = "\033[34m"
	COLOR_YELLOW = "\033[33m"
	COLOR_RED    = "\033[31m"
	COLOR_CYAN   = "\033[36m"
	COLOR_BOLD   = "\033[1m"
	COLOR_RESET  = "\033[0m"
)

var (
	Info = log.New(
		os.Stdout,
		COLOR_BLUE+"[info]"+COLOR_RESET+" ",
		log.Ltime|log.Lshortfile,
	).Printf

	Warn = log.New(
		os.Stderr,
		COLOR_YELLOW+"[warn]"+COLOR_RESET+" ",
		log.Ltime|log.Lshortfile,
	).Printf

	Fail = log.New(
		os.Stderr,
		COLOR_RED+"[fail]"+COLOR_RESET+" ",
		log.Ltime|log.Lshortfile,
	).Printf
)

func bold(msg string, f ...any) {
	fmt.Printf(COLOR_BOLD+msg+COLOR_RESET+"\n", f...)
}
