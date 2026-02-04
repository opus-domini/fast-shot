package mock

import (
	"io"

	"github.com/stretchr/testify/mock"
)

type BodyWrapper struct {
	mock.Mock
}

func (m *BodyWrapper) Read(p []byte) (n int, err error) {
	args := m.Called(p)
	return args.Int(0), args.Error(1)
}

func (m *BodyWrapper) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *BodyWrapper) ReadAsJSON(obj interface{}) error {
	args := m.Called(obj)
	return args.Error(0)
}

func (m *BodyWrapper) WriteAsJSON(obj interface{}) error {
	args := m.Called(obj)
	return args.Error(0)
}

func (m *BodyWrapper) ReadAsXML(obj interface{}) error {
	args := m.Called(obj)
	return args.Error(0)
}

func (m *BodyWrapper) WriteAsXML(obj interface{}) error {
	args := m.Called(obj)
	return args.Error(0)
}

func (m *BodyWrapper) ReadAsString() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *BodyWrapper) WriteAsString(body string) error {
	args := m.Called(body)
	return args.Error(0)
}

func (m *BodyWrapper) WriteAsFormData(fields map[string]string) (string, error) {
	args := m.Called(fields)
	return args.String(0), args.Error(1)
}

func (m *BodyWrapper) Set(body io.Reader) error {
	args := m.Called(body)
	return args.Error(0)
}

func (m *BodyWrapper) Unwrap() io.Reader {
	args := m.Called()
	return args.Get(0).(io.Reader)
}
