package fastshot

import (
	json "encoding/json"
	jsonv2 "encoding/json/v2"
	"io"
)

// JSONCodec defines how request and response bodies are encoded and decoded as JSON.
//
// It is used by BodyWrapper implementations for the WriteAsJSON and ReadAsJSON
// methods, allowing the JSON implementation to be swapped without changing the
// public API. The default codec uses the standard library encoding/json/v2.
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

// DefaultJSONCodec returns a JSONCodec backed by the standard library encoding/json/v2.
//
// Compared to NewJSONv1Codec, it applies stricter, more interoperable defaults:
// it rejects invalid UTF-8 in JSON strings, rejects duplicate names within JSON
// objects, and rejects trailing data after a top-level JSON value. Marshaled
// output carries no trailing newline. See the encoding/json/v2 documentation for
// the complete set of differences and available options.
func DefaultJSONCodec() JSONCodec {
	return JSONCodec{
		Encode: func(w io.Writer, v any) error {
			return jsonv2.MarshalWrite(w, v)
		},
		Decode: func(r io.Reader, v any) error {
			return jsonv2.UnmarshalRead(r, v)
		},
	}
}

// NewJSONv1Codec returns a JSONCodec with the lenient encoding/json (v1 API)
// semantics: duplicate object keys and invalid UTF-8 are tolerated, trailing
// data after a top-level value is ignored, and encoded output ends with a
// trailing newline.
func NewJSONv1Codec() JSONCodec {
	return JSONCodec{
		Encode: func(w io.Writer, v any) error {
			return json.NewEncoder(w).Encode(v)
		},
		Decode: func(r io.Reader, v any) error {
			return json.NewDecoder(r).Decode(v)
		},
	}
}
