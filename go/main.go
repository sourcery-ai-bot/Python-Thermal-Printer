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
	printer.LED.ToggleLoop(ctx, &wg)
	events := printer.Button.PressListener(ctx, &wg)

	// this is the event processor
	go func() {
		var lastEvent *Event
		for e := range events {
			if lastEvent == nil {
				log.Printf("first! %s", e)
				lastEvent = e
				continue
			}

			log.Printf("Got %s -- %s since last", e, e.DurationSince(lastEvent))
			lastEvent = e
		}
	}()

	// wait for everything to finish
	wg.Wait()
}
