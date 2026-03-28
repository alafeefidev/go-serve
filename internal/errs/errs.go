package errs

import (
	"fmt"
	"os"
)

func Nerr(a ...any) {
	fmt.Fprint(os.Stderr, a...)
	os.Exit(1)
}

func Nerrf(format string, a ...any) {
	fmt.Fprintf(os.Stderr, format, a...)
	os.Exit(1)
}
