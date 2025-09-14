package pokedex

import (
	"fmt"
	"pokedex/internal/pokeapi"
)

type Pokedex struct {
	Entries	map[string]PokedexEntry
}

type PokedexEntry struct {
	Name	string
	Height	int
	Weight	int
	Stats	map[string]int
	Types	[]string
}

func NewPokedex() *Pokedex {
	returnPokedex := &Pokedex{
		Entries:	make(map[string]PokedexEntry),
	}
	return returnPokedex
}

func (p *Pokedex) RegisterPokemon(data pokeapi.PokemonResponse) {
	entry := PokedexEntry {
		Name:	data.Name,
		Height:	data.Height,
		Weight:	data.Weight,
		Stats:	mapStats(data),
		Types:	getTypes(data),
	}

	p.Entries[data.Name] = entry
}

func mapStats(data pokeapi.PokemonResponse) map[string]int {
	statMap := make(map[string]int)
	for _, statEntry := range data.Stats {
		statMap[statEntry.Stat.Name] = statEntry.BaseStat
	}
	return statMap
}

func getTypes(data pokeapi.PokemonResponse) []string {
	typeSlice := make([]string, 0, 2)
	for _, typeEntry := range data.Types {
		typeSlice = append(typeSlice, typeEntry.Type.Name)
	}
	return typeSlice
}

func (p *Pokedex) PrintData(name string) {
	entry, ok := p.Entries[name]
	if !ok {
		fmt.Printf("You haven't captured %v yet!\n", name)
	} else {
		fmt.Printf("Name: %v\n", entry.Name)
		fmt.Printf("Height: %v\n", entry.Height)
		fmt.Printf("Weight: %v\n", entry.Weight)
		fmt.Printf("Stats:\n")
		for statName, statValue := range entry.Stats {
			fmt.Printf("  -%v: %v\n", statName, statValue)
		}
		fmt.Printf("Types:\n")
		for _, pokeType := range entry.Types {
			fmt.Printf("  - %v\n", pokeType)
		}
	}
}