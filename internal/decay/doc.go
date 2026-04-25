// Package decay implements an exponential-decay score tracker for hosts
// monitored by portwatch. Scores rise when new port-change events are
// recorded and fall automatically between scan cycles, giving a simple
// signal of how "active" or "noisy" a given host has been recently.
package decay
