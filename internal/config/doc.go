// Package config provides loading, saving, and validation of portwatch
// runtime configuration.
//
// A configuration file is a JSON document with the following structure:
//
//	{
//	  "hosts": ["localhost", "192.168.1.1"],
//	  "ports": [22, 80, 443],
//	  "snapshot_path": "portwatch_snapshot.json",
//	  "timeout_seconds": 2
//	}
//
// Fields:
//
//   - hosts: list of hostnames or IP addresses to scan.
//   - ports: list of TCP port numbers to check on each host.
//   - snapshot_path: path to the JSON file used to persist the last known
//     port-state snapshot between runs.
//   - timeout_seconds: per-connection dial timeout; defaults to 2 if omitted
//     or set to zero.
//
// Use [Load] to read a file from disk, or [Default] to obtain a
// ready-to-use configuration without a file.
package config
