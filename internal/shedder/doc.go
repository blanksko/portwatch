// Package shedder provides a concurrency-based load shedder for
// portwatch scan jobs. It tracks the number of active in-flight
// scans and rejects new requests once a configured ceiling is
// reached, preventing runaway resource consumption on busy hosts.
//
// Typical usage:
//
//	s := shedder.New(10)
//	if err := s.Acquire(); err != nil {
//		// drop or log the request
//		return
//	}
//	defer s.Release()
//	// … perform scan …
package shedder
