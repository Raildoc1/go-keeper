package states

import (
	"bufio"
	"fmt"
	"go-keeper/internal/client/logic/fsm"
	"google.golang.org/grpc"
	"os"
)

var _ fsm.State = (*SelectState)(nil)

const (
	quitOperationName = "quit"
)

type SelectState struct {
	dc   DependenciesContainer
	conn *grpc.ClientConn
}

func NewSelectState(dc DependenciesContainer, conn *grpc.ClientConn) *SelectState {
	return &SelectState{
		dc:   dc,
		conn: conn,
	}
}

func (s *SelectState) OnEnter() {}
func (s *SelectState) OnLeave() {}

func (s *SelectState) Process() (next fsm.State) {
	str, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return nil
	}
	cmd := str[:len(str)-1]
	switch cmd {
	case quitOperationName:
		fmt.Println("quitting...")
		return nil
	default:
		fmt.Printf("no such command '%s'\n", str)
		return NewSelectState(s.dc, s.conn)
	}
}
