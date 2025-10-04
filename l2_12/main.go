package main

import (
	"flag"
	"os"
)

func main() {
	var (
		c int
	)

	flag.IntVar(&c, "c", 1, "number of matches")
	flag.Parse()

	// file path to read from
	filePath := flag.Arg(0)

	if filePath != "" {
		file, err := os.Open(filePath)
	}
}
