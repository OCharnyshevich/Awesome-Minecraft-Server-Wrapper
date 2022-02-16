package minecraft

import (
	"context"
	"fmt"
	"github.com/OCharnyshevich/Awesome-Minecraft-Server-Wrapper/minecraft/events"
	"io"
	"time"
)

type Wrapper struct {
	Version        string
	Console        Console
	playerList     map[string]string
	ctxCancelFunc  context.CancelFunc
	gameEventsChan chan events.GameEvent
	loadedChan     chan bool
}

func NewDefaultWrapper() *Wrapper {
	config, _ := newConfig()
	cmd := NewServer(config)
	console := newConsole(cmd)
	return NewWrapper(console)
}

func NewWrapper(c Console) *Wrapper {
	wpr := &Wrapper{
		Console:    c,
		playerList: map[string]string{},
		ctxCancelFunc: func() {
			fmt.Println("Call ctxCancelFunc")
		},
		gameEventsChan: make(chan events.GameEvent, 10),
		loadedChan:     make(chan bool, 1),
	}
	return wpr
}

// Start will initialize the minecraft java process and start
// orchestrating the wrapper machine.
func (w *Wrapper) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	w.ctxCancelFunc = cancel
	go w.processLogEvents(ctx)
	return w.Console.Start()
}

// Stop pipes a 'stop' command to the minecraft java process.
func (w *Wrapper) Stop() error {
	return w.Console.WriteCmd("stop")
}

func (w *Wrapper) processLogEvents(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("processLogEvents stopped")
			return
		case <-time.After(500 * time.Millisecond):
			fmt.Println("done")
		default:
			line, err := w.Console.ReadLine()
			if err == io.EOF {
				fmt.Println("Server stopped")
				w.Kill()
				return
			}

			fmt.Println("Event: ", line)

			//ev, t := w.parseLineToEvent(line)
			//switch t {
			//case events.TypeState:
			//	w.updateState(ev.(events.StateEvent))
			//case events.TypeCmd:
			//	w.handleCmdEvent(ev.(events.GameEvent))
			//case events.TypeGame:
			//	w.handleGameEvent(ev.(events.GameEvent))
			//default:
			//}
		}
	}
}

// Kill the java process, use with caution since it will not trigger a save game.
// Kill manually perform some cleanup task and hard reset the state to 'offline'.
func (w *Wrapper) Kill() error {
	if err := w.Console.Kill(); err != nil {
		return err
	}

	// Hard reset the wrapper machine state the 'offline'.
	//w.machine.SetState(WrapperOffline)
	// Manually trigger the context cancellation since 'SetState'
	// does not trigger any callbacks on the fsm.
	w.ctxCancelFunc()
	close(w.gameEventsChan)
	return nil
}

func (w *Wrapper) GameEvents() <-chan events.GameEvent {
	return w.gameEventsChan
}

func (w *Wrapper) Loaded() <-chan bool {
	return w.loadedChan
}
