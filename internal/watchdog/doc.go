// Package watchdog coordinates a full port-watch cycle.
//
// A Watchdog instance wires together the scanner, snapshot store,
// notifier, history recorder, and metrics collector so that the
// scheduler only needs to call Run once per tick.
package watchdog
