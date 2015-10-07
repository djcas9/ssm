package ssm

// StateList type for []string
type StateList []string

// States list of states
type States []State

// Callback type
type Callback func()

// State holds state metadata and
// callbacks
type State struct {
	Name    string
	Initial bool

	To   StateList
	From StateList

	// Callbacks
	BeforeEnter Callback
	AfterEnter  Callback
	BeforeExit  Callback
	AfterExit   Callback
}

// CanHaveEvent checks if a passed event is callable
func (s *State) CanHaveEvent(event *Event) bool {
	for _, localState := range s.To {
		if localState == event.To {
			return true
		}
	}

	for _, fromState := range event.From {
		if s.Name == fromState {
			return true
		}
	}

	return false
}

// CanChangeToState checks if a passed in state passes
// validation and returns.
func (s *State) CanChangeToState(state *State) bool {

	if len(s.To) <= 0 {
		return true
	}

	for _, ls := range s.To {
		if ls == state.Name {
			return true
		}
	}

	return false
}

// CanChangeFromState checks if the passed in state passes
// validation and returns
func (s *State) CanChangeFromState(state *State) bool {
	for _, localState := range state.From {
		if localState == s.Name {
			return true
		}
	}

	return false
}
