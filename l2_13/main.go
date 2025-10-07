package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type flags struct {
	f string
	d string // new separator
	s bool   //
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

func toFindStringFiedlsLikeSliceOfStrings(indexes []int, strs []string) []string {
	if len(indexes) == 0 {
		return strs
	}
	var necessaryStrings []string = make([]string, len(indexes))
	stopInd := 0
	for i, index := range indexes {
		if index-1 >= len(strs) || index-1 < 0 {
			stopInd = i
			break
		}
		necessaryStrings[i] = strs[index-1] // to avoid an empty element, first index which program can get is 1 -> Though to substruct 1 from index value
	}
	return necessaryStrings[:stopInd]
}

func main() {
	// flags struct
	var flags flags
	//var (
	//	sentences           []string
	//	resultCuttedStrings []string
	//)
	// parse flags
	flag.StringVar(&flags.f, "f", "", "getting fields parameter")
	flag.StringVar(&flags.d, "d", "\t", "use inputted separator")
	flag.BoolVar(&flags.s, "s", false, "get parameter to ignore strings without separator")
	flag.Parse()

	if flags.f == "" {
		fmt.Println("you must specify a list fields")
		return
	}

	filePath := flag.Arg(0)
	// flag analyze:
	separator := flags.d
	isIgnoreStringsWithoutSeparator := flags.s
	// fields numbers
	fields, err := flags.parseFParameter()
	if err != nil {
		// also can get an empty array
		fmt.Println(err)
	}

	if filePath == "" {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := scanner.Text()
			lines := strings.Split(line, separator)
			if len(lines) == 1 && isIgnoreStringsWithoutSeparator { // slice of strings is not changed
				// string is ignored
				continue
			}
			// Get necessary strings as slice of strings
			strNecessaryFields := toFindStringFiedlsLikeSliceOfStrings(fields, lines)
			result := strings.Join(strNecessaryFields, separator)
			// cleaning from last separator
			//result = result[:len(result)]
			fmt.Println(result)
		}
	} else {
		// read from file
		file, err := os.OpenFile(filePath, os.O_RDONLY, 0666)
		if err != nil {
			log.Fatal(err)
		}
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			lines := strings.Split(line, separator)
			if len(lines) == 1 && isIgnoreStringsWithoutSeparator { // slice of strings is not changed
				// string is ignored
				continue
			}
			// Get necessary strings as slice of strings
			strNecessaryFields := toFindStringFiedlsLikeSliceOfStrings(fields, lines)
			result := strings.Join(strNecessaryFields, separator)
			// cleaning from last separator
			//result = result[:len(result)]
			fmt.Println(result)
		}
	}
}
