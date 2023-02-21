package app

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
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

type Route struct {
	Name      string     `json:"name"`
	Path      string     `json:"path"`
	Method    string     `json:"method"`
	Upstreams []Upstream `json:"upstreams"`
	Response  Response   `json:"response"`
}

type Upstream struct {
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Method  string            `json:"method"`
	Body    interface{}       `json:"body"`
}

// Call takes all the information in the upstream and makes a request
// based on that. URL with a value that has a env: prefix will instead use
// whatever value is in that environment variable as an URL instead for
// the upstream call
func (u Upstream) Call(client *http.Client) (*http.Response, error) {

	body, err := u.body()

	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(u.method(), u.url(), body)
	if err == nil {
		return client.Do(request)
	}
	return nil, err
}

// method returns a normalized HTTP verb in CAPS
func (u Upstream) method() string {
	return strings.ToUpper(u.Method)
}

func (u Upstream) body() (io.Reader, error) {
	if u.Method == http.MethodPost || u.Method == http.MethodPut {
		if contentType, ok := u.Headers["content-type"]; ok && contentType == "application/json" {
			body := u.Body.(map[string]interface{})
			data, err := json.Marshal(body)
			if err != nil {
				return nil, errors.New("could not marshal upstream body for " + u.URL + ", " + err.Error())
			}
			return bytes.NewReader(data), nil
		}
		return strings.NewReader(u.Body.(string)), nil
	}
	return nil, nil
}

func (u Upstream) url() string {
	url := u.URL

	if strings.HasPrefix(url, "env:") {
		url = os.Getenv(strings.TrimPrefix(url, "env:"))
	}

	return url
}

type Response struct {
	Headers                 map[string]string `json:"headers"`
	StatusCode              int               `json:"status_code"`
	Body                    interface{}       `json:"body"`
	ConcatUpstreamResponses bool              `json:"concat_upstream_responses"`
}

func (r Response) Content(upstreamResponses []*http.Response) ([]byte, error) {

	contentType, ok := r.Headers["content-type"]

	if ok && contentType == "application/json" {
		return r.json(upstreamResponses)
	}

	return r.text(upstreamResponses)
}

func (r Response) text(upstreamResponses []*http.Response) ([]byte, error) {
	content := r.Body.(string)

	if len(upstreamResponses) > 0 && r.ConcatUpstreamResponses {
		for _, upstreamResponse := range upstreamResponses {
			upstreamData, _ := io.ReadAll(upstreamResponse.Body)
			content = fmt.Sprintf("%s\n# %s %s\n%s", content, upstreamResponse.Status, upstreamResponse.Request.URL.String(), upstreamData)
		}
	}

	return []byte(content), nil
}

func (r Response) json(upstreamResponses []*http.Response) ([]byte, error) {
	body := r.Body.(map[string]interface{})
	if len(upstreamResponses) > 0 && r.ConcatUpstreamResponses {
		upstreamContents := []interface{}{}
		for _, upstreamResponse := range upstreamResponses {
			upstreamData, _ := io.ReadAll(upstreamResponse.Body)
			if upstreamResponse.Header.Get("content-type") == "application/json" {
				body := map[string]interface{}{}
				if err := json.Unmarshal(upstreamData, &body); err != nil {
					continue
				}
				upstreamContents = append(upstreamContents, map[string]interface{}{
					"url":  upstreamResponse.Request.URL.String(),
					"body": body,
				})
				continue
			}
			upstreamContents = append(upstreamContents, string(upstreamData))
		}
		body["upstreams"] = upstreamContents
	}
	return json.Marshal(body)
}

func ConfigFromFilePath(path string) (Config, error) {
	var config Config

	data, err := os.ReadFile(path)

	if err == nil {
		json.Unmarshal(data, &config)
	}

	return config, err
}
