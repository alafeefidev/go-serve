package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/alafeefidev/go-serve/internal/cli"
	"github.com/alafeefidev/go-serve/internal/errs"
)

func main() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-sig
		errs.Nerr("\nGood Bye ^_^\n")
	}()

	cli.Run()

}
