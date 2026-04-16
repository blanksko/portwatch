// Package alert provides notification functionality for portwatch.
// It formats and writes human-readable alerts when port changes are detected
// between two consecutive scans of a host.
//
// An Alerter is constructed with a target io.Writer (e.g. os.Stdout or a log
// file) and a logger for internal diagnostics. The Notify method accepts a
// snapshot.DiffResult and emits lines describing opened and closed ports.
//
// Example usage:
//
//	a := alert.New(os.Stdout, logger)
//	if err := a.Notify("192.168.1.1", diff); err != nil {
//		log.Println("alert error:", err)
//	}
package alert
