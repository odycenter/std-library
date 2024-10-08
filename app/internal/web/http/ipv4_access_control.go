package internal_http

import (
	"errors"
	"log/slog"
	"net"
)

type IPv4AccessControl struct {
	Allow *IPv4Ranges
	Deny  *IPv4Ranges
}

func (ac *IPv4AccessControl) Validate(clientIP string) error {
	address := net.ParseIP(clientIP)
	if address == nil {
		return errors.New("invalid IP address")
	}

	if ac.isLocal(address) {
		slog.Debug("allow site local client address")
		return nil
	}

	if !ac.allow(address) {
		return errors.New("access denied: IP_ACCESS_DENIED")
	}

	return nil
}

func (ac *IPv4AccessControl) isLocal(ip net.IP) bool {
	if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
		return true
	}

	return ip.IsPrivate()
}

func (ac *IPv4AccessControl) allow(address net.IP) bool {
	if address.To4() == nil { // only support ipv4
		slog.Debug("skip with ipv6 client address")
		return true
	}

	if ac.Allow != nil && ac.Allow.Matches(address) {
		slog.Debug("allow client ip within allowed ranges")
		return true
	}

	if ac.Deny == nil || ac.Deny.Matches(address) { // if deny == nil, it blocks all
		slog.Debug("deny client ip within denied ranges")
		return false
	}

	slog.Debug("allow client ip not within denied ranges")
	return true
}
