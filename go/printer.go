package main

import (
	"log"
	"os"
	"os/exec"
	"time"

	"periph.io/x/host/v3"
)

type PrinterOptions struct {
	LED    LEDOptions
	Button ButtonOptions
}

type Printer struct {
	LED    *LEDPin
	Button *ButtonPin

	isPrinting bool
}

func NewPrinter(opts PrinterOptions) (*Printer, error) {
	_, err := host.Init()
	if err != nil {
		return nil, err
	}

	led, err := NewLEDPin(opts.LED)
	if err != nil {
		return nil, err
	}

	button, err := NewButtonPin(opts.Button)
	if err != nil {
		return nil, err
	}

	return &Printer{LED: led, Button: button}, nil
}

func (p *Printer) MaybePrint(e1 *Event, e2 *Event) {
	if p.isPrinting {
		return
	}

	if e1.Is(ButtonPressed) && e2.Is(ButtonReleased) && e2.DurationSince(e1) > 1*time.Second {
		p.isPrinting = true
		go func() {
			var err error

			defer func() {
				p.isPrinting = false
				if err != nil {
					log.Printf("Print failed %v", err)
				} else {
					log.Printf("print succeeded")
				}
			}()

			time.Sleep(5 * time.Second)
			/*
				err = execPython("niceties.py")
			*/
		}()
	}
}

func execPython(name string) error {
	cmd := exec.Command("python", name)
	cmd.Dir = "/home/pi/Python-Thermal-Printer"
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
