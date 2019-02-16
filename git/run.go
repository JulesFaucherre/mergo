package git

import (
	"bytes"
	"context"
	"io"
	"os"
	"os/exec"
)

func run(ctx context.Context, args [][]string) (string, error) {
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

	return s, nil
}
