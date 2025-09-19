package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"

	//"sort"
	"strconv"
)

// Набор одноименных с флагами переменных
var (
	k                int
	n, r, u, b, c, h bool
	M                bool
)

func toReturnCorrectLessFuncs(strs []string) []func(i, j int) bool {
	var lessFuncs []func(i, j int) bool = make([]func(i int, j int) bool, 0)
	if k > 0 {
		lessFuncs = append(lessFuncs, func(i, j int) bool {
			return strs[i][k] < strs[j][k]
		})
	}
	if k < 0 {
		lessFuncs = append(lessFuncs, func(i, j int) bool {
			return strs[i][0] < strs[j][0]
		})
	}
	if n {
		lessFuncs = append(lessFuncs, func(i, j int) bool {
			n1, n2 := 0., 0.
			for _, v := range strs[i] {
				n1 += float64(v)
			}
			for _, v := range strs[j] {
				n2 += float64(v)
			}
			return n1 < n2
		})
	}
	return lessFuncs
}

func main() {
	// declare a few variables
	var (
		in             = bufio.NewScanner(os.Stdin)
		readline       = bufio.NewReader(os.Stdin)
		sentencesArray []string
	)

	// Парсинг всех параметров
	flag.IntVar(&k, "k", 0, "number of column")
	flag.BoolVar(&n, "n", false, "number of rows")
	flag.BoolVar(&r, "r", false, "number of rows")
	flag.BoolVar(&u, "u", false, "number of updates")
	flag.BoolVar(&M, "m", false, "number of updates")
	flag.BoolVar(&b, "b", false, "number of updates")
	flag.BoolVar(&c, "c", false, "number of updates")
	flag.BoolVar(&h, "h", false, "number of updates")
	flag.Parse()

	// first scanning of text for recognizing file path
	fmt.Println("Where are you want to read data from?")
	fmt.Println("If you want to write data into console, just press enter w/o any symbols")
	in.Scan()
	if path := in.Text(); path != "" {
		file, err := os.Open(path)
		if err != nil {
			log.Fatalln(err)
		}
		defer file.Close()
		var data []byte
		data, err = io.ReadAll(file)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println("Readed data:", string(data))
		arr := strings.Split(string(data), "\n")
		for _, str := range arr {
			sentencesArray = append(sentencesArray, str)
		}

	} else {
		var amountOfStrings int
		fmt.Println("And How many strings you want to write")
		fmt.Println("Please enter sentences with enter pressing at end of each sentence")
		in.Scan() // Reading of amount of strings parameter
		amountOfStrings, _ = strconv.Atoi(in.Text())

		var data []byte
		for i := 0; i < amountOfStrings; i++ {
			data, _, _ = readline.ReadLine()
			sentencesArray = append(sentencesArray, string(data))
		}
	}

	// until sorting
	fmt.Println("Before sorting")
	for _, v := range sentencesArray {
		fmt.Println(v)
	}

	// Sorting by flags
	for _, f := range toReturnCorrectLessFuncs(sentencesArray) {
		sort.Slice(sentencesArray, f)
	}

	fmt.Println("After sorting")
	for _, v := range sentencesArray {
		fmt.Println(v)
	}

}
