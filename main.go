package main

import (
	"bufio"
	"fmt"
	"os"
)

type cliCmd struct {
	name        string
	description string
	callback    func() error
}

func Scanner(s *bufio.Scanner, commands map[string]cliCmd) {
	for s.Scan() {
		text := s.Text()
		if cmd, found := commands[text]; found {
			if err := cmd.callback(); err != nil {
				fmt.Println("Error:", err)
			}
		} else {
			fmt.Println("Unknown command:", text)
		}
	}
	if err := s.Err(); err != nil {
		fmt.Println("Error reading input:", err)
	}
}

func commandHelp() error {
	fmt.Println("Available commands:")
	fmt.Println("help - displays help message")
	fmt.Println("exit - exit pokedex")
	return nil
}

func commandExit() error {
	fmt.Println("Exiting...")
	os.Exit(0)
	return nil
}

func main() {
	commands := map[string]cliCmd{
		"help": {
			name:        "help",
			description: "displays help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "exit pokedex",
			callback:    commandExit,
		},
	}

	fmt.Println("Pokedex > ")

	// Create a new scanner for standard input
	scanner := bufio.NewScanner(os.Stdin)

	// Pass the scanner and commands map to the scanner function
	Scanner(scanner, commands)
}
