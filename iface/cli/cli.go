package cli

import (
	"bufio"
	"fmt"
	"os"
	"simon/iface/cli/cmd"
	"strings"
)

type CLI int

func (CLI) Start() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(`> `) // request a command
		c, _ := reader.ReadString('\n')
		handleCmd(c)
	}
}

func handleCmd(c string) {
	c = strings.TrimSpace(c)
	args := strings.Split(c, " ")
	if len(args) < 1 {
		return
	}
	if runner, ok := cmd.Cmds[args[0]]; ok {
		runner.Run(args[1:])
	} else if len(args[0]) != 0 {
		fmt.Printf("invalid command: \"%v\"\n", args[0])
	}
}
