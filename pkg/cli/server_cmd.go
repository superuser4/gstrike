package cli

import (
	"flag"
	"fmt"
	"gstrike/pkg/core"
	"gstrike/pkg/core/comms"
	"gstrike/pkg/util"
	"os"
)

var ServerCommands = map[string]func([]string){
	"help":     help,
	"clear":    clear,
	"exit":     ExitServer,
	"use":      use,
	"https":    https,
	"generate": generate,
	"beacon":   beacon,
	"jobs":     jobs,
	"tasks":    tasks,
	"license":  license,
	"update":   update,
	"banner":   banner,
}

func help(args []string) {
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
	update				Check for updates
	banner				Displays GStrike banner
	`
	fmt.Println(cmd)
}

func clear(args []string) {
	fmt.Print("\033[H\033[2J")
}

func ExitServer(args []string) {
	fmt.Printf("%s Running cleanup..\n", util.PrintStatus)

	for i := 0; i < len(comms.Listeners); i++ {
		l := comms.Listeners[i]
		err := l.Stop()
		if err != nil {
			fmt.Printf(util.PrintBad+"Error shutting down Https Listener: %v\n", err)
		} else {
			fmt.Printf("%s Shut down listener: %s\n", util.PrintGood, l.ID)
		}
	}
	os.Exit(0)
}
func use(args []string) {
	fs := flag.NewFlagSet("use", flag.ContinueOnError)
	beacon := fs.String("beacon", "", "Starts interaction with beacon")

	fs.Usage = func() {
		fmt.Println()
		fmt.Println("Usage: use --beacon <id>")
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
	core.SelectBeacon(beacon)

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
	fmt.Printf("%s New Https listener configured\n", util.PrintStatus)
	if *startNow {
		core.StartJob(&listener.ID)
	}
}

func generate(args []string) {}

func beacon(args []string) {
	fs := flag.NewFlagSet("beacon", flag.ContinueOnError)
	list := fs.String("list", "", "Lists info about beacon(s)")

	fs.Usage = func() {
		fmt.Println()
		fmt.Println("Usage: beacon --list")
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
	core.ListBeacons(list)
}

func jobs(args []string) {
	fs := flag.NewFlagSet("jobs", flag.ContinueOnError)
	list := fs.Bool("list", false, "Lists all listeners")
	start := fs.String("start", "", "Starts configured listener")
	stop := fs.String("stop", "", "Stops configured listener")

	fs.Usage = func() {
		fmt.Println()
		fmt.Println("Usage: jobs <OPTIONS>")
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

	if *start != "" && *stop != "" && *start == *stop {
		fs.Usage()
		return
	}

	if *list {
		core.ListJobs()
	} else if *stop != "" {
		core.StopJob(stop)
	} else if *start != "" {
		core.StartJob(start)
	}
}

func tasks(args []string) {
	fs := flag.NewFlagSet("tasks", flag.ContinueOnError)
	beacon := fs.String("beacon", "", "Lists one beacon")
	list := fs.Bool("list", false, "Lists all beacons")

	fs.Usage = func() {
		fmt.Println()
		fmt.Println("Usage: tasks <OPTIONS>")
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

	if *list && *beacon == "" {
		core.PrintTasks()
	} else if !*list && *beacon != "" {
		core.PrintTask(beacon)
	} else {
		fs.Usage()
		return
	}
}

func license(args []string) {
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

func update(args []string) {}
func banner(args []string) {
	fmt.Printf(util.BANNER + "\n\n")
}
