package scanner

import (
	"fmt"
	"net"
	"time"
)

// PortState represents the state of a scanned port.
type PortState struct {
	Port     int
	Protocol string
	Open     bool
}

// ScanResult holds the results of a full scan.
type ScanResult struct {
	Host      string
	Timestamp time.Time
	Ports     []PortState
}

// Scanner scans ports on a given host.
type Scanner struct {
	Timeout time.Duration
}

// New creates a new Scanner with the given timeout.
func New(timeout time.Duration) *Scanner {
	return &Scanner{Timeout: timeout}
}

// Scan checks the given ports on host and returns a ScanResult.
func (s *Scanner) Scan(host string, ports []int) (*ScanResult, error) {
	if host == "" {
		return nil, fmt.Errorf("host must not be empty")
	}

	result := &ScanResult{
		Host:      host,
		Timestamp: time.Now(),
	}

	for _, port := range ports {
		state := PortState{
			Port:     port,
			Protocol: "tcp",
			Open:     s.isOpen(host, port),
		}
		result.Ports = append(result.Ports, state)
	}

	return result, nil
}

func (s *Scanner) isOpen(host string, port int) bool {
	address := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", address, s.Timeout)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}
