// main.go
package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"pokedexcli/internal/pokecache"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	cfg := &config{
		cache:   pokecache.NewCache(5 * time.Second),
		pokedex: make(map[string]Pokemon),
	}

	for {
		fmt.Print("Pokedex > ")

		if !scanner.Scan() {
			break
		}

		words := cleanInput(scanner.Text())
		if len(words) == 0 {
			continue
		}

		cmdName := words[0]
		args := words[1:]
		if cmd, ok := commands[cmdName]; ok {
			if err := cmd.callback(cfg, args); err != nil {
				fmt.Println("Error:", err)
			}
			continue
		}
		fmt.Println("Unknown command")
	}
	// error checker
	if err := scanner.Err(); err != nil {
		fmt.Println("read error:", err)
	}

}
