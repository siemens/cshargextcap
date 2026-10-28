#!/usr/bin/env bash
set -euo pipefail

export DEBIAN_FRONTEND=noninteractive

# Answer the "allow non-superusers to capture" prompt upfront; tests only read pipes and files.
echo "wireshark-common wireshark-common/install-setuid boolean false" | debconf-set-selections

apt-get update
apt-get install -y --no-install-recommends tshark
rm -rf /var/lib/apt/lists/*

tshark --version
