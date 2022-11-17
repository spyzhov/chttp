# Middlewares

Middleware can be created with the defining function:

```go
func(request *http.Request, next func(request *http.Request) (*http.Response, error)) (*http.Response, error)
```

**Example:**

```go
client := chttp.NewClient(nil)
client.With(func(request *http.Request, next func(request *http.Request) (*http.Response, error)) (*http.Response, error) {
    // before action
    response, err := next(request)
    // before action
    return response, err
})
```

## CustomHeaders

Adds a custom headers based on the request.

**Example:**

```go
client := chttp.NewClient(nil)
client.With(middleware.CustomHeaders(func(request *http.Request) map[string]string {
    if request.Method == http.MethodPost {
        return map[string]string{
            "Accept": "*/*",
        }
    }
    return nil
}))
```

## Debug

**NB!** Don't use it in production!

Dumps requests and responses in the logs.

**Example:**

```go
client := chttp.NewClient(nil)
client.With(middleware.Debug(true, nil))
```

## Headers

Adds a static headers.

**Example:**

```go
client := chttp.NewClient(nil)
client.With(middleware.Headers(map[string]string{
    "Accept": "*/*",
}))
```

## JSON

Adds a `Content-Type` and `Accept` headers with the `application/json` value.

**Example:**

```go
client := chttp.NewClient(nil)
client.With(middleware.JSON())
```

## OpenTracing

Adds an OpenTracing logs and headers to the request.

Source: [https://github.com/spyzhov/chttp-middleware-opentracing](https://github.com/spyzhov/chttp-middleware-opentracing)

```go
import (
	// ...
	middleware "github.com/spyzhov/chttp-middleware-opentracing"
)

client := chttp.NewClient(nil)
client.With(middleware.Opentracing())
```

## Trace

Adds short logs on each request.

**Example:**

```go
client := chttp.NewClient(nil)
client.With(middleware.Trace(nil))
```

# TBD

 - [ ] `Cache(interface{Get(string,interface{}), Set(string,interface{})}, GetKey func(*http.Request) string)`
  
   Client-wide caching layer.
 
   If `GetKey` return a blank string, then do not cache.
 - [ ] `Retry(GetCount func(*http.Request) int, BeforeRetry func (*http.Request, int) string)`

   Automatically send a retry request in case of failure.
 
   `GetCount` - returns the max retry amount. 
 
   `BeforeRetry` - should be called before retry.

   **TODO**: think about a structure as argument
   `struct {GetCount func(*http.Request) int, BeforeRetry func (*http.Request, int) string}`
