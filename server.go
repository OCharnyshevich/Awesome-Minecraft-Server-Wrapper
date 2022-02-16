package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/OCharnyshevich/Awesome-Minecraft-Server-Wrapper/minecraft"
	"github.com/urfave/cli/v2"
	"os"
)

type Server struct {
	cli            *cli.App
	wrappers       map[string]minecraft.Wrapper
	attached       string
	originalStdin  *os.File
	originalStdout *os.File
	originalStderr *os.File
}

func NewServer() *Server {
	return &Server{
		cli:            newCli(),
		originalStdin:  os.Stdin,
		originalStdout: os.Stdout,
		originalStderr: os.Stderr,
	}
}

func (s Server) Run(arguments []string) error {
	return s.cli.Run(arguments)
}

func (s *Server) hookStdin(ctx context.Context) {
	if wrp, ok := s.wrappers[s.attached]; !ok {
		fmt.Println("Active server doesn't attached: ", wrp)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			fmt.Println("Stdin hook stopped")
			return
		default:
			cmd := scanner.Text()
			fmt.Println("Typed: ", cmd)
			//wrp.Console.WriteCmd(cmd)
		}
	}
}
