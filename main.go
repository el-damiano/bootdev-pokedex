package main

import (
	"fmt"
	"strings"
)

func inputClean(text string) []string {
	textLowered := strings.ToLower(text)
	words := strings.Fields(textLowered)
	return words
}

func main() {
	fmt.Println("Hello, World!")
}
