package utils

import (
	"fmt"
	"log"
	"os"
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
	// set env variables
	os.Setenv("PORT", config.Port)
	os.Setenv("HTTPS", fmt.Sprintf("%t", config.HTTPS))
	os.Setenv("CERT", config.Cert)
	os.Setenv("KEY", config.Key)
	os.Setenv("PASSWORD", config.Password)
	os.Setenv("SECRET", config.Secret)

}
