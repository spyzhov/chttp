package chttp

import (
	"net/http"
	"net/http/cookiejar"
	"reflect"
	"testing"
	"time"
)

func TestWithTimeout(t *testing.T) {
	tests := []struct {
		name    string
		timeout time.Duration
	}{
		{
			name:    "zero",
			timeout: 0,
		},
		{
			name:    "1 sec",
			timeout: time.Second,
		},
		{
			name:    "1 min",
			timeout: time.Minute,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if client := NewClient(nil, WithTimeout(tt.timeout)); !reflect.DeepEqual(client.HTTP.Timeout, tt.timeout) {
				t.Errorf("WithTimeout() = %v, want %v", client.HTTP.Timeout, tt.timeout)
			}
		})
	}
}

func TestWithCookieJar(t *testing.T) {
	tests := []struct {
		name string
		jar  http.CookieJar
	}{
		{
			name: "nil",
			jar:  nil,
		},
		{
			name: "value",
			jar:  new(cookiejar.Jar),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if client := NewClient(nil, WithCookieJar(tt.jar)); !reflect.DeepEqual(client.HTTP.Jar, tt.jar) {
				t.Errorf("WithCookieJar()")
			}
		})
	}
}

func TestWithCheckRedirect(t *testing.T) {
	tests := []struct {
		name  string
		check func(req *http.Request, via []*http.Request) error
	}{
		{
			name:  "nil",
			check: nil,
		},
		{
			name: "value",
			check: func(req *http.Request, via []*http.Request) error {
				return nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(nil, WithCheckRedirect(tt.check))
			if reflect.ValueOf(client.HTTP.CheckRedirect).Pointer() != reflect.ValueOf(tt.check).Pointer() {
				t.Errorf("WithCheckRedirect()")
			}
		})
	}
}
