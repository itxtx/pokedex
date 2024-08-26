package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/itxtx/pokedex/pokeapi"
)

type cliCmd struct {
	name        string
	description string
	callback    func(*pokeapi.Config) error
}

func commandHelp(commands map[string]cliCmd) func(*pokeapi.Config) error {
	return func(config *pokeapi.Config) error {
		fmt.Println("Available commands:")
		for _, cmd := range commands {
			fmt.Printf("%s - %s\n", cmd.name, cmd.description)
		}
		return nil
	}
}

func commandExit() func(*pokeapi.Config) error {
	return func(config *pokeapi.Config) error {
		fmt.Println("Exiting...")
		os.Exit(0)
		return nil
	}
} // commandMap fetches and displays the next 20 location areas
func commandMap(config *pokeapi.Config) error {
	url := config.NextURL
	if url == "" {
		url = pokeapi.BaseURL + "location-area"
	}

	response, err := pokeapi.FetchLocationAreas(url)
	if err != nil {
		return fmt.Errorf("failed to fetch location areas: %v", err)
	}

	// Update the config with the new pagination URLs
	if response.Next != nil {
		config.NextURL = *response.Next
	}
	if response.Previous != nil {
		config.PreviousURL = *response.Previous
	}

	// Display the location areas
	fmt.Println("Location Areas:")
	for _, area := range response.Results {
		fmt.Println("- " + area.Name)
	}

	return nil
}

// commandMapBack fetches and displays the previous 20 location areas
func commandMapBack(config *pokeapi.Config) error {
	if config.PreviousURL == "" {
		return fmt.Errorf("no previous locations to display")
	}

	response, err := pokeapi.FetchLocationAreas(config.PreviousURL)
	if err != nil {
		return fmt.Errorf("failed to fetch location areas: %v", err)
	}

	// Update the config with the new pagination URLs
	if response.Next != nil {
		config.NextURL = *response.Next
	}
	if response.Previous != nil {
		config.PreviousURL = *response.Previous
	}

	// Display the location areas
	fmt.Println("Location Areas:")
	for _, area := range response.Results {
		fmt.Println("- " + area.Name)
	}

	return nil
}

// scanner processes user input
func Scanner(s *bufio.Scanner, commands map[string]cliCmd, config *pokeapi.Config) {
	for s.Scan() {
		text := s.Text()
		if cmd, found := commands[text]; found {
			if err := cmd.callback(config); err != nil {
				fmt.Println("Error:", err)
			}
		} else {
			fmt.Printf("Pokedex > ")
		}
	}
	if err := s.Err(); err != nil {
		fmt.Println("Error reading input:", err)
	}
}

func main() {
	// Initialize the command config
	config := &pokeapi.Config{}

	// Define the available commands
	commands := map[string]cliCmd{
		"exit": {
			name:        "exit",
			description: "exit pokedex",
			callback:    commandExit(),
		},
		"map": {
			name:        "map",
			description: "displays the next 20 location areas",
			callback:    commandMap,
		},
		"mapb": {
			name:        "map back",
			description: "displays the previous 20 location areas",
			callback:    commandMapBack,
		},
	}
	commands["help"] = cliCmd{
		name:        "help",
		description: "displays help message",
		callback:    commandHelp(commands),
	}

	fmt.Printf("Pokedex > ")

	scanner := bufio.NewScanner(os.Stdin)
	Scanner(scanner, commands, config)
}
