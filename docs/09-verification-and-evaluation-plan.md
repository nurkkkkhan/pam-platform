# 9. Verification and Evaluation Plan

## 9.1 How the aim will be demonstrated

The aim is achieved if, on a realistic test bench, PAM Platform (1) eliminates credential exposure and standing access, (2) records and controls 100% of privileged sessions, (3) produces verifiable audit evidence, and (4) does so with acceptable performance and usability — **compared to the baseline**.

## 9.2 Baseline solution for comparison

**Baseline B0 — "typical current practice":** admins connect directly with OpenSSH / psql using shared passwords and keys stored in a shared password spreadsheet or KeePass file; access is granted by adding SSH keys manually; logging = default `auth.log` / PostgreSQL logs. (To be confirmed as typical practice by interviews, assumption A1.)

**Baseline B1 (optional) — open-source alternative:** JumpServer (OSS) or Teleport Community on the same test bench, for selected criteria (setup time, connection latency, audit tamper detection).

## 9.3 Test bench (measurement conditions)

| Item | Specification |
|---|---|
| Hosts | 6 VMs (Proxmox / VirtualBox / cloud): `pam-cp` (2 vCPU, 4 GB), `pam-gw` (2 vCPU, 4 GB), `target-linux-1..3` (1 vCPU, 1 GB, Ubuntu 24.04), `target-pg-1` (PostgreSQL 16) |
| Network | Isolated LAN, 1 Gbit/s; targets accept privileged logins only from `pam-gw` |
| Clients | OpenSSH 9.x, psql 16, Chrome latest |
| Load tools | Custom Go load generator (SSH sessions with scripted commands), k6 for API |
| Data | Synthetic: 50 users, 200 targets, 500 accounts, 100k audit events |
| Repetitions | Each performance measurement repeated 5 times; report median and p95 |

## 9.4 Testing plan

| Level | What | Tools | When |
|---|---|---|---|
| Unit | vault crypto, audit chain, policy evaluation, rule matching | `go test`, coverage | every commit (CI) |
| Integration | API + DB + gateway + MinIO; rotation against real sshd/PostgreSQL containers | Testcontainers-go, Docker | every PR |
| End-to-end | Scenarios T-* below via SSH client and UI | Go e2e suite, Playwright | before each milestone |
| Security | SAST, dependency & image scan, DAST, manual pentest of own platform, audit tampering tests | gosec, golangci-lint, Trivy, OWASP ZAP | CI + before final pre-defense |
| Performance | NFR-P1…P3 | Go load generator, k6 | milestone M4 |
| Usability | SUS + task completion with 5 pilot users | questionnaire | milestone M5 |

### Key test cases (examples, full list in `tests/test-cases.md`)

| Test | Requirement | Steps | Expected |
|---|---|---|---|
| T-V1 | FR-V1, NFR-S1 | Seed 20 known secrets; `pg_dump`; grep for each secret | 0 matches |
| T-V3 | FR-V3 | Session on account with rotate_after_use; end; try old password directly | old rejected, version +1 |
| T-A3 | FR-A3 | User A requests; user A calls approve | HTTP 403, audit `approval.denied_sod` |
| T-A5 | FR-A5 | Grant 10 min; keep session open | warning at T−5, closed by T+30 s |
| T-G4 | FR-G4 | Direct `ssh root@target` from workstation | connection refused / auth fails |
| T-M4 | FR-M4 | Run `DROP DATABASE billing;` via PG gateway on prod | blocked, alert, rule hit logged |
| T-AU2 | FR-AU2 | Save the head hash outside the audit store; update 1 random audit row / delete 1 row (including the tail) / swap 2 rows; run verifier with that anchor | each tampering detected at correct seq or by anchor mismatch |
| T-R1 | NFR-R1 | Stop API; try new session; check existing | new denied, existing alive |

## 9.5 Outcome metrics

| # | Metric | Baseline B0 (expected) | Target for PAM Platform | How measured |
|---|---|---|---|---|
| M1 | Privileged credentials known to humans | all shared accounts | 0 (except approved check-outs) | Count of accounts with `checkout` events / users who saw a secret |
| M2 | Standing privileged access (accounts usable without active grant) | 100% | 0% | Attempt access without grant for all accounts |
| M3 | Share of privileged sessions recorded with attribution to a person | ~0% (shared login) | 100% | Sessions in gateway vs. sessions in target `auth.log` |
| M4 | Time to revoke access of a departing contractor | manual, minutes–days | ≤ 5 s (grant revoke) | Timed test T-A6 |
| M5 | Mean credential age | months (to confirm in interviews) | ≤ 24 h | Vault statistics |
| M6 | Audit tampering detection rate | 0% (plain logs) | 100% of injected modifications | T-AU2 with 100 random tamperings and an independently stored head hash |
| M7 | Time to answer "who ran command X on host Y" | hours (manual log search) | ≤ 1 min | Task timed with auditor role |
| M8 | Connection setup overhead | 0 (direct) | ≤ 1 s p95 | T-P1 |
| M9 | Keystroke latency overhead | 0 | ≤ 50 ms p95 at 100 sessions | T-P2 |
| M10 | Usability (SUS) | — | ≥ 70 | T-UX1 |
| M11 | Deployment time on clean VM | — (B1: measure) | ≤ 10 min | T-D1 |

## 9.6 Evaluation report structure (for the thesis)

1. Test bench description · 2. Functional results (pass/fail per Must requirement) · 3. Security results (M1–M3, M6, scans) · 4. Performance (M8, M9, charts) · 5. Usability (M10, observations) · 6. Comparison with B0/B1 · 7. Threats to validity (synthetic data, small user sample, lab network).

## 9.7 Exit criteria for MVP

All **Must** requirements pass; no open High/Critical security findings; M1, M2, M3, M6 meet targets; M8, M9 meet targets or deviations are explained.
