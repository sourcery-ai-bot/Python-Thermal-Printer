package main

import (
	"fmt"
	"time"
)

type EventKind string

const (
	ButtonPressed  EventKind = "button_pressed"
	ButtonReleased EventKind = "button_released"
	Idle           EventKind = "idle"
)

func ButtonPressedEvent(at time.Time) *Event {
	return &Event{Kind: ButtonPressed, AtTime: at}
}

func ButtonReleasedEvent(at time.Time) *Event {
	return &Event{Kind: ButtonReleased, AtTime: at}
}

func IdleEvent(at time.Time) *Event {
	return &Event{Kind: Idle, AtTime: at}
}

type Event struct {
	Kind   EventKind
	AtTime time.Time
}

func (e Event) String() string {
	return fmt.Sprintf("event[%s@%s]", e.Kind, e.AtTime)
}

func (e *Event) DurationSince(event *Event) time.Duration {
	return e.AtTime.Sub(event.AtTime)
}

func (e *Event) Is(kind EventKind) bool {
	return e.Kind == kind
}
