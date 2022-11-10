package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/spyzhov/chttp"
)

func ExampleDebug() {
	client := chttp.NewClient(nil)
	client.With(Debug(os.Getenv("DEBUG") == "true", nil))
}

func TestDebug(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	logger := &testLogger{data: make(map[string]interface{})}
	client := chttp.NewClient(nil)
	client.With(Debug(true, logger))
	resp, err := client.GET(context.Background(), server.URL)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("wrong response code: %d", resp.StatusCode)
	}
	if logger.data["Method"] != "GET" {
		t.Errorf("wrong log data: %q = %v", "method", logger.data["method"])
	}
	fields := []string{
		"Method",
		"Host",
		"Path",
	}
	for _, field := range fields {
		if logger.data[field] == nil {
			t.Errorf("not set log data: %q", field)
		}
	}
}
