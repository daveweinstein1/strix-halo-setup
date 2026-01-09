<p align="center">
  <img src="assets/logo.png" alt="Strix Halo Installer" width="120">
</p>

<h1 align="center">Strix Halo Post-Installer</h1>

<p align="center">
  <strong>Automated setup for AMD Strix Halo (gfx1151) workstations on CachyOS</strong>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/status-under%20development-yellow" alt="Status">
  <img src="https://img.shields.io/badge/license-PolyForm%20Strict-blue" alt="License">
  <img src="https://img.shields.io/badge/platform-CachyOS-green" alt="Platform">
</p>

---

> [!WARNING]
> **This project is under active development and not yet production-ready.**  
> Features may change, break, or be incomplete. Use at your own risk.

> [!NOTE]
> **Interested in this project?** I'd love to hear from you!  
> Open an issue or reach out if you're working with Strix Halo hardware.

---

## Screenshots

<p align="center">
  <strong>Terminal (TUI)</strong><br>
  <img src="assets/tui-screenshot.png" alt="TUI Screenshot" width="500"><br>
  <em>Stage-based installation with progress tracking</em>
</p>

<p align="center">
  <strong>Web UI (Browser)</strong><br>
  <em>Coming soon — graphical interface in development</em>
</p>

---

## Installation Guide

📖 **[Online Installation Guide](https://daveweinstein1.github.io/strix-halo-setup/install-guide.html)** (Recommended)

📄 [Text Version](docs/INSTALL_GUIDE.md)

---

## Quick Install

```bash
curl -fsSL https://bit.ly/strix-halo | sudo bash
```

This downloads and runs the installer binary from GitHub Releases:
```bash
# What the script does:
curl -fsSL "https://github.com/.../strix-install" -o /tmp/strix-install
chmod +x /tmp/strix-install
/tmp/strix-install "$@"
rm -f /tmp/strix-install
```

View the full script: [install.sh](install.sh)

**Direct download** (no bit.ly):
```bash
curl -fsSL https://github.com/daveweinstein1/strix-halo-setup/releases/latest/download/strix-install -o /tmp/s && chmod +x /tmp/s && sudo /tmp/s
```

**Options:** `--tui` (terminal) | `--web` (browser) | `--manual` (select stages) | `--auto` (no prompts)

---

## Stages

| Stage | Purpose |
|-------|---------|
| Kernel Config | IOMMU, device quirks (Beelink E610 fix) |
| Graphics Setup | Mesa 25.3+, LLVM 21.x, Vulkan |
| System Update | Mirrors, packages, essentials |
| LXD Setup | Containers with GPU passthrough |
| Fan Control | lm_sensors, fancontrol (optional) |
| Cleanup | Orphan removal, cache cleanup |
| Validation | Verify kernel, GPU, LXD |
| Desktop Apps | Browsers, Office (optional) |
| Workspaces | `ai-lab`, `dev-lab` containers (optional) |

---

## Why Containers?

AI and development tools are **bleeding edge** — ROCm, PyTorch, and AI coding assistants update frequently with breaking changes. We isolate these in LXD containers so:

- **Host stays stable** — Container breakage can't brick your system
- **Instant rollback** — Restore snapshots when experiments fail
- **Fresh starts** — Delete and recreate containers in minutes

📖 Full details: [Container Strategy](docs/CONTAINER_STRATEGY.md)

---

## Requirements (January 2026)

| Component | Required |
|-----------|----------|
| Kernel | **6.18+** |
| Mesa | **25.3+** |
| ROCm | **7.2+** |
| LLVM | **21.x** |

---

## Supported Hardware

| Device | Status |
|--------|--------|
| Framework Desktop | ✅ Full support |
| Beelink GTR9 Pro | ✅ E610 Ethernet fix applied |
| Minisforum MS-S1 Max | ⚠️ Advisory for Ethernet/USB4 |
| Other Strix Halo | ✅ Generic mode |

---

## License

**[Proprietary - Pre-Release](LICENSE.md)**

This software is currently under development and **not yet licensed for any use**.

- ❌ No permission to use, copy, modify, or distribute
- ❌ No warranties or guarantees
- ✅ Will be released under Apache 2.0 at v1.0

**Why this placeholder?**  
We're finalizing v1.0 before releasing under Apache 2.0. This prevents premature forks of incomplete code.

**ETA:** Apache 2.0 license coming with v1.0 release (January 2026)

**Questions?** [Open an issue](https://github.com/daveweinstein1/strix-halo-setup/issues)

---

<p align="center">
  <strong>Author:</strong> Dave Weinstein<br>
  <strong>Contact:</strong> <a href="https://github.com/daveweinstein1/strix-halo-setup/issues">GitHub Issues</a><br>
  <strong>Updated:</strong> January 2026
</p>
