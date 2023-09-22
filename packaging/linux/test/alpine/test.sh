#!/bin/bash
set -e


RED="\033[1;31m"
GREEN="\033[1;32m"
NOCOLOR="\033[0m"

docker build --pull -f Dockerfile -t cshargextcap-alpine-test-install .
docker run --rm --name cshargextcap-alpine-test-install -v ./../../../../dist:/dist cshargextcap-alpine-test-install \
    && echo -e "${GREEN}SUCCESS${NOCOLOR}" \
    || (echo -e "${RED}FAIL${NOCOLOR}"; exit 1)
