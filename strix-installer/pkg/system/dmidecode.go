package system

import (
	"context"
	"strings"
)

// DMIDecode provides hardware detection via dmidecode
type DMIDecode struct{}

// NewDMIDecode creates a new DMIDecode instance
func NewDMIDecode() *DMIDecode {
	return &DMIDecode{}
}

// GetSystemManufacturer returns the system manufacturer
func (d *DMIDecode) GetSystemManufacturer(ctx context.Context) (string, error) {
	result, err := ExecSudo(ctx, "dmidecode", "-s", "system-manufacturer")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(result.Stdout), nil
}

// GetProductName returns the product/model name
func (d *DMIDecode) GetProductName(ctx context.Context) (string, error) {
	result, err := ExecSudo(ctx, "dmidecode", "-s", "system-product-name")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(result.Stdout), nil
}

// GetSystemFamily returns the system family
func (d *DMIDecode) GetSystemFamily(ctx context.Context) (string, error) {
	result, err := ExecSudo(ctx, "dmidecode", "-s", "system-family")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(result.Stdout), nil
}

// GetBIOSVersion returns the BIOS version
func (d *DMIDecode) GetBIOSVersion(ctx context.Context) (string, error) {
	result, err := ExecSudo(ctx, "dmidecode", "-s", "bios-version")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(result.Stdout), nil
}

// GetProcessorVersion returns the CPU info
func (d *DMIDecode) GetProcessorVersion(ctx context.Context) (string, error) {
	result, err := ExecSudo(ctx, "dmidecode", "-s", "processor-version")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(result.Stdout), nil
}

// IsStrixHalo checks if the CPU is AMD Strix Halo
func (d *DMIDecode) IsStrixHalo(ctx context.Context) bool {
	proc, err := d.GetProcessorVersion(ctx)
	if err != nil {
		return false
	}
	// Strix Halo CPUs are Ryzen AI Max series
	return strings.Contains(proc, "Ryzen AI Max") ||
		strings.Contains(proc, "Ryzen AI 9 HX")
}
