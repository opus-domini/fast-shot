package fastshot

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/opus-domini/fast-shot/constant/header"
	"github.com/opus-domini/fast-shot/constant/mime"
)

func BenchmarkClientBuild(b *testing.B) {
	for b.Loop() {
		_ = NewClient("https://api.example.com").
			Auth().BearerToken("token").
			Header().Add(header.UserAgent, "bench").
			Config().SetTimeout(30_000_000_000).
			Build()
	}
}

func BenchmarkRequestBuild(b *testing.B) {
	client := DefaultClient("https://api.example.com")

	for b.Loop() {
		_ = client.GET("/users").
			Header().AddAccept(mime.JSON).
			Query().AddParam("page", "1").
			Body().AsJSON(map[string]string{"key": "value"})
	}
}

func BenchmarkSendGET(b *testing.B) {
	server := httptest.NewTestServer(b, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("hello"))
	}))
	defer server.Close()

	serverClient := server.Client()
	client := NewClient("http://example.com").
		Config().SetCustomTransport(serverClient.Transport).
		Build()
	defer serverClient.CloseIdleConnections()

	b.ReportAllocs()
	for b.Loop() {
		resp, err := client.GET("/test").Send()
		if err != nil {
			b.Fatal(err)
		}
		body, err := resp.Body().AsString()
		if err != nil {
			b.Fatal(err)
		}
		if body != "hello" {
			b.Fatalf("body got %q, want %q", body, "hello")
		}
	}
}

func BenchmarkResponseAsJSONOf(b *testing.B) {
	server := httptest.NewTestServer(b, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = DefaultJSONCodec().Encode(w, codecTestUser{Name: "Ana"})
	}))
	defer server.Close()

	serverClient := server.Client()
	client := NewClient("http://example.com").
		Config().SetCustomTransport(serverClient.Transport).
		Build()
	defer serverClient.CloseIdleConnections()

	b.ReportAllocs()
	for b.Loop() {
		user, err := client.GET("/user").Send()
		if err != nil {
			b.Fatal(err)
		}
		decoded, err := user.Body().AsJSONOf[codecTestUser]()
		if err != nil {
			b.Fatal(err)
		}
		if decoded.Name != "Ana" {
			b.Fatalf("name got %q, want %q", decoded.Name, "Ana")
		}
	}
}

func BenchmarkJSONCodec(b *testing.B) {
	benchmarks := []struct {
		name  string
		codec JSONCodec
	}{
		{"DefaultJSONCodec_v2", DefaultJSONCodec()},
		{"NewJSONv1Codec", NewJSONv1Codec()},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			payload := `{"name":"Ana","age":30,"email":"ana@example.com","active":true}`
			b.ReportAllocs()
			for b.Loop() {
				var user struct {
					Name   string `json:"name"`
					Age    int    `json:"age"`
					Email  string `json:"email"`
					Active bool   `json:"active"`
				}
				if err := bm.codec.Decode(strings.NewReader(payload), &user); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
