# Middlewares

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

## Trace

Adds short logs on each request.

**Example:**

```go
client := chttp.NewClient(nil)
client.With(middleware.Trace(nil))
```
