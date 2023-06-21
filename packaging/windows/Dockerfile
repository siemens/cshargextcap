# Simple recent makensis container image for generating Nullsoft Installers for
# Windows on Linux. Sweet. See the accompanying docker-compose.yaml for how this
# image then builds the installer binary when running the image.
FROM ubuntu:22.04
RUN apt-get update && apt-get install -y nsis
CMD ["makensis"]
