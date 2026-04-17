package metrics

import (
	"fmt"
	"io"
	"text/tabwriter"
)

// Print writes a human-readable summary of the snapshot to w.
func Print(w io.Writer, s Snapshot) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "Metric\tValue")
	fmt.Fprintln(tw, "------\t-----")
	fmt.Fprintf(tw, "Total scans\t%d\n", s.TotalScans)
	fmt.Fprintf(tw, "Open ports (last)\t%d\n", s.OpenPorts)
	if s.LastScanAt.IsZero() {
		fmt.Fprintf(tw, "Last scan\t-\n")
	} else {
		fmt.Fprintf(tw, "Last scan\t%s\n", s.LastScanAt.Format("2006-01-02 15:04:05"))
	}
	fmt.Fprintf(tw, "Last duration\t%s\n", s.LastDuration)
	return tw.Flush()
}
