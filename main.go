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

		err := command.callback(cfg)
		if err != nil {
			fmt.Printf("%v\n", err)
		}

	}
}
