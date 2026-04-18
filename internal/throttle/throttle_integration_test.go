package throttle_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/throttle"
)

// TestThrottle_ConcurrentHosts verifies that concurrent Allow calls for
// distinct hosts do not race or block each other.
func TestThrottle_ConcurrentHosts(t *testing.T) {
	th := throttle.New(10 * time.Millisecond)
	hosts := []string{"a", "b", "c", "d", "e"}

	done := make(chan struct{}, len(hosts))
	for _, h := range hosts {
		go func(host string) {
			th.Allow(host)
			done <- struct{}{}
		}(h)
	}
	for range hosts {
		<-done
	}
}

// TestThrottle_RealTime verifies throttle behaviour using actual wall-clock time.
func TestThrottle_RealTime(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping real-time throttle test in short mode")
	}
	interval := 50 * time.Millisecond
	th := throttle.New(interval)

	if !th.Allow("host") {
		t.Fatal("first allow should succeed")
	}
	if th.Allow("host") {
		t.Fatal("immediate second allow should be blocked")
	}
	time.Sleep(interval + 10*time.Millisecond)
	if !th.Allow("host") {
		t.Fatal("allow after interval should succeed")
	}
}
