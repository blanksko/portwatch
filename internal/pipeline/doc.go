// Package pipeline provides an ordered chain of scan-processing stages.
//
// Each Stage wraps the next via a next callback, allowing it to inspect,
// modify, or short-circuit the result set before it reaches the terminal
// handler (e.g. snapshot comparison and alerting).
package pipeline
