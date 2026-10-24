# ADR-004: REST for UI, gRPC+mTLS between control plane and gateways
Status: Accepted · Task: DIP-7

## Decision
REST/JSON (OpenAPI 3) + WebSocket for browser; gRPC over mutual TLS for `Authorize`, `GetCredential`, `StreamEvents`, `UploadRecording`, `PushRules`. No message broker in MVP.
## Consequences
+ Typed internal contract, low latency, browser-friendly UI API. − Two API styles to maintain; broker needed later for HA.
