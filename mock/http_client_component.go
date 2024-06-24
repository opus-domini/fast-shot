package mock

import (
	"net/http"
	"time"

	"github.com/stretchr/testify/mock"
)

type HttpClientComponent struct {
	mock.Mock
}

func (m *HttpClientComponent) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

func (m *HttpClientComponent) Transport() http.RoundTripper {
	args := m.Called()
	return args.Get(0).(http.RoundTripper)
}

func (m *HttpClientComponent) SetTransport(transport http.RoundTripper) {
	m.Called(transport)
}

func (m *HttpClientComponent) Timeout() time.Duration {
	args := m.Called()
	return args.Get(0).(time.Duration)
}

func (m *HttpClientComponent) SetTimeout(duration time.Duration) {
	m.Called(duration)
}

func (m *HttpClientComponent) SetFollowRedirects(follow bool) {
	m.Called(follow)
}
