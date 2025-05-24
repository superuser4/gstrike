package cli

import (
	"fmt"
	"strings"

	"github.com/superuser4/c2core/beaconmgr"
	"github.com/superuser4/c2core/taskmgr"
)

var BeaconCommands = func() map[string]func([]string) {
	cmds := make(map[string]func([]string))

	for k, v := range ServerCommands {
		cmds[k] = v
	}

	cmds["help"] = helpBeacon
	cmds["back"] = back
	cmds["shell"] = shell

	return cmds
}()

func helpBeacon(args []string) {
	var cmd string = `
	Beacon Specific:
	
	Commands			Description
	--------			-----------
	help				Prints help with beacon-specific commands included
	back				Exits beacon interaction context
	shell				Sends shell command for beacon to execute (eg. shell <cmd>)

	Server Specific:
	
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
	fmt.Printf("\n\n")
}

func back(args []string) {
	beaconmgr.ResetBeacon()
}

func shell(args []string) {
	var cmd string = strings.Join(args, " ")
	taskmgr.NewTask(cmd, beaconmgr.SelectedBeaconId)
}
