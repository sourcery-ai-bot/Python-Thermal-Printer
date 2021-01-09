package main

import (
	"fmt"
	"time"
)

type Event struct {
	Kind   string
	AtTime time.Time
}

func ButtonPressed(at time.Time) *Event {
	return &Event{Kind: "ButtonPressed", AtTime: at}
}

func ButtonReleased(at time.Time) *Event {
	return &Event{Kind: "ButtonReleased", AtTime: at}
}

func (e Event) String() string {
	return fmt.Sprintf("Event(%s@%s)", e.Kind, e.AtTime)
}

func (e *Event) DurationSince(event *Event) time.Duration {
	return e.AtTime.Sub(event.AtTime)
}
