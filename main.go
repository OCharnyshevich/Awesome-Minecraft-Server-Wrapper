package main

import (
	"AMSW/minecraft"
	"fmt"
	"log"
	"os"
)

func main() {
	//ctx, _ := context.WithCancel(context.Background())
	server := NewServer()
	//server.hookStdin(ctx)

	err := server.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func run() {
	wpr := minecraft.NewDefaultWrapper()
	defer wpr.Stop()

	err := wpr.Start()
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		select {
		case ev, ok := <-wpr.GameEvents():
			if !ok {
				log.Println("Game events channel closed", ev.String())
				return
			}

			log.Println("events", ev.String())
		}
	}
}
