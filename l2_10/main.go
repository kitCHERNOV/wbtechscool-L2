package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"sort"
)

// Набор одноименных с флагами переменных
var (
	k                int
	n, r, u, b, c, h bool
	M                bool
)

func toReturnCorrectLessFuncs() []func(i, j int) bool {
	var lessFuncs []func(i, j int) bool
	if k > 0 {
		return func(i, j int) bool {
			return
		}
	}
}

func main() {
	// declare a few variables
	var (
		in             = bufio.NewScanner(os.Stdin)
		readline       = bufio.NewReader(os.Stdin)
		sentencesArray []string
	)

	// Парсинг всех параметров
	flag.Parse()
	flag.IntVar(&k, "k", 0, "number of column")
	flag.BoolVar(&n, "n", false, "number of rows")
	flag.BoolVar(&r, "r", false, "number of rows")
	flag.BoolVar(&u, "u", false, "number of updates")
	flag.BoolVar(&M, "m", false, "number of updates")
	flag.BoolVar(&b, "b", false, "number of updates")
	flag.BoolVar(&c, "c", false, "number of updates")
	flag.BoolVar(&h, "h", false, "number of updates")

	// first scanning of text for recognizing file path
	fmt.Println("Where are you want to read data from?")
	fmt.Println("If you want to write data into console, just press enter w/o any symbols")
	in.Scan()
	if path := in.Text(); path != "" {

	} else {
		fmt.Println("And How many strings you want to write")
		fmt.Println("Please enter sentences with enter pressing at end of each sentence")
		var data []byte
		for {
			data, _, _ = readline.ReadLine()
		}
	}
	sort.Slice()
}
