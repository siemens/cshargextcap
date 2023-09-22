#!/bin/bash
set -e

RED="\033[1;31m"
GREEN="\033[1;32m"
NOCOLOR="\033[0m"

docker build --pull -f Dockerfile -t archpkgbuilder .
docker run --rm --name archpkgbuilder -v .:/pkg archpkgbuilder \
    && echo -e "${GREEN}SUCCESS${NOCOLOR}" \
    || (echo -e "${RED}FAIL${NOCOLOR}"; exit 1)
