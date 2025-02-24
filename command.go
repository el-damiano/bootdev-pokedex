package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type pageState struct {
	Page string
}

type commandREPL struct {
	name        string
	description string
	callback    func(*pageState) error
}

func commands() map[string]commandREPL {
	return map[string]commandREPL{
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
			callback:    commandMapForwards,
		},
		"mapb": {
			name:        "map",
			description: "Displays a list of locations",
			callback:    commandMapBackwards,
		},
	}
}

func commandExit(pageState *pageState) error {
	fmt.Printf("Closing the Pokedex... Goodbye!\n")
	os.Exit(0)
	return nil
}

func commandHelp(pageState *pageState) error {
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

func commandMapForwards(pageState *pageState) error {
	commandMap(pageState, 20, 20)
	return nil
}

func commandMapBackwards(pageState *pageState) error {
	commandMap(pageState, -20, 20)
	return nil
}

func commandMap(pageState *pageState, offset, limit int) error {
	pageState.Page = locationURL(pageState.Page, offset, limit)

	res, err := http.Get(pageState.Page)
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
