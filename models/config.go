package models

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	HTTPS    bool
	Cert     string
	Key      string
	Port     string
	Password string
	Secret   string
}

func (c *Config) SaveConfig() error {
	f, err := os.Create("./config.json")
	if err != nil {
		return err
	}
	defer f.Close()
	json, _ := json.MarshalIndent(c, "", " ")
	_, err = f.WriteString(string(json))
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) String() string {
	return fmt.Sprintf(
		`{
	HTTPS: %t
	Cert: %s
	Key: %s
	Port: %s
	Password: %s
	Secret: %s
}`, c.HTTPS, c.Cert, c.Key, c.Port, c.Password, c.Secret)
}
