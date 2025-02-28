package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/el-damiano/bootdev-pokedex/internal/pokeapi"
)

type config struct {
	pokeapiClient *pokeapi.Client
	pageNext      *string
	pagePrev      *string
	pokedex       map[string]pokeapi.Pokemon
}

func inputClean(text string) []string {
	textLowered := strings.ToLower(text)
	words := strings.Fields(textLowered)
	return words
}

func main() {
	pokeClient := pokeapi.NewClient(10*time.Second, 5*time.Minute)
	cfg := &config{
		pokeapiClient: &pokeClient,
		pokedex:       map[string]pokeapi.Pokemon{},
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
		command, ok := commands()[commandName]
		if !ok {
			fmt.Print("Unknown command\n")
			continue
		}

		args := []string{}
		if len(words) > 1 {
			args = words[1:]
		}

		err := command.callback(cfg, args...)
		if err != nil {
			fmt.Printf("%v\n", err)
		}

	}
}
