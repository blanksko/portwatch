// Package fence provides a port-range boundary guard for portwatch.
//
// A Fence is initialised with one or more [Range] values and can be used
// to filter scanner results so that only ports within the declared
// boundaries are passed downstream.  It is safe for concurrent use.
package fence
