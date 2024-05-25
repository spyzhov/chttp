package chttp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func ExampleNewGenericJSONClient() {
	type example struct {
		Data []struct {
			Fact string `json:"fact"`
		} `json:"data"`
	}

	client := NewJSON(nil)

	fact, _ := NewGenericJSONClient[example](client).GET(context.TODO(), "https://catfact.ninja/facts?limit=1&max_length=140", nil)
	fmt.Println(fact.Data[0].Fact)
}

func TestGenericJSONClient_Method(t *testing.T) {
	type example struct {
		Method string `json:"method"`
	}
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(writer).Encode(example{Method: request.Method})
	}))
	defer server.Close()
	ctx := context.Background()
	client := NewGenericJSONClient[example](nil)

	type testCase[Result any] struct {
		name string
	}
	tests := []testCase[example]{
		{name: http.MethodGet},
		// {name: http.MethodHead},
		{name: http.MethodPost},
		{name: http.MethodPut},
		{name: http.MethodPatch},
		{name: http.MethodDelete},
		{name: http.MethodConnect},
		{name: http.MethodOptions},
		{name: http.MethodTrace},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			method := client.Method(tt.name)
			want := example{Method: tt.name}
			if got, err := method(ctx, server.URL, nil); !reflect.DeepEqual(got, want) {
				t.Errorf("Method() = %v, want %v", got, want)
			} else if err != nil {
				t.Errorf("Method() error: %v", err)
			}
		})
	}
}
