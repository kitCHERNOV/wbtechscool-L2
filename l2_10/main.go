package main

import (
	"bufio"
	"flag"
	"fmt"
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

// Функция - генератор сортировки по номеру колонки и численному формату
func getUnionCompareFunc(strs []string, k int, n bool, r bool) func(i, j int) bool {
	return func(i, j int) bool {
		field1 := strings.ToLower(strs[i])
		field2 := strings.ToLower(strs[j])

		if k > 0 {
			splitArrI := strings.Split(strings.ToLower(strs[i]), "\t")
			splitArrJ := strings.Split(strings.ToLower(strs[j]), "\t")

			if len(splitArrI) >= k && len(splitArrJ) >= k {
				field1 = splitArrI[k-1]
				field2 = splitArrJ[k-1]
			}

		}

		var less bool

		if n {
			num1, err1 := strconv.Atoi(field1)
			num2, err2 := strconv.Atoi(field2)

			// Обработка ошибки так, чтобы было хоть какое-то сравнение
			if err1 != nil || err2 != nil {
				less = len(field1) < len(field2)
			} else {
				less = num1 < num2
			}
		} else {
			less = field1 < field2
		}

		if r {
			return !less
		}
		return less
	}
}

func main() {
	// declare a few variables
	var (
		sentencesArray []string
		filePath       string
	)

	// Парсинг всех параметров
	flag.IntVar(&k, "k", 0, "number of column")
	flag.BoolVar(&n, "n", false, "number of rows")
	flag.BoolVar(&r, "r", false, "number of rows")
	flag.BoolVar(&u, "u", false, "number of updates")
	flag.Parse()
	// Если существует параметр
	filePath = flag.Arg(0)

	if filePath != "" {
		file, err := os.Open(filePath)
		if err != nil {
			log.Fatalln(err)
		}
		defer func() {
			err := file.Close()
			if err != nil {
				log.Fatalln(err)
			}
		}()

		//var data []byte
		var readline = bufio.NewScanner(file)
		for readline.Scan() {
			sentencesArray = append(sentencesArray, readline.Text())
		}
		//data, err = io.ReadAll(file)
		//if err != nil {
		//	log.Fatalln(err)
		//}
		//arr := strings.Split(string(data), "\n")
		//for _, str := range arr {
		//	sentencesArray = append(sentencesArray, str)
		//}

	} else if flag.NArg() == 0 {
		var sc = bufio.NewScanner(os.Stdin)
		//var amountOfStrings int
		fmt.Println("How many strings you want to write")
		fmt.Println("Please enter sentences with enter pressing at end of each sentence")
		//for i := 0; i < amountOfStrings; i++ {
		//	in.Scan()
		//	sentencesArray = append(sentencesArray, in.Text())
		//}
		for sc.Scan() {
			sentencesArray = append(sentencesArray, sc.Text())
		}
	}

	// until sorting
	fmt.Println("Before sorting")
	for _, v := range sentencesArray {
		fmt.Println(v)
	}

	fmt.Println(
		"##=============##" +
			"\n" +
			"##=============##")

	// Sorting by flags
	compareFunc := getUnionCompareFunc(sentencesArray, k, n, r)
	sort.Slice(sentencesArray, compareFunc)

	if u {
		var m = make(map[string]bool)
		var finalArr []string
		for _, v := range sentencesArray {
			if _, ok := m[v]; !ok {
				m[v] = true
			}
		}
		for k, _ := range m {
			finalArr = append(finalArr, k)
		}

		fmt.Println("After sorting")
		for _, v := range finalArr {
			fmt.Println(v)
		}
	} else {
		fmt.Println("After sorting")
		for _, v := range sentencesArray {
			fmt.Println(v)
		}
	}
}
