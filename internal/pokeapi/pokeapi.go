package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/FJDubs/pokedexcli/internal/pokecache"
)

type Locations struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous any    `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type Area struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

type Pokemon struct {
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Name           string `json:"name"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
	Weight int `json:"weight"`
}

var cache = pokecache.NewCache(5 * time.Second)

func ListLocations(url string) (Locations, error) {
	var locs Locations
	byt, ok := cache.Get(url)
	if !ok {
		res, err := http.Get(url)
		if err != nil {
			return Locations{}, fmt.Errorf("error with html response: %v", err)
		}
		defer res.Body.Close()
		data, _ := io.ReadAll(res.Body)
		cache.Add(url, data)
		err = json.Unmarshal(data, &locs)
		if err != nil {
			return Locations{}, fmt.Errorf("error unmarshalling response data: %v", err)
		}
	} else {
		err := json.Unmarshal(byt, &locs)
		if err != nil {
			return Locations{}, fmt.Errorf("error unmarshalling from cache: %v", err)
		}
	}
	return locs, nil
}

func ListPokemonAt(url string) ([]string, error) {
	var pkmn []string
	var area Area
	byt, ok := cache.Get(url)
	if !ok {
		res, err := http.Get(url)
		if err != nil {
			return []string{}, fmt.Errorf("error with html response: %v", err)
		}
		defer res.Body.Close()
		data, _ := io.ReadAll(res.Body)
		cache.Add(url, data)
		err = json.Unmarshal(data, &area)
		if err != nil {
			return []string{}, fmt.Errorf("error unmarshalling response data: %v", err)
		}
	} else {
		err := json.Unmarshal(byt, &area)
		if err != nil {
			return []string{}, fmt.Errorf("error unmarshalling from cache: %v", err)
		}
	}
	for _, poke := range area.PokemonEncounters {
		pkmn = append(pkmn, poke.Pokemon.Name)
	}

	return pkmn, nil
}

func GetPokemonInfo(url string) (Pokemon, error) {
	var pkmn Pokemon
	byt, ok := cache.Get(url)
	if !ok {
		res, err := http.Get(url)
		if err != nil {
			return Pokemon{}, fmt.Errorf("error with html response: %v", err)
		}
		defer res.Body.Close()

		if res.StatusCode == 404 {
			return Pokemon{}, fmt.Errorf("pokemon not found")
		}
		if res.StatusCode != 200 {
			return Pokemon{}, fmt.Errorf("unexpected status code: %d", res.StatusCode)
		}

		data, _ := io.ReadAll(res.Body)
		cache.Add(url, data)
		err = json.Unmarshal(data, &pkmn)
		if err != nil {
			return Pokemon{}, fmt.Errorf("error unmarshalling response data: %v", err)
		}
	} else {
		err := json.Unmarshal(byt, &pkmn)
		if err != nil {
			return Pokemon{}, fmt.Errorf("error unmarshalling from cache: %v", err)
		}
	}

	return pkmn, nil
}
