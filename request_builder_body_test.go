package fastshot

import (
	"bytes"
	"github.com/opus-domini/fast-shot/constant"
	"strings"
	"testing"
)

func TestRequestBodyBuilder_AsReader(t *testing.T) {
	// Arrange
	client := DefaultClient("https://example.com")
	body := bytes.NewBuffer([]byte("test body"))
	// Act
	builder := client.POST("/test").
		Body().AsReader(body)
	// Assert
	if builder.request.body == nil {
		t.Errorf("Body not set correctly")
	}
}

func TestRequestBodyBuilder_AsString(t *testing.T) {
	// Arrange
	client := DefaultClient("https://example.com")
	body := "test string"
	// Act
	builder := client.POST("/path").
		Body().AsString(body)
	// Assert
	if builder.request.body == nil {
		t.Errorf("Body String not set correctly")
	}
}

func TestRequestBodyBuilder_AsJSON(t *testing.T) {
	// Arrange
	client := DefaultClient("https://example.com")
	body := map[string]string{"key": "value"}
	// Act
	builder := client.POST("/test").
		Body().AsJSON(body)
	// Assert
	if builder.request.body == nil {
		t.Errorf("Body JSON not set correctly")
	}
}

func TestRequestBodyBuilder_AsJSON_Error(t *testing.T) {
	// Arrange
	client := DefaultClient("https://example.com")
	body := func() {}
	// Act
	builder := client.POST("/path").
		Body().AsJSON(body)
	// Assert
	if len(builder.request.validations) != 1 || !strings.Contains(builder.request.validations[0].Error(), constant.ErrMsgMarshalJSON) {
		t.Errorf("Body JSON didn't capture the marshaling error")
	}
}
