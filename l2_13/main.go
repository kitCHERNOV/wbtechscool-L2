package main

import (
	"flag"
	"sort"
	"strconv"
	"strings"
)

type flags struct {
	f string
	// d string
	// s string
}

func (f *flags) parseFParameter() ([]int, error) {
	var indexes []int
	var m = make(map[int]bool)
	indexesInRowString := strings.Split(f.f, ",")
	for _, strElem := range indexesInRowString {
		separatedStrElem := strings.Split(strElem, "-")
		if len(separatedStrElem) == 1 {
			ind, err := strconv.Atoi(separatedStrElem[0])
			if err != nil {
				return []int{}, err
			}
			//indexes = append(indexes, ind)
			m[ind] = true
		} else { // else len of separatedStrElem equal 2
			start, err := strconv.Atoi(separatedStrElem[0])
			if err != nil {
				return []int{}, err
			}
			end, err := strconv.Atoi(separatedStrElem[1])
			if err != nil {
				return []int{}, err
			}
			for i := start; i <= end; i++ {
				//indexes = append(indexes, i)
				m[i] = true
			}
		}
	}
	// rewrite to indexes array from map
	for ind, _ := range m {
		indexes = append(indexes, ind) // fields are numbered since 1
	}
	sort.Ints(indexes) // sort index's order
	return indexes, nil
}

func main() {
	// flags struct
	var flags flags
	//var separator = "\t"
	// parse flags
	flag.StringVar(&flags.f, "f", "", "getting fields parameter")
	flag.Parse()

	//scanner := bufio.NewScanner(os.Stdin)
	//for scanner.Scan() {
	//	line := scanner.Text()
	//	lines := strings.Split(line, separator)
	//
	//}
}
