// Package notifier formats port-change diffs and writes human-readable
// summaries to a configurable io.Writer (stdout, file, etc.).
//
// It is intentionally decoupled from transport concerns; callers are
// responsible for wiring the writer to the desired destination.
package notifier
