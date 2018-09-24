package ticker

import "time"

type Signal int

const (
	On Signal = iota
	Off
	Kill
)

type Ticker struct {
	control <-chan Signal
	on      bool
	t       *time.Ticker
}

func New(d time.Duration) (<-chan time.Time, chan<- Signal) {
	control := make(chan Signal)
	signal := make(chan time.Time)
	ticker := time.NewTicker(d)
	t := &Ticker{control, false, ticker}
	go t.loop(signal)
	return signal, control
}

func (t *Ticker) loop(signal chan<- time.Time) {
	for {
		select {
		case tick := <- t.t.C:
			if t.on {
				signal <- tick
			}
		case s := <- t.control:
			switch s {
			case Kill:
				t.t.Stop()
				return
			default:
				t.on = s != Off
			}
		}
	}
}