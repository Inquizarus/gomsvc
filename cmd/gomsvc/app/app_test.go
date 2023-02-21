package app_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/inquizarus/gomsvc/cmd/gomsvc/app"
	"github.com/inquizarus/gomsvc/pkg/logging"
	"github.com/inquizarus/rwapper/v2/pkg/chiwrapper"
	"github.com/stretchr/testify/assert"
)

func TestThatRoutingWorksAsIntended(t *testing.T) {

	config := app.Config{
		Routes: []app.Route{
			{
				Name:      "test",
				Path:      "/test",
				Method:    "GET",
				Upstreams: []app.Upstream{},
				Response: app.Response{
					Headers: map[string]string{
						"x-test-header": "test",
					},
					StatusCode:              http.StatusOK,
					Body:                    "hello, world!",
					ConcatUpstreamResponses: true,
				},
			},
			{
				Name:      "upstream",
				Path:      "/upstream",
				Method:    "GET",
				Upstreams: []app.Upstream{},
				Response: app.Response{
					Headers:                 map[string]string{},
					StatusCode:              http.StatusOK,
					Body:                    "hello from upstream!",
					ConcatUpstreamResponses: false,
				},
			},
		},
	}

	router := chiwrapper.New(nil)

	server := httptest.NewServer(router)

	config.Routes[0].Upstreams = append(config.Routes[0].Upstreams, app.Upstream{
		URL:     server.URL + config.Routes[1].Path,
		Headers: map[string]string{},
		Method:  http.MethodGet,
	})

	app.RegisterRoutes(config, router, logging.DefaultLogger)

	request, err := http.NewRequest(http.MethodGet, server.URL+"/test", nil)

	assert.NoError(t, err)

	response, err := http.DefaultClient.Do(request)

	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)

	data, _ := io.ReadAll(response.Body)

	assert.Equal(t, "hello, world!\n# 200 OK "+server.URL+"/upstream\nhello from upstream!", string(data))
	assert.Equal(t, "test", response.Header.Get("x-test-header"))
}

func TestThatJSONRoutingWorksAsIntended(t *testing.T) {

	config := app.Config{
		Routes: []app.Route{
			{
				Name:      "test",
				Path:      "/test",
				Method:    "GET",
				Upstreams: []app.Upstream{},
				Response: app.Response{
					Headers: map[string]string{
						"content-type": "application/json",
					},
					StatusCode: http.StatusOK,
					Body: map[string]interface{}{
						"foo": "bar",
					},
					ConcatUpstreamResponses: true,
				},
			},
			{
				Name:      "upstream",
				Path:      "/upstream",
				Method:    "GET",
				Upstreams: []app.Upstream{},
				Response: app.Response{
					Headers: map[string]string{
						"content-type": "application/json",
					},
					StatusCode: http.StatusOK,
					Body: map[string]interface{}{
						"fizz": "buzz",
					},
					ConcatUpstreamResponses: false,
				},
			},
		},
	}

	router := chiwrapper.New(nil)

	server := httptest.NewServer(router)

	config.Routes[0].Upstreams = append(config.Routes[0].Upstreams, app.Upstream{
		URL:     server.URL + config.Routes[1].Path,
		Headers: map[string]string{},
		Method:  http.MethodGet,
	})

	app.RegisterRoutes(config, router, logging.DefaultLogger)

	request, err := http.NewRequest(http.MethodGet, server.URL+"/test", nil)

	assert.NoError(t, err)

	response, err := http.DefaultClient.Do(request)

	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)

	data, _ := io.ReadAll(response.Body)

	assert.Equal(t, `{"foo":"bar","upstreams":[{"body":{"fizz":"buzz"},"url":"`+server.URL+`/upstream"}]}`, string(data))
}

func TestThatMethodNotAllowedWorksForRoute(t *testing.T) {

	config := app.Config{
		Routes: []app.Route{
			{
				Name:      "test",
				Path:      "/test",
				Method:    "POST",
				Upstreams: []app.Upstream{},
				Response: app.Response{
					Headers:                 map[string]string{},
					StatusCode:              http.StatusOK,
					Body:                    "hello, world!",
					ConcatUpstreamResponses: false,
				},
			},
		},
	}

	router := chiwrapper.New(nil)

	app.RegisterRoutes(config, router, logging.DefaultLogger)

	server := httptest.NewServer(router)

	request, err := http.NewRequest(http.MethodGet, server.URL+"/test", nil)

	assert.NoError(t, err)

	response, err := http.DefaultClient.Do(request)

	assert.NoError(t, err)

	assert.Equal(t, http.StatusMethodNotAllowed, response.StatusCode)

	data, _ := io.ReadAll(response.Body)

	assert.Equal(t, "", string(data))
}
