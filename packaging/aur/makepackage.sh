#!/bin/bash
set -e

docker build --pull -f Dockerfile -t archpkgbuilder .
docker run -it --rm --name archpkgbuilder -v .:/pkg archpkgbuilder && echo "SUCCESS" || echo "FAIL"
