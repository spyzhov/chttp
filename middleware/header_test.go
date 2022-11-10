package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/spyzhov/chttp"
)

func ExampleHeaders() {
	client := chttp.NewClient(nil)
	client.With(Headers(map[string]string{
		"Accept": "*/*",
	}, false))
}

func TestHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Header.Get("Header") != "static" {
			t.Errorf("wrong header")
		}
		writer.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()
	client := chttp.NewClient(nil)
	client.With(Headers(map[string]string{
		"Header": "static",
	}, true))
	resp, err := client.GET(context.Background(), server.URL)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("wrong response code: %d", resp.StatusCode)
	}
}
