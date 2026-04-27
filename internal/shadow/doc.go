// Package shadow provides a passive secondary scan store for portwatch.
//
// A shadow store records scan results independently of the primary snapshot
// pipeline. It is intended for use cases where an operator wants to capture
// a "silent" view of the network state — for example, to compare a canary
// host against a known-good baseline without triggering the normal alert
// flow.
//
// Usage:
//
//	store, _ := shadow.New("/var/lib/portwatch/shadow")
//	store.Record("192.168.1.1", results)
//	rec, _ := store.Load("192.168.1.1")
//	delta := shadow.Compare(rec, liveResults)
package shadow
