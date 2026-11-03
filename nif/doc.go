/*
Package nif implements the “mobyshark” and “packetflix” extcap network
interfaces for Wireshark. These provide capturing network traffic from inside
stand-alone container hosts with the packetflix service deployed to them.

Both extcap network interfaces finally capture in the same way, using csharg's
container host client, but they differ in how the user selects what to capture
from:

  - [MobyNif] (“mobyshark”) lets the user select the capture target inside
    Wireshark: the user enters the URL of the Docker host in the extcap
    configuration dialog, and then reloads the “Containers” and “Interfaces”
    selectors. On reload, the extcap discovers the capture targets on that
    host, filtered by the “Show …” options. Wireshark finally passes the
    selected capture target and network interfaces as separate CLI args.
  - [PacketflixNif] (“packetflix”) instead receives a capture target already
    selected outside of Wireshark, typically in a web UI that hands over a
    “packetflix:” URL to Wireshark. This URL encodes the capture service's
    websocket URL, with the capture target, network interfaces, capture
    filter, and promiscuous mode as query parameters. The extcap converts this
    URL into the container host URL and the CLI args mobyshark would get, and
    adds the node name from the discovered capture targets. Its configuration
    dialog thus lacks the target selection on the “Standard” tab and only asks
    for the URL.
*/
package nif
