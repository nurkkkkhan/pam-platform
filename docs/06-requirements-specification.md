# 6. Requirements Specification (initial, v0.1)

Priorities use **MoSCoW**: **M** = Must (MVP, mandatory), **S** = Should, **C** = Could, **W** = Won't (this release). Every **Must** requirement has an acceptance criterion (AC) that can be checked by a test (see test IDs in `09-verification-and-evaluation-plan.md`). "Source" links each requirement to evidence (E*/A* in `02-evidence-of-relevance.md`) or literature [n].

## 6.1 Functional requirements

### Users & authentication (FR-U)

| ID | Requirement | Pri | Acceptance criterion | Source | Test |
|---|---|---|---|---|---|
| FR-U1 | Users authenticate via corporate SSO (OIDC). | M | Login via Keycloak succeeds; user without valid OIDC token gets HTTP 401 on any API call. | E7, [6] | T-U1 |
| FR-U2 | Starting a privileged session requires TOTP MFA (step-up) not older than 15 min. | M | Session start without valid TOTP is refused; with valid TOTP it starts; refusal is logged. | [6][15] | T-U2 |
| FR-U3 | Break-glass local admin account with MFA, every use triggers an alert. | S | Login as break-glass sends alert to security channel within 60 s. | [15] | T-U3 |
| FR-U4 | Users and groups synchronized from IdP claims (group → role mapping). | S | Adding user to IdP group `pam-dba` gives role `DBA` on next login. | — | T-U4 |

### Vault (FR-V)

| ID | Requirement | Pri | Acceptance criterion | Source | Test |
|---|---|---|---|---|---|
| FR-V1 | Store privileged credentials (password, SSH private key) for target accounts encrypted with AES-256-GCM envelope encryption. | M | DB dump contains no plaintext secret (automated grep of known test secrets returns 0 hits); each secret has its own DEK. | E1, [12] | T-V1 |
| FR-V2 | Automatic rotation of credentials on schedule (configurable, default 24 h) for Linux accounts and PostgreSQL roles. | M | After scheduled run, old password no longer authenticates on target; new one works via gateway. | E1, [15] IA-5 | T-V2 |
| FR-V3 | Rotation after every check-out or session end (configurable per account). | M | After session end, credential version increments and old credential is rejected by target. | E1 | T-V3 |
| FR-V4 | Controlled password check-out (show credential) only for accounts with `allow_checkout`, requires approval and is time-limited. | S | Check-out displays secret once, expires after TTL, triggers rotation, fully audited. | — | T-V4 |
| FR-V5 | Secrets API for applications (machine identities with API tokens). | C | Service retrieves a DB password via API using a scoped token; access audited. | E4, [12] | T-V5 |
| FR-V6 | Master key (KEK) rotation by re-wrapping DEKs without downtime. | S | After KEK rotation all secrets decrypt; no session interrupted. | — | T-V6 |

### Access requests & policies (FR-A)

| ID | Requirement | Pri | Acceptance criterion | Source | Test |
|---|---|---|---|---|---|
| FR-A1 | Admin defines targets (host/DB), target groups, privileged accounts and roles (RBAC eligibility). | M | CRUD via UI and API; user sees only targets their role is eligible for. | [2] | T-A1 |
| FR-A2 | User submits a JIT access request: target(s), account, reason, duration (≤ max TTL of policy). | M | Request with duration > max TTL is rejected with validation error. | E2, [5] | T-A2 |
| FR-A3 | Approval workflow: N approvers (policy-defined, 1–2), requester cannot approve own request. | M | Self-approval returns 403; grant becomes active only after required approvals. | E2 | T-A3 |
| FR-A4 | Auto-approval rules for low-risk targets (e.g. dev environment, working hours). | S | Request matching rule is granted without human approval and marked `auto`. | — | T-A4 |
| FR-A5 | Grants expire automatically; active sessions are terminated at grant expiry (with 5-min warning). | M | At T_expiry + 30 s no session for that grant exists; warning shown at T−5 min. | E2, [1] | T-A5 |
| FR-A6 | Manual revocation of a grant terminates its sessions immediately. | M | Session closes ≤ 5 s after revoke. | E3 | T-A6 |
| FR-A7 | Contextual conditions on policies: allowed time window, source IP range. | S | Request outside window is denied with reason. | [18] | T-A7 |
| FR-A8 | Notifications to approvers (e-mail, Telegram/webhook). | S | Approver receives notification ≤ 60 s after request. | — | T-A8 |

### Gateway / brokered access (FR-G)

| ID | Requirement | Pri | Acceptance criterion | Source | Test |
|---|---|---|---|---|---|
| FR-G1 | SSH gateway: user connects with standard OpenSSH client (`ssh user%account%target@gw`); gateway authenticates the user (short-lived cert or key bound to PAM identity) and connects to target injecting vaulted credential. | M | Interactive shell and `ssh … command` work; user never receives target password/key (verified by inspecting client side traffic & UI). | E1, [5] | T-G1 |
| FR-G2 | PostgreSQL gateway: user connects with `psql`/DBeaver to gateway; gateway authenticates the user with one-time token and injects DB role credential. | M | `psql` session works through gateway; target role password not known to user. | E1 | T-G2 |
| FR-G3 | Web terminal (SSH in browser) from the UI. | S | User opens shell in browser for an active grant. | — | T-G3 |
| FR-G4 | All access to targets is only possible through gateway (targets accept privileged logins only from gateway IP / keys). | M | Direct SSH from admin workstation to target with privileged account fails in test bench. | [1] complete mediation | T-G4 |
| FR-G5 | SFTP/SCP file transfers through SSH gateway are logged (file names, sizes). | C | Transfer recorded in audit with filename and size. | — | T-G5 |

### Monitoring & session control (FR-M)

| ID | Requirement | Pri | Acceptance criterion | Source | Test |
|---|---|---|---|---|---|
| FR-M1 | Record every SSH session (asciicast v2, input + output, timing) and every PostgreSQL session (queries, timestamps, rows affected). | M | 100% of sessions in test run have a recording; recording replays identically. | E3, A7 | T-M1 |
| FR-M2 | Session replay in UI with search by command text and jump-to-time. | M | Auditor finds command `rm -rf` in recording and jumps to its timestamp. | A7 | T-M2 |
| FR-M3 | Live session monitoring by security officer; manual terminate. | M | Officer sees live output ≤ 2 s delay; terminate closes session ≤ 5 s. | E3 | T-M3 |
| FR-M4 | Policy rules on commands/queries (regex): `alert`, `block` (command not sent to target), `terminate`. | M | `DROP DATABASE` on prod is blocked and alert created; rule hit recorded. | E3, E5, [13] | T-M4 |
| FR-M5 | Off-hours and unusual-target alerts. | S | Session started outside policy hours generates alert. | E5, [13] | T-M5 |
| FR-M6 | Recordings stored with object lock; retention configurable (default 365 days). | S | Attempt to delete recording before retention returns error. | [15] AU-9 | T-M6 |

### Audit & reporting (FR-AU)

| ID | Requirement | Pri | Acceptance criterion | Source | Test |
|---|---|---|---|---|---|
| FR-AU1 | Every security-relevant event (login, request, approval, grant, session start/end, rule hit, vault access, config change) is written to the audit log. | M | Test scenario produces expected list of event types; none missing. | [15] AU-2 | T-AU1 |
| FR-AU2 | Audit log is hash-chained; a verifier tool detects any modification, deletion or reordering. | M | Modifying any single row in DB → verifier reports exact broken index; untouched log verifies OK. | [3][4] | T-AU2 |
| FR-AU3 | Audit search & filter (user, target, time, event type); export CSV/JSON. | M | Search of 100k events returns ≤ 2 s; export opens in Excel. | E6 | T-AU3 |
| FR-AU4 | Reports: who has access to what; privileged sessions per user/target; rotation compliance. | S | Report matches DB state for seeded data. | E6 | T-AU4 |
| FR-AU5 | Forwarding audit events to SIEM (syslog/JSON over HTTPS). | C | Events visible in test syslog collector. | — | T-AU5 |

### Administration (FR-AD)

| ID | Requirement | Pri | Acceptance criterion | Source | Test |
|---|---|---|---|---|---|
| FR-AD1 | REST API (OpenAPI 3) for all admin operations. | M | OpenAPI spec validates; API tests cover CRUD for targets/accounts/policies. | — | T-AD1 |
| FR-AD2 | Gateway registration with one-time join token and mTLS certificate issuance. | M | Gateway without valid cert cannot connect to control plane. | [5] | T-AD2 |
| FR-AD3 | UI localization EN/RU (KZ could). | S | All UI strings switch language. | A* (interviews) | T-AD3 |

## 6.2 Non-functional requirements

| ID | Category | Requirement | Pri | Acceptance criterion / metric | Test |
|---|---|---|---|---|---|
| NFR-S1 | Security | No plaintext privileged credential at rest or in logs. | M | Secret-scanning of DB dump, logs and recordings: 0 findings. | T-S1 |
| NFR-S2 | Security | All network traffic encrypted: TLS 1.3 for UI/API, mTLS for internal gRPC, SSH for SSH. | M | `testssl.sh` shows no TLS < 1.2; internal port refuses non-mTLS. | T-S2 |
| NFR-S3 | Security | Platform passes OWASP ASVS L2 checklist for web part; no High/Critical from ZAP baseline scan, `gosec`, Trivy. | M | CI reports 0 High/Critical. | T-S3 |
| NFR-S4 | Security | STRIDE threat model of the platform maintained; each high threat has a mitigation. | M | Threat model document reviewed by supervisor. | T-S4 |
| NFR-S5 | Security | Platform admins cannot alter audit without detection (separation of duties). | M | See FR-AU2. | T-AU2 |
| NFR-P1 | Performance | Added connection setup latency via gateway ≤ 1 s (p95) vs direct SSH. | M | Benchmark 100 connections. | T-P1 |
| NFR-P2 | Performance | ≥ 100 concurrent recorded SSH sessions per gateway (2 vCPU / 4 GB) with keystroke latency overhead ≤ 50 ms (p95). | M | Load test. | T-P2 |
| NFR-P3 | Performance | API p95 latency ≤ 300 ms at 50 RPS. | S | k6 load test. | T-P3 |
| NFR-R1 | Reliability | If control plane is unavailable, new sessions are denied (fail-closed); existing sessions continue until grant expiry. | M | Stop API container → new session refused, existing continues. | T-R1 |
| NFR-R2 | Reliability | Backup & restore of DB + KEK procedure documented and tested. | S | Restore on clean host; secrets decrypt. | T-R2 |
| NFR-U1 | Usability | SUS score ≥ 70 from ≥ 5 pilot users. | S | SUS questionnaire. | T-UX1 |
| NFR-U2 | Usability | Admin connects to a target in ≤ 3 steps using existing SSH client. | M | Usability walkthrough. | T-UX2 |
| NFR-M1 | Maintainability | Unit test coverage ≥ 70% for `vault`, `audit`, `policy` packages. | M | `go test -cover`. | T-Q1 |
| NFR-D1 | Deployability | Full stack starts with one `docker compose up` in ≤ 10 min on clean VM. | M | Timed install by a person outside the team. | T-D1 |
| NFR-C1 | Compliance | Mapping of features to NIST SP 800-53 AC-2, AC-6, AU-2, AU-9, IA-2, IA-5 (and local regulation if required by client). | S | Mapping table in thesis. | — |

## 6.3 Constraints

- Self-hosted, Linux x86-64; no dependency on paid SaaS.
- Open-source dependencies with permissive or weak-copyleft licences only.
- No real production credentials or personal data are used in development or tests (synthetic test bench).

## 6.4 Traceability (summary)

| Objective (proposal §1.4) | Requirements |
|---|---|
| O5 Vault & rotation | FR-V1…V6, NFR-S1 |
| O6 Gateway | FR-G1…G5, NFR-P1, NFR-P2 |
| O7 JIT | FR-A1…A8 |
| O8 Recording & control | FR-M1…M6 |
| O9 Audit | FR-AU1…AU5, NFR-S5 |
| O10 Evaluation | NFR-*, `09-verification-and-evaluation-plan.md` |
