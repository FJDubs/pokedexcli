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
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int   `json:"chance"`
				ConditionValues []any `json:"condition_values"`
				MaxLevel        int   `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
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
