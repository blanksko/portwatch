package policy_test

import (
	"testing"

	"github.com/user/portwatch/internal/policy"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/snapshot"
)

func makeDiff(opened, closed []int) snapshot.Diff {
	d := snapshot.Diff{}
	for _, p := range opened {
		d.Opened = append(d.Opened, scanner.Result{Port: p})
	}
	for _, p := range closed {
		d.Closed = append(d.Closed, scanner.Result{Port: p})
	}
	return d
}

func TestAllow_DefaultAllowsAll(t *testing.T) {
	p := policy.Default()
	for _, port := range []int{22, 80, 443, 8080} {
		if !p.Allow(port) {
			t.Errorf("default policy should allow port %d", port)
		}
	}
}

func TestAllow_DenyRule(t *testing.T) {
	p := policy.New([]policy.Rule{
		{Ports: []int{22}, Action: "deny"},
	})
	if p.Allow(22) {
		t.Error("expected port 22 to be denied")
	}
	if !p.Allow(80) {
		t.Error("expected port 80 to be allowed by default")
	}
}

func TestAllow_CatchAllDeny(t *testing.T) {
	p := policy.New([]policy.Rule{
		{Ports: []int{443}, Action: "allow"},
		{Action: "deny"},
	})
	if !p.Allow(443) {
		t.Error("expected 443 to be allowed")
	}
	if p.Allow(80) {
		t.Error("expected 80 to be denied by catch-all")
	}
}

func TestFilter_RemovesDeniedPorts(t *testing.T) {
	p := policy.New([]policy.Rule{
		{Ports: []int{22}, Action: "deny"},
	})
	d := makeDiff([]int{22, 80}, []int{22, 443})
	out := p.Filter(d)
	if len(out.Opened) != 1 || out.Opened[0].Port != 80 {
		t.Errorf("unexpected opened ports: %v", out.Opened)
	}
	if len(out.Closed) != 1 || out.Closed[0].Port != 443 {
		t.Errorf("unexpected closed ports: %v", out.Closed)
	}
}

func TestFilter_DefaultAllowsAll(t *testing.T) {
	p := policy.Default()
	d := makeDiff([]int{80, 443}, []int{22})
	out := p.Filter(d)
	if len(out.Opened) != 2 || len(out.Closed) != 1 {
		t.Errorf("default policy should not filter any ports")
	}
}
