package ssm

// Events list of events
type Events []Event

// Event struct for added events and state transitions
type Event struct {
	Name string

	To   string
	From StateList

	Before Callback
	After  Callback
}
