package pokeapi

import (
	"net/http"
	"io"
)

type LocationAreasResponse struct {
	Count	int	`json:"count"`
	Next	string	`json:"next"`
	Previous	*string	`json:"previous"`
	Results []Location	`json:"results"`
}

type Location struct {
	Name	string	`json:"name"`
	URL	string	`json:"url"`
}

type PokemonEncountersResponse struct {
	PokemonEncounters	[]PokemonEncounter	`json:"pokemon_encounters"`
}

type PokemonEncounter struct {
	Pokemon	PokemonEncounterData	`json:"pokemon"`
}

type PokemonEncounterData struct {
	Name	string	`json:"name"`
	URL	string	`json:"url`
}

type PokemonResponse struct {
	Name	string	`json:"name"`
	BaseExperience	int	`json:"base_experience"`
}

func GetResponse(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}