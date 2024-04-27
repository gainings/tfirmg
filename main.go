package main

import (
	"fmt"
	"github.com/gainings/tfirg/cmd"
	"os"
)

func main() {
	if err := cmd.BaseCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[0], err)
		os.Exit(-1)
	}
}
