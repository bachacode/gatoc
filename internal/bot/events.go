package bot

type Event struct {
	Name    string
	Once    bool
	Handler func(cfg *BotContext) interface{}
}

var eventRegistry []Event

func RegisterEvent(event Event) {
	eventRegistry = append(eventRegistry, event)
}
