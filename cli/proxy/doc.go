/*
Package proxy implements the HTTP(S) proxy-related CLI args, which Wireshark
passes from the “Proxy” tab of the extcap configuration dialog:

  - “--proxyoff” disables all proxies, ignoring any proxy-related environment
    variables inherited from Wireshark.
  - “--httpproxy” and “--httpsproxy” specify the proxy URLs to use for HTTP
    and HTTPS connections respectively.
  - “--noproxy” specifies the destinations to connect to without any proxy.

# Note

Instead of configuring HTTP clients directly, the proxy settings are applied
to the process environment in form of the conventional “http_proxy”,
“https_proxy”, and “no_proxy” variables, in both lower and upper case. This
happens before running any extcap action, so that all HTTP clients honoring
these variables, such as those using [net/http.ProxyFromEnvironment], pick up
the proxy settings.
*/
package proxy
