// Package snapshot provides functionality for saving and loading port scan
// results to disk, as well as computing diffs between two snapshots to
// identify newly opened or closed ports on a host.
//
// Snapshots are stored as JSON files and are keyed by host. The Diff function
// compares two snapshots and returns a DiffResult describing what changed.
//
// Example usage:
//
//	snap := snapshot.New("/var/lib/portwatch")
//	if err := snap.Save("192.168.1.1", results); err != nil {
//		log.Fatal(err)
//	}
//	prev, err := snap.Load("192.168.1.1")
//	if err != nil {
//		log.Fatal(err)
//	}
//	diff := snapshot.Diff(prev, results)
package snapshot
