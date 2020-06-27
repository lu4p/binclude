package example

import (
	"github.com/lu4p/binclude"
)

var Testprg = "./testprg/testprg"

func linux() {
	binclude.Include("./testprg/testprg")
}
