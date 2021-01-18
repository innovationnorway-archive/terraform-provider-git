package provider

import (
	"net/http"
)

type userAgentTransport struct {
	transport http.RoundTripper
	userAgent string
}

func (rt *userAgentTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", rt.userAgent)
	return rt.transport.RoundTrip(req)
}
