<img align="right" width="100" height="100" src="images/csharg-icon-100x100-ltr.png" style="padding: 0 0 1ex 0.8em">

[![Siemens](https://img.shields.io/badge/github-siemens-009999?logo=github)](https://github.com/siemens)
[![Industrial Edge](https://img.shields.io/badge/github-industrial%20edge-e39537?logo=github)](https://github.com/industrial-edge)
[![Edgeshark](https://img.shields.io/badge/github-Edgeshark-003751?logo=github)](https://github.com/siemens/edgeshark)

# Containershark Extcap Plugin for Wireshark

[![PkgGoDev](https://pkg.go.dev/badge/github.com/siemens/cshargextcap)](https://pkg.go.dev/github.com/siemens/cshargextcap)
[![GitHub](https://img.shields.io/github/license/siemens/cshargextcap)](https://img.shields.io/github/license/siemens/cshargextcap)
![build and test](https://github.com/siemens/cshargextcap/workflows/build%20and%20test/badge.svg?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/siemens/cshargextcap)](https://goreportcard.com/report/github.com/siemens/cshargextcap)

> [!IMPORTANT]  
> Wireshark 4.4.0 is not supported as it breaks extcaps such as this one.
> Wireshark 4.4.1 scheduled for Oct 9th 2024 will contain two fixes so that this
> extcap plugin will be able to correctly work again.

Take a deep dive into your container host virtual networking, even if it's in a
remote location. No fiddling with special containers and juggling error-prone
CLI Docker commands. Simply click on a "fin" capture button inside one of the
containers in Ghostwire's web UI to start a Wireshark live capture session:

![Click the Fin!](images/gw-fin-ui.png)

Confirm and we're live capturing.

![](images/ws-pf-capture.png)

## What You See Is What You Get

- Capture network traffic _live_ from your containers (and pods), straight into
  your Desktop Wireshark on _Linux_ and _Windows_.

- Capture from any container _without_ preparing or modifying it for capturing.
  Just go capturing.

- Supports stand-alone container hosts, including the Siemens Industrial Edge.

- Remotely capture not only from containers, but also from the container host
  itself, process-less network namespaces, et cetera.

- this Wireshark plugin can be build for Windows 64bit (x86) as well as Linux
  64bit (x86, ARM). Currently, we support the .apk, .deb, and .rpm package
  formats on Linux.

## Installation

- **Linux**: 
  - **Alpine**, **Debian**, **Fedora**: head over to our
[releases](https://github.com/siemens/cshargextcap/releases/latest) page and
    download the package matching your CPU architecture (amd64 or arm64) and
    distro package format. Install the downloaded package as usual.
  - **Arch**:
    1. download `PKGBUILD` into a preferably clean directory:
       ```bash
       wget https://raw.githubusercontent.com/siemens/cshargextcap/main/packaging/aur/PKGBUILD
       ```
    2. `makepkg -s -r -c` in the directory you've downloaded `PKGBUILD` into.
    3. either install only the excapt plugin, or additionally the Wireshark QT
       desktop integration:
       ```bash
       # plugin only for tshark usage, without desktop Wireshark dependency
       pacman -U cshargextcap-git-cli*.zst
       # with desktop integration
       pacman -U cshargextcap-git-*.zst
       ```
    4. You can later update to new releases without the need to download the
       `PKGBUILD` file again, as it will automatically build and install from
       the latest release.

- **Mac OS**: head over to our
  [releases](https://github.com/siemens/cshargextcap/releases/latest) page and...
  1. download the `.tar.gz` archive for Darwin/macos arm64 or amd64 (Intel).
     Extract the contained `cshargextcap` plugin binary and copy/move it to
  `/Applications/Wireshark.app/Contents/MacOS/extcap`.
  2. download the `packetflix-handler.zip` archive.
  3. run the CLI command `xattr -d com.apple.quarantine packetflix-handler.zip`.
  4. unpack `packetflix-handler.zip` by double clicking it in Finder.
  5. copy `packetflix-handler`(`.app`) to your Applications folder,
     `/Applications`.
  6. go to "System Preferences" > "Security and Privacy" > tab "General" or
     "Security" section. Allow the packetflix-handler. In case you don't see
     anything here, try to start a capture from the web UI first, and as this
     will be blocked, you should now see here a notice with a button to enable
     the packetflix-handler.

- **Windows**: head over to our
  [releases](https://github.com/siemens/cshargextcap/releases/latest) page and
  download the ZIP archive for Windows amd64. Double click in file explorer to
  open its contents, then double click on the installer `.exe`. You don't need
  to extract the other files, as the installer perfectly works on its own.

See below for the [Quick Start](#quick-start).

## Project Map

The Containershark extcap plugin is part of the "Edgeshark" project that consist
of several repositories:
- [Edgeshark Hub repository](https://github.com/siemens/edgeshark)
- [G(h)ostwire discovery service](https://github.com/siemens/ghostwire)
- [Packetflix packet streaming service](https://github.com/siemens/packetflix)
- 🖝 **Containershark Extcap plugin for Wireshark** 🖜
- support modules:
  - [csharg (CLI)](https://github.com/siemens/csharg)
  - [mobydig](https://github.com/siemens/mobydig)
  - [ieddata](https://github.com/siemens/ieddata)

## Quick Start

Please deploy the [G(h)ostwire discovery
service](https://github.com/siemens/ghostwire) and [Packetflix packet streaming
service](https://github.com/siemens/packetflix) on your Docker host.

Then install this plugin: on Windows download and install the cshargextcap
installer artifact. On Linux, download and install the cshargextcap package for
your distribution (apk, deb, or rpm). In case you want to create the
installation files yourself, then simply run `make dist` in the base directory
of this repository. Please note that this will also test install the packages in distro-specific test containers to ensure the distro packages are fine. Afterwards, installation files can be found in the `dist/`
directory.

Now fire up Wireshark. If the installation went through correctly, Wireshark now
should show two new "interfaces", as shown below: 

![Container Live Capture](images/cs-docker-defaulttab.png)

It's as easy as this:

1. click the ⚙ _gear icon_ next to the network interface named "Docker host
   capture".
2. enter your Docker host's IP address or DNS name, as well as port `:5001` into
   the field "Docker host URL".
3. click the _refresh_ button to get the list of available pods (and more...).
4. pick your container.
5. click the _Start_ button.

...and your live capture starts immediately.

> **🛈** Wireshark creates the UI for our capture plugin and unfortunately we're
> therefore (quite) limited to what Wireshark has on offer. Please don't create
> UI/UX feature requests as we don't have any control over Wireshark's UI – with
> the exception of our _own_ bugs: please create issues for them in this
> project's issue tracker.

Please find more details in our [csharg Extcap ⚙ Plugin Manual](docs/manual.md).

Finally, there's also some technical background information in our
[csharg ⚙ Plugin Technical Details](docs/technical.md).

# Contributing

Please see [CONTRIBUTING.md](CONTRIBUTING.md).

## License and Copyright

(c) Siemens AG 2023-24

[SPDX-License-Identifier: MIT](LICENSE)
