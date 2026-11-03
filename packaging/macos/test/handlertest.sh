#!/usr/bin/env bash
# End-to-end test of packetflix-handler.app on macOS: opens packetflix: URLs via
# LaunchServices and checks the arguments the handler starts Wireshark with.
# A fake Wireshark gets registered under Wireshark's bundle ID outside of
# /Applications/; LaunchServices may still pick an installed real Wireshark.
set -euo pipefail

here="$(cd "$(dirname "$0")" && pwd)"
zip="${1:-${here}/../../../dist/packetflix-handler.zip}"
if [[ ! -f "${zip}" ]]; then
    echo "${zip} not found; run 'make macos-packetflix-handler' first, or pass the zip's path" >&2
    exit 1
fi
zip="$(cd "$(dirname "${zip}")" && pwd)/$(basename "${zip}")"
lsregister=/System/Library/Frameworks/CoreServices.framework/Frameworks/LaunchServices.framework/Support/lsregister

# LaunchServices doesn't seem to bind URL schemes for apps in the system temp dir.
mkdir -p "${here}/../../../dist"
work="$(cd "$(mktemp -d "${here}/../../../dist/handlertest.XXXXXX")" && pwd -P)"
handler="${work}/packetflix-handler.app"
fake="${work}/Somewhere Else/Wireshark.app"
others="${work}/other-handlers"

# lsapps JS prints the paths of the app(s) LaunchServices resolves using AppKit JS.
lsapps() {
    osascript -l JavaScript -e "ObjC.import('AppKit'); var r = \$.NSWorkspace.sharedWorkspace.$1;
        (r.count === undefined ? [r] : Array.from({length: r.count}, (_, i) => r.objectAtIndex(i)))
            .map(u => u.path.js).join('\n')" 2>/dev/null || true
}

cleanup() {
    "${lsregister}" -u "${handler}" 2>/dev/null || true
    "${lsregister}" -u "${fake}" 2>/dev/null || true
    while IFS= read -r app; do
        "${lsregister}" -f "${app}" 2>/dev/null || true
    done <"${others}"
    rm -rf "${work}"
}
# other (installed) handlers would take precedence, so set them aside until done.
lsapps 'URLsForApplicationsToOpenURL($.NSURL.URLWithString("packetflix:x"))' >"${others}"
trap cleanup EXIT
while IFS= read -r app; do
    [[ -z "${app}" ]] || "${lsregister}" -u "${app}" 2>/dev/null || true
done <"${others}"
ditto -x -k "${zip}" "${work}"
codesign --verify --verbose=2 "${handler}" || echo "(informational) bundle not fully signed"

mkdir -p "${fake}/Contents/MacOS"
cat >"${fake}/Contents/Info.plist" <<'EOF'
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>CFBundleExecutable</key><string>Wireshark</string>
	<key>CFBundleIdentifier</key><string>org.wireshark.Wireshark</string>
	<key>CFBundlePackageType</key><string>APPL</string>
</dict>
</plist>
EOF
# stays running so that its command line can be checked, like a real Wireshark.
cat >"${fake}/Contents/MacOS/Wireshark" <<'EOF'
#!/bin/sh
sleep 30
EOF
chmod +x "${fake}/Contents/MacOS/Wireshark"
"${lsregister}" -f "${fake}"
"${lsregister}" -f "${handler}"

for _ in $(seq 10); do
    urlapp="$(lsapps 'URLForApplicationToOpenURL($.NSURL.URLWithString("packetflix:x"))')"
    [[ "${urlapp}" != *"${handler#/private}" ]] || break
    sleep 1
done
echo "packetflix: URLs open: ${urlapp:-(none)}"
echo "org.wireshark.Wireshark: $(lsapps 'URLForApplicationWithBundleIdentifier("org.wireshark.Wireshark")')"
if [[ "${urlapp}" != *"${handler#/private}" ]]; then
    echo "FAIL: packetflix: URLs go to ${urlapp:-(none)} instead of ${handler}" >&2
    [[ -z "${urlapp}" ]] || echo "  unregister it: ${lsregister} -u '${urlapp}'" >&2
    echo "LaunchServices entries mentioning packetflix:" >&2
    "${lsregister}" -dump | grep -i -E '^(path|claimed schemes|bindings):.*packetflix' | sort -u >&2 ||
        echo "  (none)" >&2
    exit 1
fi
handlerlog=~/Library/Logs/packetflix-handler.log
rm -f "${handlerlog}"

# expect MARKER URL ARGS... opens the URL, waits for the Wireshark started with
# the URL (identified by the unique MARKER in it), checks its args and quits it.
expect() {
    local marker="$1" url="$2" expected pid cmdline
    shift 2
    expected="$*"
    open "${url}"
    for _ in $(seq 150); do
        if pid="$(pgrep -n -f -- "MacOS/Wireshark .*${marker}")"; then
            cmdline="$(ps -ww -o args= -p "${pid}")"
            kill "${pid}" 2>/dev/null || true
            if [[ "${cmdline}" != *" ${expected}" ]]; then
                printf 'FAIL: unexpected Wireshark args\n  got:      %s\n  expected: ... %s\n' \
                    "${cmdline}" "${expected}" >&2
                exit 1
            fi
            echo "ok: started ${cmdline%% -k *}"
            return
        fi
        sleep 0.2
    done
    echo "FAIL: no Wireshark started for ${url}; Wireshark/handler processes:" >&2
    pgrep -l -f -i 'wireshark|packetflix' >&2 || echo "  (none)" >&2
    echo "handler log ${handlerlog}:" >&2
    cat "${handlerlog}" >&2 2>/dev/null || echo "  (none, so the handler did not start)" >&2
    echo "recent handler crash reports:" >&2
    ls -t ~/Library/Logs/DiagnosticReports 2>/dev/null | grep -i packetflix | head -3 >&2 ||
        echo "  (none)" >&2
    if grep -q 'running open' "${handlerlog}" 2>/dev/null && ! grep -q 'error:' "${handlerlog}"; then
        echo "HINT: the handler started Wireshark, but Wireshark quit at once. Wireshark rejects" >&2
        echo "  -o extcap.packetflix.url:... unless cshargextcap is installed in its extcap folder." >&2
    fi
    exit 1
}

container='%7B%22name%22%3A%22foo%22%2C%22network-interfaces%22%3A%5B%22eth0%22%2C%22lo%22%5D%7D'
marker="e2e${RANDOM}${RANDOM}"
url="packetflix:ws://127.0.0.1:5001/capture?container=${container}&nif=eth0&${marker}"
expect "${marker}" "${url}" \
    -k -i packetflix -o "gui.window_title:foo: eth0, lo" -o "extcap.packetflix.url:${url}"

rm -f /tmp/packetflix-pwned
marker="e2e${RANDOM}${RANDOM}"
url="packetflix:ws://127.0.0.1:5001/\$(touch\$IFS/tmp/packetflix-pwned)?${marker}"
expect "${marker}" "${url}" -k -i packetflix -o "extcap.packetflix.url:${url}"
[[ ! -e /tmp/packetflix-pwned ]] || { echo "FAIL: command injection" >&2; exit 1; }

echo "PASS"
