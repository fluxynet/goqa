package main

import (
	"encoding/json"
	"os"
)

const (
	// configFilename is the name of the configuration file
	configFilename = "config.json"
)

// Config contains configuration variables for the app
type Config struct {
	ServerHost       string   `json:"server_host"`
	EmailHost        string   `json:"email_host"`
	EmailUsr         string   `json:"email_usr"`
	EmailPass        string   `json:"email_pass"`
	EmailFrom        string   `json:"email_from"`
	EmailSubscribers []string `json:"email_subscribers"`
	GithubSignature  string   `json:"github_signature"`
	GithubToken      string   `json:"github_token"`
}

// LoadConf from a file named config.json placed in the same directory; bleh
func LoadConf() (*Config, error) {
	var b, err = os.ReadFile(configFilename)
	if err != nil {
		return nil, err
	}

	var conf Config
	err = json.Unmarshal(b, &conf)
	if err != nil {
		return nil, err
	}

	return &conf, nil
}
