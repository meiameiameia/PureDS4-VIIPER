// Package server contains the shared PureDS4 listener policy.
package server

import (
	"fmt"
	"net"
	"strconv"
)

// LoopbackHost is also used by the local USB/IP attach client. Keep both sides
// on the same address family without depending on DNS or localhost ordering.
const LoopbackHost = "127.0.0.1"

// LocalListenAddress restricts this embedded backend to the IPv4 loopback
// endpoint used by PureDS4. Hostnames and wildcard/LAN addresses are not allowed;
// the historical localhost spelling is normalized without a DNS lookup.
func LocalListenAddress(addr string) (string, error) {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return "", fmt.Errorf("invalid local listen address %q: %w", addr, err)
	}
	if host != LoopbackHost && host != "localhost" {
		return "", fmt.Errorf("PureDS4 requires listen host %s, got %q", LoopbackHost, host)
	}
	if _, err := strconv.ParseUint(port, 10, 16); err != nil {
		return "", fmt.Errorf("invalid local listen port %q: %w", port, err)
	}
	return net.JoinHostPort(LoopbackHost, port), nil
}
