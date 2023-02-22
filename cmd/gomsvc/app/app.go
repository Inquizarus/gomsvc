package app

import (
	"net/http"
	"os"
	"strings"

	"github.com/inquizarus/gomsvc/pkg/logging"
	"github.com/inquizarus/rwapper/v2"
	"github.com/inquizarus/rwapper/v2/pkg/servemuxwrapper"
)

func Run(router rwapper.RouterWrapper, log logging.Logger) {

	if log == nil {
		log = logging.DefaultLogger
	}

	config, err := config()

	if err != nil {
		panic(err)
	}

	if router == nil {
		router = servemuxwrapper.New(nil)
	}

	RegisterRoutes(config, router, log)

	server := http.Server{
		Addr:    config.Address(),
		Handler: router,
	}

	log.Info("starting server on " + server.Addr)

	if err := server.ListenAndServe(); err != nil {
		log.Info(err)
	}
}

func config() (Config, error) {

	configPath := os.Getenv(envKeyConfigPath)
	configString := os.Getenv(envKeyConfigString)

	if configPath == "" && configString == "" {
		configPath = configPathDefault
	}

	if configPath == "" && configString != "" {
		return ConfigFromReader(strings.NewReader(configString))
	}

	return ConfigFromFilePath(configPath)
}

func RegisterRoutes(config Config, router rwapper.RouterWrapper, log logging.Logger) {
	for _, route := range config.Routes {
		log.Info("adding route " + route.Name)
		router.HandlerFunc(route.Method, route.Path, MakeHandlerFunc(route, config, log))
	}
}
