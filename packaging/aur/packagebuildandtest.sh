#!/bin/bash
set -e

cp -r /pkg /tmp/pkg
cd /tmp/pkg

makepkg -g >> PKGBUILD
makepkg --noconfirm --syncdeps --rmdeps --clean --install

# Pass back the PKGBUILD with checksum(s) included so that we can publish it as
# a release artefact from the release pipeline.
mkdir -p /pkg/dist
cp PKGBUILD /pkg/dist/

# Ask tshark to tell us the extcap interfaces it knows of: this must list the
# packetflix extcap so we know we've installed the plugin properly.
tshark -D | grep packetflix \
    && echo "OK: tshark detects extcap plugin" \
    || (echo "FAIL: tshark doesn't detect the packetflix extcap"; exit 1)

# Check that the default URL scheme handler registration is in place.
xdg-mime query default x-scheme-handler/packetflix | grep "packetflix.desktop" \
    && echo "OK: packetflix URL scheme handler registered" \
    || (echo "packetflix URL scheme handler not detected"; exit 1)
