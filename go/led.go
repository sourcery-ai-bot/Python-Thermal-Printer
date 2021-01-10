package main

import (
	"context"
	"log"
	"sync"
	"time"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
)

type LEDOptions struct {
	Pin       string
	BlinkRate time.Duration
}

type LEDPin struct {
	pin    gpio.PinIO
	level  gpio.Level
	ticker *ledTickerState
}

func NewLEDPin(opts LEDOptions) *LEDPin {
	return &LEDPin{
		pin: gpioreg.ByName(opts.Pin),
		ticker: &ledTickerState{
			blinkRate: opts.BlinkRate,
			stopped:   true,
		},
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

func (p *LEDPin) Handle(e *Event) error {
	switch e.Kind {
	case ButtonPressed:
		p.ticker.Stop()
		return p.On()
	case ButtonReleased:
		return p.Off()
	case Idle:
		p.ticker.Start()
	}
	return nil
}

func (p *LEDPin) ToggleLoop(ctx context.Context, wg *sync.WaitGroup) chan<- *Event {
	events := make(chan *Event)
	wg.Add(1)
	go func() {
		defer wg.Done()
		p.ticker.Start()
		for {
			select {
			case e := <-events:
				if err := p.Handle(e); err != nil {
					log.Printf("[Error] %s", err)
					return
				}
			case <-ctx.Done():
				if err := p.Off(); err != nil {
					log.Printf("[Error] %s", err)
				}
				return
			case <-p.ticker.t.C:
				if err := p.Toggle(); err != nil {
					log.Printf("[Error] %s", err)
					return
				}
			}
		}
	}()
	return events
}

type ledTickerState struct {
	t         *time.Ticker
	stopped   bool
	blinkRate time.Duration
}

func (t *ledTickerState) Start() {
	if t.stopped {
		t.stopped = false
		t.t = time.NewTicker(t.blinkRate)
	}
}

func (t *ledTickerState) Stop() {
	t.stopped = true
	t.t.Stop()
}
