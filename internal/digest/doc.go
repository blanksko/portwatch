// Package digest provides lightweight fingerprinting of port-scan results.
//
// A [Digest] tracks the SHA-256 hash of each host's open-port set so that
// downstream components can skip expensive diff/alert pipelines when nothing
// has changed since the previous scan cycle.
package digest
