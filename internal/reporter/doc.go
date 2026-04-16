// Package reporter provides formatting and output utilities for portwatch
// scan results.
//
// A Reporter writes human-readable summaries of port scan diffs to any
// io.Writer, such as os.Stdout or a log file. It is intended to complement
// the alert package, which handles structured notifications.
package reporter
