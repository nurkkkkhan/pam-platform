# PAM Platform

> **Diploma project:** *Development of a PAM platform for centralized management and control of privileged access in corporate infrastructure*

PAM Platform is a proposed self-hosted platform for replacing shared privileged passwords and standing admin rights. The planned system brokers SSH and PostgreSQL connections through a gateway, injects credentials from an encrypted vault, grants time-limited access after approval, records sessions, applies command rules, and writes a tamper-evident audit trail. The supplied project folder contains design materials, a clickable prototype, and selected Go modules; the full platform is not yet implemented.

## Problem
Privileged accounts (`root`, `postgres`, service accounts) are often shared between admins through chats or spreadsheets, access is permanent, credentials are rarely rotated, and nobody can prove who ran which command. Credential abuse is the most common initial access vector in breaches (22%, Verizon DBIR 2025). Commercial PAM suites solve this but are expensive and complex for SMEs; open-source tools cover only parts. → [Proposal](docs/01-project-proposal.md)

## Target users
System/DB administrators · security officers (approvers) · auditors · external contractors · platform administrators.

## Expected outcomes
- MVP: vault + rotation, SSH & PostgreSQL gateway, JIT access with approval, session recording/replay, command rules, hash-chained audit, web UI.
- Evaluation targets against current practice (shared credentials): 0 credentials known to users, 0% standing access, 100% attributed recorded sessions, 100% audit-tampering detection, ≤ 1 s connection overhead. These are targets for later measurement, not achieved results.

## Team
| Member | Responsibility |
|---|---|
| Member 1 — Ibyrkhanov Zhanibek, Backend | Technology decisions, architecture and API contracts; control plane, identity, policy, vault, rotation, SSH/PostgreSQL gateways, audit, CI, threat model and performance testing |
| Member 2 — Khassenov Nurkhan, Frontend and UI | Project coordination, research and requirements; UX prototype, React interface, live-session controls, usability study, verification plan and evaluation report |

Both members review each other's changes and integrate the system together. Member 1 owns the [OpenAPI contract](api/openapi.yaml), server-side OIDC/TOTP and authorized live-session APIs; Member 2 builds the browser OIDC flow and live-session interface against them. Member 2 leads the verification and evaluation documents using Member 1's security and performance evidence. The detailed 15-task allocation for each member is in [the backlog](project/youtrack-backlog.csv) and the separate [Member 1 plan](<project/Member 1 PAM Platform work plan.docx>) and [Member 2 plan](<project/Member 2 PAM Platform work plan.docx>). Assignment describes planned responsibility; it is not evidence that a named person performed past work.

**Supervisor:** Name and title to be supplied by the team.

## Scope & status
In scope: SSH, PostgreSQL, vault/rotation, JIT approvals, recording, command rules, audit, OIDC + TOTP. Out of scope: RDP, AD discovery, ML analytics, HA cluster.

**Current status (24 October 2026):** initial analysis and design drafts, editable diagrams, a clickable prototype, and Go `audit`, `vault`, and `policy` modules with unit tests are present. The GitHub repository is published and a YouTrack project URL is recorded. The live issue import, assignments, priorities, deadlines, user interviews, supervisor feedback, and progressive Git contributions by both members have not been verified. The gateway, web UI, database integration, and complete application remain to be built; a complete application is not required at this checkpoint.

## Technologies
Go · PostgreSQL 16 · MinIO (S3) · React + TypeScript · gRPC/mTLS · REST/OpenAPI · Keycloak (OIDC) · Docker Compose. Justification: [technology comparison](docs/05-technology-comparison.md), [ADRs](docs/adr/).

## Links
- YouTrack: [DP project](https://diploma-work.youtrack.cloud/projects/DP) and [local backlog](project/youtrack-backlog.csv).
- GitHub: [pam-platform repository](https://github.com/nurkkkkhan/pam-platform).
- Pre-defense 1 checklist: [docs/00-checklist-pre-defense-1.md](docs/00-checklist-pre-defense-1.md)
- Prototype: [prototype/index.html](prototype/index.html) (local HTML file)
- API: [api/openapi.yaml](api/openapi.yaml)

## Documentation
| # | Document |
|---|---|
| 0 | [Checklist for pre-defense 1](docs/00-checklist-pre-defense-1.md) |
| 1 | [Project proposal & scope](docs/01-project-proposal.md) |
| 2 | [Evidence of relevance](docs/02-evidence-of-relevance.md) |
| 3 | [Literature review](docs/03-literature-review.md) |
| 4 | [Existing system comparison](docs/04-existing-systems-comparison.md) |
| 5 | [Technology comparison](docs/05-technology-comparison.md) |
| 6 | [Requirements & acceptance criteria](docs/06-requirements-specification.md) |
| 7 | [Architecture & diagrams](docs/07-architecture.md) |
| 8 | [UX / prototype / API examples](docs/08-ux-prototype.md) |
| 9 | [Verification & evaluation plan](docs/09-verification-and-evaluation-plan.md) |
| 10 | [Project management & risks](docs/10-project-management-and-risks.md) |
| 11 | [Supervisor feedback log](docs/11-supervisor-feedback.md) |

## Repository structure
```
api/            OpenAPI specification
config/         configuration examples (placeholders only)
diagrams/src    editable diagram sources (Graphviz DOT)
diagrams/export SVG/PNG exports
docs/           proposal, analysis, requirements, architecture, plans, ADRs
frontend/       web UI plan (React implementation planned)
project/        local YouTrack import backlog and milestones
prototype/      clickable HTML prototype
src/            selected Go modules and audit CLI tools; gateway planned
tests/          test cases and test materials
```

## Setup, run, test
Requirements: Go 1.22+ for the current modules; Graphviz for regenerating diagrams; Docker for the optional development services.
```bash
make test        # unit tests with coverage (audit, vault, policy)
make diagrams    # regenerate diagram exports from diagrams/src
make audit-demo  # create a sample audit export and verify it
open prototype/index.html
```
`docker compose -f config/docker-compose.example.yml up` starts PostgreSQL, MinIO and Keycloak for development; application services are not included yet.

## Contributing
Each member should use their own Git identity, commit work they actually perform, and review the other member's pull requests. The [task map](scripts/tasks.tsv) and [contribution guide](CONTRIBUTING.md) are planning aids. They do not establish authorship of the files already supplied in this folder.

## Security
Never commit secrets. Use `config/*.example*` files with placeholders; copy to untracked files locally. The supplied GitHub Actions workflow includes a `gitleaks` step; no remote CI run is evidenced here.
