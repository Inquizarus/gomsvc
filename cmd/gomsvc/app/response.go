package app

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

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
