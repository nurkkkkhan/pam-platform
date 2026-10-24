// Package policy evaluates whether a user may open a privileged session
// (PDP role, see docs/07 §7.1). Requirements: FR-A2, FR-A3, FR-A5, FR-A7, FR-U2.
package policy

import (
	"errors"
	"time"
)

// Grant is an approved, time-limited permission.
type Grant struct {
	ID         string
	User       string
	TargetID   string
	Account    string
	ValidFrom  time.Time
	ValidUntil time.Time
	Revoked    bool
}

// Policy holds contextual conditions for a target group.
type Policy struct {
	MaxTTL        time.Duration
	MFAMaxAge     time.Duration
	AllowedHours  [2]int // [start,end) local hours; {0,24} = always
	RequiredVotes int
}

// SessionRequest is what the gateway asks the PDP.
type SessionRequest struct {
	User     string
	TargetID string
	Account  string
	MFAAt    time.Time
	Now      time.Time
}

var (
	ErrNoGrant      = errors.New("no active grant for user/target/account")
	ErrGrantRevoked = errors.New("grant revoked")
	ErrGrantExpired = errors.New("grant expired or not yet valid")
	ErrMFARequired  = errors.New("fresh MFA required")
	ErrOutsideHours = errors.New("outside allowed time window")
	ErrTTLTooLong   = errors.New("requested duration exceeds policy max TTL")
	ErrSelfApproval = errors.New("requester cannot approve own request")
)

// Decide returns the grant to use, or an error explaining the denial.
// It is fail-closed: anything unexpected denies.
func Decide(req SessionRequest, p Policy, grants []Grant) (*Grant, error) {
	if p.MFAMaxAge <= 0 || req.MFAAt.IsZero() || req.Now.Sub(req.MFAAt) > p.MFAMaxAge || req.MFAAt.After(req.Now) {
		return nil, ErrMFARequired
	}
	if !withinHours(req.Now, p.AllowedHours) {
		return nil, ErrOutsideHours
	}
	var lastErr error = ErrNoGrant
	for i := range grants {
		g := &grants[i]
		if g.User != req.User || g.TargetID != req.TargetID || g.Account != req.Account {
			continue
		}
		switch {
		case g.Revoked:
			lastErr = ErrGrantRevoked
		case req.Now.Before(g.ValidFrom) || !req.Now.Before(g.ValidUntil):
			lastErr = ErrGrantExpired
		default:
			return g, nil
		}
	}
	return nil, lastErr
}

// ValidateRequest checks a JIT request against policy (FR-A2).
func ValidateRequest(duration time.Duration, p Policy) error {
	if duration <= 0 || duration > p.MaxTTL {
		return ErrTTLTooLong
	}
	return nil
}

// CanApprove enforces separation of duties (FR-A3).
func CanApprove(requester, approver string) error {
	if requester == approver {
		return ErrSelfApproval
	}
	return nil
}

func withinHours(t time.Time, h [2]int) bool {
	if h[0] < 0 || h[0] >= 24 || h[1] < 0 || h[1] > 24 {
		return false
	}
	if h[0] == 0 && h[1] == 24 {
		return true
	}
	hr := t.Hour()
	if h[0] <= h[1] {
		return hr >= h[0] && hr < h[1]
	}
	return hr >= h[0] || hr < h[1] // window over midnight
}
