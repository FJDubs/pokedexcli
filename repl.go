package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
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
	}
	Config = config{
		Next:     "https://pokeapi.co/api/v2/location-area/?offset=00",
		Previous: "",
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
			err := command.callback(&Config)
			if err != nil {
				fmt.Printf("error: %s", err)
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
