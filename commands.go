// commands.go
package main

import (
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strings"

	"pokedexcli/internal/pokecache"
)

type config struct {
	next    string
	prev    string
	cache   *pokecache.Cache
	pokedex map[string]Pokemon // caught mons by name
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config, []string) error
}

// Declare first to avoid init cycles.
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
		"catch": {
			name:        "catch",
			description: "Try to catch a Pokémon: catch <name>",
			callback:    commandCatch,
		},
		"inspect": { //newly added
			name:        "inspect",
			description: "Inspect a caught Pokémon: inspect <name>",
			callback:    commandInspect,
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

// ---- NEW: catch ----

func commandCatch(cfg *config, args []string) error {
	if len(args) < 1 {
		fmt.Println("usage: catch <pokemon-name>")
		return nil
	}
	// allow "mr mime" -> "mr-mime"
	name := strings.ToLower(strings.Join(args, "-"))

	fmt.Printf("Throwing a Pokeball at %s...\n", name)

	mon, err := fetchPokemonCached(cfg.cache, name)
	if err != nil {
		return err
	}

	// Higher base experience -> harder to catch
	chance := catchChance(mon.BaseExperience)
	if rand.Float64() < chance {
		if cfg.pokedex == nil {
			cfg.pokedex = make(map[string]Pokemon)
		}
		cfg.pokedex[mon.Name] = mon
		fmt.Printf("%s was caught!\n", mon.Name)
	} else {
		fmt.Printf("%s escaped!\n", mon.Name)
	}
	return nil
}

func commandInspect(cfg *config, args []string) error {
	if len(args) < 1 {
		fmt.Println("usage: inspect <pokemon-name>")
		return nil
	}
	// Allow “mr mime” or “Mr-Mime”
	key := strings.ToLower(strings.Join(args, "-"))

	// The pokedex is keyed by the API name (lowercase with dashes)
	mon, ok := cfg.pokedex[key]
	if !ok {
		// try exact key if user typed exact API name already
		mon, ok = cfg.pokedex[strings.ToLower(args[0])]
	}
	if !ok {
		fmt.Println("you have not caught that pokemon")
		return nil
	}

	fmt.Printf("Name: %s\n", mon.Name)
	fmt.Printf("Height: %d\n", mon.Height)
	fmt.Printf("Weight: %d\n", mon.Weight)
	fmt.Println("Stats:")
	for _, s := range mon.Stats {
		fmt.Printf("  - %s: %d\n", s.Stat.Name, s.BaseStat)
	}
	fmt.Println("Types:")
	// types are already in slot order from the API
	for _, t := range mon.Types {
		fmt.Printf("  - %s\n", t.Type.Name)
	}
	return nil
}

// Simple, reasonable catch curve:
// baseExp=50  -> ~0.80
// baseExp=100 -> ~0.70
// baseExp=200 -> ~0.50
// baseExp=300 -> ~0.30
// baseExp>=400 -> clamp ~0.10
func catchChance(baseExp int) float64 {
	ch := 0.9 - (float64(baseExp) / 500.0)
	if ch < 0.10 {
		ch = 0.10
	}
	if ch > 0.90 {
		ch = 0.90
	}
	return ch
}
