# 10. Project Management and Risks

## 10.1 Tools and links

| Item | Link |
|---|---|
| Git repository | [GitHub: pam-platform](https://github.com/nurkkkkhan/pam-platform), with published `main` branch. |
| YouTrack project | [DP project](https://diploma-work.youtrack.cloud/projects/DP); live issue import and assignments still need verification. `DIP` is only the local planning prefix. |
| Documentation | `docs/` in this repository |
| Initial backlog | [`project/youtrack-backlog.csv`](../project/youtrack-backlog.csv) — 30 planned tasks with two-member ownership, priorities and acceptance criteria; import or create them in the real project. |

## 10.2 Allocation of responsibilities

The members' names are recorded in `project/team.conf`. These are planned responsibilities, not a claim about who authored existing files.

| Member | Planned responsibility | Assigned tasks |
|---|---|---|
| Member 1 — Ibyrkhanov Zhanibek, Backend | Technology comparison, architecture and API contracts; control plane, identity, policy, vault, audit, SSH/PostgreSQL gateways, CI, threat model and performance testing | DIP-7, 9, 14, 15, 16, 17, 18, 19, 20, 21, 23, 24, 25, 26, 28 |
| Member 2 — Khassenov Nurkhan, Frontend and UI | Project coordination, research, requirements, prototype, React UI, live-session interface, usability testing, verification plan and evaluation report | DIP-1, 2, 3, 4, 5, 6, 8, 10, 11, 12, 13, 22, 27, 29, 30 |

Member 1 owns the [OpenAPI contract](../api/openapi.yaml), server-side OIDC/TOTP and authorized live-event/termination APIs; Member 2 implements the browser OIDC flow and live-session controls against those contracts. Member 2 leads the verification plan and evaluation report with security and performance evidence from Member 1. Each member reviews the other's pull requests; integration, security checks, and the defense demonstration are shared work. The detailed acceptance criteria are in the [local backlog](../project/youtrack-backlog.csv).

## 10.3 Milestones

| # | Milestone | Deliverables | Date |
|---|---|---|---|
| M0 | Kick-off (proposed) | Agree topic and scope; create repo and YouTrack project; record real links | 01.10.2026 |
| **M1** | **Pre-defense 1** (tag `pre-defense-1`) | Proposal, relevance evidence, literature review, comparisons, requirements v0.1, diagrams, prototype, verification plan, backlog | 24.10.2026 |
| M2 | Core skeleton | API + DB schema + OIDC login, vault with envelope encryption, audit chain, CI | 21.11.2026 |
| M3 | First brokered session | SSH gateway with credential injection + recording; JIT request/approve | 19.12.2026 |
| M4 | **Pre-defense 2** (tag `pre-defense-2`) | PostgreSQL gateway, command rules, live view, rotation; ER diagram; performance tests | 13.02.2027 |
| M5 | Feature freeze | Reports, notifications, hardening, usability tests, security scans | 13.03.2027 |
| M6 | **Final pre-defense** (tag `final-pre-defense`) | Evaluation report, thesis draft, demo | 10.04.2027 |
| M7 | Defense | Final thesis, presentation | 29.05.2027 |

Dates other than M1 are a proposed schedule, not verified university deadlines. Align them with the official academic calendar and edit `project/milestones.csv` before importing tasks. The repository and YouTrack project URL now exist, but this does not establish that all M0 activities happened on the proposed date.

## 10.4 Branching and traceability convention

- Planned local task: `DIP-18 — Define the database schema`.
- After issue creation, use its **actual** YouTrack key in the branch name, commit message, and PR/MR description. The local `DIP-18` key is only a planning reference.
- Each PR/MR describes the change and how it was checked; the other member reviews it before merge to `main`. Record the real PR/commit link in YouTrack.

Individual contribution is established by the real authors' work and review records over time. The first two commits share one author identity on the same day; the local task assignment does not prove past authorship or progressive contributions from both members.

See `CONTRIBUTING.md`.

## 10.5 Risks

| ID | Risk | Prob. | Impact | Mitigation | Owner |
|---|---|---|---|---|---|
| R1 | Scope too large (many protocols/features) | High | High | Strict MoSCoW; RDP, AD discovery, UEBA out of scope; review scope at each milestone | Member 2 |
| R2 | Security flaws in own platform (PAM is a high-value target) | Med | High | Threat model, planned supervisor review of cryptography, scanners in CI, standard libraries only | Member 1 |
| R3 | Team inexperience with Go / SSH protocol internals | Med | Med | Spike task DIP-14 in M2, pair programming, reference permissively licensed code | Member 1 |
| R4 | PostgreSQL wire protocol proxy harder than estimated | Med | Med | Start with authentication injection and simple query capture; extended protocol later; fallback: PG access via SSH tunnel with recorded `psql` | Member 1 |
| R5 | Not enough real users for interviews / usability tests | Med | Med | Seek access to university IT, internship companies, or supervisor contacts; document actual consent and results | Member 2 |
| R6 | Test bench resources unavailable | Low | Med | Docker-based bench on laptops; identify another host if needed | Member 1 |
| R7 | Uneven contribution in team | Med | Med | Assign work evenly, hold weekly reviews, and preserve real authorship in Git and YouTrack | Member 2 |
| R8 | Accidental commit of secrets | Low | High | `.gitignore`, supplied `gitleaks` CI configuration, only `.example` configs; verify the remote CI run after publishing | Both |
| R9 | Licence conflicts of dependencies | Low | Med | Check dependency licences before use and avoid copying incompatible code | Member 1 |

## 10.6 Dependencies

| Dependency | Needed for | Status |
|---|---|---|
| Supervisor approval of scope | All | Not evidenced; request and record actual feedback |
| Access to interviewees / client | Relevance evidence, usability | Not evidenced |
| VM resources (university lab or cloud) | Integration & performance tests | Not confirmed |
| Keycloak test realm | FR-U1 | planned M2 |
| `golang.org/x/crypto/ssh`, `pgx/pgproto3`, MinIO, React | Implementation | Proposed technologies; integration remains planned |

## 10.7 Immediate next steps

1. Add the supervisor's identity and confirm repository access for the second member, supervisor and committee.
2. Confirm the live YouTrack backlog contains all 30 issues; record their actual IDs, assignments, priorities and deadlines, and correct any differences from the local CSV.
3. Conduct and document 3–5 interviews, or record a substantive client agreement; update assumptions A1–A7.
4. Request substantive supervisor feedback and record the comments, decisions, and resulting changes in `11-supervisor-feedback.md`.
5. Continue M2 work: DB schema, OIDC login, vault, audit chain and CI verification.
