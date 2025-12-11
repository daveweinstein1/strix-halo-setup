#!/bin/bash

# ============================================================================
# STAGE 4: LXD/LXC Container Setup
# ============================================================================
source "$(dirname "$0")/../lib/common.sh"

stage_start "04-lxd-setup"
check_root

# ----------------------------------------------------------------------------
# Step 1: Install LXD
# ----------------------------------------------------------------------------
step 1 "Install LXD Package"

# Only install 'lxd'. 
# Dependencies (iptables-nft, dnsmasq, etc.) are handled by pacman automatically.
# We avoid explicit 'iptables' to prevent conflict with 'iptables-nft'.
DEPS="lxd"
info "Packages: $DEPS"

if confirm_yes "Install LXD?"; then
    run_cmd "Installing LXD" pacman -S --needed --noconfirm $DEPS
fi

# ----------------------------------------------------------------------------
# Step 2: Service Setup
# ----------------------------------------------------------------------------
step 2 "Enable LXD Service"

if confirm_yes "Enable and Start lxd.service?"; then
    run_cmd "Enabling lxd.socket" systemctl enable --now lxd.socket
    
    sleep 2
    if systemctl is-active --quiet lxd.socket; then
        success "LXD socket is active"
    else
        error "LXD failed to start"
        run_cmd "Checking Status" systemctl status lxd
        stage_failed
        exit 1
    fi
fi

# ----------------------------------------------------------------------------
# Step 3: User Group
# ----------------------------------------------------------------------------
step 3 "Add User to 'lxd' Group"

CURRENT_USER=$(logname || echo $SUDO_USER)
if [ -z "$CURRENT_USER" ]; then
    warn "Could not detect SUDO_USER, assuming 'root' (skipping group add)"
else
    info "Target User: $CURRENT_USER"
    if confirm_yes "Add $CURRENT_USER to lxd group?"; then
        run_cmd "Modifying group" usermod -aG lxd "$CURRENT_USER"
        success "User added. Log out/in required."
    fi
fi

# ----------------------------------------------------------------------------
# Step 4: Initialize LXD & Networking
# ----------------------------------------------------------------------------
step 4 "Initialize LXD Network"

info "Running 'lxd init --auto' (Creates lxdbr0)"

if confirm_yes "Initialize LXD now?"; then
    if lxd init --auto &>> "$STAGE_LOG"; then
        success "LXD initialized"
        
        # FIX for systemd-resolved:
        # Instead of disabling resolved, we configure the lxdbr0 link
        if systemctl is-active --quiet systemd-resolved; then
            info "Configuring systemd-resolved for lxdbr0..."
            
            # Wait for bridge to appear
            sleep 2
            
            # Get Bridge IP
            if BRIDGE_IP=$(lxc network get lxdbr0 ipv4.address 2>/dev/null | cut -d/ -f1); then
                run_cmd "Setting DNS for lxdbr0" resolvectl dns lxdbr0 "$BRIDGE_IP"
                run_cmd "Setting Domain for lxdbr0" resolvectl domain lxdbr0 '~lxd'
                success "DNS resolution configured for .lxd domain"
            else
                warn "Could not determine lxdbr0 IP - skipping DNS config"
            fi
        fi
    else
        warn "Auto init failed. You may need to run 'lxd init' manually."
    fi
fi

# ----------------------------------------------------------------------------
# Step 6: AI/GPU Passthrough Configuration
# ----------------------------------------------------------------------------
step 6 "Configure Default Profile for AI (GPU Access)"

if systemctl is-active --quiet lxd; then
    info "Configuring 'default' profile to pass AMD GPU to all containers..."
    
    # Enable nesting (crucial for Docker-inside-LXD, common in ML workflows)
    if confirm_yes "Allow nested containers (Docker support)?"; then
        run_cmd "Enabling nesting" lxc profile set default security.nesting=true
    fi
    
    # Add GPU device
    # Strix Halo has a unified GPU, so we pass 'gpu'.
    if lxc profile device list default | grep -q "gpu0"; then
        success "GPU device already in default profile"
    else
        if confirm_yes "Pass GPU to all containers by default?"; then
            # This passes all available GPUs. 
            # gputype=physical generally gives best compute access for ROCm
            run_cmd "Adding GPU to profile" lxc profile device add default gpu0 gpu gputype=physical
        fi
    fi
    
    # Check for KFD (Compute) access
    # Recent LXD handles this with 'gpu' device, but verifying /dev/kfd permissions is key.
    # Usually requires the container user to be in 'render' group, which we can't fix from here easily
    # without mapped subgids, but passing the device is step 1.
else
    warn "LXD not running, skipping profile configuration."
fi

stage_complete
echo ""
warn "⚠️  You must LOG OUT and back in for 'lxd' group permissions to work!"
pause
