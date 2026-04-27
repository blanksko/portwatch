// Package probe manages per-host probe intervals for portwatch.
//
// A Manager holds an optional per-host duration and falls back to
// a global default, making it straightforward to scan high-risk
// hosts more frequently than the rest of the fleet.
package probe
