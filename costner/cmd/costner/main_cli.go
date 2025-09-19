//go:build cli
// +build cli

package main

import (
	"costner/internal/cli"
)

func main() {
	cli := cli.NewCLI()
	cli.Run()
}