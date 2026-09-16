package fastshot

import (
	"bytes"
	"encoding/xml"
	"io"
	"mime/multipart"
	"strings"
)

// Compile-time check that BufferedBody implements BodyWrapper.
var _ BodyWrapper = (*BufferedBody)(nil)

// Compile-time check that UnbufferedBody implements BodyWrapper.
var _ BodyWrapper = (*UnbufferedBody)(nil)

type (
	// BufferedBody is a BodyWrapper backed by an in-memory buffer.
	// It is not safe for concurrent use, matching io.ReadCloser semantics.
	BufferedBody struct {
		buffer    *bytes.Buffer
		jsonCodec JSONCodec
	}

	// UnbufferedBody is a BodyWrapper backed by a stream.
	// It is not safe for concurrent use, matching io.ReadCloser semantics.
	UnbufferedBody struct {
		reader    io.ReadCloser
		jsonCodec JSONCodec
	}
)

func (w *BufferedBody) Read(p []byte) (n int, err error) {
	return w.buffer.Read(p)
}

func (w *BufferedBody) Close() error {
	// No-op for buffered wrapper
	return nil
}

func (w *BufferedBody) ReadAsJSON(obj any) error {
	return w.jsonCodec.Decode(bytes.NewReader(w.buffer.Bytes()), obj)
}

func (w *BufferedBody) WriteAsJSON(obj any) error {
	w.buffer.Reset()
	return w.jsonCodec.Encode(w.buffer, obj)
}

func (w *BufferedBody) ReadAsXML(obj any) error {
	return xml.NewDecoder(bytes.NewReader(w.buffer.Bytes())).Decode(obj)
}

func (w *BufferedBody) WriteAsXML(obj any) error {
	w.buffer.Reset()
	return xml.NewEncoder(w.buffer).Encode(obj)
}

func (w *BufferedBody) ReadAsString() (string, error) {
	return w.buffer.String(), nil
}

func (w *BufferedBody) WriteAsString(s string) error {
	w.buffer.Reset()
	_, err := w.buffer.WriteString(s)
	return err
}

func (w *BufferedBody) WriteAsFormData(fields map[string]string) (string, error) {
	w.buffer.Reset()

	var body bytes.Buffer
	contentType, err := writeFormDataFn(&body, fields)
	if err != nil {
		return "", err
	}

	_, err = io.Copy(w.buffer, &body)
	return contentType, err
}

func (w *BufferedBody) Set(body io.Reader) error {
	w.buffer.Reset()
	_, err := io.Copy(w.buffer, body)
	return err
}

func (w *BufferedBody) Unwrap() io.Reader {
	return bytes.NewReader(w.buffer.Bytes())
}

func newBufferedBody(jsonCodec JSONCodec) *BufferedBody {
	return &BufferedBody{
		buffer:    &bytes.Buffer{},
		jsonCodec: jsonCodec,
	}
}

func (w *UnbufferedBody) Read(p []byte) (n int, err error) {
	return w.reader.Read(p)
}

func (w *UnbufferedBody) Close() error {
	return w.reader.Close()
}

func (w *UnbufferedBody) ReadAsJSON(obj any) error {
	return w.jsonCodec.Decode(w.reader, obj)
}

func (w *UnbufferedBody) WriteAsJSON(obj any) error {
	var buf bytes.Buffer
	err := w.jsonCodec.Encode(&buf, obj)
	if err != nil {
		return err
	}
	w.reader = io.NopCloser(&buf)
	return nil
}

func (w *UnbufferedBody) ReadAsXML(obj any) error {
	return xml.NewDecoder(w.reader).Decode(obj)
}

func (w *UnbufferedBody) WriteAsXML(obj any) error {
	var buf bytes.Buffer
	err := xml.NewEncoder(&buf).Encode(obj)
	if err != nil {
		return err
	}
	w.reader = io.NopCloser(&buf)
	return nil
}

func (w *UnbufferedBody) ReadAsString() (string, error) {
	stringBytes, err := io.ReadAll(w.reader)
	if err != nil {
		return "", err
	}
	return string(stringBytes), nil
}

func (w *UnbufferedBody) WriteAsString(s string) error {
	w.reader = io.NopCloser(strings.NewReader(s))
	return nil
}

func (w *UnbufferedBody) WriteAsFormData(fields map[string]string) (string, error) {
	var body bytes.Buffer
	contentType, err := writeFormDataFn(&body, fields)
	if err != nil {
		return "", err
	}

	w.reader = io.NopCloser(&body)
	return contentType, nil
}

func (w *UnbufferedBody) Set(body io.Reader) error {
	if closer, ok := body.(io.ReadCloser); ok {
		w.reader = closer
	} else {
		w.reader = io.NopCloser(body)
	}
	return nil
}

func (w *UnbufferedBody) Unwrap() io.Reader {
	return w.reader
}

// writeFormDataFn is the function used to write form data. It can be swapped in tests.
var writeFormDataFn = writeFormData

func writeFormData(dst io.Writer, fields map[string]string) (string, error) {
	writer := multipart.NewWriter(dst)

	for key, value := range fields {
		if err := writer.WriteField(key, value); err != nil {
			_ = writer.Close()
			return "", err
		}
	}

	contentType := writer.FormDataContentType()
	if err := writer.Close(); err != nil {
		return "", err
	}

	return contentType, nil
}

func newUnbufferedBody(reader io.ReadCloser, jsonCodec JSONCodec) *UnbufferedBody {
	return &UnbufferedBody{
		reader:    reader,
		jsonCodec: jsonCodec,
	}
}
