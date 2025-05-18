package main

import (
	"fmt"
	"gstrike/pkg/cli"
	"gstrike/pkg/util"
)

func main() {
	// loading db / config/ etc..
	fmt.Printf(util.BANNER + "\n\n\n")
	fmt.Println(util.PrintStatus + "GStrike v0.0.00 - CodeName")
	// executing cli
	cli.Exec()
}
