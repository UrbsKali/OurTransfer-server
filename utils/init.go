package utils

import (
	"fmt"
	"log"
	"os"

	"urbskali/file/state"
)

func StartUp() {
	fmt.Println("Starting up...")
	// check if the files directory exists
	if _, err := os.Stat("./files"); os.IsNotExist(err) {
		log.Fatal("Files directory does not exist, try running 'file setup'")
	}
	// check if the ui directory exists
	if _, err := os.Stat("./ui"); os.IsNotExist(err) {
		log.Fatal("UI directory does not exist, try running 'file setup'")
	}
	// check if the config file exists
	if _, err := os.Stat("./config.json"); os.IsNotExist(err) {
		log.Fatal("Config file does not exist, try running 'file setup'")
	}
	// load the config
	config, err := LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	// print the config
	fmt.Println(config.String())
	// set the config
	state.Config = config
}
