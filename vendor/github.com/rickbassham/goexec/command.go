package goexec

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"strings"
)

type cmd struct {
	c *exec.Cmd

	stdout chan string
	stderr chan string
}

type ExecCommand struct {
}

func (ExecCommand) Command(name string, args ...string) Command {
	return NewCommand(name, args...)
}

type Command interface {
	Start() error
	Wait() error
	Kill() error
	Signal(os.Signal) error
	Stdout() (<-chan string, error)
	Stderr() (<-chan string, error)
}

func NewCommand(name string, args ...string) Command {
	return &cmd{
		c: exec.Command(name, args...),
	}
}

func (c *cmd) Start() error {
	return c.c.Start()
}

func (c *cmd) Wait() error {
	return c.c.Wait()
}

func (c *cmd) Kill() error {
	return c.c.Process.Kill()
}

func (c *cmd) Signal(s os.Signal) error {
	return c.c.Process.Signal(s)
}

func (c *cmd) Stdout() (<-chan string, error) {
	stdout, err := c.c.StdoutPipe()
	if err != nil {
		return nil, err
	}

	c.stdout = make(chan string, 10)

	go func() {
		defer close(c.stdout)

		r := bufio.NewReader(stdout)
		var err error
		var line string

		for err == nil {
			line, err = r.ReadString('\n')
			if err == io.EOF {
				return
			}
			if err != nil {
				return
			}

			c.stdout <- strings.TrimSpace(line)
		}
	}()

	return c.stdout, nil
}

func (c *cmd) Stderr() (<-chan string, error) {
	stderr, err := c.c.StderrPipe()
	if err != nil {
		return nil, err
	}

	c.stderr = make(chan string, 10)

	go func() {
		defer close(c.stderr)

		r := bufio.NewReader(stderr)
		var err error
		var line string

		for err == nil {
			line, err = r.ReadString('\n')
			if err == io.EOF {
				return
			}
			if err != nil {
				return
			}

			c.stderr <- strings.TrimSpace(line)
		}
	}()

	return c.stderr, nil
}
