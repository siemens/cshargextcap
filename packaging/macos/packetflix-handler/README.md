# macos packetflix:// URL Handler

`packetflix-handler.app` is built from the Go sources in
[`cmd/macos-packetflix-handler`](../../../cmd/macos-packetflix-handler) by
[`../build.sh`](../build.sh) (or `make macos-packetflix-handler`), also on Linux. It starts
a new Wireshark instance located by its bundle ID `org.wireshark.Wireshark`.

[`../test/handlertest.sh`](../test/handlertest.sh) tests the built app on macOS
end-to-end via LaunchServices, using a fake Wireshark; CI runs it on Apple
Silicon and Intel runners.

[`../macos-install-handler.sh`](../macos-install-handler.sh) installs or upgrades the plugin and handler for
the current user; `macos-install-handler.sh --from dist` installs a local build after `make
dist-build macos-packetflix-handler`. The handler logs to
`~/Library/Logs/packetflix-handler.log`.

`applet.icns` is the app icon.
