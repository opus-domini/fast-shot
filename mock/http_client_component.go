package mock

import (
	"net/http"
	"time"
)

type HttpClientComponent struct {
	DoFunc                 func(req *http.Request) (*http.Response, error)
	TransportFunc          func() http.RoundTripper
	SetTransportFunc       func(transport http.RoundTripper)
	TimeoutFunc            func() time.Duration
	SetTimeoutFunc         func(duration time.Duration)
	SetFollowRedirectsFunc func(follow bool)
}

func (m *HttpClientComponent) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

func (m *HttpClientComponent) Transport() http.RoundTripper {
	return m.TransportFunc()
}

func (m *HttpClientComponent) SetTransport(transport http.RoundTripper) {
	m.SetTransportFunc(transport)
}

func (m *HttpClientComponent) Timeout() time.Duration {
	return m.TimeoutFunc()
}

func (m *HttpClientComponent) SetTimeout(duration time.Duration) {
	m.SetTimeoutFunc(duration)
}

func (m *HttpClientComponent) SetFollowRedirects(follow bool) {
	m.SetFollowRedirectsFunc(follow)
}
