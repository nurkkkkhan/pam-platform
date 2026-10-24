# ADR-005: Hash-chained audit log
Status: Accepted · Task: DIP-16

## Context
Platform admins are themselves privileged; audit must reveal tampering (FR-AU2, NFR-S5). Literature: Schneier & Kelsey 1999; Crosby & Wallach 2009.
## Decision
`hash_n = SHA-256(hash_{n-1} || canonical_json(event_n))`, genesis hash = 32 zero bytes. Head hash periodically anchored outside the DB (signed file / e-mail to auditor / SIEM). Verifier tool `pam-audit verify`.
## Consequences
+ Any modification, deletion or reordering detected. − Truncation of the tail detected only against an anchor → anchoring interval is a config parameter. History trees are future work.
