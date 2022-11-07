# cHTTP

cHTTP is a Golang HTTP client wrapper provided with middleware.

Middleware level is based on the modified [http.Client.Transport](https://pkg.go.dev/net/http#Client)
([http.RoundTripper](https://pkg.go.dev/net/http#RoundTripper)).

## Client

`chttp.Client` an HTTP client wrapper provided with the full list of unified HTTP methods and global method `Do` from
the basic [http.Client](https://pkg.go.dev/net/http#Client).

```go
func (c *Client) GET    (ctx context.Context, url string, body ...[]byte) (*http.Response, error)
func (c *Client) HEAD   (ctx context.Context, url string, body ...[]byte) (*http.Response, error)
func (c *Client) POST   (ctx context.Context, url string, body ...[]byte) (*http.Response, error)
func (c *Client) PUT    (ctx context.Context, url string, body ...[]byte) (*http.Response, error)
func (c *Client) PATCH  (ctx context.Context, url string, body ...[]byte) (*http.Response, error)
func (c *Client) DELETE (ctx context.Context, url string, body ...[]byte) (*http.Response, error)
func (c *Client) CONNECT(ctx context.Context, url string, body ...[]byte) (*http.Response, error)
func (c *Client) OPTIONS(ctx context.Context, url string, body ...[]byte) (*http.Response, error)
func (c *Client) TRACE  (ctx context.Context, url string, body ...[]byte) (*http.Response, error)
```

All methods have several `body` arguments as the latest one, but the valid is only the first one, such implementation
was done to simplify the interface of the function and calls without a body required.

Usage example:

```go
package main

import (
	"context"
	"fmt"
	"net/http/httputil"
	
	"github.com/spyzhov/chttp"
)

func main() {
    client := chttp.NewClient(nil)
	response, _ := client.HEAD(context.TODO(), "https://go.dev/")
	
	data, _ := httputil.DumpResponse(response, false)
	fmt.Println(string(data))
}
```

### JSON

`chttp.JSONClient` is a `chttp.Client` wrapper to simplify the routine of marshaling/unmarshalling structures
using `JSON` in requests. The `JSON` client takes the responsibility to marshal the request object and unmarshal the
response body to a given object structure.

If the request fails, due to (un)marshaling or HTTP request (`status code >= 300`), the result will be wrapped with 
the `*chttp.Error`, which will have the `http.Response`, its body, and the basic error if it exists.

`chttp.JSONClient` provided with the full list of unified HTTP methods:

```go
func (c *JSONClient) GET    (ctx context.Context, url string, body interface{}, result interface{}) error
func (c *JSONClient) HEAD   (ctx context.Context, url string, body interface{}, result interface{}) error
func (c *JSONClient) POST   (ctx context.Context, url string, body interface{}, result interface{}) error
func (c *JSONClient) PUT    (ctx context.Context, url string, body interface{}, result interface{}) error
func (c *JSONClient) PATCH  (ctx context.Context, url string, body interface{}, result interface{}) error
func (c *JSONClient) DELETE (ctx context.Context, url string, body interface{}, result interface{}) error
func (c *JSONClient) CONNECT(ctx context.Context, url string, body interface{}, result interface{}) error
func (c *JSONClient) OPTIONS(ctx context.Context, url string, body interface{}, result interface{}) error
func (c *JSONClient) TRACE  (ctx context.Context, url string, body interface{}, result interface{}) error
```

Usage example:

```go
package main

import (
	"context"
	"fmt"
	
	"github.com/spyzhov/chttp"
)

func main() {
	var fact struct {
		Fact string `json:"fact"`
	}
	client := chttp.NewJSON(nil) // same as chttp.NewClient(nil).JSON()
	_ = client.GET(context.TODO(), "https://catfact.ninja/facts?limit=1&max_length=140", nil, &fact)

	fmt.Println(fact.Fact)
}
```

## Middleware

Middlewares are the cHTTPs main driver. Adding various middlewares gives the ability to manage requests, adding tracing,
logs, and headers transparently into all requests through the cHTTP clients.

Current interface of the middleware is based on 2 types of functions:

```go
// RoundTripper is a RoundTrip function implementation of the http.RoundTripper interface.
type RoundTripper func(request *http.Request) (*http.Response, error)

// Middleware is an extended interface to the RoundTrip function of the http.RoundTripper interface.
type Middleware func(request *http.Request, next RoundTripper) (*http.Response, error)
```

Usage example:

```go
package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/spyzhov/chttp"
	"github.com/spyzhov/chttp/middleware"
)

func main() {
	var fact struct {
		Fact string `json:"fact"`
	}
	client := chttp.NewJSON(nil)
	client.With(middleware.JSON(), middleware.Debug(true, nil))
	client.With(func(request *http.Request, next chttp.RoundTripper) (*http.Response, error) {
		fmt.Println("Before the request")
		resp, err := next(request)
		fmt.Println("After the request")
		return resp, err
	})
	_ = client.GET(context.TODO(), "https://catfact.ninja/facts?limit=1&max_length=140", nil, &fact)

	fmt.Println(fact.Fact)
}
```

### List of middlewares

#### CustomHeaders

Adds a custom headers based on the request.

**Example:** 

```go
chttp.NewClient(nil).
	With(middleware.CustomHeaders(func(request *http.Request) map[string]string {
		if request.Method == http.MethodPost {
			return map[string]string{
				"Accept": "*/*",
			}
		}
		return nil
	}))
```

#### Debug

**NB!** Don't use it in production!

Dumps requests and responses in the logs.

**Example:** 

```go
chttp.NewClient(nil).
	With(middleware.Debug(true, nil))
```

#### Headers

Adds a static headers.

**Example:** 

```go
chttp.NewClient(nil).
	With(middleware.Headers(map[string]string{
		"Accept": "*/*",
	}))
```

#### JSON

Adds a `Content-Type` and `Accept` headers with the `application/json` value.

**Example:** 

```go
chttp.NewClient(nil).
	With(middleware.JSON())
```

#### Trace

Adds short logs on each request.

**Example:** 

```go
chttp.NewClient(nil).
	With(middleware.Trace(nil))
```

# License

MIT licensed. See the [LICENSE](LICENSE) file for details.
