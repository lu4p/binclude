package main

import "github.com/lu4p/binclude"

//go:generate go build -o=./includedPrg/includedPrg ./includedPrg
func main() {
	binclude.Include("./testdata/bench/includedPrg/includedPrg")
}
