package commandreader

import (
	"bufio"
	"context"
	"io"
)

type CommandReader struct {
	sc *bufio.Scanner
}

func NewCommandReader(in io.Reader) *CommandReader {
	sc := bufio.NewScanner(in)
	return &CommandReader{
		sc: sc,
	}
}

func (r *CommandReader) ReadNext(ctx context.Context) (cmd string, err error) {
	for !r.sc.Scan() {
		if ctx.Err() != nil {
			return "", ctx.Err()
		}
	}
	cmd = r.sc.Text()
	return cmd, nil
}
