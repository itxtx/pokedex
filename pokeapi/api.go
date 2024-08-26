package pokeapi

import (
	"encoding/json"
	"net/http"
)

const BaseURL = "https://pokeapi.co/api/v2/"

// Config struct to hold pagination information
type Config struct {
	NextURL     string
	PreviousURL string
}

// LocationAreaResponse represents the JSON response structure for location areas
type LocationAreaResponse struct {
	Count    int            `json:"count"`
	Next     *string        `json:"next"`
	Previous *string        `json:"previous"`
	Results  []LocationArea `json:"results"`
}

// LocationArea represents a single location area in the Pokemon world
type LocationArea struct {
	Name string `json:"name"`
}

// FetchLocationAreas fetches a batch of location areas from the PokeAPI
func FetchLocationAreas(url string) (*LocationAreaResponse, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var locationAreas LocationAreaResponse
	err = json.NewDecoder(resp.Body).Decode(&locationAreas)
	if err != nil {
		return nil, err
	}

	return &locationAreas, nil
}
