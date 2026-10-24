# 8. UX / Prototype

## 8.1 Clickable prototype

- Location in repo: [`prototype/index.html`](../prototype/index.html) — open in any browser, no build needed.
- Online link: `<paste the published prototype link here>`
- Covers the key scenarios below. Data in the prototype is fictional.

| Screen | Scenario | Requirements |
|---|---|---|
| Dashboard | User sees active grants, pending requests, how to connect | FR-A1, NFR-U2 |
| Request access | User requests JIT access with reason and duration | FR-A2, FR-A7 |
| Approvals | Approver reviews and approves/rejects; self-approval disabled | FR-A3 |
| Live sessions | Security officer watches live terminal, terminates | FR-M3, FR-A6 |
| Session replay | Auditor replays recording, searches commands, sees rule hits | FR-M1, FR-M2, FR-M4 |
| Vault / accounts | Admin sees accounts, rotation status, last rotated | FR-V1–V3 |
| Audit log | Auditor filters events, verifies hash chain, exports | FR-AU1–AU3 |

## 8.2 Command-line interaction (no UI needed)

```bash
# 1. Request access (CLI helper, optional) — or use the web UI
$ pam login                       # opens browser for SSO, then asks TOTP
$ pam request --target db-prod-01 --account root --duration 2h \
                --reason "INC-1042: disk full on /var"
Request REQ-2087 created. Waiting for approval (policy: 1 approver)...
Approved by s.nurlanova at 14:05. Grant G-551 valid until 16:05.

# 2. Connect with a normal OpenSSH client — no password is shown
$ ssh alice%root%db-prod-01@gw.corp.local
PAM Platform: session S-9f31 is recorded. Grant expires 16:05 (in 1h59m).
root@db-prod-01:~# df -h /var
root@db-prod-01:~# rm -rf /var/lib/postgresql
PAM Platform: command blocked by rule R-3 "Destructive delete on DB data dir". Security team notified.

# 3. PostgreSQL through the gateway
$ pam db-token --target pg-billing --account app_admin
Token (valid 5 min): pam_pg_7Hq...   # one-time password for the gateway
$ psql "host=gw.corp.local port=6432 user=alice%app_admin%pg-billing dbname=billing"
```

## 8.3 API examples (REST, see `api/openapi.yaml`)

**Create an access request**

```http
POST /api/v1/access-requests
Authorization: Bearer <OIDC access token>
Content-Type: application/json

{
  "target_id": "tgt_db-prod-01",
  "account": "root",
  "duration_minutes": 120,
  "reason": "INC-1042: disk full on /var"
}
```

```json
201 Created
{
  "id": "REQ-2087",
  "status": "pending",
  "required_approvals": 1,
  "expires_if_not_approved_at": "2026-10-24T15:00:00+05:00"
}
```

**Approve**

```http
POST /api/v1/access-requests/REQ-2087/approve
{ "comment": "OK, incident confirmed" }
```

```json
200 OK
{ "id": "REQ-2087", "status": "approved", "grant": { "id": "G-551", "valid_until": "2026-10-24T16:05:00+05:00" } }
```

**Self-approval attempt (FR-A3)**

```json
403 Forbidden
{ "error": "separation_of_duties", "message": "Requester cannot approve own request" }
```

**Audit event (as stored and exported)**

```json
{
  "seq": 18234,
  "ts": "2026-10-24T14:07:12.381+05:00",
  "actor": "alice",
  "action": "rule.blocked",
  "object": "session/S-9f31",
  "details": { "rule": "R-3", "command": "rm -rf /var/lib/postgresql", "target": "db-prod-01" },
  "prev_hash": "4be1…c9a2",
  "hash": "d07f…11e8"
}
```

**Verify audit chain**

```bash
$ pam-audit verify --from 1 --to 18234
OK: 18234 events, chain intact. Head hash d07f…11e8 matches anchor 2026-10-24T14:10Z.
```

## 8.4 Next UX steps

- Usability walkthrough of the prototype with 3–5 interviewees (record findings in `02-evidence-of-relevance.md`).
- Adjust flows (e.g., request templates, "extend grant" button) based on feedback.
