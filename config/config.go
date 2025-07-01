package config

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"time"
	"upload/cli"
)

//go:embed config.json
var configData []byte

func NewConfig() *cli.Config {
	var config struct {
		BaseURL string `json:"baseURL"`
		Timeout string `json:"timeout"`
		Token   string `json:"token"`
	}

	if err := json.Unmarshal(configData, &config); err != nil {
		panic(fmt.Sprintf("failed to parse embedded config.json: %v", err))
	}

	if config.BaseURL == "" {
		panic("baseURL is required in config.json")
	}

	if config.Timeout == "" {
		panic("timeout is required in config.json")
	}

	if config.Token == "" {
		panic("token is required in config.json")
	}

	timeout, err := time.ParseDuration(config.Timeout)
	if err != nil {
		panic(fmt.Sprintf("invalid timeout format in config.json: %v", err))
	}

	return &cli.Config{
		BaseURL: config.BaseURL,
		Timeout: timeout,
		Token:   config.Token,
	}
}
