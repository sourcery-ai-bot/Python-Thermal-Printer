package main

type Printer struct {
	LED    *LEDPin
	Button *ButtonPin
}

func NewPrinter() *Printer {
	return &Printer{
		LED:    NewLEDPin(),
		Button: NewButtonPin(),
	}
}

func (p *Printer) Shutdown() error {
	var err error
	err = p.LED.Off()
	if err != nil {
		return err
	}
	return nil
}
