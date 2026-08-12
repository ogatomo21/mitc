package main

import (
	"os"
)

// version is replaced by the release workflow with -ldflags.
var version = "v1.0"

func main() {
	os.Exit(Run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr, version))
}
