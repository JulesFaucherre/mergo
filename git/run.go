package git

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"gitlab.com/jfaucherre/mergo/tools"
)

func debug(args [][]string) string {
	a := []string{}

	for _, v := range args {
		a = append(a, strings.Join(v, " "))
	}
	return strings.Join(a, " | ")
}

func run(ctx context.Context, args [][]string) (string, error) {
	if tools.Verbose {
		fmt.Printf("Running command :\n%s\n", debug(args))
	}
	var r io.Reader
	r = os.Stdin
	var s string
	c := make(chan struct{})
	outs := make([]io.WriteCloser, len(args))
	cmds := make([]*exec.Cmd, len(args))

	for i, argLst := range args {
		cmds[i] = exec.CommandContext(ctx, argLst[0], argLst[1:]...)

		cmds[i].Stdin = r
		r, outs[i] = io.Pipe()
		cmds[i].Stdout = outs[i]
	}

	go (func() {
		buf := new(bytes.Buffer)
		buf.ReadFrom(r)
		s = buf.String()
		c <- struct{}{}
	})()

	for _, cmd := range cmds {
		if err := cmd.Start(); err != nil {
			return "", err
		}
	}

	for i, cmd := range cmds {
		if err := cmd.Wait(); err != nil {
			return "", err
		}

		outs[i].Close()
	}

	<-c

	if tools.Verbose {
		fmt.Printf("Returning :\n%s\n", s)
	}

	return s, nil
}
