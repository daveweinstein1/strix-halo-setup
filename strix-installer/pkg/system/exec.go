package system

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
)

// ExecResult contains the output of a command execution
type ExecResult struct {
	Command  string
	ExitCode int
	Stdout   string
	Stderr   string
}

// Exec runs a command and returns the result
func Exec(ctx context.Context, name string, args ...string) (*ExecResult, error) {
	cmd := exec.CommandContext(ctx, name, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	result := &ExecResult{
		Command: fmt.Sprintf("%s %s", name, strings.Join(args, " ")),
		Stdout:  stdout.String(),
		Stderr:  stderr.String(),
	}

	if cmd.ProcessState != nil {
		result.ExitCode = cmd.ProcessState.ExitCode()
	}

	return result, err
}

// ExecSudo runs a command with sudo
func ExecSudo(ctx context.Context, name string, args ...string) (*ExecResult, error) {
	sudoArgs := append([]string{name}, args...)
	return Exec(ctx, "sudo", sudoArgs...)
}

// ExecShell runs a shell command string
func ExecShell(ctx context.Context, command string) (*ExecResult, error) {
	return Exec(ctx, "bash", "-c", command)
}

// ExecShellSudo runs a shell command with sudo
func ExecShellSudo(ctx context.Context, command string) (*ExecResult, error) {
	return ExecSudo(ctx, "bash", "-c", command)
}

// CheckCommand verifies a command exists
func CheckCommand(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}
