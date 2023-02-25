package httptools

import (
	"net/http"
	"strings"
)

const (
	headerContentTypeKey  = "content-type"
	contentTypeJSON       = "application/json"
	contentTypePlainText  = "text/plain"
	contentTypeXML        = "application/xml"
	contentTypeHTML       = "text/html"
	contentTypeFormURLEnc = "application/x-www-form-urlencoded"
	contentTypeMultipart  = "multipart/form-data"
)

func IsJSON(headers http.Header) bool {
	return headers.Get(headerContentTypeKey) == contentTypeJSON
}

func IsPlainText(headers http.Header) bool {
	contentType := headers.Get(headerContentTypeKey)
	return contentType == contentTypePlainText || contentType == ""
}

func IsXML(headers http.Header) bool {
	return headers.Get(headerContentTypeKey) == contentTypeXML
}

func IsHTML(headers http.Header) bool {
	return strings.HasPrefix(headers.Get(headerContentTypeKey), contentTypeHTML)
}

func IsFormURLEncoded(headers http.Header) bool {
	return headers.Get(headerContentTypeKey) == contentTypeFormURLEnc
}

func IsMultipart(headers http.Header) bool {
	return strings.HasPrefix(headers.Get(headerContentTypeKey), contentTypeMultipart)
}
