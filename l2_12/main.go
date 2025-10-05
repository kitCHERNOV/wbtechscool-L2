package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
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
			matched = strings.Contains(sentence, searchPattern)
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
	flag.BoolVar(&flags.v, "f", false, "finding strict search line")
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
		var (
			amountOfSentences int
			err               error
		)
		fmt.Println("Enter amount of sentences: ")
		sc := bufio.NewScanner(os.Stdin)
		// scan amount of sentences
		for {
			sc.Scan()
			amountOfSentences, err = strconv.Atoi(sc.Text())
			if err != nil {
				log.Println("Error: ", err)
				fmt.Println("Please enter a valid amount of sentences")
				continue
			}
			break
		}
		sentences = make([]string, amountOfSentences)
		for i := range sentences {
			// fill sentence's array
			sc.Scan()
			sentences[i] = sc.Text()
		}
	}

	// Call grep func
	res := grepFunc(subString, sentences, flags)
	if flags.c {
		fmt.Printf("amount of matches = %d;\n", res.Amount)
	} else {
		// printing of all match strings
		for i, matchStr := range res.Strings {
			if flags.n {
				fmt.Printf("str num: %d; match: %s\n", res.NumOfStrings[i], matchStr)
				continue
			}
			fmt.Printf("match: %s\n", matchStr)
		}
	}
}
