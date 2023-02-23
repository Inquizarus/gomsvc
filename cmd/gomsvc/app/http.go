package app

import (
	"net/http"

	"github.com/inquizarus/gomsvc/pkg/logging"
)

func MakeHandlerFunc(route Route, config Config, log logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		log.Info("starting to handle request to route " + route.Name)

		// Initial checking to determine if the incoming request is a valid one according
		// to the route configuration. Usually this is already handled by a router.

		if r.Method != route.Method {
			w.WriteHeader(http.StatusMethodNotAllowed)
			log.Info("could not finish handling for request to " + route.Name + " wrong HTTP method " + r.Method)
			return
		}

		// Lets handle all potential upstreams

		upstreamResponses := []*http.Response{}
		for _, upstream := range route.Upstreams {
			upstreamResponse, err := upstream.Call(http.DefaultClient)
			if err != nil {
				log.Info("error when performing upstream request " + err.Error() + ", skipping to next upstream call")
				continue
			}
			if route.Response.IncludeUpstreamResponses {
				upstreamResponses = append(upstreamResponses, upstreamResponse)
			}
		}

		for k, v := range route.Response.Headers {
			w.Header().Set(k, v)
		}

		w.WriteHeader(route.Response.StatusCode)

		data, err := route.Response.Content(r, upstreamResponses)

		if nil != err {
			log.Error(err)
		}

		w.Write(data)

		log.Info("finished handling request to route " + route.Name)

	}
}
