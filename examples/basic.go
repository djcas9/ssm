package main

import (
	"fmt"
	"log"

	"github.com/mephux/ssm"
)

type Person struct {
	Name  string
	Age   int
	Email string
	Drunk bool
	State *ssm.SSM
}

func main() {
	person := Person{
		Name:  "Mephux",
		Age:   29,
		Drunk: true,
		Email: "cool@cool.com",
	}

	person.State = ssm.NewStateMachine(ssm.States{
		{
			Name:    "normal",
			Initial: true,
			BeforeEnter: func() {
			},
			BeforeExit: func() {
			},
		}, {
			Name: "mad",
			From: ssm.StateList{
				"happy",
				"sad",
				"normal",
			},
			To: ssm.StateList{"happy"},
			AfterEnter: func() {
				person.Drunk = true
			},
		}, {
			Name: "happy",
			From: ssm.StateList{"mad", "sad", "normal"},
			To:   ssm.StateList{"mad", "sad", "normal"},
			BeforeEnter: func() {
				fmt.Println("THINGS ARE LOOKING UP.")
			},
			AfterExit: func() {
				fmt.Println("Feeling less awesome.")
			},
		}, {
			Name: "sad",
			From: ssm.StateList{"mad", "happy"},
			To:   ssm.StateList{"happy", "normal", "mad"},
		},
	})

	if err := person.State.NewEvent(ssm.Event{
		Name: "bad_day",
		From: ssm.StateList{"happy", "normal", "sad", "mad"},
		To:   "sad",
		Before: func() {
			fmt.Println("Err, things are not going well today.")
		},
		After: func() {
			fmt.Println("super lame.")
		},
	}); err != nil {
		log.Fatal(err)
	}

	s := person.State.CurrentState()
	fmt.Println("Current State:", s.Name)

	if err := person.State.Change("mad"); err != nil {
		log.Fatal(err)
	}

	s = person.State.CurrentState()
	fmt.Println("Current State:", s.Name)

	if person.State.IsNot("mad") {
		fmt.Println("word?")
		person.State.Change("happy")
	} else {
		fmt.Println("State is mad")
	}

	if err := person.State.Change("happy"); err != nil {
		log.Fatal(err)
	}

	fmt.Println(person.State.CurrentState().Name)

	if person.State.Is("happy") && person.Drunk {
		if err := person.State.Event("bad_day"); err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println(person.State.CurrentState().Name)

}
