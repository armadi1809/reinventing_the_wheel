package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	numAgents := 5
	args := os.Args
	if len(args) > 2 {
		fmt.Printf("Usage: %s <num-of-agents>\n", args[0])
		os.Exit(1)
	}
	if len(args) == 2 {
		customNumAgents, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Printf("Invalid number for agents: %s", args[1])
		}
		numAgents = customNumAgents
	}

	fmt.Println(numAgents)

}
