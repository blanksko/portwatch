// Package rollup provides a time-windowed diff accumulator for portwatch.
//
// During periods of rapid port churn a host may produce many small diffs
// in quick succession. Rollup merges those diffs within a configurable
// window so that downstream notifiers and alerters receive a single,
// consolidated change event rather than a stream of noisy updates.
package rollup
