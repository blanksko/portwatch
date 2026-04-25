// Package checkpoint provides a persistent store for per-host scan
// checkpoints. Each checkpoint records the wall-clock time of the most recent
// completed scan cycle for a host, allowing portwatch to skip redundant
// alerting when the process is restarted.
package checkpoint
