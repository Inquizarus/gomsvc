package app

import (
	"encoding/json"
	"os"
)

type Config struct {
	Port   string  `json:"port"`
	Routes []Route `json:"routes"`
}

func (c Config) Address() string {
	port := c.Port
	if port == "" {
		port = defaultPort
	}
	return ":" + port
}

func ConfigFromFilePath(path string) (Config, error) {
	var config Config

	data, err := os.ReadFile(path)

	if err == nil {
		json.Unmarshal(data, &config)
	}

	return config, err
}
