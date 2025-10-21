package commands

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"

	"github.com/mitchellh/go-ps"
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
		fmt.Printf("DEBUG: Parsing element: [%s]\n", element)
		cmdAndArguments := strings.Fields(element)
		fmt.Printf("DEBUG: Fields result: %v\n", cmdAndArguments)
		// fmt.Println("cmdAndArgsLength:", len(cmdAndArguments))
		if len(cmdAndArguments) == 0 { // || len(cmdAndArguments[0]) == 1 {
			fmt.Println("DEBUG: Skipping empty command")
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

func (cs *commandStack) ps() ([]int, error) {
	processes, err := ps.Processes()
	if err != nil {
		return nil, err
	}
	var processPIDs []int = make([]int, len(processes))
	for i, process := range processes {
		processPIDs[i] = process.Pid()
	}
	return processPIDs, nil
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

func (cs *commandStack) Run(ctx context.Context) {

	// initialization of buffers
	for i := 0; i < len(cs.stack) - 1; i++ {
		var buffer = &bytes.Buffer{}
		cs.stack[i].Stdout = buffer
		cs.stack[i+1].Stdin = buffer
		// var reader, writer = io.Pipe()
		// cs.stack[i].Stdout = writer
		// cs.stack[i+1].Stdin = reader
	}
	// init Stdout for last command
	//var buffer bytes.Buffer
	cs.stack[len(cs.stack)-1].Stdout = os.Stdout

	// slice to store running commands
	var runningCmds []*exec.Cmd
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
			_, err = cmd.Stdout.Write([]byte(res+"\n"))
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
		case "ps":
			processesPids, err := cs.ps()
			if err != nil {
				fmt.Printf("Error: %s\n", err)
				continue
			}

			// prepare data
			var builder strings.Builder
			for _, process := range processesPids {
				strProcessPid := strconv.Itoa(process)
				builder.WriteString(strProcessPid)
				builder.WriteString("\n")
			}

			buffer := []byte(builder.String())
			_, err = cmd.Stdout.Write(buffer)
			if err != nil {
				fmt.Printf("Error: %s\n", err)
			}
		case "kill":
			if len(cmd.args) < 1 {
				fmt.Printf("Error: Not enough arguments\n")
				continue
			}
			pid := cmd.args[0]
			err := cs.kill(pid)
			if err != nil {
				fmt.Printf("Error: %s\n", err)
			}
		// handle other commands
		default:
			externalCmd := exec.CommandContext(ctx, cmd.name, cmd.args...)
			externalCmd.Stdin = cmd.Stdin
			externalCmd.Stdout = cmd.Stdout
			// err := externalCmd.Run()
			err := externalCmd.Start()
			if err != nil {
				fmt.Printf("%sError: %s\n", cmd.name, err)
				continue
			}
			
			runningCmds = append(runningCmds, externalCmd)
			
		}
	}

	// shut down all running commands
	const op = "WAITCMDS"
	for _, cmd := range runningCmds {
		err := cmd.Wait()
		if err != nil {
			// dont wtire error if context is canceled (Ctrl+C)
            if ctx.Err() == nil {
                fmt.Printf("%sError: %s\n", op, err)
            }
		}
		// // Закрыть Stdout этой команды (если это PipeWriter)
        // if w, ok := cs.stack[i].Stdout.(*io.PipeWriter); ok {
        //     w.Close()
        // }
	}

}
