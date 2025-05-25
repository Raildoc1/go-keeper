package commands

import (
	"bufio"
	"context"
	"fmt"
	"io"
)

type Commands struct {
	sc  *bufio.Scanner
	out io.Writer
}

func NewCommands(in io.Reader, out io.Writer) *Commands {
	sc := bufio.NewScanner(in)
	return &Commands{
		sc:  sc,
		out: out,
	}
}

func (r *Commands) ReadNext(ctx context.Context) (cmd string, err error) {
	for !r.sc.Scan() {
		if ctx.Err() != nil {
			return "", ctx.Err()
		}
	}
	cmd = r.sc.Text()
	return cmd, nil
}

func (r *Commands) ReadWithLabel(label string, ctx context.Context) (cmd string, err error) {
	fmt.Printf("%s: ", label)
	for !r.sc.Scan() {
		if ctx.Err() != nil {
			return "", ctx.Err()
		}
	}
	cmd = r.sc.Text()
	return cmd, nil
}

func (r *Commands) WriteWithLabel(label string, cmds []string) error {
	_, err := r.out.Write([]byte(fmt.Sprintf("%s:\n", label)))
	if err != nil {
		return fmt.Errorf("error writing commands: %w", err)
	}
	for _, cmd := range cmds {
		_, err = r.out.Write([]byte(fmt.Sprintf("--- %s\n", cmd)))
		if err != nil {
			return fmt.Errorf("error writing command: %w", err)
		}
	}
	return nil
}
