# SSM - Simple State Machine

I wanted something a little more similar to the ruby on rails plugin 
`aasm - https://github.com/aasm/aasm`. This lib is still young but 
I plan to continue my focus on validations and state enforcment.

# Installation

```
go get github.com/mephux/ssm
```

# Usage

* NewStateMachine and creating your base states

Example:

```go
person := Person{
  Name:  "Mephux",
  Age:   29,
  Drunk: true,
  Email: "cool@cool.com",
}

person.State = ssm.NewStateMachine(ssm.States{{
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
```

* Changing states and checking state

```go
s := person.State.CurrentState()

fmt.Println("Current State:", s)

if err := person.State.Change("mad"); err != nil {
  log.Fatal(err)
}

if person.State.IsNot("mad") {
  person.State.Change("happy")
} else {
  fmt.Println("State is mad")
}

if err := person.State.Change("happy"); err != nil {
  log.Fatal(err)
}
```

* Create and use Events

```go
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

if err := person.State.Event("bad_day"); err != nil {
  log.Fatal(err)
}
```

## Self-Promotion

Like ssm? Follow the repository on
[GitHub](https://github.com/mephux/ssm) and if
you would like to stalk me, follow [mephux](http://dweb.io/) on
[Twitter](http://twitter.com/mephux) and
[GitHub](https://github.com/mephux).
