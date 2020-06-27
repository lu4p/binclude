package example

import (
	"github.com/lu4p/binclude"
)

var Testprg = "./testprg/testprg_darwin"

func darwin() {
	binclude.Include("./testprg/testprg_darwin")
}
