// Package cooldown implements per-host alert cooldown tracking.
//
// A Cooldown prevents repeated notifications for the same host within a
// configurable time window. Once an alert is issued for a host, further
// calls to Allow will return false until the window elapses.
//
// Example usage:
//
//	cd := cooldown.New(10*time.Minute, time.Now)
//	if cd.Allow(host) {
//		// send alert
//	}
package cooldown
