package events

type EventHandler struct {
	Name    string
	Once    bool
	Handler interface{}
}

var registry []EventHandler

func register(event EventHandler) {
	registry = append(registry, event)
}

func All() []EventHandler {
	return registry
}
