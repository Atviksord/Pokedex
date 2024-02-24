package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*Config) error
}
type Config struct {
	Next     string
	Previous string
}
type Response struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func commandHelp(c *Config) error {
	commandmap := getCommands()
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("")

	fmt.Println("Usage: ")
	for _, cmd := range commandmap {

		fmt.Printf("\nCommand: %s", cmd.name)
		fmt.Printf("\nDescription: %s\n", cmd.description)

	}
	fmt.Println()

	return nil
}
func commandExit(c *Config) error {
	os.Exit(0)
	return nil
}
func apiCallHandler(direction bool, c *Config) error {
	var getcall string
	if c.Next == "" && c.Previous == "" {
		getcall = "https://pokeapi.co/api/v2/location-area/"
	} else if c.Next != "" && !direction {
		getcall = c.Next
	} else if c.Previous != "" && direction {
		getcall = c.Previous
	}
	response, err := http.Get(getcall)
	if err != nil {
		log.Fatal(err)
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var responseStruct Response
	err = json.Unmarshal(body, &responseStruct)
	if err != nil {
		fmt.Println(err)
	}
	c.Next = responseStruct.Next
	c.Previous = responseStruct.Previous

	for _, result := range responseStruct.Results {
		fmt.Println(result.Name) // prints each location's name and URL
	}
	defer response.Body.Close()
	return nil

}
func commandMap(c *Config) error {
	apiCallHandler(false, c)

	return nil
}
func commandMapb(c *Config) error {
	if c.Previous != "" {
		apiCallHandler(true, c)
	} else {
		fmt.Println("Nothing to go back to")
	}
	return nil
}
func getCommands() map[string]cliCommand {

	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "displays the names of 20 location areas in the Pokemon world",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "displays the previous 20 locations in the Pokemon world",
			callback:    commandMapb,
		},
	}
}

func main() {
	commandmap := getCommands()
	var c Config

	for {
		fmt.Print("Pokedex > ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		line := scanner.Text()

		if cmd, ok := commandmap[line]; ok {
			err := cmd.callback(&c)
			if err != nil {
				fmt.Println("There was an error executing the commando:", err)
			}
		}
	}

}
