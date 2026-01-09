#!/bin/bash
# Strix Halo Post-Installer Bootstrap Script
# This script downloads and runs the Strix Halo installer.
#
# Usage: curl -fsSL https://bit.ly/strix-halo | sudo bash
#
# What this script does:
# 1. Downloads the latest strix-install binary from GitHub Releases
# 2. Makes it executable
# 3. Runs it with any arguments passed to this script
# 4. Cleans up the temporary file

set -euo pipefail

REPO="daveweinstein1/strix-halo-setup"
BINARY="strix-install"
TMP_FILE="/tmp/${BINARY}"

echo "Strix Halo Post-Installer Bootstrap"
echo "===================================="
echo ""
echo "Downloading latest release from GitHub..."

# Download the binary
curl -fsSL "https://github.com/${REPO}/releases/latest/download/${BINARY}" -o "${TMP_FILE}"

# Make executable
chmod +x "${TMP_FILE}"

echo "Starting installer..."
echo ""

# Run the installer with any passed arguments
"${TMP_FILE}" "$@"

# Cleanup
rm -f "${TMP_FILE}"
