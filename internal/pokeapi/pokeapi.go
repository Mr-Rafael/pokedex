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

func GetLocations(url string) ([]byte, error) {
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