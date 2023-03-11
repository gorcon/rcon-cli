package main

import (
	"fmt"
	"os"

	"github.com/gorcon/rcon-cli/internal/executor"
)

// Version displays service version in semantic versioning (http://semver.org/).
// Can be replaced while compiling with flag `-ldflags "-X main.Version=${VERSION}"`.
var Version = "develop"

func main() {
	exec := executor.NewExecutor(os.Stdin, os.Stdout, Version)

	if err := exec.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		exec.Close()
		os.Exit(1)
	}

	exec.Close()
}
