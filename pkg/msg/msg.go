package msg

import (
	fluxevent "github.com/fluxcd/flux/pkg/event"
)

// Represents a Flux event that will get sent to an exporter
type Message struct {
	TitleLink string
	Body      string
	Type      string
	Title     string
	Event     fluxevent.Event
}
