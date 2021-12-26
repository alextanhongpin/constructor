package main

import (
	"fmt"
	"os"

	"github.com/alextanhongpin/constructor"
)

func main() {
	if err := constructor.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "constructor: %s\n", err)
		os.Exit(1)
	}
}
