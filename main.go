package main

import (
	"fmt"
	"strings"
)

func inputClean(text string) []string {
	result := []string{}

	words := strings.FieldsFunc(
		strings.ToLower(text),
		func(r rune) bool {
			return r == ' ' || r == '\n'
		})

	for _, word := range words {
		if word == "" || word == "\n" {
			continue
		}
		result = append(result, word)
	}

	return result
}

func main() {
	fmt.Println("Hello, World!")
}
