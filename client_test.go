package chttp

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestClient_Method(t *testing.T) {
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

	type args struct {
		path string
		body [][]byte
	}
	tests := []struct {
		name        string
		middlewares []Middleware
		args        args
		want        *http.Response
		wantErr     bool
	}{
		{
			name:        "default",
			middlewares: nil,
			args: args{
				path: "/",
				body: nil,
			},
			want: &http.Response{
				StatusCode:    200,
				Body:          io.NopCloser(bytes.NewBuffer(nil)),
				ContentLength: 0,
			},
			wantErr: false,
		},
		{
			name:        "with body",
			middlewares: nil,
			args: args{
				path: "/",
				body: [][]byte{
					[]byte("foo/bar"),
				},
			},
			want: &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewBufferString("bar-baz")),
			},
			wantErr: false,
		},
		{
			name: "error",
			middlewares: []Middleware{
				func(request *http.Request, next RoundTripper) (*http.Response, error) {
					return nil, fmt.Errorf("test error")
				},
			},
			args: args{
				path: "/",
				body: nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:        "broken url",
			middlewares: nil,
			args: args{
				path: "invalid url",
				body: nil,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, method := range methods {
				t.Run(method, func(t *testing.T) {
					if method == http.MethodHead {
						if len(tt.args.body) != 0 {
							t.Skipf("HEAD method with body")
						}
					}
					c := NewClient(nil)
					c.With(tt.middlewares...)
					url := "http://0.0.0.0:1234"
					body := make([]byte, 0)

					if tt.want != nil {
						body, _ = io.ReadAll(tt.want.Body)
						server := getTestServer(t, method, tt.args.path, tt.args.body, tt.want.StatusCode, body)
						defer server.Close()
						url = server.URL
					}

					got, err := c.Method(method)(context.TODO(), url+tt.args.path, tt.args.body...)
					if (err != nil) != tt.wantErr {
						t.Errorf("%s() error = %v, wantErr %v", method, err, tt.wantErr)
						return
					}
					if tt.wantErr {
						return
					}
					defer func() {
						_ = got.Body.Close()
					}()
					equalResponses(t, got, body, tt.want.StatusCode)
				})
			}
		})
	}
}

func TestClient_With(t *testing.T) {
	c := NewClient(nil)
	var index int
	for i := 0; i < 10; i++ {
		c.With((func(i int) Middleware {
			return func(request *http.Request, next RoundTripper) (*http.Response, error) {
				if index != i {
					t.Errorf("middleware called on wrong position: expected %d, actual %d", i, index)
				}
				index++
				request.Header.Add("x-index", fmt.Sprintf("%d", i))
				return next(request)
			}
		})(i))
	}
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		values := request.Header.Values("x-index")
		if !reflect.DeepEqual(values, []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}) {
			t.Errorf("invalid value in headers: %v", values)
		}
		writer.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	req, err := c.GET(context.TODO(), server.URL)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	_ = req.Body.Close()
}

func getTestServer(t *testing.T, method string, path string, bodies [][]byte, code int, result []byte) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != method {
			t.Errorf("invalid method: %v", request.Method)
		}
		if request.URL.Path != path {
			t.Errorf("invalid path: %v", request.URL.Path)
		}
		data, err := io.ReadAll(request.Body)
		_ = request.Body.Close()
		if err != nil {
			t.Errorf("error reading body: %s", err)
		}
		body := make([]byte, 0)
		if len(bodies) > 0 {
			body = bodies[0]
		}
		if !reflect.DeepEqual(data, body) {
			t.Errorf("invalid body:\nActual: %q\nWantTo: %q", string(data), string(body))
		}
		writer.WriteHeader(code)
		if result != nil {
			_, err = writer.Write(result)
			if err != nil {
				t.Errorf("error writing result: %s", err)
			}
		}
	}))
}

func equalResponses(t *testing.T, actual *http.Response, wantBody []byte, statusCode int) {
	if statusCode != actual.StatusCode {
		t.Errorf("http.Response: wrong StatusCode %d != %d", actual.StatusCode, statusCode)
		return
	}
	var (
		actualBody []byte
		err        error
	)
	if actual.Body != nil {
		if actualBody, err = io.ReadAll(actual.Body); err != nil {
			t.Errorf("http.Response error reading a actual.Body: %s", err)
			return
		}
	}
	if !reflect.DeepEqual(wantBody, actualBody) {
		t.Errorf("http.Response: wrong Body \nActual: %s\nWantTo: %s", actualBody, wantBody)
	}
}
