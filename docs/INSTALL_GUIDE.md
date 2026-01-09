# CachyOS Installation Guide for Strix Halo

**Version**: 2026.01.09 | **Author**: Dave Weinstein

## Overview

This guide covers:
1. Installing CachyOS on Strix Halo hardware
2. Running the post-install configuration

---

## Part 1: CachyOS Installation

### Requirements

- USB drive (16GB+)
- Strix Halo hardware (Framework Desktop, Beelink GTR9, Minisforum S1, etc.)
- Internet connection

### Step 1: Create Bootable USB

**Option A: Ventoy (Recommended)**
1. Download [Ventoy](https://www.ventoy.net/en/download.html)
2. Install to USB drive
3. Copy [CachyOS ISO](https://cachyos.org/download/) to USB

**Option B: Direct Write**
```bash
# Linux
sudo dd if=cachyos-desktop-linux.iso of=/dev/sdX bs=4M status=progress

# Windows
# Use Rufus or balenaEtcher
```

### Step 2: BIOS Configuration

| Setting | Required Value |
|---------|----------------|
| Secure Boot | **Disabled** |
| Boot Mode | **UEFI** |
| IOMMU | **Enabled** |
| SVM Mode | **Enabled** |

**Beelink GTR9 Pro only:**
- Disable E610 Ethernet: `Advanced` → `Demo Board` → `PCI-E Port` → `Device 3 Fun 2` → **Disabled**

### Step 3: Install CachyOS

1. Boot from USB
2. Run the Calamares installer
3. Select partitioning (BTRFS recommended)
4. Complete installation
5. Reboot

---

## Part 2: Post-Install Configuration

### Quick Install (One Command)

```bash
curl -fsSL https://bit.ly/strix-halo | sudo bash
```

Or direct from GitHub:
```bash
curl -fsSL https://github.com/daveweinstein1/strix-halo-setup/releases/latest/download/strix-install -o /tmp/s && chmod +x /tmp/s && sudo /tmp/s
```

### What the Installer Does

| Stage | Purpose |
|-------|---------|
| Kernel Config | IOMMU, device quirks |
| Graphics Stack | Mesa 25.3+, Vulkan |
| System Update | Mirror ranking, updates |
| LXD Setup | Containers, GPU passthrough |
| Thermal Control | Fan management (optional) |
| Cleanup | Remove orphans |
| Validation | Verify setup |
| Desktop Apps | Browsers, Office (optional) |
| Workspaces | ai-lab, dev-lab (optional) |

### Options

```bash
sudo ./strix-install --tui    # Terminal mode
sudo ./strix-install --web    # Browser mode
sudo ./strix-install --menu   # Select stages
```

---

## Device-Specific Notes

### Beelink GTR9 Pro
- **E610 Ethernet**: Must be disabled in BIOS *and* kernel blacklist
- Installer applies `modprobe.blacklist=ice` automatically

### Framework Desktop
- Works out of the box
- No quirks needed

### Minisforum S1 Max
- Advisory for USB4 and Ethernet edge cases
- Works for most users

---

## Troubleshooting

### No display after boot
Add `nomodeset` to kernel parameters temporarily, then run installer.

### Script fails
Check logs at `~/.config/strix-install/logs/`

### System won't boot
Boot from CachyOS USB, chroot, and fix:
```bash
sudo mount /dev/nvme0n1p2 /mnt
sudo arch-chroot /mnt
journalctl -xb
```

---

## Validation Checklist

After installation, verify:

- [ ] `glxinfo | grep -i renderer` shows RADV
- [ ] `vulkaninfo` runs without errors
- [ ] `lxc list` shows no errors
- [ ] System stable through reboots
