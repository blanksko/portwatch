// Package sampler implements a frequency-based port-result sampler for
// portwatch. It accumulates scan observations over multiple cycles and
// surfaces only those ports that appear consistently, reducing transient
// false-positives caused by ephemeral or flapping services.
package sampler
