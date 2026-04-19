// Package resolve provides hostname-to-IP resolution utilities for portwatch.
package resolve

import (
	"context"
	"fmt"
	"net"
	"time"
)

// Resolver resolves hostnames to IP addresses.
type Resolver struct {
	timeout time.Duration
	resolver *net.Resolver
}

// Result holds the resolved addresses for a hostname.
type Result struct {
	Host      string
	Addresses []string
	ResolvedAt time.Time
}

// New returns a Resolver with the given timeout.
func New(timeout time.Duration) *Resolver {
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	return &Resolver{
		timeout:  timeout,
		resolver: net.DefaultResolver,
	}
}

// Resolve returns the IP addresses for the given hostname.
// If host is already an IP address it is returned as-is.
func (r *Resolver) Resolve(host string) (Result, error) {
	if ip := net.ParseIP(host); ip != nil {
		return Result{
			Host:       host,
			Addresses:  []string{ip.String()},
			ResolvedAt: time.Now(),
		}, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	addrs, err := r.resolver.LookupHost(ctx, host)
	if err != nil {
		return Result{}, fmt.Errorf("resolve %q: %w", host, err)
	}

	return Result{
		Host:       host,
		Addresses:  addrs,
		ResolvedAt: time.Now(),
	}, nil
}

// Primary returns the first resolved address, or an error if none exist.
func (r *Resolver) Primary(host string) (string, error) {
	res, err := r.Resolve(host)
	if err != nil {
		return "", err
	}
	if len(res.Addresses) == 0 {
		return "", fmt.Errorf("resolve %q: no addresses returned", host)
	}
	return res.Addresses[0], nil
}
