package utils

import (
	"crypto/sha256"
	"fmt"
	"os"
	"os/exec"
	"urbskali/file/models"
)

func Setup() {
	config := models.Config{}
	fmt.Println("Initializing...")
	// create ./files directory
	if _, err := os.Stat("./files"); os.IsNotExist(err) {
		err := os.Mkdir("./files", 0755)
		if err != nil {
			fmt.Println("Failed to create files directory")
		}
		fmt.Println("Files directory created successfully")
	} else {
		fmt.Println("Files directory already exists")
	}

	// ask for password, hash it and save it to an environment variable
	fmt.Println("Please enter a password for the API:")
	var password string
	fmt.Scanln(&password)
	// hash the password
	h := sha256.New()
	h.Write([]byte(string(password)))
	// Convert the sha256 hash to a string
	secret := fmt.Sprintf("%x", h.Sum(nil))
	config.Secret = secret
	config.Password = password

	// ask if want to use https or not
	fmt.Println("Do you want to use HTTPS? (y/n)")
	var https string
	fmt.Scanln(&https)
	if https == "y" {
		// ask if wants to gen a self-signed cert or use a custom one
		fmt.Println("Do you want to generate a self-signed certificate? (y/n)")
		var selfSigned string
		fmt.Scanln(&selfSigned)
		if selfSigned == "y" {
			// generate a self-signed cert
			fmt.Println("Generating self-signed certificate...")
			cert, key, err := generateSelfSignedCert()
			if err != nil {
				fmt.Println("Failed to generate self-signed certificate")
			} else {
				fmt.Println("Self-signed certificate generated successfully")
				config.Cert = cert
				config.Key = key
			}
		} else {
			// ask for the cert path
			fmt.Println("Please enter the path to the certificate (relative):")
			var certPath string
			fmt.Scanln(&certPath)
			// ask for the key path
			fmt.Println("Please enter the path to the key:")
			var keyPath string
			fmt.Scanln(&keyPath)
			// save the cert path to an environment variable
			config.Cert = certPath
			config.Key = keyPath
		}
		fmt.Println("HTTPS enabled")
		config.HTTPS = true
	} else {
		fmt.Println("HTTPS disabled")
		config.HTTPS = false
	}

	// ask if wants to use a custom port
	fmt.Println("Do you want to use a custom port? (blank for default)")
	var customPort string
	fmt.Scanln(&customPort)
	if customPort != "" {
		// save the custom port to an environment variable
		config.Port = customPort
		fmt.Println("Custom port set to " + customPort)
	} else {
		if config.HTTPS {
			config.Port = "443"
		} else {
			config.Port = "80"
		}
		fmt.Println("Using default port")
	}

	//check if the ui directory exists
	if _, err := os.Stat("./ui"); os.IsNotExist(err) {
		fmt.Println("UI directory does not exist, building it...")
		BuildUI()
	} else {
		fmt.Println("UI directory already exists")
	}

	// save the config to a file
	config.SaveConfig()

	fmt.Println("Initialization complete")
}

func generateSelfSignedCert() (string, string, error) {
	// check if openssl is installed
	_, err := exec.LookPath("openssl")
	if err != nil {
		return "", "", err
	}
	// generate the cert
	cmd := exec.Command("openssl", "req", "-x509", "-newkey", "rsa:4096", "-keyout", "./cert/key.pem", "-out", "./cert/cert.pem", "-days", "365", "-nodes", "-subj", "/CN=localhost")
	err = cmd.Run()
	if err != nil {
		return "", "", err
	}
	return "./cert/cert.pem", "/cert/key", nil
}
