package main

import (
	"fmt"
	"bufio"
	"os"
	"strings"
	"time"
	"encoding/json"
	"pokedex/internal/pokeapi"
	"pokedex/internal/pokecache"
)

type cliCommand struct {
	name	string
	description	string
	callback func(*config)	error
}

type config struct {
	next	string
	previous	string
	cache	*pokecache.Cache
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
		"mapb": {
			name:	"mapb",
			description:	"Displays the previous 20 Pokemon World areas.",
			callback:	commandMapB,
		},
	}
	conf := config{
		next:	"https://pokeapi.co/api/v2/location-area?offset=0&limit=20",
		previous:	"",
		cache:	pokecache.NewCache(3 * time.Minute),
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
				command.callback(&conf)
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

func commandExit(conf *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(conf *config) error {
	fmt.Println(`Welcome to the Pokedex!
	Usage:
	
	help: Displays a help message
	exit: Exit the Pokedex`)
	return nil
}

func commandMap(conf *config) error {
	var err error
	responseBytes, ok := conf.cache.Get(conf.next)
	if !ok {
		responseBytes, err = pokeapi.GetLocations(conf.next)
		if err != nil {
			fmt.Println("Error:", err)
		}
		conf.cache.Add(conf.next, responseBytes)
	}

	var data pokeapi.LocationAreasResponse
	err = json.Unmarshal(responseBytes, &data)
	if err != nil {
		return err
	}

	conf.next = data.Next
	if data.Previous != nil {
		conf.previous = *(data.Previous)
	} else {
		conf.previous = ""
	}

	for _, location := range data.Results {
		fmt.Println(location.Name)
	}

	return nil
}

func commandMapB(conf *config) error {
	var err error
	if conf.previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}
	responseBytes, ok := conf.cache.Get(conf.previous)
	if !ok {
		responseBytes, err = pokeapi.GetLocations(conf.previous)
		if err != nil {
			fmt.Println("Error:", err)
		}
		conf.cache.Add(conf.previous, responseBytes)
	}

	var data pokeapi.LocationAreasResponse
	err = json.Unmarshal(responseBytes, &data)
	if err != nil {
		return err
	}

	conf.next = data.Next
	if data.Previous != nil {
		conf.previous = *(data.Previous)
	} else {
		conf.previous = ""
	}

	for _, location := range data.Results {
		fmt.Println(location.Name)
	}

	return nil
}

func commandUnknown() error {
	fmt.Println("Unknown command")
	return nil
}
