package commands

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"syscall"
)

type command struct {
	name   string
	args   []string
	Stdin  io.Reader
	Stdout io.Writer
}

type commandStack struct {
	stack []command
	//cmdIndexes  []int
	//argsIndexes []int
}

// Constructor of commandStack
func NewCommandStack(CommandString []string) (*commandStack, error) {
	var cmdStack commandStack = commandStack{
		stack: make([]command, 0),
		//cmdIndexes:  []int{},
		//argsIndexes: []int{},
	}
	for _, element := range CommandString {
		cmdAndArguments := strings.Fields(element)
		if len(cmdAndArguments) == 0 { // || len(cmdAndArguments[0]) == 1 {
			continue
		}
		cmdName := cmdAndArguments[0]
		cmdArgs := cmdAndArguments[1:]

		cmd := command{
			name: cmdName,
			args: cmdArgs,
		}
		cmdStack.stack = append(cmdStack.stack, cmd)
	}

	return &cmdStack, nil
}

func (cs *commandStack) echo(args []string) string {
	joinedArgs := strings.Join(args, " ")
	return joinedArgs
}

func (cs *commandStack) cd(path string) (string, error) {
	// WARN: проверить корректность восприятия пути path-ом
	err := os.Chdir(path)
	if err != nil {
		return "Not OK", fmt.Errorf("dirrectory is not found; error: %v", err)
	}
	return "OK", nil
}

func (cs *commandStack) pwd() (string, error) {
	path, err := os.Getwd()
	return path, err
}

func (cs *commandStack) ps() ([]string, error) {

}

func (cs *commandStack) kill(pidStr string) error {
	if len(pidStr) == 0 {
		return fmt.Errorf("pid is empty")
	}

	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		return fmt.Errorf("pid is not an integer; error: %v", err)
	}

	err = syscall.Kill(pid, syscall.SIGTERM)
	if err != nil {
		return fmt.Errorf("process kill error: %v", err)
	}

	return nil
}

const (
	lenOfBuffer = 50
)

func (cs *commandStack) Run(amountOfBuffers int) {

	// initialization of buffers
	for i := 0; i < amountOfBuffers-1; i++ {
		var buffer = &bytes.Buffer{}
		cs.stack[i].Stdout = buffer
		cs.stack[i+1].Stdin = buffer
	}
	// init Stdout for last command
	//var buffer bytes.Buffer
	cs.stack[len(cs.stack)-1].Stdout = os.Stdout

	// previosReturnedResult := ""
	for i, cmd := range cs.stack {
		// read data from stdin
		//inputData := make([]byte, lenOfBuffer)
		//// first command is nil pointer
		//if i != 0 {
		//	_, err := cmd.Stdin.Read(inputData)
		//	if err != nil {
		//		fmt.Printf("Error reading input data: %v\n", err)
		//	}
		//}
		_ = i

		switch cmd.name {
		case "cd":
			if len(cmd.args) == 0 {
				continue
			}
			arg := cmd.args[0]
			res, err := cs.cd(arg)
			if err != nil {
				fmt.Printf("Error: %s\n", err)
			}
			// write to buf for the next command
			// last command
			_, err = cmd.Stdout.Write([]byte(res))
			if err != nil {
				fmt.Printf("Error: %s\n", err)
			}

		case "pwd":
			res, err := cs.pwd()
			if err != nil {
				fmt.Printf("Error: %s\n", err)
			}
			// write to buf for the next command
			_, err = cmd.Stdout.Write([]byte(res))
			if err != nil {
				fmt.Printf("Error: %s\n", err)
			}
		case "echo":
			if len(cmd.args) < 1 {
				continue
			} else {
				res := cs.echo(cmd.args)
				_, err := cmd.Stdout.Write([]byte(res))
				if err != nil {
					fmt.Printf("Error: %s\n", err)
				}
			}
		}
	}
}
