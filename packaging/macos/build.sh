#!/usr/bin/env bash
# Builds packetflix-handler.app as a universal (amd64+arm64) bundle and zips it
# into OUTDIR/packetflix-handler.zip; works on Linux as well as on macOS.
#
# The arm64 binary gets ad-hoc signed by the Go linker; if rcodesign is on the
# PATH, it additionally ad-hoc signs the whole bundle.
set -euo pipefail

version="${1:?usage: build.sh VERSION [OUTDIR]}"
version="${version#v}"
version="${version%%-*}"
mkdir -p "${2:-dist}"
outdir="$(cd "${2:-dist}" && pwd)"
here="$(cd "$(dirname "$0")" && pwd)"
root="$(cd "${here}/../.." && pwd)"

work="$(mktemp -d)"
trap 'rm -rf "${work}"' EXIT
app="${work}/packetflix-handler.app"
mkdir -p "${app}/Contents/MacOS" "${app}/Contents/Resources"

cd "${root}"
for arch in amd64 arm64; do
    CGO_ENABLED=0 GOOS=darwin GOARCH="${arch}" \
        go build -trimpath -ldflags "-s -w" -o "${work}/${arch}" ./cmd/macos-packetflix-handler
done
go tool makefat "${app}/Contents/MacOS/packetflix-handler" "${work}/amd64" "${work}/arm64"
sed "s/@VERSION@/${version}/g" "${here}/Info.plist" > "${app}/Contents/Info.plist"
cp "${here}/packetflix-handler/applet.icns" "${app}/Contents/Resources/packetflix.icns"
printf 'APPL????' > "${app}/Contents/PkgInfo"

if command -v rcodesign >/dev/null; then
    rcodesign sign "${app}"
fi

rm -f "${outdir}/packetflix-handler.zip"
(cd "${work}" && zip -9 -y -r -q "${outdir}/packetflix-handler.zip" packetflix-handler.app)
echo "built ${outdir}/packetflix-handler.zip (version ${version})"
