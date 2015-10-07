package ssm

import "fmt"

// SSM simple state machine context
type SSM struct {
	States States
	Events Events
	State  *State
}

// NewStateMachine creates a new SSM context
// and starts the initial state change is specified.
func NewStateMachine(states States) *SSM {

	ssm := &SSM{
		States: states,
	}

	for i, state := range ssm.States {
		s := ssm.States[i]

		if s.Initial {

			if s.BeforeEnter != nil {
				state.BeforeEnter()
			}

			ssm.State = &s

			if s.AfterEnter != nil {
				s.AfterEnter()
			}
		}
	}

	return ssm
}

// Event allows you to trigger an event and start a state
// transition.
func (s *SSM) Event(name string) error {

	if e := s.GetEventByName(name); e != nil {

		if !s.State.CanHaveEvent(e) {
			return fmt.Errorf("Unable to call event '%s'", name)
		}

		if e.Before != nil {
			e.Before()
		}

		if err := s.Change(e.To); err != nil {
			return err
		}

		if e.After != nil {
			e.After()
		}

		return nil
	}

	return fmt.Errorf("unknown event: %s", name)
}

// GetEventByName allows you to fetch an event struct by event name.
func (s *SSM) GetEventByName(name string) *Event {
	for i, event := range s.Events {
		if event.Name == name {
			return &s.Events[i]
		}
	}

	return nil
}

// NewEvent will create a new event and add it to
// the SSM context.
func (s *SSM) NewEvent(event Event) error {
	if e := s.GetEventByName(event.Name); e == nil {
		s.Events = append(s.Events, event)

		return nil
	}

	return fmt.Errorf("An event with this name (%s) already exists.", event.Name)
}

// GetStateByName allows you to fetch a state by name
func (s *SSM) GetStateByName(name string) *State {
	for i, state := range s.States {
		if state.Name == name {
			return &s.States[i]
		}
	}

	return nil
}

// Change is used to change state. This method will
// error is validations fail.
func (s *SSM) Change(name string) error {
	if state := s.GetStateByName(name); state != nil {
		oldState := s.State

		if !s.State.CanChangeToState(state) {
			return fmt.Errorf("Unable to change state to '%s' - Allowed states: %s", name, oldState.To)
		}

		if !s.State.CanChangeFromState(state) {
			return fmt.Errorf("Unable to change state to '%s' from '%s' - Allowed from states: %s", state.Name, s.State.Name, state.From)
		}

		if oldState.BeforeExit != nil {
			oldState.BeforeExit()
		}

		if state.BeforeEnter != nil {
			state.BeforeEnter()
		}

		s.State = state

		if oldState.AfterExit != nil {
			oldState.AfterExit()
		}

		if s.State.AfterEnter != nil {
			s.State.AfterEnter()
		}

		return nil
	}

	return fmt.Errorf("Unable to change state. State '%s' does not exist.", name)
}

// CurrentState held by the SSM context
func (s *SSM) CurrentState() *State {
	return s.State
}

// Is allows you to check if the current state is X
func (s *SSM) Is(state string) bool {
	return s.State.Name == state
}

// IsNot allows you to check if the current state is not X
func (s *SSM) IsNot(state string) bool {
	return s.State.Name != state
}
