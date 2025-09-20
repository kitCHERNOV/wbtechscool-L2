package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func Parser(str string) (finstr string) {
	var preparedElem rune = 0
	for _, v := range str {
		// Первое - проверять на вход буквенные и цифровые значения.
		num, err := strconv.Atoi(string(v))
		if err != nil {
			if preparedElem != 0 {
				finstr += string(preparedElem)
			}
			preparedElem = v
		} else {
			if preparedElem != 0 {
				for i := 0; i < num; i++ {
					finstr += string(preparedElem)
				}
			}
			// Обнуляем значение
			preparedElem = 0
		}
	}
	// Если в переменной что-то осталось (последняя буква строки)
	if preparedElem != 0 {
		finstr += string(preparedElem)
	}
	return
}

func main() {
	sc := bufio.NewScanner(os.Stdin)
	sc.Scan()
	gotString := sc.Text()
	res := Parser(gotString)
	fmt.Println(res)
}
