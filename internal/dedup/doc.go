// Package dedup implements host-level result deduplication for portwatch.
//
// A Store remembers the fingerprint of the last scan result for each host.
// Callers can ask whether results have Changed since the previous scan and
// skip expensive downstream work (alerting, snapshotting, history recording)
// when the port set is identical.
package dedup
