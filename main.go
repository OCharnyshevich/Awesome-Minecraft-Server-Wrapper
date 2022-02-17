package main

import (
	"log"
	"os"
)

func main() {
	newApp := NewApp()

	NewApp().Shutdown()

	err := newApp.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
