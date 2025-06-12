package main

import (
	"fmt"
	"gstrike/pkg/config"
	"gstrike/pkg/core/transport"
	"gstrike/pkg/util"
	"os"
)

func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		fmt.Printf(util.PrintBad + "Error: %v\n",err)
		os.Exit(1)
	}
	server, err := transport.NewGStrike(*conf)
	if err != nil {
		fmt.Printf(util.PrintBad + "Error: %v\n",err)
		os.Exit(1)
	}
	err1 := server.Start()	
	if err1 != nil {
		fmt.Printf(util.PrintBad + "Error: %v\n",err1)
		os.Exit(1)
	}
}

