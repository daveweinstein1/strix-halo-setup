# CachyOS Strix Halo Post-Installer

**Automated post-install setup for AMD Strix Halo (gfx1151) workstations on CachyOS.**

## Quick Install

```bash
curl -fsSL https://github.com/daveweinstein1/strix-halo-setup/releases/latest/download/strix-install -o /tmp/strix-install && chmod +x /tmp/strix-install && sudo /tmp/strix-install
```

## What It Does

| Stage | Purpose |
|-------|---------|
| Kernel Config | IOMMU, device quirks (Beelink E610 fix) |
| Graphics Setup | Mesa 25.3+, LLVM 21.x, Vulkan |
| System Update | Mirrors, packages, essentials |
| LXD Setup | Containers with GPU passthrough |
| Cleanup | Orphan removal, cache cleanup |
| Validation | Verify kernel, GPU, LXD |
| Desktop Apps | Browsers, Office (optional) |
| Workspaces | `ai-lab`, `dev-lab` containers (optional) |

## Version Requirements (January 2026)

| Component | Required |
|-----------|----------|
| Kernel | **6.18+** |
| Mesa | **25.3+** |
| ROCm | **7.2+** |
| LLVM | **21.x** |

## Supported Hardware

- **Framework Desktop** - Full support
- **Beelink GTR9 Pro** - E610 Ethernet fix applied automatically  
- **Minisforum MS-S1 Max** - Advisory for Ethernet/USB4 quirks
- **Other Strix Halo** - Generic mode

## License

Proprietary. No use without explicit permission from the author.

---
**Author**: Dave Weinstein | **Updated**: January 2026
