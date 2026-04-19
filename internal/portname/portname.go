// Package portname maps well-known port numbers to human-readable service names.
package portname

var wellKnown = map[int]string{
	21:   "ftp",
	22:   "ssh",
	23:   "telnet",
	25:   "smtp",
	53:   "dns",
	80:   "http",
	110:  "pop3",
	143:  "imap",
	443:  "https",
	465:  "smtps",
	587:  "submission",
	993:  "imaps",
	995:  "pop3s",
	3306: "mysql",
	5432: "postgres",
	6379: "redis",
	8080: "http-alt",
	8443: "https-alt",
	27017: "mongodb",
}

// Lookup returns the service name for a port number.
// If the port is not recognised, an empty string is returned.
func Lookup(port int) string {
	return wellKnown[port]
}

// LookupWithDefault returns the service name for a port number.
// If the port is not recognised, the provided fallback is returned.
func LookupWithDefault(port int, fallback string) string {
	if name, ok := wellKnown[port]; ok {
		return name
	}
	return fallback
}

// Register adds or overwrites a port-to-name mapping at runtime.
func Register(port int, name string) {
	wellKnown[port] = name
}
