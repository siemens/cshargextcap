# makensis for building the Windows installer on Linux; used by post.sh only
# when there's no local makensis.
FROM debian:13-slim
RUN apt-get update \
    && apt-get install -y --no-install-recommends nsis \
    && rm -rf /var/lib/apt/lists/*
