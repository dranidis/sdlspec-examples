package main

import (
	"fmt"
	"time"

	"github.com/dranidis/sdlspec"
)

type HI struct{}

func helloStates(p *sdlspec.Process) {
	start := sdlspec.State(p, "start", func(s sdlspec.Signal) {
		switch s.(type) {
		case HI:
			fmt.Println("Hello SDL")
		default:
		}
	})
	go start()
}

func main() {
	die := make(chan sdlspec.Signal)
	helloProcess := sdlspec.MakeProcess(helloStates, "hello", die)
	helloProcess <- HI{}

	time.Sleep(1000 * time.Millisecond)
	close(die)
}
