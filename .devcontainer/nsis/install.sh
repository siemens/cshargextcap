#!/usr/bin/env bash
set -euo pipefail

export DEBIAN_FRONTEND=noninteractive
apt-get update
apt-get install -y --no-install-recommends nsis
rm -rf /var/lib/apt/lists/*

makensis -VERSION
