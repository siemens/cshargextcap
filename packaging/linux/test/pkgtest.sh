#!/bin/sh
# Installs, checks, and removes the cshargextcap package from /dist inside a
# distro container; supports apk, apt-get, and dnf based distributions.
set -eu

ok()   { printf '\033[1;32mOK:\033[0m %s\n' "$*"; }
fail() { printf '\033[1;31mFAIL:\033[0m %s\n' "$*" >&2; exit 1; }

case "$(uname -m)" in
    x86_64)  arch=amd64 ;;
    aarch64) arch=arm64 ;;
    *)       fail "unsupported architecture $(uname -m)" ;;
esac

pkgfile() {
    set -- /dist/cshargextcap_*_linux_"${arch}.$1"
    [ $# -eq 1 ] && [ -f "$1" ] || fail "expected exactly one package, found: $*"
    echo "$1"
}

if command -v apk >/dev/null; then
    pkg=$(pkgfile apk)
    install_harness() { apk add --no-cache xdg-utils; }
    install_pkg()     { apk add --no-cache --allow-untrusted "${pkg}"; }
    install_tshark()  { apk add --no-cache tshark; }
    remove_pkg()      { apk del cshargextcap; }
elif command -v apt-get >/dev/null; then
    export DEBIAN_FRONTEND=noninteractive
    pkg=$(pkgfile deb)
    apt_install()     { apt-get install -q -y --no-install-recommends --no-install-suggests "$@"; }
    install_harness() { apt-get update -q && apt_install xdg-utils; }
    install_pkg()     { apt_install "${pkg}"; }
    install_tshark()  { apt_install tshark; }
    remove_pkg()      { apt-get purge -q -y cshargextcap; }
elif command -v dnf >/dev/null; then
    pkg=$(pkgfile rpm)
    install_harness() { dnf install -q -y xdg-utils; }
    install_pkg()     { dnf install -q -y "${pkg}"; }
    install_tshark()  { :; } # pulled in via the wireshark-cli dependency
    remove_pkg()      { dnf remove -q -y cshargextcap; }
else
    fail "unsupported package manager"
fi

echo "== installing ${pkg}"
install_harness
install_pkg
install_tshark

extcaps() {
    ls /usr/libexec/wireshark/extcap/cshargextcap \
       /usr/lib/wireshark/extcap/cshargextcap \
       /usr/lib64/wireshark/extcap/cshargextcap \
       /usr/lib/*-linux-gnu/wireshark/extcap/cshargextcap 2>/dev/null || true
}

[ -n "$(extcaps)" ] || fail "extcap binary not installed"
for bin in $(extcaps); do
    "${bin}" --extcap-interfaces | grep -q '^extcap {version=' \
        && ok "extcap binary runs from ${bin}" \
        || fail "extcap binary at ${bin} doesn't report its interfaces"
done

tshark -D 2>/dev/null | grep -q packetflix \
    && ok "tshark detects extcap plugin" \
    || fail "tshark doesn't detect the packetflix extcap"

xdg-mime query default x-scheme-handler/packetflix | grep -q "packetflix.desktop" \
    && ok "packetflix URL scheme handler registered" \
    || fail "packetflix URL scheme handler not detected"

echo "== removing cshargextcap"
remove_pkg

[ -z "$(extcaps)" ] \
    && ok "extcap binary removed" \
    || fail "extcap binary still present after removal: $(extcaps)"

if xdg-mime query default x-scheme-handler/packetflix | grep -q "packetflix.desktop"; then
    fail "packetflix URL scheme handler still registered after removal"
fi
ok "packetflix URL scheme handler unregistered"
