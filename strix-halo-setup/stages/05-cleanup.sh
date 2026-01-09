#!/bin/bash

# ============================================================================
# STAGE 5: Cleanup & Optimization
# ============================================================================
source "$(dirname "$0")/../lib/common.sh"

stage_start "05-cleanup"
check_root

# ----------------------------------------------------------------------------
# Step 1: Orphan Cleanup
# ----------------------------------------------------------------------------
step 1 "Remove Orphaned Packages"

# Check for orphans
if [ -n "$(pacman -Qtdq)" ]; then
    ORPHANS=$(pacman -Qtdq)
    info "Found orphaned packages (dependencies no longer needed):"
    echo "$ORPHANS"
    
    if confirm_yes "Remove orphans?"; then
        run_cmd "Removing orphans" "pacman -Rns --noconfirm $ORPHANS"
    fi
else
    success "No orphans found."
fi

# ----------------------------------------------------------------------------
# Step 2: Cache Cleanup
# ----------------------------------------------------------------------------
step 2 "Clean Pacman Cache"

info "Frees disk space by removing old package versions."
if confirm_yes "Clean package cache?"; then
    run_cmd "Cleaning cache" "echo y | pacman -Scc"
fi

stage_complete
pause
