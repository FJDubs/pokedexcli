package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/FJDubs/pokedexcli/internal/pokeapi"
)

var commandList = map[string]cliCommand{}
var Config = config{}

func init() {
	commandList = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Display a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays the names of 20 location areas in the Pokemon world. Each subsequent call to map displays the next 20 locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the names of the previous 20 location areas in the Pokemon world. Each subsequent call to map displays the previous 20 locations",
			callback:    commandMapB,
		},
		"explore": {
			name:        "explore",
			description: "It takes the name of a location area as an argument, to provide a list of Pokemon in the area",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Throw a pokeball to attempt catch the pokemon you name",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Get details on a pokemon that you have already caught",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "List Pokemon you have caught",
			callback:    commandPokedex,
		},
	}
	Config = config{
		Next:          "https://pokeapi.co/api/v2/location-area/?offset=00",
		Previous:      "",
		UserArgs:      []string{},
		CaughtPokemon: make(map[string]pokeapi.Pokemon),
	}
}

func startRepl() {
	scnr := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scnr.Scan()
		usrIn := cleanInput(scnr.Text())
		command, ok := commandList[usrIn[0]]
		if ok {
			if len(usrIn) > 1 {
				Config.UserArgs = usrIn[1:]
			}
			err := command.callback(&Config)
			if err != nil {
				fmt.Printf("error: %s\n", err)
			}
		} else {
			fmt.Println("Unknown command")
		}

	}
}

func cleanInput(text string) []string {
	words := strings.Fields(strings.ToLower(text))
	return words
}
