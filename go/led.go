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
	p.level = !p.level
	return p.pin.Out(p.level)
}

func (p *LEDPin) Off() error {
	return p.pin.Out(gpio.Low)
}

func (p *LEDPin) On() error {
	return p.pin.Out(gpio.High)
}

func (p *LEDPin) ToggleLoop(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		t := time.NewTicker(500 * time.Millisecond)
		for {
			select {
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
