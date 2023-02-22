package app

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
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
		return ConfigFromReader(bytes.NewReader(data))
	}

	return config, err
}

func ConfigFromReader(r io.Reader) (Config, error) {
	var config Config
	var err error

	if r == nil {
		return config, errors.New("could not load configuration from reader, reader was nil")
	}

	data, _ := io.ReadAll(r)

	err = json.Unmarshal(data, &config)

	return config, err
}
