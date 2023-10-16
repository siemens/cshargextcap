#!/bin/bash
set -e

RED="\033[1;31m"
GREEN="\033[1;32m"
NOCOLOR="\033[0m"

arch=$(uname -m)
case ${arch} in
    x86_64)   arch="amd64";;
    aarch64) arch="arm64";;
    *)        echo -e "${RED}FAIL:${NOCOLOR} unsupported architecture ${arch}"; exit 1;;
esac

# test harness
dnf install -y xdg-utils

# the real deal...
dnf install -y /dist/cshargextcap_*_${arch}.rpm

# Ask tshark to tell us the extcap interfaces it knows of: this must list the
# packetflix extcap so we know we've installed the plugin properly.
tshark -D | grep packetflix \
    && echo -e "${GREEN}OK:${NOCOLOR} tshark detects extcap plugin" \
    || (echo -e "${RED}FAIL:${NOCOLOR} tshark doesn't detect the packetflix extcap"; exit 1)

# Check that the default URL scheme handler registration is in place.
xdg-mime query default x-scheme-handler/packetflix | grep "packetflix.desktop" \
    && echo -e "${GREEN}OK:${NOCOLOR} packetflix URL scheme handler registered" \
    || (echo -e "${RED}FAIL:${NOCOLOR} packetflix URL scheme handler not detected"; exit 1)
