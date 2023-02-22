# GOMSVC
Small project to quickly create mock service endpoints from a json configuration.

## Environment variables

**GOMSVC_CONFIG_PATH**: set this to determine which json file is loaded. By default config.json in the same directory will be loaded unless `GOMSVC_CONFIG_STRING` is set.

**GOMSVC_CONFIG_STRING**: set this with a valid JSON configuration that will be loaded. This has no effect if `GOMSVC_CONFIG_PATH` is set.

**GOMSVC_LOG_LEVEL**: set this to determine which log level to use. By default `info` will be used.

## Configuration

**port**: Determines which port the server will start on.

**routes[]**: List of all routes that should be served.

**routes[].name**: Name/Identifier of the route.

**routes[].path**: Which path this route should be served from.

**routes[].method**: Which method that should be allowed for this route.

**routes[].upstreams[]**: List of upstream calls to perform whenever this route is invoked.

**routes[].upstreams[].url**: Destination of the upstream call, if the string is prefixed with `env:`, the url value will be retrieved from the given environment variable instead.

**routes[].upstreams[].method**: Which HTTP method to use when for the upstream call. If either POST or PUT, the body of the upstream call will be sent with it.

**routes[].upstreams[].headers{}**: Object with key:value sets that are attached as headers for the upstream call. Setting `content-type` to `json/application` will trigger that any upstream body will be encoded as JSON before being sent.

**routes[].upstreams[].body**: String or object that is sent whenever the upstream call is using HTTP Method POST or PUT.

**routes[].response.headers{}**: Object with key:value sets that are attached as headers for the route response. Setting `content-type` to `json/application` will trigger the body will be encoded as JSON before being served.

**routes[].response.body**: String or object that is returned as response body.

**routes[].response.status_code**: Whatever HTTP Status Code should be used for the response.

**routes[].response.concat_upstream_responses**: If set to true, upstream responses will be injected into the response body.