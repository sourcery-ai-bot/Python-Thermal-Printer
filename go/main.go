package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	//
	// First, initialize a printer with options, we might
	// want to make it so that these can be passed in at
	// some point as CLI args, or passed in via some config
	// file.
	//
	// Later on, we might want to support dispatch actions, such as
	// reacting to how many presses, or holds sequences in a row, and
	// which script to execute with this.
	//
	printer, err := NewPrinter(PrinterOptions{
		LED: LEDOptions{
			Pin:       "18",
			BlinkRate: 500 * time.Millisecond,
		},
		Button: ButtonOptions{
			Pin:         "23",
			IdleTimeout: 2 * time.Second,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// we have our context so we can terminate the program,
	// we use this in each go-routine as well as in the interrupt
	// signal handler.
	ctx, cancel := context.WithCancel(context.Background())
	// this sets up the actual signal handler.
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-termChan
		log.Printf("captured an interrupt, shutting down")
		cancel()
	}()

	//
	// Each of our async operations will be set up with a waitgroup.
	// and we wait for the group to finish to stop the process.
	var wg sync.WaitGroup

	//
	// *The Event Loops*
	//

	//
	// 1. The button listener, which will emit events based on the above
	// configuration.
	events := printer.Button.PressListener(ctx, &wg)
	//
	// 2. We set up a another events channel that we will copy the events
	// to from the printer so that they can also be sent to the LED, as well
	// as the LED loop, which blinks the light and can turn it on/off based on
	// if the button was pressed.
	ledEvents := printer.LED.ToggleLoop(ctx, &wg)
	//
	// 3. This is a WIP event handler that reads events from the button,
	// logs them and then forwards them to the LED so that it can turn on/off
	// based on if the button was pressed.
	//
	// Eventually, this will probably dispatch the print handlers and/or figure
	// out if something was a "hold" or a sequence of presses.
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

	//
	// 4. Wait to finish everything.
	wg.Wait()
}
