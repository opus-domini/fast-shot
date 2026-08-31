package jsonv2

import (
	"encoding/json/v2"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	fastshot "github.com/opus-domini/fast-shot"
)

type testUser struct {
	Name string `json:"name"`
}

func TestCodec_RoundTrip(t *testing.T) {
	codec := New()

	var buf strings.Builder
	if err := codec.Encode(&buf, testUser{Name: "Ana"}); err != nil {
		t.Fatalf("unexpected error encoding: %v", err)
	}
	if got, want := buf.String(), `{"name":"Ana"}`; got != want {
		t.Errorf("encoded got %q, want %q", got, want)
	}

	var user testUser
	if err := codec.Decode(strings.NewReader(buf.String()), &user); err != nil {
		t.Fatalf("unexpected error decoding: %v", err)
	}
	if user.Name != "Ana" {
		t.Errorf("decoded name got %q, want %q", user.Name, "Ana")
	}
}

func TestCodec_RejectsTrailingData(t *testing.T) {
	codec := New()

	var user testUser
	err := codec.Decode(strings.NewReader(`{"name":"Ana"} trailing`), &user)
	if err == nil {
		t.Fatal("expected error decoding trailing data, got nil")
	}
}

func TestCodec_ClientIntegration(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user testUser
		if err := json.UnmarshalRead(r.Body, &user); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.MarshalWrite(w, testUser{Name: user.Name + "-echo"})
	}))
	defer server.Close()

	client := fastshot.NewClient(server.URL).
		Config().SetJSONCodec(New()).
		Build()

	resp, err := client.POST("/users").
		Body().AsJSON(testUser{Name: "Ana"}).
		Send()
	if err != nil {
		t.Fatalf("unexpected error sending: %v", err)
	}

	user, err := resp.Body().AsJSONOf[testUser]()
	if err != nil {
		t.Fatalf("unexpected error decoding response: %v", err)
	}
	if user.Name != "Ana-echo" {
		t.Errorf("name got %q, want %q", user.Name, "Ana-echo")
	}
}

func TestCodec_ImplementsFastshotCodecContract(t *testing.T) {
	codec := New()

	if codec.Encode == nil {
		t.Error("Encode func is nil")
	}
	if codec.Decode == nil {
		t.Error("Decode func is nil")
	}
}
