package log

import (
	"fmt"
	"log"
	"os"

	"github.com/ngn13/teawiki/consts"
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

func Banner() {
	fmt.Println()
	bold("teawiki (%s) simple git based wiki", consts.VERSION)
	bold("here are some links you may want to check out:")
	fmt.Println()
	bold("- readme: %s", consts.README)
	bold("- source: %s", consts.SOURCE)
	bold("- donate: %s", consts.DONATE)
	fmt.Println()
}
