#!/bin/bash
set -e

cp -r /pkg /tmp/pkg
cd /tmp/pkg

makepkg -g >> PKGBUILD
makepkg --noconfirm --syncdeps --rmdeps --install --clean

# Pass back the PKGBUILD with checksum(s) included so someone might publish it
# as a release artefact.
mkdir -p /pkg/dist
cp PKGBUILD /pkg/dist/

# Ask tshark to tell us the extcap interfaces it knows of: this must list the
# packetflix extcap so we know we've installed the plugin properly.
tshark -D | grep packetflix
