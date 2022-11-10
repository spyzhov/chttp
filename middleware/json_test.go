package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/spyzhov/chttp"
)

func ExampleJSON() {
	client := chttp.NewClient(nil)
	client.With(JSON())
}

func TestJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Header.Get("Content-Type") != "application/json" {
			t.Errorf("wrong: Content-Type")
		}
		if request.Header.Get("Accept") != "application/json" {
			t.Errorf("wrong: Accept")
		}
		writer.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()
	client := chttp.NewClient(nil)
	client.With(JSON())
	resp, err := client.GET(context.Background(), server.URL)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("wrong response code: %d", resp.StatusCode)
	}
}
