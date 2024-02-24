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
	callback    func() error
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

func commandHelp() error {
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
func commandExit() error {
	os.Exit(0)
	return nil
}
func apiCallHandler() error {
	response, err := http.Get("https://pokeapi.co/api/v2/location-area/")
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
	for _, result := range responseStruct.Results {
		fmt.Println(result.Name, result.URL) // prints each location's name and URL
	}
	defer response.Body.Close()
	return nil

}
func commandMap() error {
	apiCallHandler()
	return nil
}
func commandMapb() error {
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
	for {
		fmt.Print("Pokedex > ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		line := scanner.Text()

		if cmd, ok := commandmap[line]; ok {
			err := cmd.callback()
			if err != nil {
				fmt.Println("There was an error executing the commando:", err)
			}
		}
	}

}
