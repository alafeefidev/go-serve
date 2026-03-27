package cli

import (
	"fmt"
	"os"
)

func RunCLI() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "run":
		fmt.Fprint(os.Stdout, "sus")
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: go-serve <commands> [options]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("	run    Run the entry-point go file")
}
