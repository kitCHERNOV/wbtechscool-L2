package main

import (
	"fmt"
	"os"
)

// 1. Анализ строки на наличие команд, т.е. провека

func shellImplementation(terminalCommands []string) {
	//parameters := strings.Fields(terminalCommands)
	fmt.Println(terminalCommands)
	//syscall
}

func main() {
	// Get data from terminal
	commandLine := os.Args[1:]
	shellImplementation(commandLine)
}
