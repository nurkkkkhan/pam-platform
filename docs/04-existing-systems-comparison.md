# 4. Existing System Comparison

## 4.1 Selection of systems

Five systems were selected: three commercial **Leaders** of the 2025 Gartner Magic Quadrant for PAM [16] and two widely used **open-source / open-core** tools that a cost-sensitive organization would realistically consider.

| System | Vendor / license | Why selected |
|---|---|---|
| **CyberArk Privileged Access Manager** (Self-Hosted / Privilege Cloud) | CyberArk (Palo Alto Networks since 2026), commercial | Market leader, reference feature set |
| **BeyondTrust Password Safe + Privileged Remote Access** | BeyondTrust, commercial | Leader, strong session management |
| **Delinea Secret Server** | Delinea, commercial | Leader, popular in mid-market |
| **Teleport Community Edition** | Gravitational, AGPL-3.0 code; free binaries only for orgs < 100 employees and < US$10M revenue | Modern certificate-based infrastructure access |
| **JumpServer** | FIT2CLOUD, GPL-3.0 + paid Enterprise | Open-source bastion with vault, popular in Asia |

Excluded: HashiCorp Boundary (BSL licence since 2023, session recording only in paid tiers, depends on Vault) and Apache Guacamole (remote-desktop gateway, no vault/approval) — both are mentioned as components rather than complete PAM systems.

## 4.2 Criteria

| # | Criterion | Why it matters (link to evidence) |
|---|---|---|
| C1 | Credential vault with automatic rotation | E1, E4 |
| C2 | Brokered access without revealing credentials | E1 |
| C3 | Just-in-time access with approval workflow | E2, E7 |
| C4 | Session recording & replay (SSH + DB) | E3, A7 |
| C5 | Real-time command control (alert/block/terminate) | E3, E5 |
| C6 | Tamper-evident audit log | E6, [3] |
| C7 | SSO (OIDC/SAML) + MFA in the free/base tier | E7 |
| C8 | Self-hosted / on-prem deployment | A4 |
| C9 | Licence cost and openness for an SME | A6 |
| C10 | Deployment complexity (time-to-first-session) | A6 |

## 4.3 Comparison table

Legend: ✅ supported · ◐ partial / paid tier only · ❌ not supported · ? not confirmed in public documentation

| Criterion | CyberArk PAM | BeyondTrust | Delinea Secret Server | Teleport Community | JumpServer (OSS) | **PAM Platform (target)** |
|---|---|---|---|---|---|---|
| C1 Vault + rotation | ✅ | ✅ | ✅ | ❌ (uses short-lived certs instead of vaulting passwords) | ✅ (basic) | ✅ SSH keys & passwords, PostgreSQL |
| C2 Brokered, no credential disclosure | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| C3 JIT + approval | ✅ | ✅ | ✅ | ◐ CLI role requests only; reviews/auto-approval in Enterprise | ◐ ticket approval (more in Enterprise) | ✅ web approval, TTL, auto-revoke |
| C4 Session recording SSH + DB | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| C5 Real-time command control | ✅ | ✅ | ✅ | ◐ (enhanced recording, no inline blocking in CE) | ✅ command filters | ✅ alert/block/terminate rules |
| C6 Tamper-evident audit | ✅ (vault audit) | ? | ? | ◐ (audit events, no hash chain documented) | ❌ | ✅ hash chain + offline verifier |
| C7 SSO + MFA in base tier | ✅ | ✅ | ✅ | ◐ GitHub SSO only in CE | ✅ | ✅ OIDC + TOTP |
| C8 Self-hosted | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| C9 Licence/cost for SME | ❌ quote-based enterprise licensing | ❌ quote-based | ◐ quote-based, mid-market editions | ◐ free only below size thresholds; AGPL | ✅ GPL-3.0 | ✅ open source (licence TBD with supervisor, e.g. Apache-2.0) |
| C10 Deployment complexity | High (multiple Windows components, professional services) | High | Medium | Low–medium | Medium | Low (single `docker compose up`) — **to be measured** |

### Evidence for the claims

- Teleport Community limits (GitHub-only SSO, CLI-only access requests, free use under 100 employees / US$10M revenue): Teleport Feature Matrix, https://goteleport.com/docs/ver/17.x/feature-matrix
- JumpServer, Boundary, Guacamole features and licences: overview at https://startwithidentity.com/articles/top-7-open-source-pam-solutions.md and project repositories (https://github.com/jumpserver/jumpserver, https://github.com/hashicorp/boundary).
- Commercial vendors' positioning: Gartner MQ for PAM 2025 [16]; vendor product documentation (docs.cyberark.com, docs.beyondtrust.com, docs.delinea.com).
- Cells marked **?** must be checked against vendor documentation or a trial before the pre-defense; do not present them as facts.

## 4.4 Conclusion — why the proposed platform is needed

1. Commercial leaders cover all criteria but have quote-based enterprise pricing and high deployment complexity (C9, C10), which is a barrier for SMEs and educational/public organizations (assumption A6, to be confirmed in interviews).
2. Open-source tools each miss something important: Teleport CE lacks password vaulting and web approval workflows and is restricted by company size; JumpServer lacks verifiable (tamper-evident) audit; Boundary is no longer OSI open source and records sessions only in paid tiers.
3. **No compared open tool combines** vault + rotation, web JIT approval, inline command control and cryptographically verifiable audit. PAM Platform targets exactly this combination for SSH and PostgreSQL, with a simple deployment.

The proposed solution does **not** claim to be better than commercial suites in breadth (RDP, AD discovery, endpoint privilege management are out of scope).
