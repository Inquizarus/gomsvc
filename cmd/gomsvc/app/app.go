package app

import (
	"net/http"
	"os"

	"github.com/inquizarus/gomsvc/pkg/logging"
	"github.com/inquizarus/rwapper/v2"
	"github.com/inquizarus/rwapper/v2/pkg/servemuxwrapper"
)

func Run(router rwapper.RouterWrapper, log logging.Logger) {

	if log == nil {
		log = logging.DefaultLogger
	}

	configPath := os.Getenv(envKeyConfigPath)

	if configPath == "" {
		configPath = configPathDefault
	}

	config, err := ConfigFromFilePath(configPath)

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

func RegisterRoutes(config Config, router rwapper.RouterWrapper, log logging.Logger) {
	for _, route := range config.Routes {
		log.Info("adding route " + route.Name)
		router.HandlerFunc(route.Method, route.Path, MakeHandlerFunc(route, config, log))
	}
}
