#!/bin/bash

# ============================================================================
# STAGE 1: Kernel Configuration (E610 Fix + Optimizations)
# ============================================================================
source "$(dirname "$0")/../lib/common.sh"

stage_start "01-kernel-config"
check_root

log "Requirements: Kernel 6.14+ (6.17+ rec for Strix Halo NPU/ISP)"

# ----------------------------------------------------------------------------
# Step 1: Backup GRUB
# ----------------------------------------------------------------------------
step 1 "Backup GRUB configuration"

if confirm "Create backup of /etc/default/grub?"; then
    run_cmd "Backing up GRUB config" cp /etc/default/grub "/etc/default/grub.backup.$(date +%Y%m%d-%H%M%S)"
else
    warn "Skipped GRUB backup"
fi

# ----------------------------------------------------------------------------
# Step 2: Check Kernel Version
# ----------------------------------------------------------------------------
step 2 "Check Kernel Version"

KERNEL_CURRENT=$(uname -r)
log "Current kernel: $KERNEL_CURRENT"

# Extract major.minor
K_MAJ=$(echo "$KERNEL_CURRENT" | cut -d. -f1)
K_MIN=$(echo "$KERNEL_CURRENT" | cut -d. -f2)

# Check for 6.14+ (AMDXDNA NPU driver), warn if < 6.17
if [ "$K_MAJ" -gt 6 ] || { [ "$K_MAJ" -eq 6 ] && [ "$K_MIN" -ge 17 ]; }; then
    success "Kernel version $KERNEL_CURRENT meets recommended requirement (6.17+)"
elif [ "$K_MAJ" -eq 6 ] && [ "$K_MIN" -ge 14 ]; then
    warn "Kernel version $KERNEL_CURRENT meets minimum (6.14+) but 6.17+ recommended"
else
    warn "Kernel version $KERNEL_CURRENT is older than required (6.14+)"
    warn "Strix Halo (gfx1151) NPU support requires kernel 6.14+ (AMDXDNA driver)."
    if ! confirm "Continue anyway (Not Recommended)?"; then
        stage_failed
        exit 1
    fi
fi

# ----------------------------------------------------------------------------
# Step 3: Disable E610 (Intel Ethernet) Driver
# ----------------------------------------------------------------------------
step 3 "Configure Kernel Parameters (E610 Fix)"

info "Change: Add 'modprobe.blacklist=ice' to GRUB"
info "Reason: Fixes instability with Intel E610 (Critical for Beelink GTR9 Pro)"

if confirm "Modify GRUB to disable E610/ice driver?"; then
    if grep -q "modprobe.blacklist=ice" /etc/default/grub; then
        success "E610/ice already blacklisted in GRUB"
    else
        # Use simple sed to append to the end of the line inside the quotes
        # This is safer than regex replacing the whole line
        run_cmd "Adding modprobe.blacklist=ice" \
            "sed -i 's/GRUB_CMDLINE_LINUX_DEFAULT=\"/GRUB_CMDLINE_LINUX_DEFAULT=\"modprobe.blacklist=ice /' /etc/default/grub"
        
        # Verify
        if grep -q "modprobe.blacklist=ice" /etc/default/grub; then
            success "Parameter added successfully"
        else
            error "Failed to modify /etc/default/grub"
            start_failed
            exit 1
        fi
    fi
else
    warn "Skipped E610 disable"
fi

# ----------------------------------------------------------------------------
# Step 4: AMD Optimizations
# ----------------------------------------------------------------------------
step 4 "Add AMD Strix Halo Optimizations"

PARAMS="amd_pstate=active amdgpu.ppfeaturemask=0xffffffff amd_iommu=on iommu=pt"
info "Parameters: $PARAMS"

if confirm "Add AMD optimization parameters?"; then
    for param in $PARAMS; do
        if grep -q "$param" /etc/default/grub; then
            info "$param already present"
        else
            run_cmd "Adding $param" \
                "sed -i 's/GRUB_CMDLINE_LINUX_DEFAULT=\"/GRUB_CMDLINE_LINUX_DEFAULT=\"$param /' /etc/default/grub"
        fi
    done
    success "AMD optimizations configured"
else
    warn "Skipped AMD optimizations"
fi

# ----------------------------------------------------------------------------
# Step 5: Update GRUB
# ----------------------------------------------------------------------------
step 5 "Apply GRUB Changes"

if confirm "Run grub-mkconfig to apply changes?"; then
    run_cmd "Updating GRUB" grub-mkconfig -o /boot/grub/grub.cfg
else
    error "GRUB not updated - changes will not take effect!"
fi

# ----------------------------------------------------------------------------
# Step 6: Create Modprobe Blacklist File
# ----------------------------------------------------------------------------
step 6 "Create Modprobe Blacklist (Safety Net)"

BLOCK_FILE="/etc/modprobe.d/blacklist-e610.conf"

if confirm "Create $BLOCK_FILE?"; then
    run_cmd "Creating blacklist file" "echo 'blacklist ice' > $BLOCK_FILE"
else
    warn "Skipped blacklist file creation"
fi

# ----------------------------------------------------------------------------
# Step 7: Update Initramfs
# ----------------------------------------------------------------------------
step 7 "Update Initramfs"

if confirm "Run mkinitcpio to update initramfs?"; then
    run_cmd "Updating initramfs" mkinitcpio -P
else
    warn "Initramfs not updated"
fi

stage_complete
echo ""
warn "⚠️  REBOOT REQUIRED for kernel changes to take effect!"
pause

# ----------------------------------------------------------------------------
# Step 8: Beelink GTR9 Power/TDP Control
# ----------------------------------------------------------------------------
step 8 "Configure TDP (Beelink GTR9 Specific)"

info "This step creates a systemd service to enforce TDP limits on boot."
info "Recommended for managing heat/noise on Mini PCs."

if confirm "Configure TDP Limits?"; then
    if ! command -v ryzenadj &>/dev/null; then
        run_cmd "Installing ryzenadj" pacman -S --needed --noconfirm ryzenadj
    fi

    echo ""
    echo "Select TDP Target:"
    echo "1) 55W (Silent/Cool)"
    echo "2) 80W (Balanced)"
    echo "3) 120W+ (Max Performance - Default)"
    read -p "Selection [1-3]: " tdp_choice

    case $tdp_choice in
        1) LIMIT=55000 ;;
        2) LIMIT=80000 ;;
        *) LIMIT=0 ;;
    esac

    if [ "$LIMIT" -gt 0 ]; then
        SERVICE_FILE="/etc/systemd/system/ryzenadj.service"
        cat > /tmp/ryzenadj.service << EOF
[Unit]
Description=RyzenAdj TDP Control
After=network.target

[Service]
Type=oneshot
ExecStart=/usr/bin/ryzenadj --stapm-limit=$LIMIT --fast-limit=$LIMIT --slow-limit=$LIMIT
RemainAfterExit=yes

[Install]
WantedBy=multi-user.target
EOF
        run_cmd "Installing Service" mv /tmp/ryzenadj.service $SERVICE_FILE
        run_cmd "Enabling Service" systemctl enable ryzenadj
        success "TDP set to $(($LIMIT/1000))W (Persistent)"
    else
        info "Keeping default TDP settings."
    fi
fi

stage_complete
echo ""
warn "⚠️  REBOOT REQUIRED for changes to take effect!"
pause
