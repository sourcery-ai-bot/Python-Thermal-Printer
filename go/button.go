package main

import (
	"context"
	"sync"
	"time"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
)

type ButtonPin struct {
	pin gpio.PinIO
}

func NewButtonPin() *ButtonPin {
	pin := gpioreg.ByName(ButtonPinName)
	pin.In(gpio.PullUp, gpio.NoEdge)
	return &ButtonPin{
		pin: pin,
	}
}

type buttonState struct {
	l      gpio.Level
	t      time.Time
	update bool
}

func (s *buttonState) next(l gpio.Level, t time.Time) *Event {
	// If the level changed, we tell ourselves update the state.
	//
	// We don't immediately emit a change event because the level
	// toggles a bit between high/low.
	//
	// This will be called in a loop in the `PressListener` so this
	// early return is fine.
	if s.l != l {
		s.l = l
		s.t = t
		s.update = true
		return nil
	}

	// If the state has stabilized we can emit a change event, and make sure
	// to emit _only_ one.
	if s.update && t.Sub(s.t) > 5*time.Millisecond {
		s.update = false
		switch s.l {
		case gpio.Low:
			return ButtonPressedEvent(s.t)
		case gpio.High:
			return ButtonReleasedEvent(s.t)
		}
	}

	if s.l == gpio.High && t.Sub(s.t) > 2*time.Second {
		s.t = t
		return IdleEvent(s.t)
	}

	return nil
}

// PressListener watches for button pressers and produces a channel of Events.
func (p *ButtonPin) PressListener(ctx context.Context, wg *sync.WaitGroup) <-chan *Event {
	wg.Add(1)
	events := make(chan *Event)

	go func() {
		defer wg.Done()

		// initial state
		s := &buttonState{l: p.pin.Read(), t: time.Now()}

		// the loop
		for {
			select {
			case <-ctx.Done():
				close(events)
				return
			default:
				if event := s.next(p.pin.Read(), time.Now()); event != nil {
					events <- event
				}
			}
		}
	}()

	return events
}
