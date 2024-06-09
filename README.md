# xhttp

Extended http contains functions not included in golang standard library,
yet useful for speed of development.

## Graceful shutdown

Performs graceful shutdown of `http.Server`

```go
package main

import (
    "fmt"
    "net/http"
    "os"

    "github.com/prathoss/xhttp"
)

func main() {
    s := &http.Server{
        Addr: ":8080",
    }

    if err := xhttp.ServeWithShutdown(s); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}
```

## Handler

`xhttp.Handler` implements `http.Handler` and allows returning error which are translated to status codes.

```go
package main

import (
    "net/http"

    "github.com/prathoss/xhttp"
)

func main() {
    mux := http.NewServeMux()
    mux.Handle("GET /", xhttp.HttpHandler(func(w http.ResponseWriter, r *http.Request) (any, error) {
        return nil, nil
    }))
    mux.Handle("GET /not-found", xhttp.HttpHandler(func(w http.ResponseWriter, r *http.Request) (any, error) {
        return nil, xhttp.NewNotFoundError("not found message")
    }))
}
```

In previous example the `/` path will return status 204 No content.

The `/not-found` will return status 404 Not found
```
Content-Type: application/problem+json

{
  "status": 404,
  "type": "https://datatracker.ietf.org/doc/html/rfc7231#section-6.5.4",
  "title": "not found message"
}
```

Custom errors with status code and body can be created, by implementing `xhttp.HttpProblemWriter` and `error` interface.

## Logging middleware

Logs requests and response parameters with `log/slog` package.

```go
package main

import (
    "net/http"

    "github.com/prathoss/xhttp"
    "github.com/prathoss/xhttp/middleware"
)

func main() {
    mux := http.NewServeMux()
    mux.Handle("GET /", xhttp.HttpHandler(func(w http.ResponseWriter, r *http.Request) (any, error) {
        return nil, nil
    }))
    s := &http.Server{
        Addr: ":8080",
        Handler: middleware.LoggingHandler(mux),
    }
    s.ListenAndServe()
}
```
