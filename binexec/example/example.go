package example

import (
	"log"

	"github.com/lu4p/binclude/binexec"
)

//go:generate env GOOS=windows go build -o ./testprg/testprg.exe ./testprg
//go:generate env GOOS=linux go build -o ./testprg/testprg ./testprg
//go:generate env GOOS=darwin go build -o ./testprg/testprg_darwin ./testprg

//go:generate go run ../../cmd/binclude/

func Example() {
	cmd, err := binexec.Command(BinFS, Testprg)
	if err != nil {
		log.Fatal(err)
	}

	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
