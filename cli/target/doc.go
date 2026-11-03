/*
Package target implements the CLI args specifying what to capture from:

  - “--extcap-interface” selects the extcap network interface, that is, the
    kind of capture service to connect to, such as “mobyshark”.
  - “--container” describes the capture target, such as a pod, container, or
    process, in form of an URL-escaped JSON object; use [Unpack] to get the
    capture target details as a csharg [api.Target].
  - “--nif” lists the names of the network interfaces of the capture target to
    capture from, separated by slashes; if empty, all network interfaces get
    captured from.
  - “--extcap-capture-filter” is the capture filter expression.
  - “--noprom” avoids putting the network interfaces into promiscuous mode.

Wireshark passes “--extcap-interface” and “--extcap-capture-filter” on its own,
see “The Extcap Capture Process” in the [Wireshark Developer's Guide]. The
other args come from the “Standard” tab of the extcap configuration dialog,
except for the packetflix extcap interface, which instead takes them from the
query parameters of its packetflix: URL.

[Wireshark Developer's Guide]: https://www.wireshark.org/docs/wsdg_html_chunked/ChCaptureExtcap.html
*/
package target
