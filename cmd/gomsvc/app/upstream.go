package app

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
)

type Upstream struct {
	URL                   string            `json:"url"`
	IncludeRequestHeaders bool              `json:"include_request_headers"`
	Headers               map[string]string `json:"headers"`
	Method                string            `json:"method"`
	Body                  interface{}       `json:"body"`
}

// Call takes all the information in the upstream and makes a request
// based on that. URL with a value that has a env: prefix will instead use
// whatever value is in that environment variable as an URL instead for
// the upstream call
func (u Upstream) Call(client *http.Client, req *http.Request) (*http.Response, error) {

	body, err := u.body()

	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(u.method(), u.url(), body)

	if u.IncludeRequestHeaders && req != nil {
		for name, values := range req.Header {
			for _, v := range values {
				request.Header.Add(name, v)
			}
		}
	}

	for k, v := range u.Headers {
		request.Header.Set(k, v)
	}

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
	if (u.Method == http.MethodPost || u.Method == http.MethodPut) && u.Body != nil {
		if contentType, ok := u.Headers["content-type"]; ok && contentType == "application/json" {
			body := u.Body.(map[string]interface{})
			data, err := json.Marshal(body)
			if err != nil {
				return nil, errors.New("could not marshal upstream body for " + u.URL + ", " + err.Error())
			}
			return bytes.NewReader(data), nil
		}
		if body, ok := u.Body.(string); ok {
			return strings.NewReader(body), nil
		}
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
