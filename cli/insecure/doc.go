/*
Package insecure implements the “--insecure” CLI arg for skipping the
verification of the capture service's server certificate when connecting via
https: or wss:, shared by all extcap network interfaces connecting to capture
services. Use [SkipVerify] to check whether to skip verification.
*/
package insecure
