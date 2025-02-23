package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func inputClean(text string) []string {
	textLowered := strings.ToLower(text)
	words := strings.Fields(textLowered)
	return words
}

func main() {
	pageState := pageState{}
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()

		words := inputClean(scanner.Text())
		if len(words) == 0 {
			continue
		}

		commandName := words[0]
		command, ok := commands()[commandName]
		if !ok {
			fmt.Print("Unknown command\n")
			continue
		}

		err := command.callback(&pageState)
		if err != nil {
			fmt.Printf("%v\n", err)
		}

	}
}
