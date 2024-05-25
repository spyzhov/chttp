package chttp

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestJSONClient_Method(t *testing.T) {
	methods := [...]string{
		http.MethodGet,
		http.MethodHead,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodConnect,
		http.MethodOptions,
		http.MethodTrace,
	}

	type TestStruct struct {
		Foo string `json:"foo"`
		Bar string `json:"bar"`
	}

	type args struct {
		path    string
		request interface{}
		result  interface{}
	}
	tests := []struct {
		name        string
		middlewares []Middleware
		args        args
		want        interface{}
		wantStatus  int
		wantErr     bool
	}{
		{
			name:        "default",
			middlewares: nil,
			args: args{
				path:    "/",
				request: nil,
				result:  nil,
			},
			want:       nil,
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name:        "with request",
			middlewares: nil,
			args: args{
				path: "/",
				request: TestStruct{
					Foo: "bar",
					Bar: "baz",
				},
				result: new(TestStruct),
			},
			want: &TestStruct{
				Foo: "bar",
				Bar: "baz",
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name:        "marshaling error",
			middlewares: nil,
			args: args{
				path:    "/",
				request: math.Inf(0),
				result:  nil,
			},
			want:       nil,
			wantStatus: 0,
			wantErr:    true,
		},
		{
			name: "error",
			middlewares: []Middleware{
				func(request *http.Request, next func(request *http.Request) (*http.Response, error)) (*http.Response, error) {
					return nil, fmt.Errorf("test error")
				},
			},
			args: args{
				path:    "/",
				request: nil,
			},
			want:       nil,
			wantStatus: http.StatusOK,
			wantErr:    true,
		},
		{
			name:        "not found",
			middlewares: nil,
			args: args{
				path:    "/",
				request: nil,
			},
			want:       nil,
			wantStatus: http.StatusNotFound,
			wantErr:    true,
		},
		{
			name:        "broken url",
			middlewares: nil,
			args: args{
				path:    "invalid url",
				request: nil,
			},
			want:       nil,
			wantStatus: http.StatusOK,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, method := range methods {
				t.Run(method, func(t *testing.T) {
					if method == http.MethodHead {
						if tt.args.request != nil {
							t.Skipf("HEAD method with request")
						}
					}
					c := NewJSON(nil)
					c.With(tt.middlewares...)
					body := make([]byte, 0)

					if tt.want != nil {
						var err error
						body, err = json.Marshal(tt.want)
						if err != nil {
							t.Errorf("json.Marshal() error = %v", err)
						}
					}
					server := getTestServer(t, method, tt.args.path, [][]byte{body}, tt.wantStatus, body)
					defer server.Close()
					url := server.URL

					err := c.Method(method)(context.TODO(), url+tt.args.path, tt.args.request, &tt.args.result)
					if (err != nil) != tt.wantErr {
						t.Errorf("%s() error = %v, wantErr %v", method, err, tt.wantErr)
						return
					}
					if tt.wantErr {
						return
					}
					if !reflect.DeepEqual(tt.args.result, tt.want) {
						t.Errorf("%s() wrong response \nactual: %v\n want: %v", method, tt.args.result, tt.want)
					}
				})
			}
		})
	}
}

func TestJSONClient_Request(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write([]byte(`"broken`))
	}))
	defer server.Close()
	var result int
	err := NewJSON(nil).Request(context.TODO(), http.MethodGet, server.URL, nil, &result)
	if err == nil {
		t.Errorf("Request() error wanted")
		return
	}
	if err.Error() != "unmarshaling response error: unexpected end of JSON input" {
		t.Errorf("Request() wrong error: %q", err.Error())
	}
}

func TestJSONClient_Request_status(t *testing.T) {
	server := testServerJSON(http.StatusInternalServerError, 123)
	defer server.Close()
	var result int
	err := NewJSON(nil).Request(context.TODO(), http.MethodGet, server.URL, nil, &result)
	if err == nil {
		t.Errorf("Request() error wanted")
		return
	}
	if err.Error() != "wrong status code, status_code=500" {
		t.Errorf("Request() wrong error: %q", err.Error())
	}
}

func ExampleJSONClient_GET() {
	var fact struct {
		Data []struct {
			Fact string `json:"fact"`
		} `json:"data"`
	}
	client := NewJSON(nil)
	client.With(func(request *http.Request, next func(request *http.Request) (*http.Response, error)) (*http.Response, error) {
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Accept", "application/json")
		return next(request)
	})
	_ = client.GET(context.TODO(), "https://catfact.ninja/facts?limit=1&max_length=140", nil, &fact)
	fmt.Println(fact.Data[0].Fact)
}

func TestJSONClient_Clone(t *testing.T) {
	type example struct {
		Foo string `json:"foo"`
	}
	server := testServerJSON(http.StatusOK, example{Foo: "bar"})
	defer server.Close()

	ctx := context.Background()
	var result example
	first := 0
	second := 0

	client := NewJSON(
		nil,
		WithMiddleware(func(request *http.Request, next func(request *http.Request) (*http.Response, error)) (*http.Response, error) {
			first++
			return next(request)
		}),
	)

	clone := client.Clone()

	err := clone.GET(ctx, server.URL, nil, &result)
	if err != nil {
		t.Errorf("GET() error: %v", err)
		return
	}
	if first != 1 {
		t.Errorf("first middleware called wrong times: %v", first)
		return
	}

	another := client.Clone()
	another.With(func(request *http.Request, next func(request *http.Request) (*http.Response, error)) (*http.Response, error) {
		second++
		return next(request)
	})

	err = another.GET(ctx, server.URL, nil, &result)
	if err != nil {
		t.Errorf("GET() error: %v", err)
		return
	}
	if first != 2 {
		t.Errorf("first middleware called wrong times: %v != 2", first)
		return
	}
	if second != 1 {
		t.Errorf("second middleware called wrong times: %v != 1", second)
		return
	}
}

func testServerJSON(status int, response interface{}) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		writer.WriteHeader(status)
		_ = json.NewEncoder(writer).Encode(response)
	}))
}
