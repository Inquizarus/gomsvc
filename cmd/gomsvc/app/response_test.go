package app_test

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/inquizarus/gomsvc/cmd/gomsvc/app"
	"github.com/inquizarus/gomsvc/internal/pkg/httptools"
	"github.com/stretchr/testify/assert"
)

func TestResponseContentJSON(t *testing.T) {
	headers := map[string]string{"content-type": "application/json"}
	body := map[string]interface{}{
		"name": "John",
		"age":  30,
	}
	resp := app.Response{
		Headers: headers,
		Body:    body,
	}
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", bytes.NewBufferString("hello, world"))
	req.Header.Set("content-type", "application/json")
	upstreamResponses := []*http.Response{}
	expected, _ := httptools.FormatJSON(body)

	result, err := resp.Content(req, upstreamResponses)
	if err != nil {
		t.Errorf("Error calling Content method: %v", err)
	}
	if !bytes.Equal(result, expected) {
		t.Errorf("Returned byte array does not match expected byte array:\nExpected: %s\nReturned: %s", string(expected), string(result))
	}
}

func TestJSONWithUpstream(t *testing.T) {
	body := map[string]interface{}{
		"name": "John",
		"age":  30,
	}
	url, _ := url.Parse("http://example.com")
	upstreamResponse := &http.Response{
		Status:     "200 OK",
		StatusCode: http.StatusOK,
		Header:     http.Header{},
		Body:       io.NopCloser(strings.NewReader(`{"city": "New York"}`)),
		Request: &http.Request{
			URL: url,
		},
	}
	upstreamResponse.Header.Set("content-type", "application/json")
	upstreamResponses := []*http.Response{upstreamResponse}
	response := app.Response{
		Headers:                   map[string]string{"content-type": "application/json"},
		Body:                      body,
		IncludeUpstreamResponses:  true,
		IncludeRequestInformation: false,
	}

	expectedJSON := `{
 "age": 30,
 "name": "John",
 "upstreams": [
  {
   "body": {
    "city": "New York"
   },
   "headers": {
    "Content-Type": [
     "application/json"
    ]
   },
   "status_code": 200,
   "url": "http://example.com"
  }
 ]
}`
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	jsonBytes, err := response.Content(req, upstreamResponses)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	x, _ := httptools.FormatJSONData(jsonBytes)

	assert.Equal(t, expectedJSON, string(x))
}
