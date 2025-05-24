package jobs

import (
	"fmt"
	"gstrike/pkg/core/transport"
	"gstrike/pkg/util"
)

func ListJobs() {
	util.ListDisplay([]string{"Listeners", "Port", "Status"})
	for i := 0; i < len(transport.Listeners); i++ {
		current := transport.Listeners[i]
		fmt.Printf("%s		%d		%s\n", current.ID, current.Port, current.Status)
	}
	fmt.Printf("\n\n")
}

func StopJob(stop *string) {
	for i := 0; i < len(transport.Listeners); i++ {
		c := transport.Listeners[i]
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
}

func StartJob(start *string) {
	for i := 0; i < len(transport.Listeners); i++ {
		c := transport.Listeners[i]
		if c.ID == *start {
			go c.Start()
			return
		}
	}
	fmt.Printf("%s No such listener ID found...\n", util.PrintBad)
}
