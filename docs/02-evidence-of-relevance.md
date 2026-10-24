# 2. Evidence of Relevance

> **Status: data collection in progress.** The committee checks *what evidence informs the requirements* and whether *verified findings are distinguished from assumptions*. This file therefore keeps three things separate: (A) verified secondary data with sources, (B) primary data from interviews — **to be filled with real interview notes only**, and (C) working assumptions that the interviews must confirm or reject.
>
> ⚠️ Do not invent interview results. If an interview has not happened yet, leave its row as "planned".

## 2.1 (A) Verified secondary evidence

| # | Finding | Source | Relevance to requirements |
|---|---|---|---|
| E1 | Credential abuse was the most common initial access vector — 22% of breaches | Verizon DBIR 2025 [14] | Vault + no credential disclosure to users (FR-V*, FR-G*) |
| E2 | Third-party involvement in breaches doubled from 15% to 30% | Verizon DBIR 2025 [14] | Time-limited contractor access (FR-A3, FR-A5) |
| E3 | Human element present in ~60% of breaches | Verizon DBIR 2025 [14] | Monitoring, command rules, session termination (FR-M*) |
| E4 | Hard-coded / shared secrets in repositories and artifacts are widespread; practitioners lack consistent secret-management practices | Basak et al., 2022 [12] | Central vault, API for applications (FR-V5) |
| E5 | Insider threats include *privilege exploitation*, data theft and sabotage; indicators include off-hours access and unexpected privilege changes | Singh & Sharma [13]; Yuan & Wu [10] | Off-hours rule, alerting, audit search (FR-M4, FR-AU3) |
| E6 | Least privilege and accountable use of privileged accounts are mandatory controls (AC-6, AC-2, AU-9) | NIST SP 800-53 Rev.5 [15] | Compliance-oriented reports (FR-AU4) |
| E7 | Zero-trust architectures require per-session, policy-based access instead of network-perimeter trust | NIST SP 800-207 [5]; Syed et al. [6] | JIT, per-session authorization (FR-A*) |
| E8 | Leading commercial PAM products (CyberArk, BeyondTrust, Delinea) are recognised as market Leaders but target enterprise budgets | Gartner MQ for PAM, Oct 2025 [16] (as reported by vendors) | Gap for SME, self-hosted alternative |

## 2.2 (B) Primary evidence — interviews with potential users

**Plan:** 3–5 semi-structured interviews (30–45 min) with people who administer or audit infrastructure: system administrators, DevOps engineers, information-security officers, an IT auditor. At least one interview with the client organization if a client is confirmed.

**Consent:** participants are told the purpose, that notes are anonymized (role + company size only), and that they can stop at any time. No personal data or real credentials are recorded.

| ID | Role | Org. type / size | Status | Notes file |
|---|---|---|---|---|
| INT-1 | `<e.g., Senior Linux administrator>` | `<e.g., bank, 1000+ employees>` | planned | `interviews/INT-1.md` |
| INT-2 | `<Information security officer>` | | planned | `interviews/INT-2.md` |
| INT-3 | `<DevOps engineer>` | | planned | `interviews/INT-3.md` |
| INT-4 | `<IT auditor / compliance>` | | planned | `interviews/INT-4.md` |
| INT-5 | `<Client representative>` | | planned | `interviews/INT-5.md` |

### Interview guide

1. **Context.** What infrastructure do you administer (number of servers, DBs, network devices, cloud)? Which protocols do you use daily (SSH, RDP, DB clients, web consoles)?
2. **Accounts.** How many people know the `root`/`Administrator`/DB superuser passwords? Where are they stored today?
3. **Granting access.** Describe the last time a new admin or contractor needed access. Who approved it, how long did it take, how was it revoked?
4. **Rotation.** How often are privileged passwords/keys changed? What happens when an admin leaves?
5. **Accountability.** If something broke on a server last night, how would you find out who did it and what commands were run?
6. **Audit / compliance.** What do auditors or regulators ask you to prove? How long does preparing that evidence take?
7. **Incidents.** Have you had incidents related to shared or leaked credentials (without confidential details)?
8. **Existing tools.** Have you used or evaluated any PAM product (CyberArk, BeyondTrust, Teleport, JumpServer…)? Why adopted / rejected? (cost, complexity, licensing, local support)
9. **Priorities.** Rank: vault & rotation · session recording · approval workflow · command blocking · reports · SSO/MFA.
10. **Constraints.** On-prem only? Air-gapped segments? Required languages (KZ/RU/EN)? Acceptable latency when connecting?

### Interview notes template (`interviews/INT-n.md`)

```markdown
# INT-n — <role>, <org type/size>
Duration: · Interviewer: · Format: (in person / online)

## Current practice
## Pain points (quotes, anonymized)
## Priorities (ranking from Q9)
## Constraints
## Confirms / rejects assumptions: A1 ✅/❌, A2 ...
## New requirements discovered
```

### Synthesis (fill after interviews)

| Theme | INT-1 | INT-2 | INT-3 | INT-4 | INT-5 | Conclusion |
|---|---|---|---|---|---|
| Shared passwords in use | | | | | | |
| Standing access for contractors | | | | | | |
| No command-level audit | | | | | | |
| SSH + PostgreSQL most used | | | | | | |
| Commercial PAM too expensive | | | | | | |

## 2.3 (C) Working assumptions to verify

| ID | Assumption | How verified | Status |
|---|---|---|---|
| A1 | Privileged credentials are shared among ≥3 people in typical SMEs | Interviews Q2 | ⏳ unverified |
| A2 | SSH and PostgreSQL cover most of the daily privileged sessions | Interviews Q1 | ⏳ unverified |
| A3 | Contractor access is revoked late or not at all | Interviews Q3 | ⏳ unverified |
| A4 | Organizations need on-prem / self-hosted deployment | Interviews Q10 | ⏳ unverified |
| A5 | Users accept ≤ 1 s extra connection time | Interviews Q10 + measurement | ⏳ unverified |
| A6 | Commercial PAM cost or complexity is a barrier for SMEs | Interviews Q8, vendor pricing | ⏳ unverified |
| A7 | Auditors need replayable session evidence | Interviews Q6 | ⏳ unverified |

When an assumption is confirmed or rejected, update this table **and** the linked requirement in `06-requirements-specification.md` (and record the change in a commit referencing the YouTrack task).
