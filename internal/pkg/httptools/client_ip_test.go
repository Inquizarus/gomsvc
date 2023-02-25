package httptools_test

import (
	"net/http"
	"testing"

	"github.com/inquizarus/gomsvc/internal/pkg/httptools"
)

func TestClientIP_WithXRealIPHeader(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("X-Real-IP", "192.0.2.1")

	ip := httptools.ClientIP(req)
	if ip != "192.0.2.1" {
		t.Errorf("Expected 192.0.2.1, but got %s", ip)
	}
}

func TestClientIP_WithXForwardedForHeader(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("X-Forwarded-For", "192.0.2.2, 192.0.2.1")

	ip := httptools.ClientIP(req)
	if ip != "192.0.2.1" {
		t.Errorf("Expected 192.0.2.1, but got %s", ip)
	}
}

func TestClientIP_WithXForwardedForHeaderWithoutComma(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("X-Forwarded-For", "192.0.2.1")

	ip := httptools.ClientIP(req)
	if ip != "192.0.2.1" {
		t.Errorf("Expected 192.0.2.1, but got %s", ip)
	}
}

func TestClientIP_WithoutHeaders(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.RemoteAddr = "127.0.0.1"
	ip := httptools.ClientIP(req)
	if ip != "127.0.0.1" {
		t.Errorf("Expected 127.0.0.1, but got %s", ip)
	}
}

func TestClientIP_WithNilRequest(t *testing.T) {
	ip := httptools.ClientIP(nil)
	if ip != "" {
		t.Errorf("Expected empty string, but got %s", ip)
	}
}
