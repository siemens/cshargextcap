<img align="right" width="100" height="100" src="../images/csharg-icon-100x100-ltr.png" style="padding: 0 0 1ex 0.8em">

[![Siemens](https://img.shields.io/badge/github-siemens-009999?logo=github)](https://github.com/siemens)
[![Industrial Edge](https://img.shields.io/badge/github-industrial%20edge-e39537?logo=github)](https://github.com/industrial-edge)
[![Edgeshark](https://img.shields.io/badge/github-Edgeshark-003751?logo=github)](https://github.com/siemens/edgeshark)

# Containershark ⚙ Plugin Technical Details

Please️ note that function names such as `ExtcapConfigure` are high-level
functions of the extcap pluging. These then rely on
[csharg](https://github.com/siemens/csharg) library functions for the dirty
work.

## Wireshark Startup

When Wireshark starts, it looks into its `extcap` directories (system-wide and
per user) searching for plugin executables, which it then starts one after
another in order to gather information:

- the network interfaces provided by these extcap plugins,
- the UI configuration parameters used by these plugins,
- and finally the data link layer encapsulation used by these network interfaces.

The important takeaway here is that extcap plugins are very "_RESTless_" because
Wireshark calls them _multiple_ times. This does not exactly improve performance
due to having to initialize multiple times. However, the compiled code and the
Go runtime do a well performing job here in our experience, so this issue is
mitigated to a good extend.

```mermaid
sequenceDiagram
    actor user
    participant Wireshark
    participant cshargextcap as cshargextcap<br/>«extcap plugin»

    user->>Wireshark: start
    Note over Wireshark,cshargextcap: each call runs a new plugin process

    Wireshark->>+cshargextcap: --extcap-interfaces
    cshargextcap->>cshargextcap: ExtcapInterfaces
    Note over cshargextcap: supported version, help URL:<br/>extcap {version=...}{help=...}<br/>list of "network interfaces":<br/>interface {value=mobyshark}{display=...}<br/>interface ...
    cshargextcap-->>-Wireshark: exit

    Wireshark->>+cshargextcap: --extcap-config --extcap-interface mobyshark
    cshargextcap->>+cshargextcap: ExtcapConfigure
    cshargextcap->>cshargextcap: MobyNif.Configure
    deactivate cshargextcap
    Note over cshargextcap: list of config UI arguments:<br/>arg {number=...}...
    cshargextcap-->>-Wireshark: exit

    Wireshark->>+cshargextcap: --extcap-dlts --extcap-interface mobyshark
    cshargextcap->>cshargextcap: ExtcapDlts
    Note over cshargextcap: dlt {number=147}{name=...}{display=...}
    cshargextcap-->>-Wireshark: exit
```

## Interface Configuration

The user clicks on the gear icon next to an extcap plugin network interface.
And yes, Wireshark asks for the UI configuration parameters _again_. Then, the
user presses a "reload" button in order to update, for instance, the list of
containers.

To simplify single-host service configuration the plugin only talks to the
Packetflix service on a standalone container host, and Packetflix implements a
forwarding path to reach the associated Ghostwire service instance.

```mermaid
sequenceDiagram
    actor user
    participant Wireshark
    participant cshargextcap as cshargextcap<br/>«extcap plugin»
    participant Packetflix as Packetflix<br/>«service»
    participant Ghostwire as Ghostwire<br/>«service»

    user->>Wireshark: configure
    Wireshark->>+cshargextcap: --extcap-config<br/>--extcap-interface mobyshark
    cshargextcap->>cshargextcap: MobyNif.Configure
    Note over cshargextcap: list of config UI arguments:<br/>arg {number=...}...
    cshargextcap-->>-Wireshark: exit

    user->>Wireshark: reload (containers)
    Wireshark->>+cshargextcap: --extcap-config<br/>--extcap-interface mobyshark<br/>--extcap-reload-option container
    cshargextcap->>+cshargextcap: ExtcapReloadOption
    cshargextcap->>+cshargextcap: MobyNif.ReloadOption
    cshargextcap->>+Packetflix: GET /discover/mobyshark
    Packetflix->>+Ghostwire: GET /mobyshark
    Ghostwire-->>-Packetflix: discovery results
    Packetflix-->>-cshargextcap: discovery results
    deactivate cshargextcap
    deactivate cshargextcap
    Note over cshargextcap: list of config UI values, Wireshark only accepts values<br/>for the arg {number=...} whose "reload" button was<br/>pressed, other value {number=...} will be ignored:<br/>value {number=...}...
    cshargextcap-->>-Wireshark: exit
```

## Capture from Paketflix Service via External Capture Plugin

After all the hassle, let's capture packets. Here's the basic picture – where
we've left out all the complex details about updating the toolbar and statusbar,
and then even more.

```mermaid
sequenceDiagram
    actor user
    participant Wireshark
    participant cshargextcap as cshargextcap<br/>«extcap plugin»
    participant Packetflix as Packetflix<br/>«service»

    user->>Wireshark: start (capture)
    Wireshark->>+cshargextcap: --capture<br/>--extcap-interface mobyshark<br/>--fifo ...
    cshargextcap->>Wireshark: fifo.open('wb')
    cshargextcap->>+Packetflix: ws://... GET /capture?...
    Packetflix-->>cshargextcap: connected
    alt stream packet capture data
        loop each message
            Packetflix->>cshargextcap: message(data)
            cshargextcap->>Wireshark: fifo.write(data)
        end
    else stop (capture)
        user->>Wireshark: stop
        Wireshark->>Wireshark: fifo.close
        Wireshark-)cshargextcap: fifo closed
        cshargextcap->>Packetflix: close
    end
    Packetflix-->>-cshargextcap: closed
    deactivate cshargextcap
```
