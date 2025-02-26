package main

import (
	"errors"
	"fmt"
	"os"
)

type commandREPL struct {
	name        string
	description string
	callback    func(*config) error
}

func commands() map[string]commandREPL {
	return map[string]commandREPL{
		// TODO: add quit
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"quit": {
			name:        "quit",
			description: "Quit Pokedex",
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
			callback:    commandMapf,
		},
		"mapb": {
			name:        "map",
			description: "Displays a list of locations",
			callback:    commandMapb,
		},
	}
}

func commandExit(cfg *config) error {
	fmt.Printf("Closing the Pokedex... Goodbye!\n")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config) error {
	fmt.Print(`Welcome to the Pokedex!
Usage:

help: Displays a help message
exit: Exit the Pokedex
quit: Quit Pokedex
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

func commandMapf(cfg *config) error {
	locationsResp, err := cfg.pokeapiClient.ListLocations(cfg.pageNext)
	if err != nil {
		return err
	}

	cfg.pageNext = locationsResp.Next
	cfg.pagePrev = locationsResp.Previous

	for _, location := range locationsResp.Results {
		fmt.Println(location.Name)
	}

	return nil
}

func commandMapb(cfg *config) error {
	if cfg.pagePrev == nil {
		return errors.New("you're on the first page")
	}

	locationsResp, err := cfg.pokeapiClient.ListLocations(cfg.pagePrev)
	if err != nil {
		return err
	}

	cfg.pageNext = locationsResp.Next
	cfg.pagePrev = locationsResp.Previous

	for _, location := range locationsResp.Results {
		fmt.Println(location.Name)
	}

	return nil
}
