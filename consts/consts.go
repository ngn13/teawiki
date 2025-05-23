package consts

import "fmt"

// VERSION is overwritten by the commit ID when "make" is run with "RELEASE=1"

var (
	VERSION = "dev"
	DOCS    = fmt.Sprintf("https://github.com/ngn13/teawiki/tree/%s/docs", VERSION)
)

const (
	SOURCE = "https://github.com/ngn13/teawiki"
	DONATE = "https://ngn.tf/donate"
)
