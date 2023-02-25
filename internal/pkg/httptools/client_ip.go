package httptools

import (
	"net/http"
	"strings"
)

func ClientIP(request *http.Request) string {
	if request == nil {
		return ""
	}
	if clientIP := request.Header.Get("X-Real-IP"); clientIP != "" {
		return clientIP
	}
	if clientIP := request.Header.Get("X-Forwarded-For"); clientIP != "" {
		ips := strings.Split(clientIP, ", ")
		return ips[len(ips)-1]
	}
	return request.RemoteAddr
}
