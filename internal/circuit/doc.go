// Package circuit implements a lightweight per-host circuit breaker for
// portwatch scan pipelines.
//
// A Breaker tracks consecutive scan failures per host. Once the failure
// count reaches the configured threshold the circuit opens, causing
// subsequent Allow calls to return false and skip the scan. After the
// cooldown period the circuit transitions to half-open, allowing a single
// probe attempt. A successful scan closes the circuit and resets the
// counter; another failure re-opens it immediately.
//
// Usage:
//
//	br := circuit.New(5, 30*time.Second)
//	if !br.Allow(host) {
//	    // skip scan
//	}
//	// ... perform scan ...
//	if err != nil {
//	    br.RecordFailure(host)
//	} else {
//	    br.RecordSuccess(host)
//	}
package circuit
