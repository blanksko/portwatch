// Package backoff implements per-host exponential backoff with
// configurable base duration, ceiling, growth factor, and jitter.
//
// Typical usage:
//
//	b := backoff.Default()
//
//	// on scan failure:
//	wait := b.Next(host)
//	time.Sleep(wait)
//
//	// on scan success:
//	b.Reset(host)
package backoff
