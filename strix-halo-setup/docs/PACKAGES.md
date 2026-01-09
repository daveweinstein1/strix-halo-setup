# Package Inventory
*Strix Halo Post-Install - January 2026*

## Host System

### Graphics (Stage 02)
```
mesa lib32-mesa mesa-utils
vulkan-radeon lib32-vulkan-radeon vulkan-tools
linux-firmware
llvm lib32-llvm
```

### Essentials (Stage 03)
```
base-devel git wget curl vim neovim btop neofetch fastfetch
```

### Containers (Stage 04)
```
lxd
```

### User Apps (Stage 07)

**From official repos:**
```
firefox signal-desktop vlc yay
```

**From AUR (installed using `yay` helper):**
```
google-chrome
ungoogled-chromium-bin
helium
onlyoffice-bin
```

---

## LXD Containers (Stage 08)

### ai-lab
For AI/ML workloads with ROCm GPU acceleration:
```
rocm-hip-sdk python-pytorch-rocm python-numpy python-pip
git base-devel fastfetch vim
```

### dev-lab
For general development:
```
base-devel git rust go nodejs npm
python python-pip vim neovim fastfetch
```

---

## Notes

- **yay**: AUR helper tool that builds packages from the Arch User Repository
- **Antigravity IDE**: Install manually or in dev-lab container as needed (not included by default)
