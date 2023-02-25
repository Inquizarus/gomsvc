package app

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/inquizarus/gomsvc/internal/pkg/httptools"
)

type Response struct {
	Headers                   map[string]string `json:"headers"`
	StatusCode                int               `json:"status_code"`
	Body                      interface{}       `json:"body"`
	IncludeUpstreamResponses  bool              `json:"concat_upstream_responses"`
	IncludeRequestInformation bool              `json:"include_request_information"`
}

func (r Response) Content(request *http.Request, upstreamResponses []*http.Response) ([]byte, error) {

	if request == nil {
		return nil, errors.New("request was nil")
	}

	contentType, ok := r.Headers["content-type"]

	if ok && contentType == "application/json" {
		return r.json(request, upstreamResponses)
	}

	return r.text(request, upstreamResponses)
}

func (r Response) text(request *http.Request, upstreamResponses []*http.Response) ([]byte, error) {
	var buf bytes.Buffer

	if r.shouldIncludeRequestInformation(request) {
		buf.WriteString("######################\n")
		buf.WriteString("#   Request headers  #\n")
		buf.WriteString("######################\n\n")

		var headerBuilder strings.Builder
		for k, v := range request.Header {
			headerBuilder.Reset()
			headerBuilder.WriteString(k)
			headerBuilder.WriteString(":")
			headerBuilder.WriteString(strings.Join(v, ","))
			headerBuilder.WriteString("\n")
			buf.WriteString(headerBuilder.String())
		}
	}

	buf.WriteString("\n--------------------------------------------------\n\n")
	buf.WriteString(r.Body.(string))
	buf.WriteString("\n\n--------------------------------------------------\n")

	if len(upstreamResponses) > 0 && r.IncludeUpstreamResponses {
		buf.WriteString("\n#####################\n")
		buf.WriteString("#   Upstream calls  #\n")
		buf.WriteString("#####################\n\n")

		for _, upstreamResponse := range upstreamResponses {
			if upstreamResponse == nil {
				continue
			}
			upstreamData, err := io.ReadAll(upstreamResponse.Body)
			if err != nil {
				return nil, err
			}

			if httptools.IsJSON(upstreamResponse.Header) {
				upstreamData, err = httptools.FormatJSONData(upstreamData)
				if err != nil {
					return nil, err
				}
			}

			buf.WriteString(fmt.Sprintf(
				"\n\t%s - %s - %s\n\tFROM %s \n\n\t%s",
				upstreamResponse.Request.Method,
				upstreamResponse.Request.URL.String(),
				upstreamResponse.Status,
				httptools.ClientIP(request),
				bytes.ReplaceAll(upstreamData, []byte{'\n'}, []byte{'\n', '\t'}), // Keeps everything indented
			))
		}
	}

	return buf.Bytes(), nil
}

func (r Response) json(request *http.Request, upstreamResponses []*http.Response) ([]byte, error) {

	body := r.copyBody()

	if r.shouldIncludeRequestInformation(request) {

		body["request"] = map[string]interface{}{
			"client_ip": httptools.ClientIP(request),
			"method":    request.Method,
			"headers":   r.Headers,
		}
	}

	if len(upstreamResponses) > 0 && r.IncludeUpstreamResponses {
		upstreamContents := []interface{}{}
		for _, upstreamResponse := range upstreamResponses {
			upstreamData, _ := io.ReadAll(upstreamResponse.Body)
			if httptools.IsJSON(upstreamResponse.Header) {
				container := map[string]interface{}{}
				if err := json.Unmarshal(upstreamData, &container); err != nil {
					continue
				}
				upstreamContents = append(upstreamContents, map[string]interface{}{
					"url":         upstreamResponse.Request.URL.String(),
					"headers":     upstreamResponse.Header,
					"status_code": upstreamResponse.StatusCode,
					"body":        container,
				})
				continue
			}
			upstreamContents = append(upstreamContents, string(upstreamData))
		}
		body["upstreams"] = upstreamContents
	}

	return httptools.FormatJSON(body)
}

func (r Response) shouldIncludeRequestInformation(request *http.Request) bool {
	if r.IncludeRequestInformation || request.Header.Get(httpHeaderAddRequestHeadersInResponse) != "" {
		return true
	}
	return false
}

func (r Response) copyBody() map[string]interface{} {
	body := map[string]interface{}{}

	for k, v := range r.Body.(map[string]interface{}) {
		body[k] = v
	}

	return body
}
