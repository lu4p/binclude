package example

//go:generate env GOOS=windows go build -o ./testprg/testprg.exe ./testprg
//go:generate go run ../../cmd/binclude/
import (
	"github.com/lu4p/binclude"
)

var Testprg = "./testprg/testprg.exe"

func windows() {
	binclude.Include("./testprg/testprg.exe")
}
