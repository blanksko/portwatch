// Package filter provides port filtering for portwatch scan results.
//
// A Filter can be configured with an inclusion list, an exclusion list, or both.
// When an inclusion list is provided, only ports in that list are kept.
// Exclusions are always applied and take precedence over inclusions.
package filter
