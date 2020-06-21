package example

//go:generate go build -o ./testprg/testprg ./testprg
//go:generate binclude
import (
	"github.com/lu4p/binclude"
)

func Example() {
	binclude.Include("./testprg")
}
