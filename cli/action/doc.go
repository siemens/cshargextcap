/*
Package action implements the handling of Wireshark's extcap CLI action flag.

# Note

Wireshark tells external capture plugins which action to perform not through
action verbs (as it is good practise nowadays) but instead through a set of
mutually exclusive flags, see “Extcap Command Line Interface” in the [Wireshark
Developer's Guide]:

  - “--extcap-interfaces”: “Query For Available Interfaces”,
  - “--extcap-dlts”: “Ask For DLTs For Each Interface”,
  - “--extcap-config”: “The Extcap Configuration Interface”,
  - “--capture”: “The Extcap Capture Process”.

[Wireshark Developer's Guide]: https://www.wireshark.org/docs/wsdg_html_chunked/ChCaptureExtcap.html
*/
package action
