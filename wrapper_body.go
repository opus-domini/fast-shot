package fastshot

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"strings"
	"sync"
)

// BufferedBody implements BodyWrapper interface and provides a default HTTP context.
var _ BodyWrapper = (*BufferedBody)(nil)

// UnbufferedBody implements BodyWrapper interface and provides a default HTTP context.
var _ BodyWrapper = (*UnbufferedBody)(nil)

type (
	BufferedBody struct {
		buffer *bytes.Buffer
		mutex  sync.RWMutex
	}

	UnbufferedBody struct {
		reader io.ReadCloser
		mutex  sync.RWMutex
	}
)

func (w *BufferedBody) Read(p []byte) (n int, err error) {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	return w.buffer.Read(p)
}

func (w *BufferedBody) Close() error {
	// No-op for buffered wrapper
	return nil
}

func (w *BufferedBody) ReadAsJSON(obj interface{}) error {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	return json.NewDecoder(bytes.NewReader(w.buffer.Bytes())).Decode(obj)
}

func (w *BufferedBody) WriteAsJSON(obj interface{}) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.buffer.Reset()
	return json.NewEncoder(w.buffer).Encode(obj)
}

func (w *BufferedBody) ReadAsXML(obj interface{}) error {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	return xml.NewDecoder(bytes.NewReader(w.buffer.Bytes())).Decode(obj)
}

func (w *BufferedBody) WriteAsXML(obj interface{}) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.buffer.Reset()
	return xml.NewEncoder(w.buffer).Encode(obj)
}

func (w *BufferedBody) ReadAsString() (string, error) {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	return w.buffer.String(), nil
}

func (w *BufferedBody) WriteAsString(s string) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.buffer.Reset()
	_, err := w.buffer.WriteString(s)
	return err
}

func (w *BufferedBody) Set(body io.Reader) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.buffer.Reset()
	_, err := io.Copy(w.buffer, body)
	return err
}

func (w *BufferedBody) Unwrap() io.Reader {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	return bytes.NewReader(w.buffer.Bytes())
}

func newBufferedBody() *BufferedBody {
	return &BufferedBody{
		buffer: &bytes.Buffer{},
	}
}

func (w *UnbufferedBody) Read(p []byte) (n int, err error) {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	return w.reader.Read(p)
}

func (w *UnbufferedBody) Close() error {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	return w.reader.Close()
}

func (w *UnbufferedBody) ReadAsJSON(obj interface{}) error {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	return json.NewDecoder(w.reader).Decode(obj)
}

func (w *UnbufferedBody) WriteAsJSON(obj interface{}) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(obj)
	if err != nil {
		return err
	}
	w.reader = io.NopCloser(&buf)
	return nil
}

func (w *UnbufferedBody) ReadAsXML(obj interface{}) error {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	return xml.NewDecoder(w.reader).Decode(obj)
}

func (w *UnbufferedBody) WriteAsXML(obj interface{}) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	var buf bytes.Buffer
	err := xml.NewEncoder(&buf).Encode(obj)
	if err != nil {
		return err
	}
	w.reader = io.NopCloser(&buf)
	return nil
}

func (w *UnbufferedBody) ReadAsString() (string, error) {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	stringBytes, err := io.ReadAll(w.reader)
	if err != nil {
		return "", err
	}
	return string(stringBytes), nil
}

func (w *UnbufferedBody) WriteAsString(s string) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.reader = io.NopCloser(strings.NewReader(s))
	return nil
}

func (w *UnbufferedBody) Set(body io.Reader) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	if closer, ok := body.(io.ReadCloser); ok {
		w.reader = closer
	} else {
		w.reader = io.NopCloser(body)
	}
	return nil
}

func (w *UnbufferedBody) Unwrap() io.Reader {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	return w.reader
}

func newUnbufferedBody(reader io.ReadCloser) *UnbufferedBody {
	return &UnbufferedBody{
		reader: reader,
	}
}
