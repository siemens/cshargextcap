/*
Package e2e contains opt-in end-to-end tests that run the cshargextcap plugin
binary the way Wireshark does against a real Packetflix capture service. They
also drive Wireshark's GUI on Xvfb, inspecting the extcap options dialog via
AT-SPI, and start captures from Edgeshark's web UI in Firefox through the
packaged packetflix: URL handler. They only run when explicitly selected using
the Ginkgo label filter "e2e".
*/
package e2e
