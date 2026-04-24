// Package limiter caps the number of concurrent port-scan goroutines
// that may run against a single host at any one time.
//
// Use New to create a Limiter, then wrap each scan goroutine with
// Acquire / Release pairs to stay within the configured concurrency
// budget.
package limiter
