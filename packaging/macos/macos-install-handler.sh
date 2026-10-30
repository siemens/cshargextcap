#!/bin/bash
# Installs or upgrades the cshargextcap Wireshark plugin and the packetflix: URL
# handler on macOS for the current user, without needing admin rights, and
# cleans up older installations:
#
#   /bin/bash -c "$(curl -fsSL https://github.com/siemens/cshargextcap/releases/latest/download/macos-install-handler.sh)"
#
# "macos-install-handler.sh --from DIR" installs from a local build in DIR instead, that is,
# packetflix-handler.zip and cshargextcap_*_darwin_ARCH.tar.gz.
set -euo pipefail

repo=https://github.com/siemens/cshargextcap
extcapdir="${HOME}/.local/lib/wireshark/extcap"
handler="${HOME}/Applications/packetflix-handler.app"
handlerids="com.github.siemens.packetflix-handler com.apple.ScriptEditor.id.packetflix-handler"
lsregister=/System/Library/Frameworks/CoreServices.framework/Frameworks/LaunchServices.framework/Support/lsregister

[[ "$(uname -s)" == Darwin ]] || { echo "this installer is for macOS only" >&2; exit 1; }
arch=amd64
[[ "$(sysctl -n hw.optional.arm64 2>/dev/null)" != 1 ]] || arch=arm64

# lsapps JS prints the paths of the app(s) LaunchServices resolves using AppKit JS.
lsapps() {
    osascript -l JavaScript -e "ObjC.import('AppKit'); var r = \$.NSWorkspace.sharedWorkspace.$1;
        (r.count === undefined ? [r] : Array.from({length: r.count}, (_, i) => r.objectAtIndex(i)))
            .map(u => u.path.js).join('\n')" 2>/dev/null || true
}

# trash PATH removes PATH, falling back to Finder for app bundles that macOS
# protects from Terminal; it returns non-zero if neither works.
trash() {
    rm -rf "$1" 2>/dev/null && [[ ! -e "$1" ]] && return
    osascript -e 'on run argv' -e 'tell application "Finder" to delete (POSIX file (item 1 of argv) as alias)' \
        -e 'end run' "$1" >/dev/null 2>&1
}

work="$(mktemp -d)"
trap 'rm -rf "${work}"' EXIT

if [[ "${1:-}" == "--from" ]]; then
    from="$(cd "${2:?usage: macos-install-handler.sh [--from DIR]}" && pwd)"
    echo "installing from ${from}"
    tarball="$(ls -t "${from}"/cshargextcap_*_darwin_"${arch}".tar.gz 2>/dev/null | head -n 1 || true)"
    [[ -n "${tarball}" ]] || { echo "no cshargextcap_*_darwin_${arch}.tar.gz in ${from}" >&2; exit 1; }
    zip="${from}/packetflix-handler.zip"
else
    tag="$(curl -fsSLI -o /dev/null -w '%{url_effective}' "${repo}/releases/latest")"
    tag="${tag##*/}"
    echo "installing ${tag} (${arch})"
    tarball="${work}/cshargextcap.tar.gz"
    zip="${work}/packetflix-handler.zip"
    curl -fsSL -o "${tarball}" "${repo}/releases/download/${tag}/cshargextcap_${tag#v}_darwin_${arch}.tar.gz"
    curl -fsSL -o "${zip}" "${repo}/releases/download/${tag}/packetflix-handler.zip"
fi

tar -xzf "${tarball}" -C "${work}" cshargextcap
ditto -x -k "${zip}" "${work}"
if [[ ! -x "${work}/packetflix-handler.app/Contents/MacOS/packetflix-handler" ]]; then
    echo "this release still has the old AppleScript packetflix handler; please use a newer release" >&2
    exit 1
fi

# Wireshark plugin into the personal extcap folder, as macOS doesn't allow
# modifying Wireshark.app.
mkdir -p "${extcapdir}"
install -m 0755 "${work}/cshargextcap" "${extcapdir}/cshargextcap"
"${extcapdir}/cshargextcap" --extcap-interfaces 2>/dev/null | grep packetflix >/dev/null ||
    { echo "installed plugin ${extcapdir}/cshargextcap doesn't work" >&2; exit 1; }
echo "installed Wireshark plugin ${extcapdir}/cshargextcap"

wireshark="$(lsapps 'URLForApplicationWithBundleIdentifier("org.wireshark.Wireshark")')"
if [[ -z "${wireshark}" ]]; then
    echo "NOTE: Wireshark not found; please install it from https://www.wireshark.org/download.html"
elif [[ -e "${wireshark}/Contents/MacOS/extcap/cshargextcap" ]]; then
    if trash "${wireshark}/Contents/MacOS/extcap/cshargextcap"; then
        echo "removed old plugin from ${wireshark}"
    else
        echo "NOTE: please remove the outdated plugin from ${wireshark}: in Finder, right-click Wireshark,"
        echo "  Show Package Contents, open Contents/MacOS/extcap, and move cshargextcap to the Trash."
    fi
fi

# Remove old packetflix: handlers from the usual places, and unregister any
# other ones, such as in download folders, so that they don't take precedence.
for app in "/Applications/packetflix-handler.app" "${handler}"; do
    [[ -e "${app}" ]] || continue
    id="$(/usr/libexec/PlistBuddy -c 'Print :CFBundleIdentifier' "${app}/Contents/Info.plist" 2>/dev/null || true)"
    [[ " ${handlerids} " == *" ${id} "* ]] || continue
    "${lsregister}" -u "${app}" 2>/dev/null || true
    if trash "${app}"; then
        [[ "${app}" == "${handler}" ]] || echo "removed old handler ${app}"
    else
        echo "NOTE: please move the old ${app} to the Trash" >&2
    fi
done
lsapps 'URLsForApplicationsToOpenURL($.NSURL.URLWithString("packetflix:x"))' | while IFS= read -r app; do
    [[ -n "${app}" && -e "${app}" ]] || continue
    "${lsregister}" -u "${app}" 2>/dev/null || true
    echo "unregistered other handler ${app}; consider deleting it"
done

mkdir -p "$(dirname "${handler}")"
mv "${work}/packetflix-handler.app" "${handler}"
"${lsregister}" -f "${handler}"
for _ in $(seq 10); do
    urlapp="$(lsapps 'URLForApplicationToOpenURL($.NSURL.URLWithString("packetflix:x"))')"
    [[ "${urlapp}" != "${handler}" ]] || break
    sleep 1
done
if [[ "${urlapp}" != "${handler}" ]]; then
    echo "FAIL: packetflix: URLs open ${urlapp:-nothing} instead of ${handler}" >&2
    exit 1
fi
echo "installed packetflix: URL handler ${handler}"
echo "done; please restart Wireshark if it is running."
