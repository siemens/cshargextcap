/*
Package cshargextcap is an external capture plugin for Wireshark for capturing network packets from
inside containers (Docker, others) without having to prepare these containers
for capturing. These containers can be deployed on single hosts, but also in
Kubernetes clusters.

This OpenSource module implements capture clients for capturing from container
hosts (including a KinD host). The following two extcap interfaces are
implemented:
  - [MobyNif] connects to a capture service via http: and https: protocol URLs.
  - [PacketflixNif] connects to a capture service described by a packetflix:
    protocol URL. Under the hood, the packetflix: protocol encodes http: or
    https: URLs. The use case for packetflix: URLs is to allow easy hand-over
    from web-browser based UIs to Wireshark.

Wireshark's extcap interface is documented in the section “Adding Capture
Interfaces And Log Sources Using Extcap” of the [Wireshark Developer's Guide],
with the grammar of the extcap output in the [extcap(4)] manual page.

[Wireshark Developer's Guide]: https://www.wireshark.org/docs/wsdg_html_chunked/ChCaptureExtcap.html
[extcap(4)]: https://www.wireshark.org/docs/man-pages/extcap.html
*/
package cshargextcap
