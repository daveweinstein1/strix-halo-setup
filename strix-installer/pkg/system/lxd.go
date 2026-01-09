package system

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// LXD provides container management operations
type LXD struct{}

// NewLXD creates a new LXD instance
func NewLXD() *LXD {
	return &LXD{}
}

// Init initializes LXD with automatic defaults
func (l *LXD) Init(ctx context.Context) error {
	result, err := ExecSudo(ctx, "lxd", "init", "--auto")
	if err != nil {
		return fmt.Errorf("lxd init failed: %s\n%s", err, result.Stderr)
	}
	return nil
}

// AddUserToGroup adds a user to the lxd group
func (l *LXD) AddUserToGroup(ctx context.Context, user string) error {
	result, err := ExecSudo(ctx, "usermod", "-aG", "lxd", user)
	if err != nil {
		return fmt.Errorf("failed to add user to lxd group: %s\n%s", err, result.Stderr)
	}
	return nil
}

// IsUserInGroup checks if a user is in the lxd group
func (l *LXD) IsUserInGroup(ctx context.Context, user string) bool {
	result, err := Exec(ctx, "groups", user)
	if err != nil {
		return false
	}
	return strings.Contains(result.Stdout, "lxd")
}

// CreateContainer creates a new container from an image
func (l *LXD) CreateContainer(ctx context.Context, name, image string) error {
	result, err := Exec(ctx, "lxc", "launch", image, name)
	if err != nil {
		return fmt.Errorf("failed to create container %s: %s\n%s", name, err, result.Stderr)
	}
	return nil
}

// ContainerExists checks if a container exists
func (l *LXD) ContainerExists(ctx context.Context, name string) bool {
	result, _ := Exec(ctx, "lxc", "info", name)
	return result.ExitCode == 0
}

// DeleteContainer removes a container
func (l *LXD) DeleteContainer(ctx context.Context, name string, force bool) error {
	args := []string{"delete", name}
	if force {
		args = append(args, "--force")
	}
	result, err := Exec(ctx, "lxc", args...)
	if err != nil {
		return fmt.Errorf("failed to delete container %s: %s\n%s", name, err, result.Stderr)
	}
	return nil
}

// ExecInContainer runs a command inside a container
func (l *LXD) ExecInContainer(ctx context.Context, name string, command ...string) (*ExecResult, error) {
	args := append([]string{"exec", name, "--"}, command...)
	return Exec(ctx, "lxc", args...)
}

// SetProfileConfig sets a configuration on the default profile
func (l *LXD) SetProfileConfig(ctx context.Context, key, value string) error {
	result, err := Exec(ctx, "lxc", "profile", "set", "default", key, value)
	if err != nil {
		return fmt.Errorf("failed to set profile config %s=%s: %s\n%s", key, value, err, result.Stderr)
	}
	return nil
}

// AddGPUDevice adds a GPU device to the default profile
func (l *LXD) AddGPUDevice(ctx context.Context) error {
	// Add GPU device with full access
	result, err := Exec(ctx, "lxc", "profile", "device", "add", "default", "gpu", "gpu", "gid=110")
	if err != nil && !strings.Contains(result.Stderr, "already exists") {
		return fmt.Errorf("failed to add GPU device: %s\n%s", err, result.Stderr)
	}
	return nil
}

// EnableNesting enables container nesting (for Docker-in-LXD etc)
func (l *LXD) EnableNesting(ctx context.Context) error {
	return l.SetProfileConfig(ctx, "security.nesting", "true")
}

// ListContainers returns a list of container names
func (l *LXD) ListContainers(ctx context.Context) ([]string, error) {
	result, err := Exec(ctx, "lxc", "list", "--format=json")
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %s", err)
	}

	var containers []struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal([]byte(result.Stdout), &containers); err != nil {
		return nil, err
	}

	names := make([]string, len(containers))
	for i, c := range containers {
		names[i] = c.Name
	}
	return names, nil
}

// WaitForNetwork waits for a container to have network connectivity
func (l *LXD) WaitForNetwork(ctx context.Context, name string) error {
	// Simple ping test
	for i := 0; i < 30; i++ {
		result, _ := l.ExecInContainer(ctx, name, "ping", "-c1", "-W1", "1.1.1.1")
		if result.ExitCode == 0 {
			return nil
		}
		// Wait a bit before retrying
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
	}
	return fmt.Errorf("container %s did not get network connectivity", name)
}
