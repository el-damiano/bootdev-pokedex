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

type commandREPL struct {
	name        string
	description string
	callback    func() error
}

func commandExit() error {
	fmt.Printf("Closing the Pokedex... Goodbye!\n")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Print(`Welcome to the Pokedex!
Usage:

help: Displays a help message
exit: Exit the Pokedex
`)
	return nil
}

func main() {

	commands := map[string]commandREPL{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Print the Help section",
			callback:    commandHelp,
		},
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()

		words := inputClean(scanner.Text())
		if len(words) == 0 {
			continue
		}

		commandName := words[0]
		command, ok := commands[commandName]
		if !ok {
			fmt.Print("Unknown command\n")
			continue
		}

		err := command.callback()
		if err != nil {
			fmt.Printf("%v", err)
		}

	}

}
