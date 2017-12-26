package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

func main() {
	host, _ := os.Hostname()
	prompt := fmt.Sprintf("[barryz@%s]$ ", host)
	r := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print(prompt)
		if !r.Scan() {
			break
		}
		line := r.Text()
		if len(line) == 0 {
			continue
		}
		cmds := strings.Split(line, "|")
		switch len(cmds) {
		case 1:
			args := strings.Fields(line)
			cmd := exec.Command(args[0], args[1:]...)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err := cmd.Run()
			if err != nil {
				fmt.Println(err)
			}
		case 2:
			args1, args2 := strings.Fields(cmds[0]), strings.Fields(cmds[1])
			cmd1, cmd2 := exec.Command(args1[0], args1[1:]...), exec.Command(args2[0], args2[1:]...)
			r, w := io.Pipe()
			cmd1.Stdout = w
			cmd2.Stdin = r
			cmd2.Stdout = os.Stdout

			cmd1.Start()
			cmd2.Start()

			go func() {
				defer w.Close()

				cmd1.Wait()
			}()
			cmd2.Wait()
		}
	}
}
