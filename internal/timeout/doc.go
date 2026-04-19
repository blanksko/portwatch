// Package timeout manages per-host TCP scan timeouts for portwatch.
//
// A Manager is created with a default duration and allows individual
// hosts to override that value. This is useful when some targets are
// known to be slow or unreliable and require a longer dial window.
package timeout
