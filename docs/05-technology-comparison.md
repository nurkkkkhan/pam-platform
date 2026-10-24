# 5. Technology Comparison

Each decision is scored against criteria derived from the requirements (`06-requirements-specification.md`). Scores: 1 (poor) – 5 (excellent). Weights reflect importance for this project. Detailed decisions are recorded as ADRs in `docs/adr/`.

## 5.1 Decision 1 — Server platform (backend + gateway language)

| Criterion (weight) | **Go 1.22+** | Java 21 + Spring Boot | Python 3.12 + FastAPI |
|---|---|---|---|
| Mature SSH / network protocol libraries (×3) | 5 — `golang.org/x/crypto/ssh` supports server and client side, used by Teleport/Boundary | 4 — Apache MINA SSHD | 3 — Paramiko/asyncssh, slower |
| Concurrency for many long-lived sessions (×3) | 5 — goroutines, low memory per session | 4 — virtual threads (Java 21) | 3 — asyncio, GIL limits CPU work |
| Performance / latency overhead (×2) | 5 | 4 | 2 |
| Deployment simplicity (×2) | 5 — single static binary | 3 — JVM, larger images | 3 — interpreter + deps |
| Team experience (×1) | Not assessed | Not assessed | Not assessed |
| Ecosystem for web API (×1) | 4 | 5 | 5 |
| **Weighted total (excl. team experience)** | **54** | 43 | 33 |

**Decision:** Go for the control-plane API and the access gateway. Justification: NFR-P1 (connection overhead ≤ 1 s), NFR-P2 (≥ 100 concurrent sessions) and the need for SSH server+client in one process. Risk: team experience with Go — mitigated by training tasks in sprint 1 (see risks R3).

## 5.2 Decision 2 — Data storage

### 2a. Primary database

| Criterion (weight) | **PostgreSQL 16** | MySQL 8 | MongoDB 7 |
|---|---|---|---|
| ACID transactions for grants/approvals (×3) | 5 | 5 | 4 |
| Append-only audit table enforcement (triggers, REVOKE UPDATE/DELETE) (×2) | 5 | 3 | 2 |
| Semi-structured audit details (JSONB / indexes) (×2) | 5 | 3 | 5 |
| Row-level security (×1) | 5 | 1 | 2 |
| Operational familiarity in target orgs (×1) | 5 | 5 | 3 |
| **Weighted total** | **45** | 33 | 31 |

### 2b. Session recordings

| Criterion (weight) | **S3-compatible object storage (MinIO)** | Local filesystem | Database BLOB |
|---|---|---|---|
| Scales to GBs of recordings (×3) | 5 | 3 | 1 |
| Immutability (object lock / WORM) (×2) | 5 | 2 | 2 |
| Simplicity (×1) | 3 | 5 | 4 |
| **Total** | **28** | 18 | 11 |

**Decision:** PostgreSQL for metadata, policies, grants and the hash-chained audit log; MinIO (S3 API) for recordings with object lock. The S3 API keeps cloud storage possible later.

## 5.3 Decision 3 — Secret protection (vault encryption)

| Criterion (weight) | **Built-in envelope encryption (AES-256-GCM, per-secret DEK, KEK from KMS/file)** | External HashiCorp Vault / OpenBao | pgcrypto column encryption |
|---|---|---|---|
| Keys never stored next to ciphertext (×3) | 5 — KEK outside DB | 5 | 2 — key passed in SQL |
| Key rotation without re-encrypting all data (×2) | 5 — re-wrap DEKs only | 5 | 2 |
| Operational dependencies (×2) | 5 — none | 2 — extra HA cluster to run & unseal | 5 |
| Licence (×1) | 5 | 3 — Vault BSL; OpenBao MPL-2.0 | 5 |
| Demonstrates engineering work (×1) | 5 | 2 | 3 |
| **Total** | **45** | 34 | 28 |

**Decision:** built-in envelope encryption behind a `KeyProvider` interface (file/env KEK for the test bench; interface allows HSM/cloud KMS/OpenBao Transit later). Reference implementation: `src/internal/vault/envelope.go`.

## 5.4 Decision 4 — Method of component interaction

| Criterion (weight) | REST/JSON only | **REST (UI ↔ API) + gRPC with mTLS (control plane ↔ gateways)** | Message broker (RabbitMQ/NATS) for all |
|---|---|---|---|
| Browser compatibility for UI (×2) | 5 | 5 (REST part) | 2 |
| Strongly typed contract for gateway authorization calls (×2) | 3 | 5 — protobuf | 3 |
| Low latency per-session authorization (×2) | 3 | 5 — HTTP/2, persistent streams | 3 |
| Streaming live session to observers (×2) | 2 — needs WebSocket add-on | 4 — gRPC streams + WebSocket to browser | 5 |
| Operational simplicity (×1) | 5 | 4 | 2 |
| **Total** | 31 | **42** | 28 |

**Decision:** REST + OpenAPI 3 for the web UI and external integrations; gRPC over mutual TLS between the control plane and gateways (authorization checks, credential injection, event streaming); WebSocket from API to browser for live session view. A broker is deferred until multi-node HA is needed.

## 5.5 Decision 5 — Web frontend

| Criterion (weight) | **React 18 + TypeScript + Vite** | Angular 17 | Vue 3 |
|---|---|---|---|
| Terminal player component availability (xterm.js, asciinema-player) (×2) | 5 | 4 | 4 |
| Team experience (×2) | Not assessed | Not assessed | Not assessed |
| Ecosystem / UI kits (×1) | 5 | 4 | 4 |
| Type safety (×1) | 5 | 5 | 4 |

**Provisional decision:** React + TypeScript. The two members should assess their actual experience with all three options before confirming this choice.

## 5.6 Decision 6 — User authentication

| Criterion | **OIDC to corporate IdP (Keycloak in test bench) + TOTP MFA** | Direct LDAP/AD bind | Local accounts only |
|---|---|---|---|
| Fits SSO in corporations | ✅ | ◐ | ❌ |
| MFA support | ✅ (IdP or built-in TOTP) | ❌ needs extra | built-in only |
| Password never handled by PAM | ✅ | ❌ | ❌ |

**Decision:** OIDC (Authorization Code + PKCE); PAM Platform additionally enforces step-up TOTP before privileged session start (FR-U2). A local break-glass admin account exists with mandatory MFA.

## 5.7 Summary of selected stack

| Layer | Choice |
|---|---|
| Backend / gateway | Go |
| Primary DB | PostgreSQL 16 |
| Recording storage | MinIO (S3, object lock) |
| Secrets | Built-in envelope encryption, AES-256-GCM, pluggable KEK provider |
| Interaction | REST/OpenAPI (UI), gRPC + mTLS (internal), WebSocket (live view) |
| Frontend | React + TypeScript |
| Identity | OIDC (Keycloak) + TOTP |
| Deployment | Docker Compose (dev/test), Helm chart (stretch) |
| CI | GitHub Actions / GitLab CI: `go test`, `golangci-lint`, `gosec`, `npm test`, Trivy image scan |
