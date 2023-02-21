package main

import (
	"os"

	"github.com/inquizarus/gomsvc/cmd/gomsvc/app"
	"github.com/inquizarus/gomsvc/pkg/logging"
)

const (
	envKeyLogLevel  = "GOMSVC_LOG_LEVEL"
	defaultLogLevel = "info"
)

func main() {
	logLevel := os.Getenv(envKeyLogLevel)

	if logLevel == "" {
		logLevel = defaultLogLevel
	}

	app.Run(nil, logging.NewLogrusLogger(nil, logLevel, nil))
}
