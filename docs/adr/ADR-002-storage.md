# ADR-002: PostgreSQL for metadata & audit, MinIO for recordings
Status: Accepted · Task: DIP-7

## Context
Need ACID for grants/approvals, append-only audit with JSONB details, large immutable recordings (docs/05 §5.2).
## Decision
PostgreSQL 16 (API role: INSERT-only on `audit_events`); MinIO with object lock for recordings.
## Consequences
+ One well-known DB; S3 API portable to cloud. − Two stateful services to back up.
