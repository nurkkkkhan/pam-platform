// Command pam-demo-audit prints a small synthetic audit export (JSON Lines)
// used to demonstrate the offline verifier.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/nurkkkkhan/pam-platform/internal/audit"
)

func main() {
	c := audit.NewChain()
	t := time.Date(2026, 10, 24, 13, 58, 44, 0, time.FixedZone("ALMT", 5*3600))
	steps := []struct {
		actor, action, object string
		details               any
	}{
		{"alice", "request.created", "REQ-2087", map[string]any{"target": "db-prod-01", "account": "root", "duration_min": 120}},
		{"s.nurlanova", "request.approved", "REQ-2087", map[string]any{"grant": "G-551"}},
		{"alice", "session.started", "session/S-9f31", map[string]any{"gateway": "gw-1"}},
		{"alice", "rule.blocked", "session/S-9f31", map[string]any{"rule": "R-3", "command": "rm -rf /var/lib/postgresql"}},
		{"gateway-1", "session.ended", "session/S-9f31", map[string]any{"recording_sha256": "a91c7e02"}},
		{"rotation", "secret.rotated", "root@db-prod-01", map[string]any{"version": 42}},
	}
	enc := json.NewEncoder(os.Stdout)
	for i, s := range steps {
		e, err := c.Append(t.Add(time.Duration(i)*time.Minute), s.actor, s.action, s.object, s.details)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		_ = enc.Encode(e)
	}
}
