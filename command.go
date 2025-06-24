package main

import (
	"fmt"
	"os"

	"github.com/FJDubs/pokedexcli/internal/pokeapi"
)

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
	locs, err := pokeapi.ListLocations(Conf.Next)

	Conf.Next = locs.Next
	if locs.Previous != nil {
		Conf.Previous = locs.Previous.(string)
	}

	for _, loc := range locs.Results {
		fmt.Println(loc.Name)
	}

	return err
}

func commandMapB(Conf *config) error {
	if Conf.Previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}
	locs, err := pokeapi.ListLocations(Conf.Previous)
	if err != nil {
		return err
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

	return err
}
