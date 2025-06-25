package main

import (
	"fmt"
	"math/rand/v2"
	"os"

	"github.com/FJDubs/pokedexcli/internal/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

type config struct {
	Next          string
	Previous      string
	UserArgs      []string
	CaughtPokemon map[string]pokeapi.Pokemon
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

func commandExplore(Conf *config) error {
	exploreLocation := Conf.UserArgs[0]
	searchUrl := "https://pokeapi.co/api/v2/location-area/" + exploreLocation + "/"
	fmt.Printf("Exploring %s...\n", Conf.UserArgs[0])
	pkmn, err := pokeapi.ListPokemonAt(searchUrl)
	if err != nil {
		return err
	}
	if len(pkmn) == 0 {
		fmt.Println("No pokemon Found")
	} else {
		fmt.Println("Found Pokemon:")
		for _, name := range pkmn {
			fmt.Printf("- %s\n", name)
		}
	}

	return err
}

func commandCatch(Conf *config) error {
	pkmnName := Conf.UserArgs[0]
	const highestPokemonBaseExperience = 608
	searchUrl := "https://pokeapi.co/api/v2/pokemon/" + pkmnName
	pkmn, err := pokeapi.GetPokemonInfo(searchUrl)
	if err != nil {
		return err
	}
	fmt.Printf("Throwing a Pokeball at %s...\n", pkmnName)
	catchRate := 70 - ((pkmn.BaseExperience * 40) / highestPokemonBaseExperience)
	randNum := rand.IntN(100)
	if randNum <= catchRate {
		fmt.Printf("%s was caught!\n", pkmnName)
		fmt.Println("You may now inspect it with the inspect command.")
		Conf.CaughtPokemon[pkmnName] = pkmn
	} else {
		fmt.Printf("%s escaped\n", pkmnName)
	}

	return nil
}

func commandInspect(Conf *config) error {
	pkmnName := Conf.UserArgs[0]
	pkmn, ok := Conf.CaughtPokemon[pkmnName]
	if !ok {
		fmt.Println("You have not caught that Pokemon")
		return nil
	}
	fmt.Printf("Name: %s\n", pkmn.Name)
	fmt.Printf("Height: %v\n", pkmn.Height)
	fmt.Printf("Weight: %v\n", pkmn.Weight)
	fmt.Println("Stats:")
	for _, stat := range pkmn.Stats {
		fmt.Printf(" -%s: %v\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, t := range pkmn.Types {
		fmt.Printf(" -%s\n", t.Type.Name)
	}
	return nil
}

func commandPokedex(Conf *config) error {
	if len(Conf.CaughtPokemon) < 1 {
		fmt.Println("You have not caught any Pokemon")
		return nil
	}
	fmt.Println("Your Pokedex:")
	for pkmn := range Conf.CaughtPokemon {
		fmt.Printf(" - %s\n", pkmn)
	}
	return nil
}
