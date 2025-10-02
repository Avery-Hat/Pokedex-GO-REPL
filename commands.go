// commands.go
package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"pokedexcli/internal/pokecache"
)

type config struct {
	next  string
	prev  string
	cache *pokecache.Cache
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config, []string) error // <-- now accepts args
}

// Declare first (no initializer) to avoid init cycles.
var commands map[string]cliCommand

func init() {
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
		"map": {
			name:        "map",
			description: "List the next 20 location areas",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "List the previous 20 location areas",
			callback:    commandMapBack,
		},
		"explore": {
			name:        "explore",
			description: "Explore a location area: explore <area_name>",
			callback:    commandExplore,
		},
	}
}

func commandExit(_ *config, _ []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(_ *config, _ []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()

	// Print "help" first, then others alphabetically (stable output).
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

func commandMap(cfg *config, _ []string) error {
	url := cfg.next
	resp, err := fetchLocationAreasCached(cfg.cache, url)
	if err != nil {
		return err
	}
	for _, r := range resp.Results {
		fmt.Println(r.Name)
	}
	cfg.next = strptr(resp.Next)
	cfg.prev = strptr(resp.Previous)
	return nil
}

func commandMapBack(cfg *config, _ []string) error {
	if cfg.prev == "" {
		fmt.Println("you're on the first page")
		return nil
	}
	resp, err := fetchLocationAreasCached(cfg.cache, cfg.prev)
	if err != nil {
		return err
	}
	for _, r := range resp.Results {
		fmt.Println(r.Name)
	}
	cfg.next = strptr(resp.Next)
	cfg.prev = strptr(resp.Previous)
	return nil
}

func commandExplore(cfg *config, args []string) error {
	if len(args) < 1 {
		fmt.Println("usage: explore <location-area-name>")
		return nil
	}
	// Input already lowercased by cleanInput; join in case user wrote with spaces.
	area := strings.TrimSpace(strings.Join(args, " "))

	fmt.Printf("Exploring %s...\n", area)
	detail, err := fetchLocationAreaDetailCached(cfg.cache, area)
	if err != nil {
		return err
	}

	fmt.Println("Found Pokemon:")
	for _, enc := range detail.PokemonEncounters {
		fmt.Printf(" - %s\n", enc.Pokemon.Name)
	}
	return nil
}
