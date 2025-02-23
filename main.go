package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
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
map: Displays a list of map locations
`)
	return nil
}

type LocationsPage struct {
	Count     int    `json:"count"`
	Next      string `json:"next"`
	Previous  string `json:"previous"`
	Locations []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func commandMap() error {
	urlFull := "https://pokeapi.co/api/v2/location-area"

	res, err := http.Get(urlFull)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("non-OK HTTP status: %s", res.Status)
	}

	var locationsPage LocationsPage
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&locationsPage)
	if err != nil {
		return err
	}

	for _, location := range locationsPage.Locations {
		fmt.Printf("%v\n", location.Name)
	}

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
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays a list of locations",
			callback:    commandMap,
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
			fmt.Printf("%v\n", err)
		}

	}

}
