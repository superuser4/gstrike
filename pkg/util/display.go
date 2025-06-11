package util

import (
	"fmt"
	"strings"
)

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
