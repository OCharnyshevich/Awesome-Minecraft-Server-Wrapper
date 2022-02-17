package minecraft

/* #include <unistd.h> */
import "C"
import (
	"fmt"
	"github.com/c9s/goprocinfo/linux"
	"io"
	"log"
	"os/exec"
)

type Process interface {
	Stat() (*Stat, error)
	Stdout() io.ReadCloser
	Stdin() io.WriteCloser
	Start() error
	Kill() error
}

type ProcessExec struct {
	cmd *exec.Cmd
}

type Stat struct {
	PID    int
	Memory uint // In megabytes
}

func NewProcess(config *Config) *ProcessExec {
	cmd := exec.Command(config.JavaPath, fmt.Sprintf("-Xmx%dM", config.RAMMax), fmt.Sprintf("-Xms%dM", config.RAMMin), "-jar", config.GetPathJar(), "nogui") //nolint:gosec
	log.Println(cmd.Args)
	cmd.Dir = config.GetPath()
	return &ProcessExec{cmd: cmd}
}

func (e *ProcessExec) Stdout() io.ReadCloser {
	r, _ := e.cmd.StdoutPipe()
	return r
}

func (e *ProcessExec) Stdin() io.WriteCloser {
	w, _ := e.cmd.StdinPipe()
	return w
}

func (e *ProcessExec) Start() error {
	return e.cmd.Start()
}

func (e *ProcessExec) Kill() error {
	return e.cmd.Process.Kill()
}

func (e ProcessExec) readProcessStatus() (*linux.ProcessStatus, error) {
	return linux.ReadProcessStatus(fmt.Sprintf("/proc/%d/status", e.cmd.Process.Pid))
}

func (e ProcessExec) Stat() (*Stat, error) {
	status, err := linux.ReadProcessStatus(fmt.Sprintf("/proc/%d/status", e.cmd.Process.Pid))
	if err != nil {
		return nil, err
	}

	return &Stat{
		PID:    e.cmd.Process.Pid,
		Memory: uint(status.VmSize / 1024),
	}, err
}
