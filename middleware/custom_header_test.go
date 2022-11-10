package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/spyzhov/chttp"
)

func ExampleCustomHeaders() {
	client := chttp.NewClient(nil)
	client.With(CustomHeaders(func(request *http.Request) map[string]string {
		if request.Method == http.MethodPost {
			return map[string]string{
				"Accept": "*/*",
			}
		}
		return nil
	}))
}

func TestCustomHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Header.Get("CustomHeader") != "GET" {
			t.Errorf("wrong header")
		}
		writer.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()
	client := chttp.NewClient(nil)
	client.With(CustomHeaders(func(request *http.Request) map[string]string {
		return map[string]string{
			"CustomHeader": request.Method,
		}
	}))
	resp, err := client.GET(context.Background(), server.URL)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("wrong response code: %d", resp.StatusCode)
	}
}
