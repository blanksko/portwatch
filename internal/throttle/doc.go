// Package throttle provides per-host scan throttling for portwatch.
//
// It ensures that no single host is scanned more frequently than a
// configured minimum interval, reducing noise and network load.
package throttle
