package system

import (
	"context"
	"fmt"
	"strings"
)

// Pacman provides package management operations
type Pacman struct{}

// NewPacman creates a new Pacman instance
func NewPacman() *Pacman {
	return &Pacman{}
}

// Install installs packages
func (p *Pacman) Install(ctx context.Context, packages ...string) error {
	args := append([]string{"-S", "--needed", "--noconfirm"}, packages...)
	result, err := ExecSudo(ctx, "pacman", args...)
	if err != nil {
		return fmt.Errorf("pacman install failed: %s\n%s", err, result.Stderr)
	}
	return nil
}

// Update performs a full system update
func (p *Pacman) Update(ctx context.Context) error {
	result, err := ExecSudo(ctx, "pacman", "-Syu", "--noconfirm")
	if err != nil {
		return fmt.Errorf("pacman update failed: %s\n%s", err, result.Stderr)
	}
	return nil
}

// Remove removes packages
func (p *Pacman) Remove(ctx context.Context, packages ...string) error {
	args := append([]string{"-Rns", "--noconfirm"}, packages...)
	result, err := ExecSudo(ctx, "pacman", args...)
	if err != nil {
		return fmt.Errorf("pacman remove failed: %s\n%s", err, result.Stderr)
	}
	return nil
}

// IsInstalled checks if a package is installed
func (p *Pacman) IsInstalled(ctx context.Context, pkg string) bool {
	result, err := Exec(ctx, "pacman", "-Q", pkg)
	return err == nil && result.ExitCode == 0
}

// GetVersion returns the installed version of a package
func (p *Pacman) GetVersion(ctx context.Context, pkg string) (string, error) {
	result, err := Exec(ctx, "pacman", "-Q", pkg)
	if err != nil {
		return "", fmt.Errorf("package not installed: %s", pkg)
	}
	parts := strings.Fields(result.Stdout)
	if len(parts) >= 2 {
		return parts[1], nil
	}
	return "", fmt.Errorf("could not parse version for %s", pkg)
}

// CleanOrphans removes orphaned packages
func (p *Pacman) CleanOrphans(ctx context.Context) error {
	// First check if there are orphans
	orphans, err := Exec(ctx, "pacman", "-Qtdq")
	if err != nil || strings.TrimSpace(orphans.Stdout) == "" {
		// No orphans
		return nil
	}

	// Remove orphans
	_, err = ExecShellSudo(ctx, "pacman -Rns --noconfirm $(pacman -Qtdq)")
	return err
}

// CleanCache cleans the package cache
func (p *Pacman) CleanCache(ctx context.Context) error {
	_, err := ExecShellSudo(ctx, "echo y | pacman -Scc")
	return err
}

// Yay provides AUR package management (runs as user, not root)
type Yay struct {
	user string
}

// NewYay creates a new Yay instance for the specified user
func NewYay(user string) *Yay {
	return &Yay{user: user}
}

// Install installs AUR packages
func (y *Yay) Install(ctx context.Context, packages ...string) error {
	args := append([]string{"-u", y.user, "yay", "-S", "--needed", "--noconfirm"}, packages...)
	result, err := Exec(ctx, "sudo", args...)
	if err != nil {
		return fmt.Errorf("yay install failed: %s\n%s", err, result.Stderr)
	}
	return nil
}

// IsInstalled checks if a package is installed (works for AUR too)
func (y *Yay) IsInstalled(ctx context.Context, pkg string) bool {
	result, err := Exec(ctx, "pacman", "-Q", pkg)
	return err == nil && result.ExitCode == 0
}
