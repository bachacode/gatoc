package events

import "github.com/bachacode/go-discord-bot/internal/config"

type Event struct {
	Name    string
	Once    bool
	Handler func(cfg *config.BotConfig) interface{}
}

var registry []Event

func register(event Event) {
	registry = append(registry, event)
}

func All() []Event {
	return registry
}
