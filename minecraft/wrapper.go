package minecraft

import (
	"context"
	"errors"
	"fmt"
	"github.com/OCharnyshevich/Awesome-Minecraft-Server-Wrapper/app"
	"github.com/OCharnyshevich/Awesome-Minecraft-Server-Wrapper/minecraft/events"
	"io"
	"log"
	"strings"
)

type Wrappers interface {
	AddWrapper(wrapper *Wrapper) error
	GetWrapper(name string) (*Wrapper, error)
	HookStdin(name string)
}

// Config TODO: Create interface for config
type Config struct {
	Name      string `yaml:"name"`
	Version   string `yaml:"version"`
	RAMMin    int    `yaml:"ram-min"`
	RAMMax    int    `yaml:"ram-max"`
	ServerDir string `yaml:"server-dir"`
	JavaPath  string `yaml:"java-path"`
	JarName   string `yaml:"jar-name"`
	JavaFlags string `yaml:"java-params"`
	RootDir   string
}

func (c Config) GetPath() string {
	return fmt.Sprintf("%s/%s/%s", c.RootDir, c.ServerDir, c.Name)
}

func (c Config) GetPathJar() string {
	return fmt.Sprintf("%s/%s", c.GetPath(), c.JarName)
}

type Wrapper struct {
	Config         *Config
	Console        Console
	playerList     map[string]string
	ctxCancelFunc  context.CancelFunc
	gameEventsChan chan events.GameEvent
	loadedChan     chan bool
}

func NewDefaultWrapper(c *app.Config, name string, version string, javaPath string) *Wrapper {
	wrapper := newWrapper()
	wrapper.Config = &Config{
		Name:      name,
		Version:   version,
		RAMMin:    c.RAMMin,
		RAMMax:    c.RAMMax,
		RootDir:   c.RootDir,
		ServerDir: c.ServerDir,
		JavaPath:  javaPath,
		JarName:   c.JarName,
		JavaFlags: c.JavaFlags,
	}

	wrapper.Console = newConsole(NewProcess(wrapper.Config))
	return wrapper
}

func newWrapper() *Wrapper {
	wpr := &Wrapper{
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
func (a *Wrapper) Start() error {
	_, cancel := context.WithCancel(context.Background())
	a.ctxCancelFunc = cancel
	//go a.processLogEvents(ctx)
	//defer a.Stop() //TODO: bug(me) somehow call this function in the end of the loading
	return a.Console.Start()
}

// Stop pipes a 'stop' command to the minecraft java process.
func (a *Wrapper) Stop() error {
	log.Println("Stopping")
	return a.Console.WriteCmd("stop")
}

func (a *Wrapper) processLogEvents(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("processLogEvents stopped")
			return
		default:
			line, err := a.Console.ReadLine()
			if errors.Is(err, io.EOF) {
				fmt.Println("Process stopped")
				a.Kill()
				return
			}

			fmt.Printf("event: %s\n", strings.TrimSpace(line))

			// ev, t := w.parseLineToEvent(line)
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
func (a *Wrapper) Kill() error {
	if err := a.Console.Kill(); err != nil {
		return err
	}

	// Hard reset the wrapper machine state the 'offline'.
	//w.machine.SetState(WrapperOffline)
	// Manually trigger the context cancellation since 'SetState'
	// does not trigger any callbacks on the fsm.
	a.ctxCancelFunc()
	close(a.gameEventsChan)
	return nil
}

func (a *Wrapper) GameEvents() <-chan events.GameEvent {
	return a.gameEventsChan
}

func (a *Wrapper) Loaded() <-chan bool {
	return a.loadedChan
}
