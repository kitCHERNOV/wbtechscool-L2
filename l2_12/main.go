package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"sort"
	"strings"
)

// type for returning result
type grepResult struct {
	Strings      []string
	NumOfStrings []int
	Amount       int // if flag c is set
}

type Flags struct {
	c bool // flag to write amount of desired lines
	i bool // flag to ignore register of word
	v bool // inverse of filters
	f bool // string should be strict equal the string is inputted
	n bool // print number of each string
	A int  // return A elems before match
	B int  // return B elems after match
	C int  // return C elems before and after match
}

func (gr *grepResult) addElems(str string, ind int) {
	gr.Strings = append(gr.Strings, str)
	gr.NumOfStrings = append(gr.NumOfStrings, ind)
	gr.Amount++
}

func grepFunc(substring string, sentences []string, flags Flags) grepResult {
	var result grepResult
	// prefer to -F flag solution
	searchPattern := substring
	if flags.f && flags.i {
		searchPattern = strings.ToLower(substring)
	}
	// prefer to no -f flag solution
	var regExp *regexp.Regexp
	var pattern string
	if !flags.f {
		pattern = substring
		if flags.i {
			pattern = "(?i)" + substring // (?i) делает regex нечувствительным к регистру
		}
		regExp = regexp.MustCompile(pattern)
	}

	for i, sentence := range sentences {
		// adding
		var matched bool
		if flags.f {
			if flags.i {
				matched = strings.Contains(strings.ToLower(sentence), searchPattern)
			} else {
				matched = strings.Contains(sentence, searchPattern)
			}
		} else {
			if regExp != nil {
				matched = regExp.MatchString(sentence)
			}
		}
		// matched = true and v flag = false or reverse params
		if matched != flags.v {
			result.addElems(sentences[i], i)
		}

	}
	return result
}

func main() {
	// flag's parameters
	var flags Flags
	// sentences
	var sentences []string

	flag.BoolVar(&flags.c, "c", false, "number of matches")
	flag.BoolVar(&flags.i, "i", false, "ignore of register")
	flag.BoolVar(&flags.v, "v", false, "reverse flag parameter")
	flag.BoolVar(&flags.n, "n", false, "enumerate each string")
	flag.BoolVar(&flags.f, "F", false, "finding strict search line")
	flag.IntVar(&flags.A, "A", 0, "amount of words after match")
	flag.IntVar(&flags.B, "B", 0, "amount of words before match")
	flag.IntVar(&flags.C, "C", 0, "amount of words after and before match")
	flag.Parse()

	// file path to read from
	subString := flag.Arg(0)
	filePath := flag.Arg(1)

	if filePath != "" {
		file, err := os.Open(filePath)
		if err != nil {
			log.Fatalln(err)
		}
		defer func() {
			if err := file.Close(); err != nil {
				log.Fatalln(err)
			}
		}()
		// Reading of file
		readline := bufio.NewScanner(file)
		for readline.Scan() {
			sentences = append(sentences, readline.Text())
		}
	} else {
		sc := bufio.NewScanner(os.Stdin)
		sentences = make([]string, 0)
		for sc.Scan() {
			sentences = append(sentences, sc.Text())
		}
	}

	// Call grep func
	res := grepFunc(subString, sentences, flags)

	if flags.c {
		fmt.Printf("amount of matches = %d;\n", res.Amount)
	} else {
		// Alternative solution
		var (
			A       int          = 0
			B       int          = 0
			indexes []int        // indexes for printing
			m       map[int]bool = make(map[int]bool)
		)
		if flags.A >= 1 {
			A = flags.A
		}
		if flags.B >= 1 {
			B = flags.B
		}
		if flags.C >= 1 {
			A = flags.C
			B = flags.C
		}
		for _, numOfStr := range res.NumOfStrings {
			// calculate difference parameters
			differenceAfter := int(math.Min(float64(numOfStr+A), float64(len(sentences)-1)))
			differenceBefore := int(math.Max(float64(numOfStr-B), 0))
			//if flags.n {
			if B >= 1 {
				for i := differenceBefore; i <= numOfStr; i++ {
					m[i] = true
				}
			}
			if A >= 1 {
				for i := numOfStr; i <= differenceAfter; i++ {
					m[i] = true
				}
			}
			if flags.B == 0 && flags.A == 0 && flags.C == 0 {
				m[numOfStr] = true
				//continue
			}
		}

		// print left strings
		for elem := range m {
			indexes = append(indexes, elem)
		}
		sort.Ints(indexes)
		for _, elem := range indexes {
			if flags.n {
				fmt.Printf("str num: %d; %s\n", elem, sentences[elem])
			} else {
				fmt.Printf("match: %s\n", sentences[elem])
			}
		}
	}
}
