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
				func(request *http.Request, next RoundTripper) (*http.Response, error) {
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
					url := "http://0.0.0.0:1234"
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
					url = server.URL

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
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusInternalServerError)
		_, _ = writer.Write([]byte(`123`))
	}))
	defer server.Close()
	var result int
	err := NewJSON(nil).Request(context.TODO(), http.MethodGet, server.URL, nil, &result)
	if err == nil {
		t.Errorf("Request() error wanted")
		return
	}
	if err.Error() != "http error, status_code=500" {
		t.Errorf("Request() wrong error: %q", err.Error())
	}
}
