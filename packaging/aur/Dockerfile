FROM archlinux:base-devel

RUN useradd -m build && \
    pacman -Syu --noconfirm && \
    pacman -Sy --noconfirm git && \
    echo "build ALL=(ALL) NOPASSWD: ALL" > /etc/sudoers.d/build
USER build
WORKDIR /pkg
COPY packagebuild.sh /
CMD ["/bin/bash", "packagebuildandtest.sh"]
