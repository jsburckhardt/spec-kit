// Package main provides the gospecify CLI entry point
package main

import (
	"os"

	"github.com/jsburckhardt/spec-kit/gospecify/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
