package system

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// PackageInfo contains metadata about a package
type PackageInfo struct {
	Name           string
	Version        string
	ExpectedMinVer string
	Source         string // "core", "extra", "cachyos", "aur"
	Description    string
	Installed      bool
}

// VersionCheck represents the result of a version comparison
type VersionCheck struct {
	Package      string
	Current      string
	Expected     string
	Status       VersionStatus
	CanProceed   bool
	UserDecision bool // true if user approved despite mismatch
}

type VersionStatus int

const (
	VersionOK VersionStatus = iota
	VersionNewer
	VersionOlder
	VersionMissing
)

func (v VersionStatus) String() string {
	switch v {
	case VersionOK:
		return "✓ OK"
	case VersionNewer:
		return "↑ Newer"
	case VersionOlder:
		return "↓ Older"
	case VersionMissing:
		return "✗ Missing"
	default:
		return "?"
	}
}

// ExpectedVersions defines the versions we expect for Jan 2026
var ExpectedVersions = map[string]string{
	// Graphics
	"mesa":           "25.3.0",
	"vulkan-radeon":  "25.3.0",
	"llvm":           "21.0.0",
	"linux-firmware": "20250108",

	// Kernel
	"linux-cachyos": "6.18.0",

	// ROCm
	"rocm-hip-sdk":        "7.2.0",
	"python-pytorch-rocm": "2.5.1",

	// LXD
	"lxd": "6.2",

	// Tools
	"git":    "2.47.0",
	"neovim": "0.10.0",
}

// CheckPackageVersion compares installed version against expected
func CheckPackageVersion(ctx context.Context, pkg string, expectedMin string) (*VersionCheck, error) {
	pacman := NewPacman()

	check := &VersionCheck{
		Package:  pkg,
		Expected: expectedMin,
	}

	// Get installed version
	installed, err := pacman.GetVersion(ctx, pkg)
	if err != nil || installed == "" {
		check.Current = "not installed"
		check.Status = VersionMissing
		check.CanProceed = true // Can install
		return check, nil
	}

	check.Current = installed

	// Compare versions
	cmp := CompareVersions(installed, expectedMin)
	switch {
	case cmp == 0:
		check.Status = VersionOK
		check.CanProceed = true
	case cmp > 0:
		check.Status = VersionNewer
		check.CanProceed = true // Newer is fine
	case cmp < 0:
		check.Status = VersionOlder
		check.CanProceed = false // User should confirm
	}

	return check, nil
}

// CompareVersions compares two version strings
// Returns: -1 if a < b, 0 if a == b, 1 if a > b
func CompareVersions(a, b string) int {
	partsA := parseVersion(a)
	partsB := parseVersion(b)

	maxLen := len(partsA)
	if len(partsB) > maxLen {
		maxLen = len(partsB)
	}

	for i := 0; i < maxLen; i++ {
		var numA, numB int
		if i < len(partsA) {
			numA = partsA[i]
		}
		if i < len(partsB) {
			numB = partsB[i]
		}

		if numA < numB {
			return -1
		}
		if numA > numB {
			return 1
		}
	}

	return 0
}

// parseVersion extracts numeric parts from version string
func parseVersion(v string) []int {
	// Remove common suffixes
	v = strings.Split(v, "-")[0]
	v = strings.Split(v, "+")[0]
	v = strings.Split(v, "_")[0]

	// Extract numbers
	re := regexp.MustCompile(`\d+`)
	matches := re.FindAllString(v, -1)

	parts := make([]int, len(matches))
	for i, m := range matches {
		num, _ := strconv.Atoi(m)
		parts[i] = num
	}

	return parts
}

// SummarizeVersionChecks returns human-readable summary
func SummarizeVersionChecks(checks []*VersionCheck) string {
	var sb strings.Builder

	ok, newer, older, missing := 0, 0, 0, 0

	for _, c := range checks {
		switch c.Status {
		case VersionOK:
			ok++
		case VersionNewer:
			newer++
		case VersionOlder:
			older++
		case VersionMissing:
			missing++
		}
	}

	sb.WriteString(fmt.Sprintf("Version Check Summary:\n"))
	sb.WriteString(fmt.Sprintf("  ✓ Expected: %d\n", ok))
	if newer > 0 {
		sb.WriteString(fmt.Sprintf("  ↑ Newer:    %d (will use current)\n", newer))
	}
	if older > 0 {
		sb.WriteString(fmt.Sprintf("  ↓ Older:    %d (may need confirmation)\n", older))
	}
	if missing > 0 {
		sb.WriteString(fmt.Sprintf("  ✗ Missing:  %d (will install)\n", missing))
	}

	return sb.String()
}

// =============================================================================
// Phase 11: Pre-Install Version Verification
// =============================================================================

// CheckAllVersions checks all packages in ExpectedVersions map
func CheckAllVersions(ctx context.Context) ([]*VersionCheck, error) {
	checks := make([]*VersionCheck, 0, len(ExpectedVersions))

	for pkg, expectedVer := range ExpectedVersions {
		check, err := CheckPackageVersion(ctx, pkg, expectedVer)
		if err != nil {
			// Log but continue checking other packages
			check = &VersionCheck{
				Package:  pkg,
				Expected: expectedVer,
				Current:  "error",
				Status:   VersionMissing,
			}
		}
		checks = append(checks, check)
	}

	return checks, nil
}

// FormatVersionTable returns a formatted table of version checks
func FormatVersionTable(checks []*VersionCheck) string {
	var sb strings.Builder

	// Header
	sb.WriteString(fmt.Sprintf("%-20s %-15s %-15s %s\n", "Package", "Expected", "Installed", "Status"))
	sb.WriteString(strings.Repeat("-", 60) + "\n")

	// Sort by status (problems first)
	sortedChecks := make([]*VersionCheck, len(checks))
	copy(sortedChecks, checks)

	// Simple sort: older/missing first, then OK/newer
	for i := 0; i < len(sortedChecks); i++ {
		for j := i + 1; j < len(sortedChecks); j++ {
			if sortedChecks[j].Status > sortedChecks[i].Status {
				sortedChecks[i], sortedChecks[j] = sortedChecks[j], sortedChecks[i]
			}
		}
	}

	for _, c := range sortedChecks {
		current := c.Current
		if len(current) > 15 {
			current = current[:12] + "..."
		}
		expected := c.Expected
		if len(expected) > 15 {
			expected = expected[:12] + "..."
		}
		sb.WriteString(fmt.Sprintf("%-20s %-15s %-15s %s\n", c.Package, expected, current, c.Status.String()))
	}

	return sb.String()
}

// HasCriticalMismatches returns true if any package has older version
func HasCriticalMismatches(checks []*VersionCheck) bool {
	for _, c := range checks {
		if c.Status == VersionOlder {
			return true
		}
	}
	return false
}

// GetMismatches returns only packages that need attention (older versions)
func GetMismatches(checks []*VersionCheck) []*VersionCheck {
	mismatches := make([]*VersionCheck, 0)
	for _, c := range checks {
		if c.Status == VersionOlder {
			mismatches = append(mismatches, c)
		}
	}
	return mismatches
}
