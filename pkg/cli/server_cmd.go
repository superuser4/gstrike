package cli

import (
	"flag"
	"fmt"
	"gstrike/pkg/core"
	"gstrike/pkg/core/comms"
	"gstrike/pkg/util"
	"os"
	"strconv"
	"strings"
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

	if *beacon != "" {
		var exist bool
		for i := 0; i < len(core.Beacons); i++ {
			if core.Beacons[i].ID == *beacon {
				exist = true
				break
			}
		}
		if exist {
			fmt.Printf("%s Beacon selected: %s\n", util.PrintStatus, *beacon)
			core.SelectedBeaconId = *beacon
		}
	}
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
		listener.Status = "running"
		go listener.Start()
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

	var lol = []string{"ID", "External IP", "Internal IP", "Hostname", "Username", "Domain", "OS", "Arch", "PID", "LastSeen", "FirstSeen"}
	ListDisplay(lol)

	for i := 0; i < len(core.Beacons); i++ {
		c := core.Beacons[i]

		if *list != "" && c.ID != *list {
			continue
		}
		fmt.Printf("%s\t\t%s\t\t%s\t\t%s\t\t%s\t\t%s\t\t%s\t\t%s\t\t%s\t\t%s\t\t%s\t\t", c.ID, c.ExternalIP, c.InternalIp, c.Hostname, c.Username, c.Domain, c.OS, c.Arch, strconv.Itoa(c.PID), c.LastSeen, c.FirstSeen)
	}
	fmt.Printf("\n")
}

func ListDisplay(list []string) {
	for i := 0; i < len(list); i++ {
		fmt.Printf("%s\t\t", list[i])
	}
	fmt.Printf("\n")
	for i := 0; i < len(list); i++ {
		fmt.Printf("%s\t\t", strings.Repeat("-", len(list[i])))
	}
	fmt.Printf("\n")
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
		ListDisplay([]string{"Listeners", "Port", "Status"})
		for i := 0; i < len(comms.Listeners); i++ {
			current := comms.Listeners[i]
			fmt.Printf("%s		%d		%s\n", current.ID, current.Port, current.Status)
		}
		fmt.Printf("\n\n")

	} else if *stop != "" {
		for i := 0; i < len(comms.Listeners); i++ {
			c := comms.Listeners[i]
			if c.ID == *stop {
				if c.Status == "stopped" {
					fmt.Printf("%s Listener already at a stopped state\n", util.PrintBad)
					return
				}
				err := c.Stop()
				if err != nil {
					fmt.Printf("%s Error while stopping listener <%s>: %v\n", util.PrintBad, c.ID, err)
					return
				}
				fmt.Printf("%s Stopped listener %s\n", util.PrintStatus, *stop)
				return
			}
		}
		fmt.Printf("%s No such listener ID found...\n", util.PrintBad)

	} else if *start != "" {
		for i := 0; i < len(comms.Listeners); i++ {
			c := comms.Listeners[i]
			if c.ID == *start {
				go c.Start()
				return
			}
		}
		fmt.Printf("%s No such listener ID found...\n", util.PrintBad)
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
		ListDisplay([]string{"Task ID", "Beacon ID", "Command", "Status", "Created At", "Finished At", "Output"})

		for i := 0; i < len(core.Tasks); i++ {
			t := core.Tasks[i]
			fmt.Printf("%s\t\t%s\t\t%s\t\t%s\t\t%s\t\t%s\t\t%s\n", t.TaskID, t.BeaconID, t.Command, t.Status, t.CreatedAt, t.FinishedAt, t.Output)
		}
	} else if !*list && *beacon != "" {
		ListDisplay([]string{"Task ID", "Beacon ID", "Command", "Status", "Created At", "Finished At", "Output"})
		var exist bool

		for i := 0; i < len(core.Tasks); i++ {
			t := core.Tasks[i]
			if t.BeaconID == *beacon {
				exist = true
				fmt.Printf("%s\t\t%s\t\t%s\t\t%s\t\t%s\t\t%s\t\t%s\n", t.TaskID, t.BeaconID, t.Command, t.Status, t.CreatedAt, t.FinishedAt, t.Output)
			}
		}
		if !exist {
			fmt.Printf("%s No such beacon ID found...\n", util.PrintBad)
		}
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
