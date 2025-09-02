package main

import (
	"fmt"
	"bufio"
	"os"
	"strings"
	"pokedex/internal/pokeapi"
)

type cliCommand struct {
	name	string
	description	string
	callback func(config)	error
}

type config struct {
	next	string
	previous	string
}

func main() {
	supportedCommands := map[string]cliCommand{
		"exit": {
			name:	"exit",
			description:	"Exit the Pokedex",
			callback:	commandExit,
		},
		"help": {
			name:	"help",
			description: "Displays a help message",
			callback:	commandHelp,
		},
		"map": {
			name:	"map",
			description: "Displays the next 20 Pokemon World areas.",
			callback:	commandMap,
		},
	}
	conf := config{
		next:	"",
		previous:	"",
	}

	scanner := bufio.NewScanner(os.Stdin)
	for ;; {
		fmt.Print("Pokedex > ")
		if scanner.Scan() {
			input := scanner.Text()
			command, ok := supportedCommands[input]
			if !ok {
				commandUnknown()
			} else {
				command.callback(conf)
			}
		}
	}
	return
}

func cleanInput(text string) []string {
	input := strings.ToLower(strings.TrimSpace(text))
	words := strings.Fields(input)
	return words
}

func commandExit(conf config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(conf config) error {
	fmt.Println(`Welcome to the Pokedex!
	Usage:
	
	help: Displays a help message
	exit: Exit the Pokedex`)
	return nil
}

func commandMap(conf config) error {
	fmt.Println(pokeapi.ImportTest())
	return nil
}

func commandUnknown() error {
	fmt.Println("Unknown command")
	return nil
}
