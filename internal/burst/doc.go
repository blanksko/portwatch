// Package burst implements a sliding-window burst detector for portwatch.
//
// Use Detector.Record to register each change event for a host. Record returns
// true the moment the event count inside the configured window exceeds the
// burst threshold, allowing callers to suppress noisy hosts or raise a
// dedicated high-frequency alert.
package burst
