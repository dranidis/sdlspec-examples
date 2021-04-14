package main

import (
	"time"

	"github.com/dranidis/sdlspec"
)

type UP struct {
	n int
}
type DN struct{}
type OVER struct{}

var out chan sdlspec.Signal

func Counter(p *sdlspec.Process) {
	var goingDn func() // for mutual definition

	counter := 0
	goingUp := sdlspec.State(p, "goingUp", func(s sdlspec.Signal) {
		switch v := s.(type) {
		case UP:
			counter += v.n
			if counter > 4 {
				out <- OVER{}
				defer goingDn()
				return
			}
		default:
			sdlspec.DefaultMessage(p, s)
		}
	})
	goingDn = sdlspec.State(p, "goingDn", func(s sdlspec.Signal) {
		switch s.(type) {
		case DN:
			counter -= 1
			if counter == 0 {
				defer goingUp()
				return
			}
		default:
			sdlspec.DefaultMessage(p, s)
		}
	})

	go goingUp()
}

func main() {
	//sdlspec.DisableLogging()
	die := make(chan sdlspec.Signal)

	out = sdlspec.MakeBuffer()
	counterChan := sdlspec.MakeProcess(Counter, "Counter", die)

	go sdlspec.ChannelConsumer(die, "ENV", out)

	sdlspec.Execute(
		sdlspec.Transmission{MsDelay: 10, Receiver: counterChan, Signal: UP{}},
		sdlspec.Transmission{MsDelay: 10, Receiver: counterChan, Signal: DN{}},
		sdlspec.Transmission{MsDelay: 10, Receiver: counterChan, Signal: UP{4}},
		sdlspec.Transmission{MsDelay: 10, Receiver: counterChan, Signal: DN{}},
		sdlspec.Transmission{MsDelay: 10, Receiver: counterChan, Signal: DN{}},
		sdlspec.Transmission{MsDelay: 10, Receiver: counterChan, Signal: DN{}},
		sdlspec.Transmission{MsDelay: 10, Receiver: counterChan, Signal: DN{}},
		sdlspec.Transmission{MsDelay: 10, Receiver: counterChan, Signal: DN{}},
		sdlspec.Transmission{MsDelay: 10, Receiver: counterChan, Signal: UP{}},
	)

	time.Sleep(2000 * time.Millisecond)
	close(die)
}
