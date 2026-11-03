/*
Package wireshark implements the Wireshark-specific extcap CLI args that are
not tied to a particular extcap action or capture target:

  - “--extcap-reload-option” names the single configuration arg whose values
    Wireshark wants to refresh, in combination with “--extcap-config”; see
    “Reload A Selector” in the [Wireshark Developer's Guide].
  - “--extcap-version” is the version of Wireshark invoking this plugin; see
    “Query For Available Interfaces”.
  - “--fifo” is the path of the named pipe to write the captured packet stream
    into, in combination with “--capture”; see “The Extcap Capture Process”.

[Wireshark Developer's Guide]: https://www.wireshark.org/docs/wsdg_html_chunked/ChCaptureExtcap.html
*/
package wireshark
