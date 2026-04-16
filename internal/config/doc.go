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
// Use [Load] to read a file from disk, or [Default] to obtain a
// ready-to-use configuration without a file.
package config
