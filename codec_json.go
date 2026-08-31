package fastshot

import (
	"encoding/json"
	"io"
)

// JSONCodec defines how request and response bodies are encoded and decoded as JSON.
//
// It is used by BodyWrapper implementations for the WriteAsJSON and ReadAsJSON
// methods, allowing the JSON implementation to be swapped without changing the
// public API. The default codec uses the standard library encoding/json.
//
// Example usage with a custom codec:
//
//	codec := fastshot.DefaultJSONCodec()
//	codec.Decode = myStrictDecode
//	client := fastshot.NewClient("https://api.example.com").
//		Config().SetJSONCodec(codec).
//		Build()
type JSONCodec struct {
	// Encode writes the JSON representation of v to w.
	Encode func(w io.Writer, v any) error
	// Decode reads JSON from r and stores it in v.
	Decode func(r io.Reader, v any) error
}

// DefaultJSONCodec returns a JSONCodec backed by the standard library encoding/json.
func DefaultJSONCodec() JSONCodec {
	return JSONCodec{
		Encode: func(w io.Writer, v any) error {
			return json.NewEncoder(w).Encode(v)
		},
		Decode: func(r io.Reader, v any) error {
			return json.NewDecoder(r).Decode(v)
		},
	}
}
