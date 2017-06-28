package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// Config - data structure for configuration
type Config struct {
	Vars   map[string]string `json:"vars"`
	Parsed bool
}

var config Config

// parseConfig - parses FLOWCONFIG=./config.json
func parseConfig() {
	configFile := os.Getenv("FLOWCONFIG")
	config = Config{}

	if configFile == "" {
		panic("Configuration file path is not defined. FLOWCONFIG environment variable must be set to config.json. Ex: export FLOWCONFIG=$GOPATH/config.json")
	}

	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(file, &config)
	if err != nil {
		panic(fmt.Sprintf("Error reading configuration file: %s", err))
	}
	config.Parsed = true
}

// GetConfigVar - gets configuration setting defined in CONFIG=./config.json
func GetConfigVar(key string) string {
	if !config.Parsed {
		parseConfig()
	}

	val, ok := config.Vars[key]
	if !ok {
		return ""
	}
	return val
}
