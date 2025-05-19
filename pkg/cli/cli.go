package cli

import (
	"bufio"
	"fmt"
	"gstrike/pkg/core"
	"gstrike/pkg/util"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func PrintPrompt() {
	fmt.Printf("%s [%s] > ", util.PROMPT, core.SelectedBeaconId)
}

func Exec() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		for range sigCh {
			fmt.Printf("\nUse 'exit' or Ctrl-D to exit GStrike\n\n")
			PrintPrompt()
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		PrintPrompt()
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				fmt.Fprintf(os.Stderr, "%s IO Error: %v\n\n", util.PrintBad, err)
			} else {
				fmt.Printf("\n")
				ExitServer()
			}
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
	case "update":
		update()
	case "banner":
		banner()

	default:
		fmt.Println(util.PrintBad + "Unknown command")
	}
}
