#!/bin/bash
set -e
DIR=$(dirname -- "$0")

PLUGIN_VERSION="$(cat defs_version.go | sed -n -e 's/^.* SemVersion\s*=\s*"\(.*\)"$/\1/p')"
IFS='-' read -ra FILE_VERSION <<< "$PLUGIN_VERSION"
FILE_VERSION="${FILE_VERSION[0]}.0"
BINARYDIRNAME=/cshargextcap/$(realpath --relative-to="$(pwd)" $(dirname -- "$1"))
BINARYBASENAME=$(basename -- "$1")
BINARYNAME=${BINARYBASENAME%-installer.exe}.exe

echo "!define VERSION \"$PLUGIN_VERSION\"" > ${DIR}/pluginversion.nsh
echo "!define FILEVERSION \"${FILE_VERSION}\"" >> ${DIR}/pluginversion.nsh
echo "!define COPYRIGHT \"Copyright © Siemens $(date +%Y)\"" >> ${DIR}/pluginversion.nsh
echo "!define BINARYPATH \"${BINARYDIRNAME}\"" >> ${DIR}/pluginversion.nsh
echo "!define BINARYNAME \"${BINARYNAME}\"" >> ${DIR}/pluginversion.nsh
echo "!define INSTALLERBINARY \"${BINARYDIRNAME}/${BINARYBASENAME}\"" >> ${DIR}/pluginversion.nsh

sed -i -e "s/(c) 2019\(-[[:digit:]]\+\) \?/(c) 2019-$(date +%Y) /" ${DIR}/license.txt

docker build -t makensis ${DIR}

mv "$1" "$(dirname -- $1)/$BINARYNAME"
docker run -t --rm --network host --volume .:/cshargextcap -w /cshargextcap/packaging/windows makensis makensis cshargextcap.nsi
rm "$(dirname -- $1)/$BINARYNAME"
