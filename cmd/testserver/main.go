package main

import (
	"fmt"
	"net/http"

	"github.com/alafeefidev/go-serve/middleware"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<html><head><title>sussy</title></head><body><h1>Hello!</h1></body></html>`))
	})
	serve := middleware.Serve(mux)

	fmt.Println("Server startd at localhost:6961")
	http.ListenAndServe(":6961", serve)
}
