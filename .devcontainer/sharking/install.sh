#!/usr/bin/env bash
set -euo pipefail

export DEBIAN_FRONTEND=noninteractive

# Answer the "allow non-superusers to capture" prompt upfront; tests only read pipes and files.
echo "wireshark-common wireshark-common/install-setuid boolean false" | debconf-set-selections

apt-get update
# Wireshark's GUI runs on Xvfb and gets inspected via AT-SPI on a private D-Bus session.
apt-get install -y --no-install-recommends tshark wireshark xvfb dbus-daemon at-spi2-core
# Firefox opens Edgeshark's packetflix: capture links via the packaged URL handler.
apt-get install -y --no-install-recommends desktop-file-utils curl ca-certificates
# Ubuntu's firefox package only installs the snap, which doesn't work in containers.
if ! firefox --version >/dev/null 2>&1; then
    install -d -m 0755 /etc/apt/keyrings
    curl -fsSL https://packages.mozilla.org/apt/repo-signing-key.gpg -o /etc/apt/keyrings/packages.mozilla.org.asc
    echo "deb [signed-by=/etc/apt/keyrings/packages.mozilla.org.asc] https://packages.mozilla.org/apt mozilla main" \
        > /etc/apt/sources.list.d/mozilla.list
    printf 'Package: *\nPin: origin packages.mozilla.org\nPin-Priority: 1000\n' > /etc/apt/preferences.d/mozilla
    apt-get update
    apt-get install -y --no-install-recommends firefox
fi
rm -rf /var/lib/apt/lists/*

tshark --version
