#!/bin/bash

# ============================================================================
# STAGE 7: User Applications (Dev, Office, Internet)
# ============================================================================
source "$(dirname "$0")/../lib/common.sh"

stage_start "07-user-apps"
check_root

REAL_USER=$(logname || echo $SUDO_USER)
if [ -z "$REAL_USER" ]; then
    error "Cannot detect non-root user (SUDO_USER) required for AUR operations."
    exit 1
fi
info "AUR operations will run as user: $REAL_USER"

# ----------------------------------------------------------------------------
# Step 0: AUR Helper (Required for some apps)
# ----------------------------------------------------------------------------
step 0 "Setup AUR Helper"

if ! command -v yay &>/dev/null; then
    info "Installing 'yay' for AUR support..."
    run_cmd "Installing yay" pacman -S --needed --noconfirm yay
fi

# ----------------------------------------------------------------------------
# Step 1: Internet & Browsers
# ----------------------------------------------------------------------------
step 1 "Install Web Browsers"
info "Note: You can install multiple browsers independently."

if confirm_yes "Install Firefox?"; then
    run_cmd "Installing Firefox" pacman -S --needed --noconfirm firefox
fi

if confirm "Install Google Chrome?"; then
    run_cmd "Installing Chrome" sudo -u "$REAL_USER" yay -S --needed --noconfirm google-chrome
fi

if confirm "Install Ungoogled Chromium?"; then
    run_cmd "Installing Ungoogled Chromium" sudo -u "$REAL_USER" yay -S --needed --noconfirm ungoogled-chromium-bin
fi

if confirm "Install Helium Browser?"; then
    # Helium often carbon-based or minimal floating browser.
    # Assuming 'helium-browser' or 'helium'. 
    # If not found, manual build required.
    run_cmd "Installing Helium" sudo -u "$REAL_USER" yay -S --needed --noconfirm helium
fi

# ----------------------------------------------------------------------------
# Step 3: Office Suite
# ----------------------------------------------------------------------------
step 3 "Install Office Suite"

if confirm "Install OnlyOffice?"; then
    run_cmd "Installing OnlyOffice" sudo -u "$REAL_USER" yay -S --needed --noconfirm onlyoffice-bin
fi

# ----------------------------------------------------------------------------
# Step 4: Communication
# ----------------------------------------------------------------------------
step 4 "Install Communication Tools"

if confirm "Install Signal Desktop?"; then
    run_cmd "Installing Signal" pacman -S --needed --noconfirm signal-desktop
fi

# ----------------------------------------------------------------------------
# Step 5: Media
# ----------------------------------------------------------------------------
step 5 "Install Media Tools"

if confirm "Install VLC Media Player?"; then
    run_cmd "Installing VLC" pacman -S --needed --noconfirm vlc
fi

stage_complete
pause
