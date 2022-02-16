package minecraft

import (
	"fmt"
	"io"
	"os/exec"
)

type Server interface {
	Stdout() io.ReadCloser
	Stdin() io.WriteCloser
	Start() error
	Kill() error
}

type ServerExec struct {
	cmd *exec.Cmd
}

func NewServer(config *Config) *ServerExec {
	serverPath := fmt.Sprintf("%s/%s", config.RootDir, config.ServerDir)

	cmd := exec.Command("java", fmt.Sprintf("-Xmx%dM", config.RAMMax), fmt.Sprintf("-Xms%dM", config.RAMMin), "-jar", fmt.Sprintf("%s/%s", serverPath, config.JarName), "nogui")
	cmd.Dir = serverPath
	return &ServerExec{cmd: cmd}
}

func (j *ServerExec) Stdout() io.ReadCloser {
	r, _ := j.cmd.StdoutPipe()
	return r
}

func (j *ServerExec) Stdin() io.WriteCloser {
	w, _ := j.cmd.StdinPipe()
	return w
}

func (j *ServerExec) Start() error {
	return j.cmd.Start()
}

func (j *ServerExec) Kill() error {
	return j.cmd.Process.Kill()
}
