// Package quota enforces per-host scan rate limits using a rolling time window.
//
// A Quota instance tracks how many times each host has been scanned within a
// configurable duration and rejects requests that exceed the configured maximum.
//
// Example usage:
//
//	q := quota.New(10, time.Minute)
//	if err := q.Allow(host); err != nil {
//		log.Printf("skipping %s: %v", host, err)
//		return
//	}
package quota
