package cli

import (
	"bufio"
	"fmt"
	"gstrike/pkg/core"
	"gstrike/pkg/util"
	"os"
	"strings"
)

func Exec() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("%s [%s] > ", util.PROMPT, core.SelectedBeaconId)
		if !scanner.Scan() {
			break
		}
		input := scanner.Text()
		if input == "" {
			continue
		}
		split := strings.Split(input, " ")
		handleCmd(split)
	}
}

func handleCmd(commands []string) {

	cmd, args := commands[0], commands[1:]
	switch cmd {
	case "help":
		help()
	case "clear":
		clear()
	case "exit":
		ExitServer()
	case "use":
		use(args)
	case "https":
		https(args)
	case "generate":
		generate(args)
	case "beacon":
		beacon(args)
	case "jobs":
		jobs(args)
	case "tasks":
		tasks(args)
	case "license":
		license()
	case "version":
		version()
	case "update":
		update()
	case "banner":
		banner()
	default:
		fmt.Println(util.PrintBad + "Unknown command")
	}
}
