# Test cases (v0.1)
Linked to requirements in `docs/06-requirements-specification.md`. Status: ✅ automated & passing · ⏳ planned.

| Test | Req. | Level | Description | Status |
|---|---|---|---|---|
| T-AU2a | FR-AU2 | unit | Modified event detected at correct seq (`TestDetectModification`) | ✅ |
| T-AU2b | FR-AU2 | unit | Deleted event detected (`TestDetectDeletion`) | ✅ |
| T-AU2c | FR-AU2 | unit | Reordered events detected (`TestDetectReorder`) | ✅ |
| T-AU2d | FR-AU2 | unit | Re-hashed forgery detected at next link (`TestDetectRehashedForgery`) | ✅ |
| T-AU2e | FR-AU2 | unit | Tail truncation detected against anchor | ✅ |
| T-AU2f | FR-AU2 | CLI | `pam-audit verify` on `data/sample-audit.jsonl` → OK; on `data/tampered-audit.jsonl` → FAIL at seq 4 | ✅ |
| T-V1 | FR-V1 | unit | Plaintext absent from stored ciphertext/wrapped DEK | ✅ |
| T-V1b | FR-V1 | unit | Ciphertext bound to account+version (swapped row fails) | ✅ |
| T-V6 | FR-V6 | unit | KEK rotation by rewrap; old KEK cannot open | ✅ |
| T-A2 | FR-A2 | unit | Duration > max TTL rejected | ✅ |
| T-A3 | FR-A3 | unit | Self-approval denied | ✅ |
| T-U2 | FR-U2 | unit | Stale / missing MFA denies session | ✅ (logic) / ⏳ (e2e) |
| T-A5 | FR-A5 | unit | Expired / revoked grant denied | ✅ (logic) / ⏳ (e2e termination) |
| T-A7 | FR-A7 | unit | Time window incl. overnight window | ✅ |
| T-G1, T-G2, T-G4 | FR-G* | e2e | Brokered SSH/PG, no credential disclosure, direct access blocked | ⏳ M3/M4 |
| T-M1…T-M6 | FR-M* | e2e | Recording, replay, live, rules, retention | ⏳ M3/M4 |
| T-P1, T-P2, T-P3 | NFR-P* | perf | Latency and concurrency | ⏳ M4 |
| T-S1…T-S4 | NFR-S* | security | Secret scan, TLS, scanners, threat model | ⏳ M5 |
| T-UX1, T-UX2 | NFR-U* | usability | SUS, steps to connect | ⏳ M5 |

Run: `make test` and `make audit-demo` from the repository root.
