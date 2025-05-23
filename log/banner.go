package log

import (
	"fmt"

	"github.com/ngn13/teawiki/consts"
)

func Banner() {
	fmt.Println()
	bold("teawiki (%s) - simple git based wiki", consts.VERSION)
	bold("here are some links you may want to check out:")
	fmt.Println()
	bold("- docs  : %s", consts.DOCS)
	bold("- source: %s", consts.SOURCE)
	bold("- donate: %s", consts.DONATE)
	fmt.Println()
}
