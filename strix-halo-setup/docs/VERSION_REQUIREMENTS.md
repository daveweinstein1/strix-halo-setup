# Version Requirements (Strix Halo)
*Updated: January 9, 2026*

| Component | Minimum | Recommended | Reason |
|-----------|---------|-------------|--------|
| **Kernel** | 6.14 | **6.17+** | 6.14 = AMDXDNA NPU driver, 6.17+ = stable performance |
| **Mesa** | 25.0 | **25.4+** | GFX1151 Graphics Support |
| **ROCm** | 7.1 | **7.2+** | **CRITICAL**: Full Strix Halo compute support (CES 2026) |
| **LLVM** | 19.x | **21.x** | Shader Compiler |

> [!IMPORTANT]
> Use CachyOS `v3` or `v4` repositories to ensure you get these bleeding-edge versions. Standard Arch repos may trail behind.

> [!NOTE]
> **Strix Halo = gfx1151** (Ryzen AI Max+). Strix Point = gfx1150 (different chip).
> Linux 6.18 has a known bug affecting gfx1150, but gfx1151 appears unaffected.
