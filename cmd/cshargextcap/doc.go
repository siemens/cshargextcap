/*
Cshargextcap is an external capture (extcap) plugin for Wireshark that remotely
captures network packets from containers, k8s pods, and other capture targets.

Wireshark runs cshargextcap itself, passing exactly one of the extcap action
flags together with further flags. Before parsing, cshargextcap fixes up
Wireshark's way of passing boolean flags as "--flag true" or "--flag false".

If an extcap action fails, cshargextcap exits with the action's non-zero exit
code; other errors result in exit code 1.

Usage:

	cshargextcap [flags]

The extcap action flags are:

	--extcap-interfaces
		List the available extcap network interfaces.
	--extcap-dlts
		List the data link types of the extcap network interface.
	--extcap-config
		List the configuration options of the extcap network interface.
	--capture
		Capture packets and write them into the fifo.

The flags selecting what to capture from are:

	--extcap-interface name
		The extcap network interface to use.
	--container details
		Pod, container, or process reference details to capture from.
	--nif names
		Network interface name(s) to capture from, separated by slashes.
	--extcap-capture-filter expr
		Capture filter expression.
	--noprom
		Don't put the network interface(s) into promiscuous mode.

The flags selecting which kinds of capture targets to offer are:

	--showpods
		Discover pods.
	--showcontainers
		Discover stand-alone containers.
	--showprocs
		Discover non-containers with virtual IP stacks.
	--showemptynetns
		Discover IP stacks without any attached processes.

The further Wireshark extcap flags are:

	--extcap-reload-option name
		Reload the values of a specific configuration option.
	--extcap-version version
		Version of Wireshark invoking cshargextcap.
	--fifo path
		Wireshark fifo pathname to write the captured packets to.

The connection flags are:

	--proxyoff
		Never use any HTTP/HTTPS proxy.
	--httpproxy url
		HTTP proxy URL.
	--httpsproxy url
		HTTPS proxy URL.
	--noproxy destinations
		Destination(s) to which no proxy should be used.
	--insecure
		Skip server certificate verification.
	--timeout seconds
		Time limit for discovery and establishing capture connections
		(default 30).

The debugging flags are:

	--debug
		Enable debug logging.
	--debug-file path
		Enable debug logging to the specified file (implies --debug).
*/
package main
