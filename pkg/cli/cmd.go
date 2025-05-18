package cli

import (
	"flag"
	"fmt"
	"gstrike/pkg/comms"
	"gstrike/pkg/util"
	"os"
)

func help() {
	var cmd string = `
	Commands			Description
	--------			-----------
	help				Prints help, run "<command> --help" for more
	clear				Clears screen
	exit				Exits GStrike
	use				Selects beacon for interaction
	https				Configures an https listener
	generate			Generates a beacon payload
	beacon				Beacon management
	jobs				Lists running listeners
	tasks				Beacon Task manangement
	license				Prints GStrike's license 
	version				Display version information
	update				Check for updates
	banner				Displays GStrike banner
	`
	fmt.Println(cmd)
}

func clear() {
	fmt.Print("\033[H\033[2J")
}

func ExitServer() {
	for i := 0; i < len(comms.Listeners); i++ {
		l := comms.Listeners[i]
		err := l.Stop()
		if err != nil {
			fmt.Printf(util.PrintBad+"Error shutting down Https Listener: %v\n", err)
		}
	}
	os.Exit(0)
}
func use(args []string) {
	beacon := args[0]
}
func https(args []string) {
	fs := flag.NewFlagSet("https", flag.ContinueOnError)
	port := fs.Int("port", 0, "port to serve on")
	startNow := fs.Bool("start-now", false, "starts the listener now")

	fs.Usage = func() {
		fmt.Println()
		fmt.Println("Usage: https --port <port>")
		fs.PrintDefaults()
	}

	err := fs.Parse(args)
	if err != nil {
		return
	}

	if len(args) == 0 {
		fs.Usage()
		return
	}

	if *port == 0 {
		fmt.Println("Error: --port is required")
		fs.Usage()
		return
	}
	listener := comms.NewHttps(*port)
	if *startNow {
		listener.Status = "running"
		go listener.Start()
	}
}

func generate(args []string) {}
func beacon(args []string)   {}
func jobs(args []string) {
	fmt.Println(`
Listeners				Port	Status
---------				----	------`)
	for i := 0; i < len(comms.Listeners); i++ {
		current := comms.Listeners[i]
		fmt.Printf("%s	%d	%s\n", current.ID, current.Port, current.Status)
	}
	fmt.Printf("\n\n")
}
func tasks(args []string) {}

func license() {
	var license string = `
	Copyright (c) 2025 superuser4

	Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"),
	to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, 
	distribute, sublicense, and/or sell copies of the Software, 
	and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

	The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, 
	INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. 
	IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
	WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
	ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
`
	fmt.Println(license)
}

func version() {}
func update()  {}
func banner() {
	fmt.Printf(util.BANNER + "\n\n")
}
