package fsm

type State interface {
	OnEnter()
	OnLeave()
	Process() (next State)
}

type StateMachine struct {
	currentState State
}

func NewStateMachine(initialState State) *StateMachine {
	initialState.OnEnter()
	return &StateMachine{
		currentState: initialState,
	}
}

func (s *StateMachine) Process() {
	for s.currentState != nil {
		nextState := s.currentState.Process()
		s.currentState.OnLeave()
		s.currentState = nextState
		if nextState != nil {
			s.currentState.OnEnter()
		}
	}
}
