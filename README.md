<p align="center">
  <img src="assets/logo.png" alt="StrixForge" width="120">
</p>

<h1 align="center">StrixForge</h1>

<p align="center">
  <strong>Forge your Strix Halo into a Professional-Grade AI Workstation.</strong>
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
  <strong>Container Hub (TUI)</strong><br>
  <img src="assets/marketplace-screenshot.png" alt="Marketplace Screenshot" width="500"><br>
  <em>Browse and install community containers (kyuz0, AMD, etc.)</em>
</p>

<p align="center">
  <strong>Web UI (Browser)</strong><br>
  <em>Coming soon — graphical interface in development</em>
</p>

---

## The Container Hub

Strixforge includes a curated "Container Hub" that allows you to pull pre-built, optimized environments directly into your LXD setup.

**Launch the Hub:**
```bash
strixforge --hub
```

**Available Sources:**
* **kyuz0:** Highly optimized Strix Halo toolboxes (ROCm-7rc, LLaMA-Vulkan, PyTorch-2.5).
* **AMD Official:** Standard AMD-maintained AI containers (ComfyUI-ROCm).
* **Community:** Verified contributions for specific workflows.

## Why Containers?

AI and development tools are **bleeding edge** — ROCm, PyTorch, and AI coding assistants update frequently with breaking changes. We isolate these in LXD containers so:

- **Host stays stable** — Container breakage can't brick your system
- **Instant rollback** — Restore snapshots when experiments fail
- **Fresh starts** — Delete and recreate containers in minutes

📖 Full details: [Container Strategy](docs/CONTAINER_STRATEGY.md)

---

## Documentation & Deep Dives

We have written detailed guides to help you understand the architecture and get the most out of your hardware.

* **[The Unified Memory Advantage](https://daveweinstein1.github.io/strixforge/Unified-Memory-Advantage.html)**, Why Strix Halo allows you to run models that even an RTX 5090 cannot touch.
* **[Zero-Risk Architecture: Understanding LXD](https://daveweinstein1.github.io/strixforge/Zero-Risk-Architecture.html)**, How we use containers to isolate the fragile AI stack from your stable OS.
* **[Supported Hardware & Known Quirks](https://daveweinstein1.github.io/strixforge/Supported-Hardware-and-Quirks.html)**, Specific notes for Framework, Minisforum, and Beelink users.
* **[Benchmarks vs. Reality](https://daveweinstein1.github.io/strixforge/Benchmarks-vs-Reality.html)**, Moving past TFLOPS to measure real-world compile times and tokens/sec.
* **[Community Recipes](https://daveweinstein1.github.io/strixforge/Community-Recipes.html)**, Verified stacks for Oobabooga, ComfyUI, and more.
* **[Full Installation Guide](https://daveweinstein1.github.io/strixforge/install-guide.html)**, A step-by-step HTML guide for printing.

---

## The Origin Story

The name **StrixForge** wasn't chosen by accident. In ancient mythology, the *Strix* was a bird of ill omen, a screeching owl that brought terror in the night. When we first got our hands on the AMD Strix Halo hardware, that's exactly what it felt like. The hardware was a beast, massive unified memory, incredible potential, but the software stack was a nightmare of broken dependencies, kernel panics, and fragmented documentation. It screeched at us every time we tried to run a simple inference, "GGGGGGGGG....", over and over. ;-)

This project is to silence the screeching, to take the raw, chaotic potential of the "Strix Halo" APU and put it through the Forge, hammer out the imperfections, tame the drivers, and sharpen the software stack until it became a precise, reliable tool.

## Why We Built This

Just a few months ago, the drivers and core libraries were so disjointed that even a seasoned software engineer would give up after a week of frustration and go shopping for a DGX Spark. While there has been significant effort from AMD and others to get the low-level "plumbing" working, the reality on the ground is still a dependency nightmare.

**StrixForge simplifies this chaos in two specific ways:**

1.  **Automated Configuration:** It handles the tedious configuration and installation of all the base packages you need on the host. It applies critical kernel patches (like the Beelink E610 fix) automatically.
2.  **Containerized Safety:** It sets up LXD containers that allow you to leverage the work of elite developers (like `kyuz0`) or to experiment with bleeding-edge software. We help you simple spin up a container so you can safely try out a new app or technology without putting your development machine at risk. If it makes a mess inside the container, your host system remains pristine.

### Who This Is For
It's for the **"People in the Middle."**

You are a developer, a data scientist, or a power user. You know what a tensor is, and you know why you want local inference. But you **don't** want to be a Linux kernel maintainer.
* You want to write code, not debug `make` files.
* You want to run Llama 3 70B, not re-compile LLVM.
* You want the power of Linux without the "Linux tax" on your time.

## Quick Start

**Prerequisites:**
* A machine with an AMD Strix Halo (Ryzen AI Max+) processor.
* A fresh installation of **CachyOS** (recommended) or Arch Linux. (more Linux distros coming soon!)
* Secure Boot disabled in BIOS.

**Installation:**
Run this single command to download and launch the installer:


📖 **[Online Installation Guide](https://daveweinstein1.github.io/strixforge/install-guide.html)** (Recommended)

📄 [Text Version](docs/INSTALL_GUIDE.md)

---

## Quick Install

```bash
curl -fsSL https://bit.ly/strixforge | sudo bash
```

This downloads and runs the installer binary from GitHub Releases:
```bash
# What the script does:
curl -fsSL "https://github.com/.../strixforge" -o /tmp/strixforge
chmod +x /tmp/strixforge
/tmp/strixforge "$@"
rm -f /tmp/strixforge
```

View the full script: [install.sh](install.sh)

**Direct download** (no bit.ly):
```bash
curl -fsSL https://github.com/daveweinstein1/strixforge/releases/latest/download/strixforge -o /tmp/s && chmod +x /tmp/s && sudo /tmp/s
```

**Installer Options:**

| Flag | Description |
|------|-------------|
| `(none)` | Auto-detects best UI (GUI, then fallback to TUI) |
| `--tui` | Forces Terminal UI |
| `--gui` | Forces Graphical UI (GUI window pops up on localhost) |
| `--auto` | Runs all stages without prompts (Unattended) |
| `--manual` | Interactive stage selection |
| `--hub` | Browse and install community containers (Container Hub) |
| `--check-versions` | Verify package versions |
| `--dry-run` | Simulate without changes |

*Auto-detects GUI if `$DISPLAY` or `$WAYLAND_DISPLAY` is set, otherwise uses TUI.*

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

**Questions?** [Open an issue](https://github.com/daveweinstein1/strixforge/issues)

---

<p align="center">
  <strong>Author:</strong> Dave Weinstein<br>
  <strong>Contact:</strong> <a href="https://github.com/daveweinstein1/strixforge/issues">GitHub Issues</a><br>
  <strong>Updated:</strong> January 2026
</p>
