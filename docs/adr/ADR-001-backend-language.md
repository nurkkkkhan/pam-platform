# ADR-001: Go for control plane and gateway
Status: Accepted (pre-defense-1) · Task: DIP-7

## Context
Gateway must act as SSH server and client, hold ≥100 long-lived sessions, add ≤1 s setup latency (NFR-P1/P2). See comparison in docs/05 §5.1.
## Decision
Use Go 1.22+ for API, workers and gateway; single repository, separate binaries.
## Consequences
+ Native `x/crypto/ssh`, single static binaries, low memory per session.
− Team must learn Go (risk R3) → spike task DIP-14.
