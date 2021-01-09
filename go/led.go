package main

import (
	"context"
	"log"
	"sync"
	"time"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
)

type LEDPin struct {
	pin   gpio.PinIO
	level gpio.Level
}

func NewLEDPin() *LEDPin {
	return &LEDPin{
		pin: gpioreg.ByName(LEDPinName),
	}
}

func (p *LEDPin) Toggle() error {
	return p.Out(!p.level)
}

func (p *LEDPin) Off() error {
	return p.Out(gpio.Low)
}

func (p *LEDPin) On() error {
	return p.Out(gpio.High)
}

func (p *LEDPin) Out(l gpio.Level) error {
	p.level = l
	return p.pin.Out(l)
}

func (p *LEDPin) ToggleLoop(ctx context.Context, wg *sync.WaitGroup, events <-chan *Event) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		stopped := false
		t := ticker()
		for {
			select {
			case e := <-events:
				// we send a copy of the events here,
				// and this only handles the press, release, and idle.
				switch e.Kind {
				case ButtonPressed:
					t.Stop()
					stopped = true
					_ = p.On()
				case ButtonReleased:
					_ = p.Off()
				case Idle:
					if stopped {
						t = ticker()
						stopped = false
					}
				}
			case <-ctx.Done():
				if err := p.Off(); err != nil {
					log.Printf("[Error] %s", err)
				}
				return
			case <-t.C:
				if err := p.Toggle(); err != nil {
					log.Printf("[Error] %s", err)
					return
				}
			}
		}
	}()
}

func ticker() *time.Ticker {
	return time.NewTicker(500 * time.Millisecond)
}
