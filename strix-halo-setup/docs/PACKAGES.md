# Package Inventory

*Strix Halo Post-Install — January 2026*

> [!NOTE]
> Version numbers reflect packages available in CachyOS repos as of January 2026.
> Actual versions may vary based on your system's repository state.

---

## Graphics Stack (Stage 2)

| Package | Source | Version | Description |
|---------|--------|---------|-------------|
| `mesa` | CachyOS | 25.3.0 | Open-source OpenGL/Vulkan implementation for AMD GPUs |
| `lib32-mesa` | CachyOS | 25.3.0 | 32-bit Mesa libraries for legacy/Wine applications |
| `mesa-utils` | Extra | 9.0.0 | Mesa utilities including `glxinfo` and `glxgears` |
| `vulkan-radeon` | CachyOS | 25.3.0 | Vulkan driver for AMD RDNA/RDNA2/RDNA3 GPUs |
| `lib32-vulkan-radeon` | CachyOS | 25.3.0 | 32-bit Vulkan driver for Wine/Proton gaming |
| `vulkan-tools` | Extra | 1.3.290 | Vulkan utilities including `vulkaninfo` |
| `linux-firmware` | Core | 20250108 | Firmware files for AMD GPUs and other hardware |
| `llvm` | CachyOS | 21.0.0 | LLVM compiler infrastructure (required for Mesa) |
| `lib32-llvm` | CachyOS | 21.0.0 | 32-bit LLVM libraries |

---

## System Essentials (Stage 3)

| Package | Source | Version | Description |
|---------|--------|---------|-------------|
| `base-devel` | Core | meta | Meta-package for building software (gcc, make, etc.) |
| `git` | Extra | 2.47.1 | Distributed version control system |
| `wget` | Extra | 1.25 | Network downloader for HTTP/HTTPS/FTP |
| `curl` | Core | 8.11.1 | Command-line tool for transferring data with URLs |
| `vim` | Extra | 9.1.0950 | Vi IMproved - classic terminal text editor |
| `neovim` | Extra | 0.10.3 | Hyperextensible Vim-based text editor |
| `btop` | Extra | 1.4.0 | Resource monitor with modern terminal UI |
| `neofetch` | Extra | 7.1.0 | System information tool with ASCII art |
| `fastfetch` | Extra | 2.34.0 | Faster alternative to neofetch |

---

## Containerization (Stage 4)

| Package | Source | Version | Description |
|---------|--------|---------|-------------|
| `lxd` | Extra | 6.2 | System container and VM manager (Canonical) |

---

## Fan & Thermal Control (Stage 5)

| Package | Source | Version | Description |
|---------|--------|---------|-------------|
| `lm_sensors` | Extra | 3.6.0 | Hardware monitoring tools for temperature/voltage/fans |
| `fancontrol` | Extra | 3.6.0 | Fan speed control daemon (part of lm_sensors) |

---

## Desktop Applications (Stage 7) — Optional

### Official Repositories

| Package | Source | Version | Description |
|---------|--------|---------|-------------|
| `firefox` | Extra | 134.0 | Mozilla Firefox web browser |
| `signal-desktop` | Extra | 7.38.0 | Signal encrypted messaging client |
| `vlc` | Extra | 3.0.21 | VLC media player - plays almost anything |
| `yay` | AUR Helper | 12.4.2 | Yet Another Yogurt - AUR helper for installing community packages |

### AUR (Arch User Repository)

| Package | Source | Version | Description |
|---------|--------|---------|-------------|
| `google-chrome` | AUR | 131.0 | Google Chrome web browser |
| `ungoogled-chromium-bin` | AUR | 131.0 | Chromium without Google services/tracking |
| `helium` | AUR | 1.1.0 | Minimal Chromium-based browser |
| `onlyoffice-bin` | AUR | 8.2.2 | OnlyOffice Desktop - Microsoft Office compatible suite |

---

## LXD Container Packages (Stage 8) — Optional

### ai-lab Container

For AI/ML workloads with ROCm GPU acceleration:

| Package | Source | Version | Description |
|---------|--------|---------|-------------|
| `rocm-hip-sdk` | CachyOS | 7.2.0 | AMD ROCm HIP SDK for GPU computing |
| `python-pytorch-rocm` | CachyOS | 2.5.1 | PyTorch with ROCm backend for AMD GPUs |
| `python-numpy` | Extra | 2.2.1 | Numerical computing library for Python |
| `python-pip` | Extra | 24.3.1 | Python package installer |
| `git` | Extra | 2.47.1 | Version control |
| `base-devel` | Core | meta | Build tools |

### dev-lab Container

For general software development:

| Package | Source | Version | Description |
|---------|--------|---------|-------------|
| `rust` | Extra | 1.84.0 | Rust programming language and cargo |
| `go` | Extra | 1.23.4 | Go programming language |
| `nodejs` | Extra | 23.5.0 | JavaScript runtime built on V8 |
| `npm` | Extra | 10.9.2 | Node.js package manager |
| `python` | Extra | 3.13.1 | Python programming language |
| `python-pip` | Extra | 24.3.1 | Python package installer |
| `git` | Extra | 2.47.1 | Version control |
| `base-devel` | Core | meta | Build tools |

---

## Package Dependencies

Some packages have dependencies on others:

```
Graphics Stack
├── mesa ← requires llvm
├── vulkan-radeon ← requires mesa
└── lib32-* packages ← require 32-bit base libraries

Containers
├── lxd ← requires lxd group membership
└── GPU passthrough ← requires graphics stack

AUR Packages
└── all ← require yay (or another AUR helper)

ROCm/AI
├── rocm-hip-sdk ← requires kernel 6.18+ with AMDGPU
└── python-pytorch-rocm ← requires rocm-hip-sdk
```

---

## Notes

- **CachyOS repos** provide optimized builds with PGO/LTO for better performance
- **AUR packages** are built from source and require `yay` or similar helper
- **Version numbers** are approximate and will be updated by pacman to latest
- **Antigravity IDE**: Install separately as needed (not included by default)
