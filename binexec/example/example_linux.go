package example

import (
	"github.com/lu4p/binclude"
)

// Testprg the program to run for the Example
var Testprg = "./testprg/testprg"

func linux() {
	binclude.Include("./testprg/testprg")
}
