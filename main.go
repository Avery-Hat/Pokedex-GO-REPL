// main.go (only the call changes)
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
		cache: pokecache.NewCache(5 * time.Second),
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
		args := words[1:] // <-- pass the rest to the callback
		if cmd, ok := commands[cmdName]; ok {
			if err := cmd.callback(cfg, args); err != nil {
				fmt.Println("Error:", err)
			}
			continue
		}
		fmt.Println("Unknown command")
	}
}
