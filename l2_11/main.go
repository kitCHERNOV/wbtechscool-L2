package main

import (
	"fmt"
	"sort"
	"strings"
)

func anagramFind(set []string) map[string][]string {
	sort.Slice(set, func(i, j int) bool {
		s1 := strings.ToLower(set[i])
		s2 := strings.ToLower(set[j])
		return s1 < s2
	})
	var dictionary = make(map[string][]string)
	for _, v := range set {
		newStr := []rune(v)
		sort.Slice(newStr, func(i, j int) bool {
			return newStr[i] < newStr[j]
		})
		dictionary[string(newStr)] = append(dictionary[string(newStr)], v)
	}

	var finalMap = make(map[string][]string)
	for _, v := range dictionary {
		if len(v) == 1 {
			continue
		}
		finalMap[v[0]] = v
	}

	return finalMap
}

func main() {
	var set = []string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик", "стол"}
	m := anagramFind(set)
	for k, v := range m {
		fmt.Printf("%s:%v\n", k, v)
	}
}
