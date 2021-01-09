package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"periph.io/x/host/v3"
)

const (
	LEDPinName    = "18"
	ButtonPinName = "23"
)

func main() {
	var err error

	_, err = host.Init()
	if err != nil {
		log.Fatal(err)
	}

	printer := NewPrinter()

	// we have our context so we can terminate the program
	ctx, cancel := context.WithCancel(context.Background())

	// Handle sigterm and await termChan signal
	termChan := make(chan os.Signal)
	signal.Notify(termChan, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-termChan
		log.Printf("captured an interrupt, shutting down")
		cancel()
	}()

	// set up our actions/listeners given the context
	var wg sync.WaitGroup

	// this looks for button presses and releases, and will
	// emit an Idle event if there was no press for 2 seconds.
	events := printer.Button.PressListener(ctx, &wg)

	ledEvents := make(chan *Event)
	printer.LED.ToggleLoop(ctx, &wg, ledEvents)

	// this is the event processor
	go func() {
		var lastEvent *Event
		for e := range events {
			if lastEvent == nil {
				log.Printf("%s", e)
			} else {
				log.Printf("%s | %s", e, e.DurationSince(lastEvent))
			}
			lastEvent = e
			ledEvents <- e
		}
	}()

	// wait for everything to finish
	wg.Wait()
}
