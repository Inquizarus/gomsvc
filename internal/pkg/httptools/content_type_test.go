package httptools_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/inquizarus/gomsvc/internal/pkg/httptools"
)

func TestIsJSON(t *testing.T) {
	jsonReq, _ := http.NewRequest(http.MethodGet, "/", strings.NewReader(`{"name": "John", "age": 30}`))
	jsonReq.Header.Set("Content-Type", "application/json")
	if !httptools.IsJSON(jsonReq.Header) {
		t.Error("expected true, got false")
	}

	xmlReq, _ := http.NewRequest(http.MethodGet, "/", strings.NewReader(`<?xml version="1.0" encoding="UTF-8"?><note><to>Tove</to><from>Jani</from><heading>Reminder</heading><body>Don't forget me this weekend!</body></note>`))
	xmlReq.Header.Set("Content-Type", "application/xml")
	if httptools.IsJSON(xmlReq.Header) {
		t.Error("expected false, got true")
	}

	textReq, _ := http.NewRequest(http.MethodGet, "/", strings.NewReader("this is plain text"))
	textReq.Header.Set("Content-Type", "text/plain")
	if httptools.IsJSON(textReq.Header) {
		t.Error("expected false, got true")
	}
}

func TestIsXML(t *testing.T) {
	jsonReq, _ := http.NewRequest(http.MethodGet, "/", strings.NewReader(`{"name": "John", "age": 30}`))
	jsonReq.Header.Set("Content-Type", "application/json")
	if httptools.IsXML(jsonReq.Header) {
		t.Error("expected false, got true")
	}

	xmlReq, _ := http.NewRequest(http.MethodGet, "/", strings.NewReader(`<?xml version="1.0" encoding="UTF-8"?><note><to>Tove</to><from>Jani</from><heading>Reminder</heading><body>Don't forget me this weekend!</body></note>`))
	xmlReq.Header.Set("Content-Type", "application/xml")
	if !httptools.IsXML(xmlReq.Header) {
		t.Error("expected true, got false")
	}

	textReq, _ := http.NewRequest(http.MethodGet, "/", strings.NewReader("this is plain text"))
	textReq.Header.Set("Content-Type", "text/plain")
	if httptools.IsXML(textReq.Header) {
		t.Error("expected false, got true")
	}
}

func TestIsPlainText(t *testing.T) {
	jsonReq, _ := http.NewRequest(http.MethodGet, "/", strings.NewReader(`{"name": "John", "age": 30}`))
	jsonReq.Header.Set("Content-Type", "application/json")
	if httptools.IsPlainText(jsonReq.Header) {
		t.Error("expected false, got true")
	}

	xmlReq, _ := http.NewRequest(http.MethodGet, "/", strings.NewReader(`<?xml version="1.0" encoding="UTF-8"?><note><to>Tove</to><from>Jani</from><heading>Reminder</heading><body>Don't forget me this weekend!</body></note>`))
	xmlReq.Header.Set("Content-Type", "application/xml")
	if httptools.IsPlainText(xmlReq.Header) {
		t.Error("expected false, got true")
	}

	textReq, _ := http.NewRequest(http.MethodGet, "/", strings.NewReader("this is plain text"))
	textReq.Header.Set("Content-Type", "text/plain")
	if !httptools.IsPlainText(textReq.Header) {
		t.Error("expected true, got false")
	}
}

func TestIsHTML(t *testing.T) {
	tests := []struct {
		name           string
		contentType    string
		expectedResult bool
	}{
		{
			name:           "HTML content type with utf-8",
			contentType:    "text/html; charset=utf-8",
			expectedResult: true,
		},
		{
			name:           "HTML content type without charset",
			contentType:    "text/html",
			expectedResult: true,
		},
		{
			name:           "Non-HTML content type",
			contentType:    "application/json",
			expectedResult: false,
		},
	}

	for k, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			req.Header.Set("Content-Type", tt.contentType)

			if httptools.IsHTML(req.Header) != tt.expectedResult {
				t.Errorf("expected %v but got %v -> %d", tt.expectedResult, !tt.expectedResult, k)
			}
		})
	}
}

func TestIsFormURLEncoded(t *testing.T) {
	tests := []struct {
		name           string
		contentType    string
		expectedResult bool
	}{
		{
			name:           "Form URL-encoded content type",
			contentType:    "application/x-www-form-urlencoded",
			expectedResult: true,
		},
		{
			name:           "Plain text content type",
			contentType:    "text/plain",
			expectedResult: false,
		},
		{
			name:           "JSON content type",
			contentType:    "application/json",
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			req.Header.Set("Content-Type", tt.contentType)

			if httptools.IsFormURLEncoded(req.Header) != tt.expectedResult {
				t.Errorf("expected %v but got %v", tt.expectedResult, !tt.expectedResult)
			}
		})
	}
}

func TestIsMultipart(t *testing.T) {
	tests := []struct {
		name           string
		contentType    string
		expectedResult bool
	}{
		{
			name:           "Multipart content type",
			contentType:    "multipart/form-data; boundary=------------------------78c4658e492de3c4",
			expectedResult: true,
		},
		{
			name:           "Plain text content type",
			contentType:    "text/plain",
			expectedResult: false,
		},
		{
			name:           "JSON content type",
			contentType:    "application/json",
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			req.Header.Set("Content-Type", tt.contentType)

			if httptools.IsMultipart(req.Header) != tt.expectedResult {
				t.Errorf("expected %v but got %v", tt.expectedResult, !tt.expectedResult)
			}
		})
	}
}
