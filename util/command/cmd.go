package command

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"syscall"
	"time"
)

func NewCmdAsUser(command string) *exec.Cmd {
	if runtime.GOOS == "windows" {
		return exec.Command("cmd", "/C", command)
	}

	shell, ok := os.LookupEnv("SHELL")
	if !ok {
		shell = "sh"
	}

	cmd := exec.Command(shell, "-c", command)

	if os.Geteuid() == 0 {
		uidStr, ok := os.LookupEnv("SUDO_UID")
		if ok {
			cmd = exec.Command("sudo", "-E", "-u", "#"+uidStr, shell, "-c", command)
		}
	}

	cmd.SysProcAttr = defaultSysProcAttr()

	return cmd
}

type Cmd struct {
	cmd                    *exec.Cmd
	done                   chan struct{}
	err                    error
	stdoutPipe, stderrPipe io.ReadCloser

	env []string
	dir string
}

type CmdOption func(*Cmd)

func WithEnv(env []string) CmdOption {
	return func(c *Cmd) {
		c.env = env
	}
}

func WithDir(dir string) CmdOption {
	return func(c *Cmd) {
		c.dir = dir
	}
}

func New(cmd *exec.Cmd, opts ...CmdOption) (*Cmd, error) {
	c := &Cmd{
		done: make(chan struct{}),
		cmd:  cmd,
	}

	for _, opt := range opts {
		opt(c)
	}

	// Env.
	cmd.Env = append(os.Environ(),
		c.env...,
	)

	cmd.Dir = c.dir

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	c.stdoutPipe = stdoutPipe
	c.stderrPipe = stderrPipe

	return &Cmd{
		done:       make(chan struct{}),
		stdoutPipe: stdoutPipe,
		stderrPipe: stderrPipe,
		cmd:        cmd,
	}, nil
}

func (i *Cmd) Run() error {
	err := i.cmd.Start()
	if err != nil {
		return err
	}

	go func() {
		err := i.cmd.Wait()
		if err != nil {
			i.err = fmt.Errorf("exited: %s", err)
		}

		close(i.done)
	}()

	return nil
}

func (i *Cmd) IsRunning() bool {
	if i.cmd == nil || i.cmd.Process == nil {
		return false
	}

	select {
	case <-i.done:
		return false
	default:
		return true
	}
}

func (i *Cmd) Stdout() io.ReadCloser {
	return i.stdoutPipe
}

func (i *Cmd) SetStdin(r io.Reader) {
	i.cmd.Stdin = r
}

func (i *Cmd) Stderr() io.ReadCloser {
	return i.stderrPipe
}

func (i *Cmd) Wait() error {
	if i.cmd == nil || i.cmd.Process == nil {
		return nil
	}

	<-i.done

	return i.err
}

func (i *Cmd) WaitChannel() <-chan struct{} {
	return i.done
}

func (i *Cmd) Stop(cleanupTimeout time.Duration) error {
	if i.cmd == nil || i.cmd.Process == nil {
		return nil
	}

	_ = CmdSignal(i.cmd, syscall.SIGINT)

	go func() {
		time.Sleep(cleanupTimeout)

		if i.IsRunning() {
			_ = CmdSignal(i.cmd, syscall.SIGKILL)
		}
	}()

	return i.Wait()
}
