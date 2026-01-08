#!/bin/bash

# ============================================================================
# STAGE 2: Graphics Stack Setup (Mesa/Vulkan/Firmware)
# ============================================================================
source "$(dirname "$0")/../lib/common.sh"

stage_start "02-graphics-setup"
check_root

log "Target: AMD Strix Halo (gfx1151)"
log "Requirements (Jan 2026): Mesa 25.0+ (25.4+ rec), LLVM 19+ (21+ rec)"

# ----------------------------------------------------------------------------
# Step 1: Sync Database
# ----------------------------------------------------------------------------
step 1 "Sync Package Databases"

if confirm_yes "Sync pacman databases (pacman -Sy)?"; then
    run_cmd "Syncing databases" pacman -Sy --noconfirm
fi

# ----------------------------------------------------------------------------
# Step 2: Mesa Installation
# ----------------------------------------------------------------------------
step 2 "Install/Upgrade Mesa (OpenGL/Vulkan)"

PACKAGES="mesa lib32-mesa mesa-utils"
info "Packages: $PACKAGES"

if confirm_yes "Install Mesa packages?"; then
    run_cmd "Installing Mesa" pacman -S --needed --noconfirm $PACKAGES
    
    # Version Check
    CURRENT_VER=$(get_version mesa)
    log "Installed Mesa: $CURRENT_VER"
    
    # Strip Epoch (e.g., 1:25.3.1 -> 25.3.1) for comparison
    CLEAN_VER=$(echo "$CURRENT_VER" | cut -d: -f2)
    MAJ=$(echo "$CLEAN_VER" | cut -d. -f1)
    MIN=$(echo "$CLEAN_VER" | cut -d. -f2)
    
    if [ "$MAJ" -ge 25 ]; then
        success "Version $CURRENT_VER is EXCELLENT (25.x series)"
    elif [ "$MAJ" -eq 24 ] && [ "$MIN" -ge 1 ]; then
        success "Version $CURRENT_VER is ADEQUATE (24.1+)"
    else
        warn "Version $CURRENT_VER is potentially too old for gfx1150!"
        warn "Consider enabling CachyOS testing repos or manual update."
        if ! confirm "Continue despite version warning?"; then
            stage_failed
            exit 1
        fi
    fi
fi

# ----------------------------------------------------------------------------
# Step 3: Vulkan Setup
# ----------------------------------------------------------------------------
step 3 "Install Vulkan Drivers (RADV)"

PACKAGES="vulkan-radeon lib32-vulkan-radeon vulkan-tools"
info "Packages: $PACKAGES"

if confirm_yes "Install Vulkan packages?"; then
    run_cmd "Installing Vulkan" pacman -S --needed --noconfirm $PACKAGES
fi

# ----------------------------------------------------------------------------
# Step 4: Firmware
# ----------------------------------------------------------------------------
step 4 "Install AMD GPU Firmware"

PACKAGES="linux-firmware"
info "Packages: $PACKAGES"

if confirm_yes "Install/Upgrade linux-firmware?"; then
    run_cmd "Installing Firmware" pacman -S --needed --noconfirm $PACKAGES
fi

# ----------------------------------------------------------------------------
# Step 5: LLVM
# ----------------------------------------------------------------------------
step 5 "Install LLVM (Shader Compiler)"

PACKAGES="llvm lib32-llvm"
info "Packages: $PACKAGES"

if confirm_yes "Install LLVM packages?"; then
    run_cmd "Installing LLVM" pacman -S --needed --noconfirm $PACKAGES
    
    # Check version
    LLVM_VER=$(llvm-config --version 2>/dev/null || echo "unknown")
    log "Installed LLVM: $LLVM_VER"
    
    MAJ=$(echo "$LLVM_VER" | cut -d. -f1)
    if [ "$MAJ" -ge 21 ]; then
        success "LLVM $LLVM_VER is Current Stable (Dec 2025)"
    elif [ "$MAJ" -ge 17 ]; then
        success "LLVM $LLVM_VER meets minimum requirements"
    else
        warn "LLVM $LLVM_VER is old! Needs 17+ for gfx1150 backend."
    fi
fi

# ----------------------------------------------------------------------------
# Step 6: GPU Detection
# ----------------------------------------------------------------------------
step 6 "Verify GPU Detection"

if confirm "Check lspci for GPU?"; then
    run_cmd "Checking VGA devices" "lspci | grep -i vga"
    run_cmd "Checking AMD devices" "lspci | grep -i amd"
fi

if [ -n "$DISPLAY" ]; then
    if confirm "Run glxinfo test?"; then
        run_cmd "Checking OpenGL Renderer" "glxinfo | grep -E 'OpenGL vendor|OpenGL renderer|OpenGL version'"
    fi
fi

# ----------------------------------------------------------------------------
# Step 7: Environment Config
# ----------------------------------------------------------------------------
step 7 "Configure Environment Optimizations"

ENV_FILE="/etc/environment.d/10-amd-gpu.conf"

if confirm "Create $ENV_FILE?"; then
    mkdir -p /etc/environment.d
    cat > /tmp/amd-env.conf << 'EOF'
# AMD GPU optimizations for Strix Halo (gfx1150)
AMD_VULKAN_ICD=RADV
RADV_PERFTEST=gpl,nggc,rt
RADV_DEBUG=noatocdithering
EOF
    
    run_cmd "Installing environment config" "mv /tmp/amd-env.conf $ENV_FILE"
    log "Config created at $ENV_FILE"
fi

stage_complete
pause
