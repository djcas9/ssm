package ssm

import (
	"os"
	"testing"
)

type Person struct {
	Name  string
	Age   int
	Email string
	Drunk bool
	State *SSM
}

var (
	p *Person
)

func TestMain(m *testing.M) {

	stateMachine := NewStateMachine(States{{
		Name:    "normal",
		Initial: true,
		BeforeEnter: func() {
		},
		BeforeExit: func() {
		},
	}, {
		Name: "mad",
		From: StateList{"happy", "sad", "normal"},
		To:   StateList{"happy"},
		AfterEnter: func() {
		},
	}, {
		Name: "happy",
		From: StateList{"mad", "sad", "normal"},
		To:   StateList{"mad", "sad", "normal"},
		BeforeEnter: func() {
		},
		AfterExit: func() {
		},
	}, {
		Name: "sad",
		From: StateList{"mad", "happy"},
		To:   StateList{"happy", "normal", "mad"},
	}})

	p = &Person{
		Name:  "Alice",
		Age:   29,
		Drunk: false,
		Email: "b@c.com",
	}

	p.State = stateMachine

	code := m.Run()
	os.Exit(code)
}

func TestCurrentState(t *testing.T) {
	if p.State.CurrentState().Name != "normal" {
		t.Error("Initial state should be: normal")
	}
}

func TestChangeState(t *testing.T) {
	if err := p.State.Change("mad"); err != nil {
		t.Error(err)
	}
}

func TestChangeStateFail(t *testing.T) {
	if err := p.State.Change("sad"); err == nil {
		t.Error("Should not be able to change state from nnormal to sad")
	}
}

func TestNewEvent(t *testing.T) {

	if err := p.State.Change("happy"); err != nil {
		t.Error(err)
	}

	if err := p.State.NewEvent(Event{
		Name: "bad_day",
		From: StateList{"happy", "normal", "sad", "mad"},
		To:   "sad",
		Before: func() {
			p.Drunk = true
		},
		After: func() {
		},
	}); err != nil {
		t.Fatal(err)
	}

	if err := p.State.Event("bad_day"); err != nil {
		t.Error(err)
	}

	if !p.Drunk {
		t.Error("callback fail")
	}
}
