package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Flags struct {
	c bool // flag to write amount of desired lines
}

func (f *Flags) getCParameter() bool {
	return f.c
}

func grepFunc(substring string, sentences []string, flags Flags) []string {

	for _, sentence := range sentences {
		if strings.Contains(sentence, substring) {

		}
	}
	return sentences
}

func main() {
	// flag's parameters
	var flags Flags
	// sentences
	var sentences []string

	flag.BoolVar(&flags.c, "c", false, "number of matches")
	flag.Parse()

	// file path to read from
	filePath := flag.Arg(0)

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
	grepFunc(sentences, flags)
}
