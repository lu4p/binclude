package example

//go:generate go build -o ./testprg/testprg ./testprg
//go:generate go run ../../cmd/binclude/
import (
	"github.com/lu4p/binclude"
)

func Example() {
	binclude.Include("./testprg")
}
