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
				ExitServer([]string{})
			}
			break
		}
		input := scanner.Text()
		if input == "" {
			continue
		}
		split := strings.Split(input, " ")
		if core.SelectedBeaconId == "" {
			dispatchCmd(ServerCommands, split)
		} else {
			dispatchCmd(BeaconCommands, split)
		}
	}
}

func dispatchCmd(cmdMap map[string]func([]string), input []string) {
	if len(input) == 0 {
		return
	}
	cmd, args := input[0], input[1:]
	if handler, ok := cmdMap[cmd]; ok {
		handler(args)
	} else {
		fmt.Println(util.PrintBad + "Unknown command...")
	}
}
