<p align="center">
  <img src="strix-halo-setup/assets/logo.png" alt="Strix Halo Installer" width="120">
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
  <img src="strix-halo-setup/assets/tui-screenshot.png" alt="TUI Screenshot" width="500"><br>
  <em>Simulated screenshot — actual UI coming soon</em>
</p>

<p align="center">
  <strong>Web UI (Browser)</strong><br>
  <img src="strix-halo-setup/assets/gui-screenshot.png" alt="GUI Screenshot" width="500"><br>
  <em>Simulated screenshot — actual UI coming soon</em>
</p>

---

## Quick Install

**Short version** (via bit.ly — *coming soon*):
```bash
curl -fsSL https://bit.ly/strix-halo | sudo bash
```

<details>
<summary><strong>What does this do?</strong> (click to expand)</summary>

The short URL redirects to this bootstrap script:
```
https://raw.githubusercontent.com/daveweinstein1/strix-halo-setup/main/install.sh
```

Which does the following:
```bash
#!/bin/bash
# Downloads the installer binary from GitHub Releases
curl -fsSL "https://github.com/daveweinstein1/strix-halo-setup/releases/latest/download/strix-install" -o /tmp/strix-install

# Makes it executable
chmod +x /tmp/strix-install

# Runs the installer
/tmp/strix-install "$@"

# Cleans up
rm -f /tmp/strix-install
```

View the full script: [install.sh](strix-halo-setup/install.sh)

</details>

**Direct download** (no bit.ly):
```bash
curl -fsSL https://github.com/daveweinstein1/strix-halo-setup/releases/latest/download/strix-install -o /tmp/s && chmod +x /tmp/s && sudo /tmp/s
```

**Options:** `--tui` (terminal) | `--web` (browser) | `--menu` (select stages)

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

📖 Full details: [Container Strategy](strix-halo-setup/docs/CONTAINER_STRATEGY.md)

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

**[PolyForm Strict 1.0.0](strix-halo-setup/LICENSE)**

- ✅ Personal/noncommercial use allowed
- ✅ Educational and research use allowed
- ❌ Commercial use requires separate license
- ❌ Redistribution not permitted

**Hardware manufacturers & integrators:** Interested in bundling this with your Strix Halo products? [Open an issue](https://github.com/daveweinstein1/strix-halo-setup/issues) to discuss licensing.

---

<p align="center">
  <strong>Author:</strong> Dave Weinstein<br>
  <strong>Contact:</strong> <a href="https://github.com/daveweinstein1/strix-halo-setup/issues">GitHub Issues</a><br>
  <strong>Updated:</strong> January 2026
</p>

