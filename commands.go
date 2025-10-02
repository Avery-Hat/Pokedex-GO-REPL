package main

import (
	"fmt"
	"os"
	"sort"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

// declare first, without an initializer to avoid cycles
var commands map[string]cliCommand

func init() {
	// populate after functions are declared
	commands = map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
	}
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()

	// Print "help" first, then the rest alphabetically for stable output.
	if c, ok := commands["help"]; ok {
		fmt.Printf("%s: %s\n", c.name, c.description)
	}
	keys := make([]string, 0, len(commands))
	for k := range commands {
		if k == "help" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		c := commands[k]
		fmt.Printf("%s: %s\n", c.name, c.description)
	}
	return nil
}
