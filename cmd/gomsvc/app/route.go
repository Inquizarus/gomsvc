package app

type Route struct {
	Name      string     `json:"name"`
	Path      string     `json:"path"`
	Method    string     `json:"method"`
	Upstreams []Upstream `json:"upstreams"`
	Response  Response   `json:"response"`
}
