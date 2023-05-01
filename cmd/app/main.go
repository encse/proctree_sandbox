package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/encse/proctree_sandbox/internal/hosts"
	"github.com/encse/proctree_sandbox/internal/pids"
)

func main() {
	fmt.Println("Welcome to Gorluin! 👋")
	fmt.Println("")
	fmt.Println("This program simulates a process tree running in some virtual environment.")
	fmt.Println("- Most lines you enter create a new 'shell' using the argument as the prompt.")
	fmt.Println("- Entering 'x' exits the given shell and returns to the parent.")
	fmt.Println("- You can list processes with 'ps' and kill processes with 'kill <pid>'.")
	fmt.Println("  Killing a process will kill its whole process tree with all of its children.")
	fmt.Println("  E.g. 'kill 1' exits the program immediately.")
	fmt.Println("")

	host := hosts.New("Gorluin")
	conn := hosts.NewConn(os.Stdin, os.Stdout)
	hosts.Exec(conn, host, Shell{Prompt: "home"})

	fmt.Println("Have a nice day.")
}

// commands
type Shell struct {
	Prompt string
}

func (s Shell) Name() string {
	return "shell"
}

func (s Shell) Run(env hosts.Env) error {

	reader := bufio.NewReader(env.Stdin)
	for {
		fmt.Fprintf(env.Stdout, "%s (pid: %d)> ", s.Prompt, env.Pid)
		cmd, err := reader.ReadString('\n')

		if err != nil {
			return fmt.Errorf("could not read input %w", err)
		}

		fields := strings.Fields(cmd)
		if len(fields) > 0 {
			cmd = fields[0]
			args := fields[1:]
			if cmd == "x" {
				return nil
			} else if cmd == "ps" {
				env.Exec(env.Host, Ps{})
			} else if cmd == "kill" {
				env.Exec(env.Host, Kill{Args: args})
			} else {
				env.Exec(env.Host, Shell{Prompt: cmd})
			}
		}
	}
}

type Ps struct {
	Prompt string
}

func (cmd Ps) Name() string {
	return "ps"
}

func (cmd Ps) Run(env hosts.Env) error {
	fmt.Fprintf(env.Stdout, "Processes on %s:\n", env.Host.Name)
	for _, vproc := range env.Host.GetVprocs() {
		fmt.Fprintln(env.Stdout,
			vproc.Pid, "\t",
			vproc.Executable.Name(), "\t",
			vproc.Started.Format(time.TimeOnly),
		)
	}
	return nil
}

type Kill struct {
	Args []string
}

func (cmd Kill) Name() string {
	return "kill"
}

func (cmd Kill) Run(env hosts.Env) error {
	pid, _ := strconv.Atoi(cmd.Args[0])
	for _, vproc := range env.Host.GetVprocs() {
		if vproc.Pid == pids.Pid(pid) {
			env.Host.Terminate(vproc)
		}
	}
	return nil
}
