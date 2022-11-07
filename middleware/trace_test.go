package middleware

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/spyzhov/chttp"
)

func ExampleTrace() {
	client := chttp.NewClient(nil)
	client.With(Trace(nil))
}

func TestTrace(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	logger := &testLogger{data: make(map[string]interface{})}
	client := chttp.NewClient(nil)
	client.With(Trace(logger))
	resp, err := client.GET(context.Background(), server.URL+"/test")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("wrong response code: %d", resp.StatusCode)
	}
	if logger.data["method"] != "GET" {
		t.Errorf("wrong log data: %q = %v", "method", logger.data["method"])
	}
	if logger.data["host"] == nil {
		t.Errorf("not set log data: %q", "host")
	}
	if logger.data["path"] != "/test" {
		t.Errorf("wrong log data: %q = %v", "path", logger.data["path"])
	}
	if logger.data["request_time"] == nil {
		t.Errorf("not set log data: %q", "request_time")
	}
	if logger.data["status_code"] != http.StatusOK {
		t.Errorf("wrong log data: %q = %v", "status_code", logger.data["status_code"])
	}
	if logger.data["content_length"] == nil {
		t.Errorf("not set log data: %q", "content_length")
	}
}

type testLogger struct {
	data  map[string]interface{}
	print string
}

func (t *testLogger) WithContext(_ context.Context) Logger {
	return t
}

func (t *testLogger) WithField(name string, value interface{}) Logger {
	t.data[name] = value
	return t
}

func (t *testLogger) Printf(format string, args ...interface{}) {
	t.print = fmt.Sprintf(format, args...)
}
