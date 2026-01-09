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

// =============================================================================
// Phase 10: Container Lifecycle Management
// =============================================================================

// Snapshot represents an LXD container snapshot
type Snapshot struct {
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	Stateful  bool   `json:"stateful"`
}

// CreateSnapshot creates a snapshot of a container
func (l *LXD) CreateSnapshot(ctx context.Context, container, snapshotName string) error {
	result, err := Exec(ctx, "lxc", "snapshot", container, snapshotName)
	if err != nil {
		return fmt.Errorf("failed to create snapshot %s/%s: %s\n%s", container, snapshotName, err, result.Stderr)
	}
	return nil
}

// RestoreSnapshot restores a container to a previous snapshot
func (l *LXD) RestoreSnapshot(ctx context.Context, container, snapshotName string) error {
	result, err := Exec(ctx, "lxc", "restore", container, snapshotName)
	if err != nil {
		return fmt.Errorf("failed to restore snapshot %s/%s: %s\n%s", container, snapshotName, err, result.Stderr)
	}
	return nil
}

// DeleteSnapshot removes a snapshot from a container
func (l *LXD) DeleteSnapshot(ctx context.Context, container, snapshotName string) error {
	result, err := Exec(ctx, "lxc", "delete", fmt.Sprintf("%s/%s", container, snapshotName))
	if err != nil {
		return fmt.Errorf("failed to delete snapshot %s/%s: %s\n%s", container, snapshotName, err, result.Stderr)
	}
	return nil
}

// ListSnapshots returns all snapshots for a container
func (l *LXD) ListSnapshots(ctx context.Context, container string) ([]Snapshot, error) {
	result, err := Exec(ctx, "lxc", "info", container, "--format=json")
	if err != nil {
		return nil, fmt.Errorf("failed to get container info: %s", err)
	}

	var info struct {
		Snapshots []Snapshot `json:"snapshots"`
	}
	if err := json.Unmarshal([]byte(result.Stdout), &info); err != nil {
		return nil, err
	}

	return info.Snapshots, nil
}

// ContainerStatus represents the state of a container
type ContainerStatus struct {
	Name    string
	Status  string // "Running", "Stopped", "Frozen"
	Created string
}

// GetContainerStatus returns detailed status of a container
func (l *LXD) GetContainerStatus(ctx context.Context, name string) (*ContainerStatus, error) {
	result, err := Exec(ctx, "lxc", "info", name, "--format=json")
	if err != nil {
		return nil, fmt.Errorf("failed to get container status: %s", err)
	}

	var info struct {
		Name      string `json:"name"`
		Status    string `json:"status"`
		CreatedAt string `json:"created_at"`
	}
	if err := json.Unmarshal([]byte(result.Stdout), &info); err != nil {
		return nil, err
	}

	return &ContainerStatus{
		Name:    info.Name,
		Status:  info.Status,
		Created: info.CreatedAt,
	}, nil
}

// StopContainer stops a running container
func (l *LXD) StopContainer(ctx context.Context, name string, force bool) error {
	args := []string{"stop", name}
	if force {
		args = append(args, "--force")
	}
	result, err := Exec(ctx, "lxc", args...)
	if err != nil {
		return fmt.Errorf("failed to stop container %s: %s\n%s", name, err, result.Stderr)
	}
	return nil
}

// StartContainer starts a stopped container
func (l *LXD) StartContainer(ctx context.Context, name string) error {
	result, err := Exec(ctx, "lxc", "start", name)
	if err != nil {
		return fmt.Errorf("failed to start container %s: %s\n%s", name, err, result.Stderr)
	}
	return nil
}

// RecreateContainer deletes and recreates a container from image
func (l *LXD) RecreateContainer(ctx context.Context, name, image string) error {
	// Stop if running
	status, err := l.GetContainerStatus(ctx, name)
	if err == nil && status.Status == "Running" {
		if err := l.StopContainer(ctx, name, true); err != nil {
			return err
		}
	}

	// Delete
	if err := l.DeleteContainer(ctx, name, true); err != nil {
		return err
	}

	// Create fresh
	return l.CreateContainer(ctx, name, image)
}

