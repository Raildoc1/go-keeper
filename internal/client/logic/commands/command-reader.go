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

func (r *Commands) ReadWithLabel(label string, ctx context.Context) (cmd string, err error) {
	fmt.Printf("%s: ", label)

	in := make(chan string)
	errc := make(chan error)

	defer close(in)
	defer close(errc)

	go func() {
		var s string
		_, err := fmt.Scan(&s)
		if err != nil {
			errc <- err
		}
		in <- s
	}()

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case err := <-errc:
		return "", err
	case cmd = <-in:
		return cmd, nil
	}
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
