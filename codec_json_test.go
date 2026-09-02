package fastshot

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type codecTestUser struct {
	Name string `json:"name"`
}

func TestDefaultJSONCodec(t *testing.T) {
	codec := DefaultJSONCodec()

	t.Run("round trip", func(t *testing.T) {
		var buf strings.Builder
		if err := codec.Encode(&buf, codecTestUser{Name: "Ana"}); err != nil {
			t.Fatalf("unexpected error encoding: %v", err)
		}
		if got, want := buf.String(), `{"name":"Ana"}`; got != want {
			t.Errorf("encoded got %q, want %q", got, want)
		}

		var user codecTestUser
		if err := codec.Decode(strings.NewReader(buf.String()), &user); err != nil {
			t.Fatalf("unexpected error decoding: %v", err)
		}
		if user.Name != "Ana" {
			t.Errorf("decoded name got %q, want %q", user.Name, "Ana")
		}
	})

	t.Run("rejects trailing data", func(t *testing.T) {
		var user codecTestUser
		err := codec.Decode(strings.NewReader(`{"name":"Ana"} trailing`), &user)
		if err == nil {
			t.Error("expected error decoding trailing data, got nil")
		}
	})

	t.Run("rejects duplicate object keys", func(t *testing.T) {
		var user map[string]string
		err := codec.Decode(strings.NewReader(`{"name":"Ana","name":"Bia"}`), &user)
		if err == nil {
			t.Error("expected error decoding duplicate keys, got nil")
		}
	})

	t.Run("codec funcs are never nil", func(t *testing.T) {
		if codec.Encode == nil {
			t.Error("Encode func is nil")
		}
		if codec.Decode == nil {
			t.Error("Decode func is nil")
		}
	})
}

func TestNewJSONv1Codec(t *testing.T) {
	codec := NewJSONv1Codec()

	var buf strings.Builder
	if err := codec.Encode(&buf, codecTestUser{Name: "Ana"}); err != nil {
		t.Fatalf("unexpected error encoding: %v", err)
	}
	// encoding/json Encoder appends a trailing newline.
	if got, want := buf.String(), "{\"name\":\"Ana\"}\n"; got != want {
		t.Errorf("encoded got %q, want %q", got, want)
	}

	var user codecTestUser
	if err := codec.Decode(strings.NewReader(`{"name":"Ana"} trailing`), &user); err != nil {
		t.Fatalf("unexpected error decoding: %v", err)
	}
	if user.Name != "Ana" {
		t.Errorf("decoded name got %q, want %q", user.Name, "Ana")
	}
}

func TestJSONCodec_ClientIntegration(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user codecTestUser
		if err := DefaultJSONCodec().Decode(r.Body, &user); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = DefaultJSONCodec().Encode(w, codecTestUser{Name: user.Name + "-echo"})
	}))
	defer server.Close()

	client := NewClient(server.URL).Build()

	resp, err := client.POST("/users").
		Body().AsJSON(codecTestUser{Name: "Ana"}).
		Send()
	if err != nil {
		t.Fatalf("unexpected error sending: %v", err)
	}

	user, err := resp.Body().AsJSON[codecTestUser]()
	if err != nil {
		t.Fatalf("unexpected error decoding response: %v", err)
	}
	if user.Name != "Ana-echo" {
		t.Errorf("name got %q, want %q", user.Name, "Ana-echo")
	}
}

func TestJSONCodec_V1ClientIntegration(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user codecTestUser
		if err := NewJSONv1Codec().Decode(r.Body, &user); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = NewJSONv1Codec().Encode(w, codecTestUser{Name: user.Name + "-echo"})
	}))
	defer server.Close()

	client := NewClient(server.URL).
		Config().SetJSONCodec(NewJSONv1Codec()).
		Build()

	resp, err := client.POST("/users").
		Body().AsJSON(codecTestUser{Name: "Ana"}).
		Send()
	if err != nil {
		t.Fatalf("unexpected error sending: %v", err)
	}

	user, err := resp.Body().AsJSON[codecTestUser]()
	if err != nil {
		t.Fatalf("unexpected error decoding response: %v", err)
	}
	if user.Name != "Ana-echo" {
		t.Errorf("name got %q, want %q", user.Name, "Ana-echo")
	}
}
