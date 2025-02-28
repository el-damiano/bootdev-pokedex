package main

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"os"
	"sort"
)

type commandREPL struct {
	name        string
	description string
	callback    func(*config, ...string) error
}

func commands() map[string]commandREPL {
	return map[string]commandREPL{
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
		"pokedex": {
			name:        "pokedex",
			description: "Display a list of caught Pokemon",
			callback:    commandPokedex,
		},
		"map": {
			name:        "map",
			description: "Get the next page of locations",
			callback:    commandMapf,
		},
		"mapb": {
			name:        "mapb",
			description: "Get the previous page of locations",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Display list of Pokemons at a specific location",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch [name]",
			description: "Attempt catching a Pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect [name]",
			description: "Display information about a caught Pokemon",
			callback:    commandInspect,
		},
		"find": {
			name:        "find [name]",
			description: "Find locations where [name] can be found",
			callback:    commandFind,
		},
	}
}

func commandExit(cfg *config, args ...string) error {
	fmt.Printf("Closing the Pokedex... Goodbye!\n")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config, args ...string) error {
	fmt.Println("Summary of Pokedex Commands")
	fmt.Println()

	commands := commands()
	keys := make([]string, len(commands))

	i := 0
	for key := range commands {
		keys[i] = key
		i++
	}

	sort.Strings(keys)
	for _, key := range keys {
		cmd := commands[key]
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}

	fmt.Println()
	return nil
}

func commandPokedex(cfg *config, args ...string) error {
	fmt.Println("Your Pokedex:")
	for _, pokemon := range cfg.pokedex {
		fmt.Printf("  - %s\n", pokemon.Name)
	}
	return nil
}

func commandMapf(cfg *config, args ...string) error {
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

func commandMapb(cfg *config, args ...string) error {
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

func commandExplore(cfg *config, args ...string) error {
	if len(args) == 0 {
		return errors.New("explore command expects a location name")
	}
	name := args[0]

	location, err := cfg.pokeapiClient.GetLocation(name)
	if err != nil {
		return err
	}

	fmt.Printf("Exploring %s...\n", location.Name)
	fmt.Println("Found Pokemon:")
	for _, encounter := range location.PokemonEncounters {
		fmt.Println(encounter.Pokemon.Name)
	}

	return nil
}

func commandCatch(cfg *config, args ...string) error {
	if len(args) == 0 {
		return errors.New("catch command expects a Pokemon name")
	}
	name := args[0]
	pokemon, err := cfg.pokeapiClient.GetPokemon(name)
	if err != nil {
		return err
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon.Name)

	maxExp := 400
	catchRate := rand.IntN(pokemon.BaseExperience) * 100 / maxExp
	if catchRate < 15 {
		cfg.pokedex[pokemon.Name] = pokemon
		fmt.Printf("%s was caught!\n", pokemon.Name)
	} else {
		fmt.Printf("%s escaped!\n", pokemon.Name)
	}

	return nil
}

func commandInspect(cfg *config, args ...string) error {
	if len(args) == 0 {
		return errors.New("catch command expects a Pokemon name")
	}
	name := args[0]
	pokemon, ok := cfg.pokedex[name]
	if !ok {
		return fmt.Errorf("%s has not been caught yet!", name)
	}

	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)

	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf("  - %s: %v\n", stat.Stat.Name, stat.BaseStat)
	}

	fmt.Println("Types:")
	for _, pokeType := range pokemon.Types {
		fmt.Printf("  - %s\n", pokeType.Type.Name)
	}

	fmt.Printf("Total moves: %d\n", len(pokemon.Moves))

	fmt.Println("Abilities:")
	for _, ability := range pokemon.Abilities {
		fmt.Printf("  - %s\n", ability.Ability.Name)
	}

	return nil
}

func commandFind(cfg *config, args ...string) error {
	if len(args) == 0 {
		return errors.New("find command expects a Pokemon name")
	}
	name := args[0]

	encounters, err := cfg.pokeapiClient.GetPokemonEncounters(name)
	if err != nil {
		return err
	}
	fmt.Printf("%s can be found at:\n", name)
	for _, enc := range encounters {
		fmt.Printf("  - %s\n", enc.LocationArea.Name)
	}

	return nil
}
