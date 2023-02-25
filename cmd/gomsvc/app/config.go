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

	if err == nil {
		routes, err := LoadRoutesFromDir()
		if err == nil {
			config.Routes = append(config.Routes, routes...)
		}
	}
	return config, err
}

func LoadRoutesFromDir() ([]Route, error) {
	routes := []Route{}
	dir := os.Getenv(envKeyRoutesDir)
	if dir != "" {
		entries, err := os.ReadDir(dir)
		if err != nil {
			return nil, err
		}
		for _, file := range entries {
			data, err := os.ReadFile(dir + "/" + file.Name())
			if err != nil {
				return nil, err
			}
			route := Route{}
			err = json.Unmarshal(data, &route)
			if err != nil {
				return nil, err
			}
			routes = append(routes, route)
		}
	}
	return routes, nil
}
