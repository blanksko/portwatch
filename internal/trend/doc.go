// Package trend provides a sliding-window change-frequency tracker for
// monitored hosts. It records opened/closed port events and exposes
// counts within a configurable time window, enabling callers to detect
// hosts whose port sets are changing unusually often.
package trend
