#!/usr/bin/env bash
# Replaces a Windows cshargextcap binary built by goreleaser with an NSIS
# installer containing it; called as: post.sh <binary path> <version> <goarch>
set -euo pipefail

here=$(dirname -- "$(realpath -- "$0")")
root=$(realpath -- "${here}/../..")
installer=$(realpath -- "$1")
version="$2"
arch="$3"

stage="${root}/dist/nsis_${arch}"
rm -rf "${stage}" && mkdir -p "${stage}"
mv "${installer}" "${stage}/cshargextcap.exe"
# Unicode NSIS installers only recognize UTF-8 license texts with a BOM.
{ printf '\xEF\xBB\xBF'; cat "${root}/LICENSE"; } > "${stage}/license.txt"

makensis=(makensis -V2 -WX -INPUTCHARSET UTF8
    "-DVERSION=${version}"
    "-DFILEVERSION=${version%%[-+]*}.0"
    "-DARCH=${arch}"
    "-DBINARY=${stage}/cshargextcap.exe"
    "-DLICENSE=${stage}/license.txt"
    "-DOUTFILE=${installer}"
    "${here}/cshargextcap.nsi")

if command -v makensis >/dev/null; then
    "${makensis[@]}"
else
    docker build -q -t cshargextcap-makensis "${here}" >/dev/null
    docker run --rm --user "$(id -u):$(id -g)" \
        -v "${root}:${root}" -w "${here}" \
        cshargextcap-makensis "${makensis[@]}"
fi
rm -rf "${stage}"
echo "built ${installer#"${root}/"}"
