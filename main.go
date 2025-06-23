package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Print("Hello, World!")
}

func cleanInput(text string) []string {
	words := strings.Fields(text)
	return words
}
