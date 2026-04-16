package reporter_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/reporter"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/snapshot"
)

func makeDiff(opened, closed []int) snapshot.Diff {
	toResults := func(ports []int) []scanner.Result {
		var rs []scanner.Result
		for _, p := range ports {
			rs = append(rs, scanner.Result{Port: p, Protocol: "tcp", Open: true})
		}
		return rs
	}
	return snapshot.Diff{
		Opened: toResults(opened),
		Closed: toResults(closed),
	}
}

func TestReport_NoChanges(t *testing.T) {
	var buf bytes.Buffer
	r, _ := reporter.New(&buf)
	r.Report(makeDiff(nil, nil))
	if !strings.Contains(buf.String(), "No changes detected") {
		t.Errorf("expected no-changes message, got: %s", buf.String())
	}
}

func TestReport_OpenedPorts(t *testing.T) {
	var buf bytes.Buffer
	r, _ := reporter.New(&buf)
	r.Report(makeDiff([]int{80, 443}, nil))
	out := buf.String()
	if !strings.Contains(out, "Opened ports (2)") {
		t.Errorf("expected opened ports summary, got: %s", out)
	}
	if !strings.Contains(out, "+ 80/tcp") {
		t.Errorf("expected port 80 in output, got: %s", out)
	}
}

func TestReport_ClosedPorts(t *testing.T) {
	var buf bytes.Buffer
	r, _ := reporter.New(&buf)
	r.Report(makeDiff(nil, []int{22}))
	out := buf.String()
	if !strings.Contains(out, "Closed ports (1)") {
		t.Errorf("expected closed ports summary, got: %s", out)
	}
	if !strings.Contains(out, "- 22/tcp") {
		t.Errorf("expected port 22 in output, got: %s", out)
	}
}

func TestNew_NilWriter(t *testing.T) {
	_, err := reporter.New(nil)
	if err == nil {
		t.Error("expected error for nil writer")
	}
}
