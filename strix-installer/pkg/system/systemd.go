package system

import (
	"context"
	"fmt"
	"strings"
)

// Systemd provides service management operations
type Systemd struct{}

// NewSystemd creates a new Systemd instance
func NewSystemd() *Systemd {
	return &Systemd{}
}

// Enable enables a service (starts on boot)
func (s *Systemd) Enable(ctx context.Context, service string) error {
	result, err := ExecSudo(ctx, "systemctl", "enable", service)
	if err != nil {
		return fmt.Errorf("failed to enable %s: %s\n%s", service, err, result.Stderr)
	}
	return nil
}

// Start starts a service
func (s *Systemd) Start(ctx context.Context, service string) error {
	result, err := ExecSudo(ctx, "systemctl", "start", service)
	if err != nil {
		return fmt.Errorf("failed to start %s: %s\n%s", service, err, result.Stderr)
	}
	return nil
}

// EnableAndStart enables and starts a service
func (s *Systemd) EnableAndStart(ctx context.Context, service string) error {
	if err := s.Enable(ctx, service); err != nil {
		return err
	}
	return s.Start(ctx, service)
}

// Stop stops a service
func (s *Systemd) Stop(ctx context.Context, service string) error {
	result, err := ExecSudo(ctx, "systemctl", "stop", service)
	if err != nil {
		return fmt.Errorf("failed to stop %s: %s\n%s", service, err, result.Stderr)
	}
	return nil
}

// Disable disables a service
func (s *Systemd) Disable(ctx context.Context, service string) error {
	result, err := ExecSudo(ctx, "systemctl", "disable", service)
	if err != nil {
		return fmt.Errorf("failed to disable %s: %s\n%s", service, err, result.Stderr)
	}
	return nil
}

// IsActive checks if a service is running
func (s *Systemd) IsActive(ctx context.Context, service string) bool {
	result, _ := Exec(ctx, "systemctl", "is-active", service)
	return strings.TrimSpace(result.Stdout) == "active"
}

// IsEnabled checks if a service is enabled
func (s *Systemd) IsEnabled(ctx context.Context, service string) bool {
	result, _ := Exec(ctx, "systemctl", "is-enabled", service)
	return strings.TrimSpace(result.Stdout) == "enabled"
}

// Status returns the status of a service
func (s *Systemd) Status(ctx context.Context, service string) (string, error) {
	result, err := Exec(ctx, "systemctl", "status", service)
	if err != nil && result.ExitCode != 3 { // Exit 3 means service is stopped (valid)
		return "", err
	}
	return result.Stdout, nil
}

// DaemonReload reloads systemd configuration
func (s *Systemd) DaemonReload(ctx context.Context) error {
	_, err := ExecSudo(ctx, "systemctl", "daemon-reload")
	return err
}
