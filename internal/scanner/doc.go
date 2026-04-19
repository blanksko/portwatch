// Package scanner provides functionality to probe TCP ports on a target host
// and report their open/closed state.
//
// Basic usage:
//
//	s := scanner.New(500 * time.Millisecond)
//	result, err := s.Scan("localhost", []int{22, 80, 443, 8080})
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, p := range result.Ports {
//		if p.Open {
//			fmt.Printf("Port %d/%s is OPEN\n", p.Port, p.Protocol)
//		}
//	}
//
// The Scanner respects a configurable dial timeout to avoid blocking
// indefinitely on filtered ports. Scans are performed concurrently across
// all requested ports, with results collected and returned together once
// all probes have completed.
package scanner
