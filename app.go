package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/OCharnyshevich/Awesome-Minecraft-Server-Wrapper/app"
	"github.com/OCharnyshevich/Awesome-Minecraft-Server-Wrapper/cmd"
	"github.com/OCharnyshevich/Awesome-Minecraft-Server-Wrapper/minecraft"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	config          *app.Config
	cli             *cli.App
	versionManifest *minecraft.VersionManifest
	wrappers        map[string]*minecraft.Wrapper
	attached        string
	originalStdin   *os.File
	originalStdout  *os.File
	originalStderr  *os.File
	context         context.Context
	cancelFunc      context.CancelFunc
}

func NewApp() *App {
	ctx, cancel := context.WithCancel(context.Background())

	a := &App{
		config:          app.NewConfig(),
		cli:             cmd.NewCli(),
		versionManifest: &minecraft.VersionManifest{},
		wrappers:        map[string]*minecraft.Wrapper{},
		originalStdin:   os.Stdin,
		originalStdout:  os.Stdout,
		originalStderr:  os.Stderr,
		context:         ctx,
		cancelFunc:      cancel,
	}

	a.cli.Before = func(c *cli.Context) error {
		minecraft.PreloadManifest(a.versionManifest)
		err := a.readWrappersConfig()
		return err
	}

	a.AddCommand(cmd.AvailableVersions(a.versionManifest))
	a.AddCommand(cmd.Server(a.config, a.versionManifest, a))

	return a
}

func (a *App) Run(arguments []string) error {
	return a.cli.Run(arguments)
}

func (a *App) Shutdown() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigs
		for _, wrapper := range a.wrappers {
			err := wrapper.Stop()
			if err != nil {
				log.Fatal(err)
			}
		}

		log.Printf("\nEverything is shuted down")

		time.Sleep(10 * time.Second)
		os.Exit(2)
	}()
}

func (a *App) AddCommand(commands ...*cli.Command) {
	a.cli.Commands = append(a.cli.Commands, commands...)
}

func (a App) AddWrapper(wrapper *minecraft.Wrapper) error {
	a.wrappers[wrapper.Config.Name] = wrapper
	return a.saveWrappersConfig()
}

func (a App) saveWrappersConfig() error {
	file, _ := json.MarshalIndent(a.wrappers, "", " ")

	return ioutil.WriteFile(fmt.Sprintf("%s/.wrappers.json", a.config.GetPath()), file, 0700)
}

func (a App) GetWrapper(name string) (*minecraft.Wrapper, error) {
	wrp, ok := a.wrappers[name]

	if !ok {
		return nil, fmt.Errorf("server '%s' not found", name)
	}

	return wrp, nil
}

func (a *App) HookStdin(name string) {
	var (
		wrp *minecraft.Wrapper
		ok  bool
	)

	if wrp, ok = a.wrappers[name]; !ok {
		log.Println("Active server doesn't attached: ", wrp)
	}

	a.attached = name

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		select {
		case <-a.context.Done():
			log.Println("Stdin hook stopped")
			return
		default:
			text := scanner.Text()
			log.Println("Typed: ", text)
			stat, err := wrp.Console.Stat()
			fmt.Printf("PID: %d MEM: %d M %v", stat.PID, stat.Memory, err)
			if err := wrp.Console.WriteCmd(text); err != nil {
				log.Fatal(err)
			}
		}
	}
}

func (a App) readWrappersConfig() error {
	file, err := ioutil.ReadFile(fmt.Sprintf("%s/.wrappers.json", a.config.GetPath()))

	if err != nil {
		if !os.IsExist(err) {
			return a.saveWrappersConfig()
		}

		return err
	}

	data := a.wrappers

	err = json.Unmarshal(file, &data)

	for name, wrapper := range data {
		wrp := minecraft.NewDefaultWrapper(a.config, wrapper.Config.Name, wrapper.Config.Version, wrapper.Config.JavaPath)
		wrp.Config = wrapper.Config
		wrp.Config.RootDir = a.config.RootDir
		a.wrappers[name] = wrp
	}

	return err
}
