package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var commandList = map[string]cliCommand{}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Printf("Usage:\n\n")
	for _, cmd := range commandList {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

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
			err := command.callback()
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

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}
