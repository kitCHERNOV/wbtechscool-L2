package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"shell/commands"
	"strings"
)

// 1. Анализ строки на наличие команд, т.е. провека слайса команд/аргументов
// 2. Рекурсивный вызов позволит последовательно обрабатывать последовательность команд,
// так чтобы каждый результат мог быть обработан следующей командой

// TODO: Сделать реализацию каждой команды shell оболочки
// TODO: Функция вызова этих команд по необходимости (добавить pipeline - пока опционально)

func shellManager(text string) {
	// prepare get interupt syscall
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer func() {
		cancel()
	}()

	// Get sentences of commands with their arguments
	commandLine := strings.Split(text, "|")

	cmdStack, err := commands.NewCommandStack(commandLine)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}


	cmdStack.Run(ctx)
}

func shellLaunch() {
	// intercept interupt signal
	sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, os.Interrupt)
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
