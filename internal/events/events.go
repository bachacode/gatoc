package events

type Event struct {
	Name    string
	Once    bool
	Handler interface{}
}

var registry []Event

func register(event Event) {
	registry = append(registry, event)
}

func All() []Event {
	return registry
}
