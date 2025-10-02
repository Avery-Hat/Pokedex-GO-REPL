package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

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
		if cmd, ok := commands[cmdName]; ok {
			if err := cmd.callback(); err != nil {
				fmt.Println("Error:", err)
			}
			continue
		}

		fmt.Println("Unknown command")
	}

	// check scanner.Err()?
}
