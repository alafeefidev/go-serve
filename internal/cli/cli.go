package cli

import (
	"fmt"
	"os"
)

func Run() {
	// [1] go-serve
	// [2] run/init
	// [3] file/-
	printSplash()
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, "please provide a command")
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "run":
		runCMD(os.Args[2:])
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

func printSplash() {
	fmt.Fprintln(os.Stdout,`
   ▄██████▄   ▄██████▄     ▄████████    ▄████████    ▄████████  ▄█    █▄     ▄████████ 
  ███    ███ ███    ███   ███    ███   ███    ███   ███    ███ ███    ███   ███    ███ 
  ███    █▀  ███    ███   ███    █▀    ███    █▀    ███    ███ ███    ███   ███    █▀  
 ▄███        ███    ███   ███         ▄███▄▄▄      ▄███▄▄▄▄██▀ ███    ███  ▄███▄▄▄     
▀▀███ ████▄  ███    ███ ▀███████████ ▀▀███▀▀▀     ▀▀███▀▀▀▀▀   ███    ███ ▀▀███▀▀▀     
  ███    ███ ███    ███          ███   ███    █▄  ▀███████████ ███    ███   ███    █▄  
  ███    ███ ███    ███    ▄█    ███   ███    ███   ███    ███ ███    ███   ███    ███ 
  ████████▀   ▀██████▀   ▄████████▀    ██████████   ███    ███  ▀██████▀    ██████████                       
	`)
}