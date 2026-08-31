// Package jsonv2 provides a fastshot.JSONCodec backed by the standard library
// encoding/json/v2, available since Go 1.27.
//
// Opt in by setting it on the client builder:
//
//	import "github.com/opus-domini/fast-shot/jsonv2"
//
//	client := fastshot.NewClient("https://api.example.com").
//		Config().SetJSONCodec(jsonv2.New()).
//		Build()
//
// Compared to the default codec (encoding/json), this codec applies stricter,
// more interoperable defaults: it rejects invalid UTF-8 in JSON strings,
// rejects duplicate names within JSON objects, and rejects trailing data after
// a top-level JSON value. Marshaled output also carries no trailing newline.
// See the encoding/json/v2 documentation for the complete set of differences
// and available options.
package jsonv2

import (
	json "encoding/json/v2"
	"io"

	fastshot "github.com/opus-domini/fast-shot"
)

// New returns a fastshot.JSONCodec using encoding/json/v2.
func New() fastshot.JSONCodec {
	return fastshot.JSONCodec{
		Encode: func(w io.Writer, v any) error {
			return json.MarshalWrite(w, v)
		},
		Decode: func(r io.Reader, v any) error {
			return json.UnmarshalRead(r, v)
		},
	}
}
