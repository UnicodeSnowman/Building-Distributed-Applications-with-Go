package coordinator

import "time"

type EventRaiser interface {
	AddListener(eventName string, f func(interface{}))
}

type EventData struct {
	Name      string
	Value     float64
	Timestamp time.Time
}

type EventAggregator struct {
	listeners map[string][]func(interface{})
}

func NewEventAggregator() *EventAggregator {
	ea := EventAggregator{
		listeners: make(map[string][]func(interface{})),
	}

	return &ea
}

func (ea *EventAggregator) AddListener(name string, f func(interface{})) {
	ea.listeners[name] = append(ea.listeners[name], f)
}

func (ea *EventAggregator) PublishEvent(name string, eventData interface{}) {
	listeners := ea.listeners[name]
	if listeners != nil {
		for _, listener := range listeners {
			listener(eventData)
		}
	}
}
