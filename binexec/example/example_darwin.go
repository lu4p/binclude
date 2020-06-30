package example

import (
	"github.com/lu4p/binclude"
)

// Testprg the program to run for the Example
var Testprg = "./testprg/testprg_darwin"

func darwin() {
	binclude.Include("./testprg/testprg_darwin")
}
