package mock

import (
	"io"
)

type BodyWrapper struct {
	ReadFunc            func(p []byte) (int, error)
	CloseFunc           func() error
	ReadAsJSONFunc      func(obj any) error
	WriteAsJSONFunc     func(obj any) error
	ReadAsXMLFunc       func(obj any) error
	WriteAsXMLFunc      func(obj any) error
	ReadAsStringFunc    func() (string, error)
	WriteAsStringFunc   func(body string) error
	WriteAsFormDataFunc func(fields map[string]string) (string, error)
	SetFunc             func(body io.Reader) error
	UnwrapFunc          func() io.Reader
}

func (m *BodyWrapper) Read(p []byte) (int, error) {
	return m.ReadFunc(p)
}

func (m *BodyWrapper) Close() error {
	return m.CloseFunc()
}

func (m *BodyWrapper) ReadAsJSON(obj any) error {
	return m.ReadAsJSONFunc(obj)
}

func (m *BodyWrapper) WriteAsJSON(obj any) error {
	return m.WriteAsJSONFunc(obj)
}

func (m *BodyWrapper) ReadAsXML(obj any) error {
	return m.ReadAsXMLFunc(obj)
}

func (m *BodyWrapper) WriteAsXML(obj any) error {
	return m.WriteAsXMLFunc(obj)
}

func (m *BodyWrapper) ReadAsString() (string, error) {
	return m.ReadAsStringFunc()
}

func (m *BodyWrapper) WriteAsString(body string) error {
	return m.WriteAsStringFunc(body)
}

func (m *BodyWrapper) WriteAsFormData(fields map[string]string) (string, error) {
	return m.WriteAsFormDataFunc(fields)
}

func (m *BodyWrapper) Set(body io.Reader) error {
	return m.SetFunc(body)
}

func (m *BodyWrapper) Unwrap() io.Reader {
	return m.UnwrapFunc()
}
