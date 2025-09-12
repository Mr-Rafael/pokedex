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
	"pokedex/internal/pokeball"
)

type cliCommand struct {
	name	string
	description	string
	callback func(*config, string)	error
}

type config struct {
	next	string
	previous	string
	exploreURL	string
	pokemonURL	string
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
		"explore": {
			name:	"explore",
			description:	"Displays the pokemon in an area. Receives the area name as argument.",
			callback:	commandExplore,
		},
		"catch": {
			name:	"catch",
			description:	"Throws a pokeball at the specified pokemon, with a probability of success",
			callback:	commandCatch,
		},
		"cache": {
			name:	"cache",
			description:	"Prints the currently cached pages.",
			callback:	commandCache,
		},
	}
	conf := config{
		next:	"https://pokeapi.co/api/v2/location-area?offset=0&limit=20",
		previous:	"",
		exploreURL:	"https://pokeapi.co/api/v2/location-area/",
		pokemonURL: "https://pokeapi.co/api/v2/pokemon/",
		cache:	pokecache.NewCache(3 * time.Minute),
	}

	scanner := bufio.NewScanner(os.Stdin)
	for ;; {
		fmt.Print("Pokedex > ")
		if scanner.Scan() {
			input := scanner.Text()
			arguments := strings.Fields(input)
			command, ok := supportedCommands[arguments[0]]
			if !ok {
				commandUnknown()
			} else {
				if len(arguments) == 1 {
					command.callback(&conf, "")
				} else {
					command.callback(&conf, arguments[1])
				}
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

func commandExit(conf *config, arg1 string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(conf *config, arg1 string) error {
	fmt.Println(`Welcome to the Pokedex!
	Usage:
	
	help: Displays a help message
	exit: Exit the Pokedex`)
	return nil
}

func commandMap(conf *config, arg1 string) error {
	var err error
	responseBytes, ok := conf.cache.Get(conf.next)
	if !ok {
		responseBytes, err = pokeapi.GetResponse(conf.next)
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

func commandMapB(conf *config, arg1 string) error {
	var err error
	if conf.previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}
	responseBytes, ok := conf.cache.Get(conf.previous)
	if !ok {
		responseBytes, err = pokeapi.GetResponse(conf.previous)
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

func commandExplore(conf *config, arg1 string) error {
	if len(arg1) <= 0 {
		fmt.Println("Please specify the area to explore.")
		return nil
	}
	var err error
	fullURL := conf.exploreURL + arg1
	responseBytes, ok := conf.cache.Get(fullURL)
	if !ok {
		responseBytes, err = pokeapi.GetResponse(fullURL)
		if err != nil {
			fmt.Println("Error:", err)
		}
		conf.cache.Add(fullURL, responseBytes)
	}

	var data pokeapi.PokemonEncountersResponse
	err = json.Unmarshal(responseBytes, &data)
	if err != nil {
		return err
	}

	for _, pokemon := range data.PokemonEncounters {
		fmt.Println(pokemon.Pokemon.Name)
	}

	return nil
}

func commandCatch(conf *config, arg1 string) error {
	if len(arg1) <= 0 {
		fmt.Println("Please specify the Pokémon you're trying to catch!")
		return nil
	}
	var err error
	fullURL := conf.pokemonURL + arg1
	responseBytes, ok := conf.cache.Get(fullURL)
	if !ok {
		responseBytes, err = pokeapi.GetResponse(fullURL)
		if err != nil {
			fmt.Println("Error:", err)
		}
		conf.cache.Add(fullURL, responseBytes)
	}

	var data pokeapi.PokemonResponse
	err = json.Unmarshal(responseBytes, &data)
	if err != nil {
		return err
	}

	fmt.Printf("Throwing a Pokeball at %v...\n", data.Name)

	if pokeball.Throw(data) {
		fmt.Printf("%v was caught!\n", data.Name)
	} else {
		fmt.Printf("%v escaped!\n", data.Name)
	}

	return nil
}

func commandCache(conf *config, arg1 string) error {
	conf.cache.PrintStatus()
	return nil
}

func commandUnknown() error {
	fmt.Println("Unknown command")
	return nil
}
