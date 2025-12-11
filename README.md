# CachyOS Strix Halo Automation Scripts

**Automated post-install setup for AMD Strix Halo (gfx1150) workstations on CachyOS.**

> [!TIP]
> **Start Here**: Read [docs/INSTALL_GUIDE.md](docs/INSTALL_GUIDE.md) for the complete end-to-end guide, including creating the USB, BIOS settings (Beelink/Framework), and manual OS installation.

## 🚀 Quick Start

Once CachyOS is installed and you are logged in:

```bash
# 1. Copy scripts to your home directory (if from USB)
cp -r /path/to/usb/strix-halo-setup ~/
cd ~/strix-halo-setup

# 2. Make executable
chmod +x master-control.sh stages/*.sh

# 3. Designate yourself as the Master Control Program
./master-control.sh
```

## 🛠️ Automation Stages

The `master-control.sh` script guides you through these verified stages:

| Stage | Name | Purpose |
|-------|------|---------|
| **1** | **Kernel Config** | Enables IOMMU, applies critical E610 Ethernet fix (Beelink), helps set TDP. |
| **2** | **Graphics Setup** | Verifies Mesa 25.3+, LLVM 21.x, and Vulkan drivers for Strix Halo. |
| **3** | **System Update** | Configures CachyOS repositories (v3/v4) and updates base packages. |
| **4** | **LXD Containerization** | Installs LXD, configures networking, and sets up GPU Passthrough. |
| **5** | **Cleanup** | Removes conflicting "AI Sysadmin" packages and unnecessary bloat. |
| **6** | **Validation** | Verifies Kernel, GPU, and IOMMU status before proceeding. |
| **7** | **User Applications** | (Optional) Installs Browsers, IDEs (Antigravity), and Office tools. |
| **8** | **Workspaces** | (New) Provisions `ai-lab` (ROCm/PyTorch) and `dev-lab` (Rust/Go) containers. |

## 📚 Documentation

All detailed documentation is consolidated in:
*   **[docs/INSTALL_GUIDE.md](docs/INSTALL_GUIDE.md)**: The primary manual.
    *   **Appendix A**: Software Package Manifest.
    *   **Appendix B**: Troubleshooting Guide.
*   **[docs/VERSION_REQUIREMENTS.md](docs/VERSION_REQUIREMENTS.md)**: Specific version pins for Kernel (6.18+), ROCm (7.1+), and Mesa.

## ⚠️ Key Features

*   **Beelink GTR9 Pro Support**: Stage 1 includes a Fix for the E610 Ethernet crash and a TDP control tool.
*   **Framework Desktop Support**: Verified "Golden Path" configuration.
*   **AI-Ready**: Stage 8 sets up "Bleeding Edge" AI containers with ROCm 7.1 pre-installed.
*   **Resume Capability**: Failed scripts can be resumed mid-stream (see Install Guide).
*   **AI Error Logs**: Failures generate context-rich logs in `logs/` optimized for AI diagnosis.

---
**Version**: 2025.12.11 | **Author**: Dave Weinstein (@daveweinstein1)

## Hardware Target
- **CPU**: AMD Strix Halo
- **GPU**: gfx1150 (integrated)
- **Special Requirements**: Latest kernel, Mesa, and firmware for proper gfx1150 support

## Installation Stages

### Stage 0: Pre-Installation Preparation
- [ ] Backup existing data
- [ ] Create bootable CachyOS USB
- [ ] Document current system (if migrating)
- [ ] **Script**: `00-pre-install.sh`

### Stage 1: Base System Installation
- [ ] Boot from USB
- [ ] Partition disk
- [ ] Install base CachyOS
- [ ] Configure bootloader
- [ ] **Script**: `01-base-install.sh`

### Stage 2: Kernel Configuration
- [ ] Disable E610 driver
- [ ] Add kernel parameters
- [ ] Configure GRUB
- [ ] **Script**: `02-kernel-config.sh`

### Stage 3: Graphics Stack (Critical for gfx1150)
- [ ] Install latest Mesa
- [ ] Install AMD firmware
- [ ] Install Vulkan drivers
- [ ] **Script**: `03-graphics-setup.sh`

### Stage 4: Package Updates & Repository Configuration
- [ ] Configure CachyOS repositories
- [ ] Update system to latest packages
- [ ] Install build tools
- [ ] **Script**: `04-system-update.sh`

### Stage 5: LXD/LXC Container System
- [ ] Fix LXD dependencies
- [ ] Configure container networking
- [ ] Test container creation
- [ ] **Script**: `05-lxd-setup.sh`

### Stage 6: Cleanup & Optimization
- [ ] Remove AI sysadmin components
- [ ] Clean unnecessary packages
- [ ] Configure system services
- [ ] **Script**: `06-cleanup.sh`

### Stage 7: Validation & Testing
- [ ] GPU rendering test
- [ ] Container functionality test
- [ ] Boot stability test
- [ ] **Script**: `07-validation.sh`

## Script Features

Each installation script includes:
- ✅ **Breakpoints**: Confirm before each major command
- 🤖 **AI Help Prompts**: Copy-paste error messages for AI assistance
- 📋 **Validation**: Check if each step succeeded
- 🔄 **Rollback**: Instructions if a step fails
- 📝 **Logging**: All output saved to `logs/` directory

## Quick Start

```bash
# Clone/download this guide
cd ~/strix-halo-install

# Make scripts executable
chmod +x *.sh

# Start with pre-installation prep
./00-pre-install.sh
```

## Critical Notes for Strix Halo

> [!IMPORTANT]
> **Version Requirements (VERIFIED Dec 2025):**
> See [VERSION_REQUIREMENTS.md](./VERSION_REQUIREMENTS.md) for **verified current package versions**.
> All version numbers have been researched and confirmed as of December 11, 2025.

> [!IMPORTANT]
> **gfx1150 requires recent packages:**
> - Kernel: 6.15+ (6.18 LTS recommended - released Nov 2025)
> - Mesa: 24.1+ (25.3.1 current - released Dec 2025)
> - LLVM: 17+ (21.1.7 current - released Dec 2025)
> - Firmware: linux-firmware latest from repos

> [!WARNING]
> **E610 driver conflict:**
> The E610 network driver must be disabled via kernel parameters to avoid system instability on some Strix Halo configurations.

> [!CAUTION]
> **LXD dependency issues:**
> CachyOS may have conflicts between LXD and certain kernel modules. Stage 5 addresses these systematically.

## Troubleshooting

See [TROUBLESHOOTING.md](./TROUBLESHOOTING.md) for common issues and solutions.

## Architecture

See [ARCHITECTURE.md](./ARCHITECTURE.md) for technical details about the installation approach.
