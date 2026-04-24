// Package healthcheck provides lightweight TCP liveness probing for monitored
// hosts. It distinguishes hosts that are genuinely unreachable from hosts that
// simply have no open ports at a given moment, helping portwatch avoid
// generating spurious alerts when a host is temporarily offline.
//
// Results are cached per-host so that callers can query the last known status
// without incurring a network round-trip on every access.
package healthcheck
