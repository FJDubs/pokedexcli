package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

type config struct {
	Next     string
	Previous string
}

func commandExit(Conf *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(Conf *config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Printf("Usage:\n\n")
	for _, cmd := range commandList {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func commandMap(Conf *config) error {
	res, err := http.Get(Conf.Next)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	defer res.Body.Close()

	var locs Locations
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&locs)
	if err != nil {
		return fmt.Errorf("error decoding request: %w", err)
	}

	Conf.Next = locs.Next
	if locs.Previous != nil {
		Conf.Previous = locs.Previous.(string)
	}

	for _, loc := range locs.Results {
		fmt.Println(loc.Name)
	}

	return nil
}

func commandMapB(Conf *config) error {
	if Conf.Previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}
	res, err := http.Get(Conf.Previous)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	defer res.Body.Close()

	var locs Locations
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&locs)
	if err != nil {
		return fmt.Errorf("error decoding request: %w", err)
	}

	Conf.Next = locs.Next
	if locs.Previous != nil {
		Conf.Previous = locs.Previous.(string)
	} else {
		Conf.Previous = ""
	}

	for _, loc := range locs.Results {
		fmt.Println(loc.Name)
	}

	return nil
}
