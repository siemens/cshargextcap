FROM archlinux:base-devel

RUN useradd -m build && \
    pacman -Syu --noconfirm && \
    pacman -Sy --noconfirm git xdg-utils && \
    echo "build ALL=(ALL) NOPASSWD: ALL" > /etc/sudoers.d/build
USER build
WORKDIR /pkg
COPY packagebuildandtest.sh /
CMD ["/bin/bash", "packagebuildandtest.sh"]
