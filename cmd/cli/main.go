package main

import (
	"blockchain/internal/cli"
)

func main() {
	// defer os.Exit(0)
	c := cli.CommandLine{}
	c.Run()
}
