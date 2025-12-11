#!/bin/bash

# ============================================================================
# STAGE 8: Workspace Provisioning (AI & Dev Containers)
# ============================================================================
source "$(dirname "$0")/../lib/common.sh"

stage_start "08-workspace-setup"
check_root

# ----------------------------------------------------------------------------
# Step 1: Verify LXD Status
# ----------------------------------------------------------------------------
step 1 "Check LXD Status"

if ! systemctl is-active --quiet lxd; then
    error "LXD is not running. Please run Stage 4 first."
    start_failed
    exit 1
fi
success "LXD is active"

# ----------------------------------------------------------------------------
# Step 2: AI Lab Provisioning
# ----------------------------------------------------------------------------
step 2 "Provision 'ai-lab' Container"

info "This container is for AI/ML workloads (ROCm, PyTorch)."

if confirm_yes "Create/Setup 'ai-lab'?"; then
    if lxc list | grep -q "ai-lab"; then
        info "'ai-lab' already exists."
    else
        run_cmd "Launching 'ai-lab' (Arch Linux)" lxc launch images:archlinux/current ai-lab
        spin "Waiting for network..." "while [ -z \"\$(lxc list ai-lab -c 4 --format csv)\" ]; do sleep 1; done"
    fi

    if confirm_yes "Install AI Stack (ROCm/PyTorch) inside 'ai-lab'?"; then
        info "Installing packages inside container..."
        # Update first
        run_cmd "Updating container" lxc exec ai-lab -- pacman -Syu --noconfirm
        
        # Install Strix Halo AI Stack
        # python-pytorch-rocm: Official Arch package for ROCm PyTorch
        # rocm-hip-sdk: Core SDK for development
        PACKAGES="rocm-hip-sdk python-pytorch-rocm python-numpy python-pip git base-devel fastfetch vim"
        
        run_cmd "Installing AI Packages (this may take time)" \
            lxc exec ai-lab -- pacman -S --needed --noconfirm $PACKAGES
            
        # Verify ROCm Version (Critical for Strix Halo)
        # Use pacman -Q (short) and strip epoch (cut -d: -f2)
        ROCM_VER=$(lxc exec ai-lab -- pacman -Q rocm-core 2>/dev/null | awk '{print $2}' | cut -d: -f2)
        
        if [[ "$ROCM_VER" == 7.1* ]]; then
            success "ROCm 7.1 detected ($ROCM_VER) - Strix Halo Ready"
        else
            warn "ROCm version is $ROCM_VER. Strix Halo requires 7.1+!"
            warn "Ensure 'cachyos-extra-v3' repo is active inside container."
        fi

        success "AI Stack installed"
    fi
fi

# ----------------------------------------------------------------------------
# Step 3: Dev Lab Provisioning
# ----------------------------------------------------------------------------
step 3 "Provision 'dev-lab' Container"

info "This container is for General Development (Rust, Go, Node, etc)."

if confirm "Create/Setup 'dev-lab'?"; then
    if lxc list | grep -q "dev-lab"; then
        info "'dev-lab' already exists."
    else
        run_cmd "Launching 'dev-lab' (Arch Linux)" lxc launch images:archlinux/current dev-lab
        spin "Waiting for network..." "while [ -z \"\$(lxc list dev-lab -c 4 --format csv)\" ]; do sleep 1; done"
    fi

    if confirm_yes "Install Dev Tools (Rust/Go/Node) inside 'dev-lab'?"; then
        info "Installing packages inside container..."
        run_cmd "Updating container" lxc exec dev-lab -- pacman -Syu --noconfirm
        
        PACKAGES="base-devel git rust go nodejs npm python python-pip vim neovim fastfetch"
        
        run_cmd "Installing Dev Packages" \
            lxc exec dev-lab -- pacman -S --needed --noconfirm $PACKAGES
            
        success "Dev Tools installed"
    fi
fi

# ----------------------------------------------------------------------------
# Step 4: Verification
# ----------------------------------------------------------------------------
step 4 "Verify Workspaces"

if lxc list | grep -q "ai-lab"; then
    info "Verifying AI Lab GPU Access..."
    # Check for KFD accessible
    if lxc exec ai-lab -- ls -l /dev/kfd &>/dev/null; then
        success "ai-lab can see /dev/kfd (GPU Compute OK)"
    else
        warn "ai-lab cannot see /dev/kfd. Check LXD profile in Stage 4."
    fi
fi

stage_complete
pause
