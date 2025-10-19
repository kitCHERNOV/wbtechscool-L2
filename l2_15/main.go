package main

import (
	"bufio"
	"fmt"
	"os"
	"shell/commands"
	"strings"
)

// 1. Анализ строки на наличие команд, т.е. провека слайса команд/аргументов
// 2. Рекурсивный вызов позволит последовательно обрабатывать последовательность команд,
// так чтобы каждый результат мог быть обработан следующей командой

// TODO: Сделать реализацию каждой команды shell оболочки
// TODO: Функция вызова этих команд по необходимости (добавить pipeline - пока опционально)

func shellManager(text string) {
	// Get sentences of commands with their arguments
	commandLine := strings.Split(text, "|")

	cmdStack, err := commands.NewCommandStack(commandLine)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}

	cmdStack.Run(len(commandLine) - 1)
}

func shellLaunch() {
	// get data from terminal
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		if text == "exit" {
			// shout down program
			break
		}
		if len(text) < 1 {
			continue
		}
		shellManager(text)
	}
}

func main() {
	// launch shell
	shellLaunch()
}
