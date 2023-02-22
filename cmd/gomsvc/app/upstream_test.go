package app_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/inquizarus/gomsvc/cmd/gomsvc/app"
	"github.com/inquizarus/rwapper/v2/pkg/servemuxwrapper"
	"github.com/stretchr/testify/assert"
)

func TestThatPlainTextUpstreamPostWorksAsIntended(t *testing.T) {
	router := servemuxwrapper.New(nil)
	called := false

	expectedBody := "hello, world!"

	router.Handle(http.MethodPost, "/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		payload, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		assert.Equal(t, expectedBody, string(payload))
	}))

	server := httptest.NewServer(router)

	client := server.Client()

	upstream := app.Upstream{
		URL:    server.URL,
		Method: http.MethodPost,
		Body:   expectedBody,
	}

	upstream.Call(client)

	assert.True(t, called)
}

func TestThatJSONUpstreamPostWorksAsIntended(t *testing.T) {
	router := servemuxwrapper.New(nil)
	called := false

	expectedBody := map[string]interface{}{
		"foo": "bar",
	}

	router.Handle(http.MethodPost, "/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		payload, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		expectedPayload, _ := json.Marshal(expectedBody)
		assert.Equal(t, string(expectedPayload), string(payload))
	}))

	server := httptest.NewServer(router)

	client := server.Client()

	upstream := app.Upstream{
		URL:    server.URL,
		Method: http.MethodPost,
		Headers: map[string]string{
			"content-type": "application/json",
		},
		Body: expectedBody,
	}

	upstream.Call(client)

	assert.True(t, called)
}
