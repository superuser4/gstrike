package main

import (
	"fmt"
	"gstrike/pkg/config"
	"gstrike/pkg/util"
)

func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		fmt.Printf(util.PrintBad + "Error: %v\n",err)
	}
	// Start main server

}

