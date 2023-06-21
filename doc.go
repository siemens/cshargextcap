/*
An external capture plugin for Wireshark for capturing network packets from
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
*/
package cshargextcap
