package cli

import (
	"flag"
	"fmt"
	"gstrike/pkg/core/beaconmgr"
	"gstrike/pkg/core/jobs"
	"gstrike/pkg/core/taskmgr"
	"gstrike/pkg/core/transport"
	"gstrike/pkg/util"
	"os"
)

var ServerCommands = map[string]func([]string){
	"help":     help_cmd,
	"clear":    clear_cmd,
	"exit":     ExitServer_cmd,
	"use":      use_cmd,
	"https":    https_cmd,
	"generate": generate_cmd,
	"beacon":   beacon_cmd,
	"jobs":     jobs_cmd,
	"tasks":    tasks_cmd,
	"license":  license_cmd,
	"update":   update_cmd,
	"banner":   banner_cmd,
}

func help_cmd(args []string) {
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

func clear_cmd(args []string) {
	fmt.Print("\033[H\033[2J")
}

func ExitServer_cmd(args []string) {
	fmt.Printf("%s Running cleanup..\n", util.PrintStatus)

	for i := 0; i < len(transport.Listeners); i++ {
		l := transport.Listeners[i]
		err := l.Stop()
		if err != nil {
			fmt.Printf(util.PrintBad+"Error shutting down Https Listener: %v\n", err)
		} else {
			fmt.Printf("%s Shut down listener: %s\n", util.PrintGood, l.ID)
		}
	}
	os.Exit(0)
}
func use_cmd(args []string) {
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
	beaconmgr.SelectBeacon(beacon)

}
func https_cmd(args []string) {
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
	listener := transport.NewHttps(*port)
	fmt.Printf("%s New Https listener configured\n", util.PrintStatus)
	if *startNow {
		jobs.StartJob(&listener.ID)
	}
}

func generate_cmd(args []string) {}

func beacon_cmd(args []string) {
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
	beaconmgr.ListBeacons(list)
}

func jobs_cmd(args []string) {
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
		jobs.ListJobs()
	} else if *stop != "" {
		jobs.StopJob(stop)
	} else if *start != "" {
		jobs.StartJob(start)
	}
}

func tasks_cmd(args []string) {
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
		taskmgr.PrintTasks()
	} else if !*list && *beacon != "" {
		taskmgr.PrintTask(beacon)
	} else {
		fs.Usage()
		return
	}
}

func license_cmd(args []string) {
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

func update_cmd(args []string) {}
func banner_cmd(args []string) {
	fmt.Printf(util.BANNER + "\n\n")
}
