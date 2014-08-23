package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Database struct {
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"database_name"`
}

// Config represents the configuration information.
type Config struct {
	AppURL   string   `json:"app_url"`
	Database Database `json:"database"`
}

var Conf Config

func init() {
	// Get the config file
	config_file, err := ioutil.ReadFile("./config.json")
	if err != nil {
		fmt.Printf("File error: %v\n", err)
	}
	json.Unmarshal(config_file, &Conf)
}
