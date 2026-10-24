# ADR-003: Built-in envelope encryption with pluggable KEK provider
Status: Proposed — to confirm with supervisor · Task: DIP-15

## Context
Secrets must never be stored or decryptable from DB alone (FR-V1, NFR-S1); KEK rotation without re-encrypting all data (FR-V6). Alternatives: external Vault/OpenBao, pgcrypto (docs/05 §5.3).
## Decision
Each secret version gets a random 256-bit DEK; secret encrypted with AES-256-GCM (DEK), DEK wrapped with KEK (AES-256-GCM). KEK provided by `KeyProvider` interface (file/env now; HSM, cloud KMS, OpenBao Transit later). Associated data = account ID + version, preventing ciphertext swapping between rows. Only Go standard library crypto is used — no custom primitives.
## Consequences
+ No extra cluster to operate; KEK rotation = re-wrap DEKs. − We own the correctness → unit tests, supervisor review (risk R2).
