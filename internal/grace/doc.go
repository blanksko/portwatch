// Package grace provides graceful shutdown coordination for portwatch.
//
// A Coordinator listens for SIGINT / SIGTERM and waits up to a configurable
// timeout for any in-flight scan cycles to complete before allowing the
// process to exit, preventing partial writes to snapshot and history files.
package grace
