package chttp

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestError_Error(t *testing.T) {
	var err *Error
	if err.Error() != "" {
		t.Errorf("broken error")
	}
}

func TestErrorUnmarshalTo(t *testing.T) {
	type example struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	expected := example{
		Code:    123,
		Message: "example data",
	}

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		writer.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(writer).Encode(expected)
	}))
	defer server.Close()

	err := NewJSON(nil).GET(context.Background(), server.URL, nil, nil)
	if err == nil {
		t.Errorf("should be error")
		return
	}
	result, pErr := ErrorUnmarshalTo[example](err)
	if pErr != nil {
		t.Errorf("ErrorUnmarshalTo() error: %v", pErr)
		return
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("ErrorUnmarshalTo() = %v, want %v", result, expected)
	}
}
