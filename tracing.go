package main

import (
	"fmt"
	"os"
	"sync"
)

type TraceEvent struct {
	Node  int
	State string
	Msg   string
}

var traceCh = make(chan TraceEvent, 10000)
var traceOnce sync.Once

func startTrace() {
	traceOnce.Do(func() {
		f, err := os.Create("trace.log")
		if err != nil {
			panic(err)
		}
		go func() {
			for ev := range traceCh {
				fmt.Fprintf(f, "%d %s %s\n", ev.Node, ev.State, ev.Msg)
			}
			f.Close()
		}()
	})
}

func trace(node *Node, msg string) {
	startTrace()
	select {
	case traceCh <- TraceEvent{Node: node.name, State: node.state, Msg: msg}:
	default:
		// drop if buffer full to avoid blocking
	}
}
