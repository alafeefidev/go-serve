# Go-Serve

A middleware and CLI to test your go websites with auto refresh

## Installation
1. Install the middleware library
```bash
go get ...
```
2. Install the CLI tool
```bash
go install ...
```

## Running It
1. Use the middleware with any package that supports the net/http middlewares like chi or the net/http ServeMux
```go
// With stdlib
import "github.com/alafeefidev/go-serve/middleware"
mux := http.NewServeMux()
serve := middleware.Serve(mux)
```
```go
// With chi
#TODO
```
2. Run the CLI
```cmd
# you can write the filename with or without the .go extension
go-serve main.go
```
## How It Works
You run with `go-serve cmd/main.go`, before runnig the script it spins a go routine with a websocket server with a non taken port also can be specified in `goserve.conf`, then it runs the script which needs to use the `Serve` middleware for routes that serves html and they will automatically update on change.

