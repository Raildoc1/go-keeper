package states

import (
	"fmt"
	"go-keeper/internal/client/logic/fsm"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var _ fsm.State = (*ConnectingState)(nil)

type ConnectingState struct {
	dc DependenciesContainer
}

func NewConnectingState(dc DependenciesContainer) *ConnectingState {
	return &ConnectingState{
		dc: dc,
	}
}

func (s *ConnectingState) OnEnter() {}
func (s *ConnectingState) OnLeave() {}

func (s *ConnectingState) Process() (next fsm.State) {
	fmt.Println("Connecting to " + s.dc.Config.Address + "...")
	options := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err := grpc.NewClient(s.dc.Config.Address, options)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	fmt.Println("Connected successfully!")
	return NewSelectState(s.dc, conn)
}
