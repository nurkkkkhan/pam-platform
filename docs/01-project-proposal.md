# 1. Project Proposal

| Field | Value |
|---|---|
| **Project title** | Development of a PAM platform for centralized management and control of privileged access in corporate infrastructure |
| **Working name** | PAM Platform |
| **Team** | Member 1 (name pending), Backend: architecture and API contracts, control plane, identity, policy, vault, gateways, audit, CI, security and performance. Member 2 (name pending), Frontend and UI: project coordination, research, requirements, prototype, React interface, usability and evaluation documents. Both members review and integrate each other's work. |
| **Supervisor** | Name and academic title to be supplied by the team; review has not been evidenced in this folder. |
| **Client / partner organization** | Not confirmed; do not present an organization as a project partner without agreement. |
| **Project type** | Software engineering project (security platform, web + network services) |

## 1.1 Problem statement

Corporate infrastructure (Linux servers, databases, network devices, cloud consoles) is administered through **privileged accounts**: `root`, `postgres`, `Administrator`, service accounts, and similar. In many small and medium organizations these accounts are managed informally:

- passwords and SSH keys are shared between administrators through chats, spreadsheets, or wiki pages;
- access is **standing** (permanent), so a former contractor or a compromised laptop keeps working credentials;
- there is no record of **who** used a shared account, **when**, and **what commands** they ran;
- credentials are rarely rotated because rotation is manual and risky;
- auditors (internal, ISO/IEC 27001, regulators) cannot get evidence of least privilege.

Industry data confirms this is the main attack path: in the Verizon 2025 DBIR, credential abuse was the most common initial access vector (22% of breaches), and breaches involving third parties doubled to 30% [14]. Commercial PAM suites (CyberArk, BeyondTrust, Delinea) solve the problem but are expensive and complex to deploy for SMEs, while open-source tools each cover only part of the problem (see `04-existing-systems-comparison.md`).

## 1.2 Target users

| User role | Need |
|---|---|
| **System / DB administrator** (primary) | Connect to servers and databases without knowing the password; request elevated access quickly |
| **Security officer / CISO** | Approve access, watch live sessions, get alerts on dangerous commands, terminate sessions |
| **Auditor / compliance** | Search and replay sessions; export tamper-evident audit evidence |
| **External contractor** | Time-limited access to specific hosts only |
| **Platform administrator** | Register targets, accounts, policies; integrate with corporate SSO |

## 1.3 Aim

To design and implement a self-hostable PAM platform that removes standing privileged access and shared passwords in a corporate infrastructure by brokering all privileged connections through a controlled gateway, storing credentials in an encrypted vault, granting access just-in-time with approval, and producing a complete, tamper-evident audit trail.

## 1.4 Objectives

1. Analyse privileged-access practices in target organizations (3–5 interviews) and the state of the art (≥12 sources).
2. Compare 4–5 existing PAM systems on a consistent set of criteria and identify the gap.
3. Specify functional and non-functional requirements with verifiable acceptance criteria.
4. Design the architecture (control plane, access gateway, vault, audit subsystem, web UI).
5. Implement an **encrypted credential vault** with automatic password/key rotation for Linux (SSH) and PostgreSQL accounts.
6. Implement an **access gateway** (SSH proxy and PostgreSQL proxy) that injects credentials so the user never sees them.
7. Implement **just-in-time (JIT) access requests** with approval workflow, time limits and automatic revocation.
8. Implement **session recording and replay**, command logging and rule-based alerting/termination for risky commands.
9. Implement a **tamper-evident (hash-chained) audit log** and audit export.
10. Evaluate the platform against a baseline (manual shared-credential administration) using the metrics in `09-verification-and-evaluation-plan.md`.

## 1.5 Expected outcomes

- Working MVP deployable with Docker Compose in a test corporate network (≥3 Linux hosts, 1 PostgreSQL server).
- Source code, automated tests, OpenAPI specification, deployment guide.
- Evaluation report: security metrics (0 plaintext credentials exposed to users, 100% sessions recorded), performance metrics (connection overhead, concurrent sessions), usability results (SUS score from pilot users).
- Diploma thesis text and defense materials.

## 1.6 Scope

**In scope (MVP)**

| Area | Included |
|---|---|
| Protocols | SSH (interactive shell, exec), PostgreSQL wire protocol |
| Vault | Passwords and SSH private keys; AES-256-GCM envelope encryption; check-out; scheduled and post-session rotation |
| Access model | RBAC (roles → target groups → accounts) + JIT requests with approval, TTL and auto-revocation |
| Authentication | Corporate SSO via OIDC (Keycloak in test bench) + mandatory TOTP MFA |
| Monitoring | Live session view, session recording (asciicast v2 for SSH; query log for PostgreSQL), replay, command/query policy rules (alert / block / terminate) |
| Audit | Hash-chained audit log, search, CSV/JSON export, verification tool |
| UI | Web console (React) for users, approvers, auditors, admins |
| Integration | REST API (OpenAPI 3), webhook/Telegram/e-mail notifications |

**Out of scope (possible future work)**

- RDP / VNC graphical sessions (could be added later via Apache Guacamole integration).
- Automatic discovery of privileged accounts in Active Directory / cloud IAM.
- ML-based user behaviour analytics (UEBA) — only rule-based detection in the MVP; ML is discussed in the literature review as future work.
- High availability clustering and hardware security module (HSM) integration — designed for (interfaces) but not implemented.
- Endpoint privilege management (local admin rights on workstations).

## 1.7 Engineering complexity

The project involves non-trivial engineering beyond CRUD:

- implementing **protocol-level proxies** (SSH man-in-the-middle by design with host-key management; PostgreSQL wire protocol parsing for authentication injection and query capture);
- **cryptography**: envelope encryption, key hierarchy, key rotation without downtime;
- **concurrency**: many simultaneous long-lived sessions with live streaming to observers;
- **security engineering**: threat modelling (STRIDE) of the platform itself, since a PAM system is a high-value target;
- **integrity guarantees**: hash-chained audit records verifiable offline;
- **distributed components**: control plane ↔ gateway communication over mutually authenticated gRPC.

## 1.8 Realism of scope

The two-member team limits the MVP to SSH and PostgreSQL as representative privileged-access protocols; interviews must test whether that choice fits the target users. RDP, AD discovery and UEBA remain outside the MVP. The project has a proposed milestone schedule in `10-project-management-and-risks.md`; capacity and final scope should be reviewed with the supervisor as implementation results arrive.

References: see `03-literature-review.md`.
