package portname_test

import (
	"testing"

	"github.com/user/portwatch/internal/portname"
)

func TestLookup_KnownPort(t *testing.T) {
	cases := []struct {
		port int
		want string
	}{
		{22, "ssh"},
		{80, "http"},
		{443, "https"},
		{3306, "mysql"},
		{6379, "redis"},
	}
	for _, tc := range cases {
		got := portname.Lookup(tc.port)
		if got != tc.want {
			t.Errorf("Lookup(%d) = %q; want %q", tc.port, got, tc.want)
		}
	}
}

func TestLookup_UnknownPort(t *testing.T) {
	got := portname.Lookup(9999)
	if got != "" {
		t.Errorf("Lookup(9999) = %q; want empty string", got)
	}
}

func TestLookupWithDefault_Known(t *testing.T) {
	got := portname.LookupWithDefault(22, "unknown")
	if got != "ssh" {
		t.Errorf("LookupWithDefault(22) = %q; want %q", got, "ssh")
	}
}

func TestLookupWithDefault_Unknown(t *testing.T) {
	got := portname.LookupWithDefault(9999, "unknown")
	if got != "unknown" {
		t.Errorf("LookupWithDefault(9999) = %q; want %q", got, "unknown")
	}
}

func TestRegister_AddsMapping(t *testing.T) {
	portname.Register(12345, "myservice")
	got := portname.Lookup(12345)
	if got != "myservice" {
		t.Errorf("Lookup(12345) after Register = %q; want %q", got, "myservice")
	}
}

func TestRegister_OverwritesMapping(t *testing.T) {
	portname.Register(80, "custom-http")
	got := portname.Lookup(80)
	if got != "custom-http" {
		t.Errorf("Lookup(80) after overwrite = %q; want %q", got, "custom-http")
	}
	// restore
	portname.Register(80, "http")
}
