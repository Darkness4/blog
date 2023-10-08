package mermaid

import (
	"os/exec"

	"go.abhg.dev/goldmark/mermaid"
)

type CLI struct {
	mermaid.CLI
}

// Command builds an exec.Cmd to run 'mmdc' with the given arguments.
func (c *CLI) Command(args ...string) *exec.Cmd {
	args = append(
		[]string{"-b", "transparent"},
		args...,
	)
	return c.CLI.Command(args...)
}
