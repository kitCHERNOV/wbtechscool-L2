package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"shell/commands"
	"strings"
	"syscall"
)

// 1. Анализ строки на наличие команд, т.е. провека слайса команд/аргументов
// 2. Рекурсивный вызов позволит последовательно обрабатывать последовательность команд,
// так чтобы каждый результат мог быть обработан следующей командой

// TODO: Сделать реализацию каждой команды shell оболочки
// TODO: Функция вызова этих команд по необходимости (добавить pipeline - пока опционально)

func shellManager(closeChan chan os.Signal, text string) (error) {
	// Get sentences of commands with their arguments
	commandLine := strings.Split(text, "|")
	// init shut down channel
	shutDownChannel := initCloseChannel(closeChan)
	select {
	case <-shutDownChannel:
		return fmt.Errorf("shut down error signal")
	default:
		cmdStack, err := commands.NewCommandStack(commandLine)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		}
	
		cmdStack.Run(shutDownChannel)
	}

	return nil
}

func initCloseChannel(signCh <-chan os.Signal) <-chan struct{} {
	var closeChan = make(chan struct{})
	go func() {
		select {
		case <-signCh:
			closeChan <- struct{}{}
		}
	}()

	return closeChan
}

func shellLaunch() {
	// prepare ctrl+c syscall 
	ctrlcSignalChannel := make(chan os.Signal, 1)
	signal.Notify(ctrlcSignalChannel, syscall.SIGINT)

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

		err := shellManager(ctrlcSignalChannel, text)
		if err != nil {
			return
		}
	}
}

func main() {
	// launch shell
	shellLaunch()
}
