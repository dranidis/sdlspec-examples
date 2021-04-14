package main

import (
	"time"

	"github.com/dranidis/sdlspec"
)

// Signals
type DEPOSIT struct {
	amount int
}
type WITHDRAW struct {
	amount int
}
type OPEN struct{}
type CLOSE struct{}

var out chan sdlspec.Signal

func Account(p *sdlspec.Process) {
	// states
	var initial func()
	var open func()
	var active func()
	var closed func()

	// memory
	balance := 0

	// states, incoming signals and responses
	initial = sdlspec.State(p, "initial", func(s sdlspec.Signal) {
		switch s.(type) {
		case OPEN:
			// next state
			defer open()
		default:
			sdlspec.DefaultMessage(p, s)
		}
	})
	open = sdlspec.State(p, "open", func(s sdlspec.Signal) {
		switch signal := s.(type) {
		case DEPOSIT:
			balance += signal.amount
			defer active()
		case CLOSE:
			defer closed()
		default:
			sdlspec.DefaultMessage(p, s)
		}
	})
	active = sdlspec.State(p, "active", func(s sdlspec.Signal) {
		switch signal := s.(type) {
		case DEPOSIT:
			balance += signal.amount
			defer active()
		case WITHDRAW:
			if balance > signal.amount {
				balance -= signal.amount
				defer active()
			} else if balance == signal.amount {
				balance -= signal.amount
				defer open()
			}
		default:
			sdlspec.DefaultMessage(p, s)
		}
	})
	closed = sdlspec.State(p, "closed", func(s sdlspec.Signal) {
		switch s.(type) {
		default:
			sdlspec.DefaultMessage(p, s)
		}
	})
	// Initial state
	go initial()
}

func main() {
	//sdlspec.DisableLogging()
	die := make(chan sdlspec.Signal)

	out = sdlspec.MakeBuffer()
	accountChan := sdlspec.MakeProcess(Account, "Account", die)

	go sdlspec.ChannelConsumer(die, "ENV", out)

	sdlspec.Execute(
		sdlspec.Transmission{MsDelay: 10, Receiver: accountChan, Signal: OPEN{}},
		sdlspec.Transmission{MsDelay: 10, Receiver: accountChan, Signal: DEPOSIT{10}},
		sdlspec.Transmission{MsDelay: 10, Receiver: accountChan, Signal: DEPOSIT{20}},
		sdlspec.Transmission{MsDelay: 10, Receiver: accountChan, Signal: WITHDRAW{10}},
		sdlspec.Transmission{MsDelay: 10, Receiver: accountChan, Signal: WITHDRAW{20}},
		sdlspec.Transmission{MsDelay: 10, Receiver: accountChan, Signal: CLOSE{}},
	)

	time.Sleep(2000 * time.Millisecond)
	close(die)
}
