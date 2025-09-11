package pokeball

import (
	"fmt"
	"math/rand"
	"pokedex/internal/pokeapi"
)

func Throw(pokemon pokeapi.PokemonResponse) {
	randomNumber := rand.Intn(324)
	fmt.Printf("\nBase exp of %v: %v\n", pokemon.Name, pokemon.BaseExperience)
	fmt.Printf("\nGenerated a random number: %v\n", randomNumber)
}