package app

const (
	envKeyConfigPath   = "GOMSVC_CONFIG_PATH"
	envKeyConfigString = "GOMSVC_CONFIG_STRING"
	envKeyRoutesDir    = "GOMSVC_ROUTES_DIR"
	configPathDefault  = "config.json"
	defaultPort        = "8080"

	httpHeaderAddRequestHeadersInResponse = "X-GOMSVC-Add-Request-Headers-In-Response"
	httpHeaderAddUpstreamsInResponse      = "X-GOMSVC-Add-Upstreams-In-Response"
)
