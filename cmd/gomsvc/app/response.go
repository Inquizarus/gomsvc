package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Response struct {
	Headers                   map[string]string `json:"headers"`
	StatusCode                int               `json:"status_code"`
	Body                      interface{}       `json:"body"`
	IncludeUpstreamResponses  bool              `json:"concat_upstream_responses"`
	IncludeRequestInformation bool              `json:"include_request_information"`
}

func (r Response) Content(request *http.Request, upstreamResponses []*http.Response) ([]byte, error) {

	contentType, ok := r.Headers["content-type"]

	if ok && contentType == "application/json" {
		return r.json(request, upstreamResponses)
	}

	return r.text(request, upstreamResponses)
}

func (r Response) text(request *http.Request, upstreamResponses []*http.Response) ([]byte, error) {
	content := ""

	if r.shouldIncludeRequestInformation(request) {
		content = content + `######################
#   Request headers  #
######################
` + "\n"
		for k, v := range request.Header {
			content = content + k + ":" + strings.Join(v, ",") + "\n"
		}
	}

	content = content + "\n--------------------------------------------------\n\n" + r.Body.(string)

	if len(upstreamResponses) > 0 && r.IncludeUpstreamResponses {
		content = content + "\n\n" + `
#####################
#   Upstream calls  #
#####################
` + "\n"
		for _, upstreamResponse := range upstreamResponses {
			upstreamData, _ := io.ReadAll(upstreamResponse.Body)
			content = fmt.Sprintf(
				"%s\n\t%s - %s - %s\n\tFROM %s \n\n\t%s",
				content,
				upstreamResponse.Request.Method,
				upstreamResponse.Request.URL.String(),
				upstreamResponse.Status,
				r.getClientIP(request),
				bytes.ReplaceAll(upstreamData, []byte{'\n'}, []byte{'\n', '\t'}),
			)
		}
	}

	return []byte(content), nil
}

func (r Response) json(request *http.Request, upstreamResponses []*http.Response) ([]byte, error) {
	body := r.Body.(map[string]interface{})
	if r.shouldIncludeRequestInformation(request) {

		body["request"] = map[string]interface{}{
			"client_ip": r.getClientIP(request),
			"method":    request.Method,
			"headers":   r.Headers,
		}
	}
	if len(upstreamResponses) > 0 && r.IncludeUpstreamResponses {
		upstreamContents := []interface{}{}
		for _, upstreamResponse := range upstreamResponses {
			upstreamData, _ := io.ReadAll(upstreamResponse.Body)
			if upstreamResponse.Header.Get("content-type") == "application/json" {
				body := map[string]interface{}{}
				if err := json.Unmarshal(upstreamData, &body); err != nil {
					continue
				}
				upstreamContents = append(upstreamContents, map[string]interface{}{
					"url":         upstreamResponse.Request.URL.String(),
					"headers":     upstreamResponse.Header,
					"status_code": upstreamResponse.StatusCode,
					"body":        body,
				})
				continue
			}
			upstreamContents = append(upstreamContents, string(upstreamData))
		}
		body["upstreams"] = upstreamContents
	}
	return json.Marshal(body)
}

func (r Response) shouldIncludeRequestInformation(request *http.Request) bool {
	if r.IncludeRequestInformation || request.Header.Get(httpHeaderAddRequestHeadersInResponse) != "" {
		return true
	}
	return false
}

func (r Response) getClientIP(request *http.Request) string {
	if clientIP := request.Header.Get("X-Forwarded-For"); clientIP != "" {
		return clientIP
	}
	return request.RemoteAddr
}
