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
