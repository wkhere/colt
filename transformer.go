package colt

import (
	"bytes"
	"io"
	"os"
	"os/exec"
)

type Transformer interface {
	Copy([]byte) error
	Transform([]byte) error
}

type CommandT struct {
	Command        []string
	Stdout, Stderr io.Writer
}

func (c *CommandT) Copy(b []byte) error {
	_, err := c.Stdout.Write(b)
	return err
}

func (c *CommandT) Transform(d []byte) error {
	var b bytes.Buffer
	cmd := exec.Command(c.Command[0], append(c.Command[1:], string(d))...)
	cmd.Env = append(os.Environ(), "COLOR=1")
	cmd.Stdout = &b
	cmd.Stderr = c.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}
	_, err = c.Stdout.Write(chomp(b.Bytes()))
	return err
}
