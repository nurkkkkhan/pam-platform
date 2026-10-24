package policy

import (
	"errors"
	"testing"
	"time"
)

var (
	now = time.Date(2026, 10, 24, 14, 0, 0, 0, time.UTC)
	pol = Policy{MaxTTL: 4 * time.Hour, MFAMaxAge: 15 * time.Minute, AllowedHours: [2]int{0, 24}, RequiredVotes: 1}
	g   = Grant{ID: "G-1", User: "alice", TargetID: "db-prod-01", Account: "root",
		ValidFrom: now.Add(-time.Hour), ValidUntil: now.Add(time.Hour)}
	req = SessionRequest{User: "alice", TargetID: "db-prod-01", Account: "root", MFAAt: now.Add(-5 * time.Minute), Now: now}
)

func TestAllowWithActiveGrant(t *testing.T) {
	got, err := Decide(req, pol, []Grant{g})
	if err != nil || got.ID != "G-1" {
		t.Fatalf("want allow, got %v", err)
	}
}

func TestDenyCases(t *testing.T) {
	cases := []struct {
		name   string
		mutate func(*SessionRequest, *Policy, *Grant)
		want   error
	}{
		{"no grant for other user", func(r *SessionRequest, _ *Policy, _ *Grant) { r.User = "bob" }, ErrNoGrant},
		{"stale MFA", func(r *SessionRequest, _ *Policy, _ *Grant) { r.MFAAt = now.Add(-20 * time.Minute) }, ErrMFARequired},
		{"no MFA", func(r *SessionRequest, _ *Policy, _ *Grant) { r.MFAAt = time.Time{} }, ErrMFARequired},
		{"expired", func(_ *SessionRequest, _ *Policy, gr *Grant) { gr.ValidUntil = now }, ErrGrantExpired},
		{"revoked", func(_ *SessionRequest, _ *Policy, gr *Grant) { gr.Revoked = true }, ErrGrantRevoked},
		{"outside hours", func(_ *SessionRequest, p *Policy, _ *Grant) { p.AllowedHours = [2]int{9, 13} }, ErrOutsideHours},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r, p, gr := req, pol, g
			c.mutate(&r, &p, &gr)
			if _, err := Decide(r, p, []Grant{gr}); !errors.Is(err, c.want) {
				t.Fatalf("want %v, got %v", c.want, err)
			}
		})
	}
}

func TestOvernightWindow(t *testing.T) {
	p := pol
	p.AllowedHours = [2]int{22, 6}
	r := req
	r.Now = time.Date(2026, 10, 24, 23, 0, 0, 0, time.UTC)
	r.MFAAt = r.Now.Add(-time.Minute)
	gr := g
	gr.ValidUntil = r.Now.Add(time.Hour)
	if _, err := Decide(r, p, []Grant{gr}); err != nil {
		t.Fatalf("23:00 inside 22–06 window: %v", err)
	}
}

func TestInvalidAllowedHoursDeny(t *testing.T) {
	for _, hours := range [][2]int{{-1, 24}, {0, 25}, {24, 24}} {
		p := pol
		p.AllowedHours = hours
		if _, err := Decide(req, p, []Grant{g}); !errors.Is(err, ErrOutsideHours) {
			t.Errorf("hours %v: want ErrOutsideHours, got %v", hours, err)
		}
	}
}

func TestValidateRequestTTL(t *testing.T) {
	if err := ValidateRequest(5*time.Hour, pol); !errors.Is(err, ErrTTLTooLong) {
		t.Fatal("5h > max 4h must be rejected")
	}
	if err := ValidateRequest(2*time.Hour, pol); err != nil {
		t.Fatal(err)
	}
}

func TestSeparationOfDuties(t *testing.T) {
	if !errors.Is(CanApprove("alice", "alice"), ErrSelfApproval) {
		t.Fatal("self-approval must be denied")
	}
	if CanApprove("alice", "bob") != nil {
		t.Fatal("other approver allowed")
	}
}
