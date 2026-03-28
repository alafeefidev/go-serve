# Go-Serve

A middleware and CLI to test your go websites with auto refresh

## How It Works
You run with `go-serve cmd/main.go`, before runnig the script it spins a go routine with a websocket server with a non taken port also can be specified in `goserve.conf`, then it runs the script which needs to use the `Serve` middleware for routes that serves html and they will automatically update on change.

