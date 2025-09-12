package pokeball

import (
	"math/rand"
	"pokedex/internal/pokeapi"
)

const BlisseyBaseExperience = 608
const MewTwoFailChance = 0.94

func Throw(pokemon pokeapi.PokemonResponse) bool {
	maxBaseExp := BlisseyBaseExperience/MewTwoFailChance
	randomNumber := rand.Intn(int(maxBaseExp))
	return randomNumber > pokemon.BaseExperience
}
