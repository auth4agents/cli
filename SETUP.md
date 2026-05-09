# Auth4Agent CLI Setup Guide

## Requirements

- Go 1.23+
- Running Auth4Agent server
- Verified operator domain
- Network access to Auth4Agent server

---

# 1. Clone Repository

## Windows

```powershell
git clone https://github.com/auth4agents/auth4agent-cli.git
cd auth4agent-cli
```

---

## macOS/Linux

```bash
git clone https://github.com/auth4agents/auth4agent-cli.git
cd auth4agent-cli
```

---

# 2. Install Dependencies

## Windows

```powershell
go mod tidy
```

---

## macOS/Linux

```bash
go mod tidy
```

---

# 3. Build CLI

## Windows

```powershell
go build -o build/auth4agent.exe
```

Run:

```powershell
./build/auth4agent.exe --help
```

---

## macOS/Linux

```bash
go build -o build/auth4agent
```

Run:

```bash
./build/auth4agent --help
```

---

# 4. Check Server Health

Before starting, verify the server is running:

## Windows

```powershell
./build/auth4agent.exe health
```

## macOS/Linux

```bash
./build/auth4agent health
```

Expected:

```json
{
  "status": "ok",
  "timestamp": 1234567890
}
```

---

# 5. Operator Setup

Operator setup establishes:
- trust domain
- root identity
- authorization authority

This is typically performed once per organization/domain.

---

# 6. Initialize Operator Identity

Generates:
- operator keypair
- operator identity config

---

## Windows

```powershell
./build/auth4agent.exe init --operator --domain example.com --server http://localhost:8080
```

---

## macOS/Linux

```bash
./build/auth4agent init --operator --domain example.com --server http://localhost:8080
```

Expected output:

```json
{
  "mode": "operator",
  "domain": "example.com",
  "root_public_key": "base64...",
  "server": "http://localhost:8080"
}
```

---

# 7. Register Operator

Registers operator with server.

---

## Windows

```powershell
./build/auth4agent.exe register operator
```

---

## macOS/Linux

```bash
./build/auth4agent register operator
```

Expected:

```text
operator registered

id: 550e8400-e29b-41d4-a716-446655440000
domain: example.com
status: active
```

**IMPORTANT**: Save the operator ID for agent registration.

---

# 8. Get DNS Verification Instructions

---

## Windows

```powershell
./build/auth4agent.exe verify-operator instructions
```

---

## macOS/Linux

```bash
./build/auth4agent verify-operator instructions
```

Expected:

```text
dns verification instructions

domain: example.com

record_type: TXT
name: _auth4agents.example.com
value: abc123def456...

after propagation run:
auth4agents verify-operator confirm
```

---

# 9. Add DNS TXT Record

Example:

```text
Type: TXT
Name: _auth4agents.example.com
Value: abc123def456...
TTL: 300
```

Wait for DNS propagation (usually 5-30 minutes).

---

# 10. Confirm Operator Verification

---

## Windows

```powershell
./build/auth4agent.exe verify-operator confirm
```

---

## macOS/Linux

```bash
./build/auth4agent verify-operator confirm
```

Expected:

```text
operator verified

verified_at: 2024-01-15T10:30:00Z
```

---

# 11. Check Verification Status

---

## Windows

```powershell
./build/auth4agent.exe verify-operator status
```

---

## macOS/Linux

```bash
./build/auth4agent verify-operator status
```

Expected:

```text
operator verification status

id: 550e8400-e29b-41d4-a716-446655440000
domain: example.com
status: active
verified_at: 2024-01-15T10:30:00Z
```

---

# 12. Agent Setup

Agents are autonomous machine identities.

Each agent:
- owns its own DID
- owns its own private key
- authenticates independently

---

# 13. Initialize Agent Identity

Generates:
- DID
- agent keypair
- DID document

---

## Windows

```powershell
./build/auth4agent.exe init --domain example.com --server http://localhost:8080
```

---

## macOS/Linux

```bash
./build/auth4agent init --domain example.com --server http://localhost:8080
```

Expected output:

```json
{
  "mode": "agent",
  "did": "did:agent:example.com:5e5009ae48d7e3ea",
  "public_key": "base64...",
  "did_document": {...},
  "server": "http://localhost:8080"
}
```

---

# 14. Register Agent

Use operator ID from step 7.

---

## Windows

```powershell
./build/auth4agent.exe register agent --operator-id YOUR_OPERATOR_ID
```

---

## macOS/Linux

```bash
./build/auth4agent register agent --operator-id YOUR_OPERATOR_ID
```

Expected:

```text
agent registered

did: did:agent:example.com:5e5009ae48d7e3ea
status: active
```

---

# 15. Check Agent Status

## Windows

```powershell
./build/auth4agent.exe whoami
```

## macOS/Linux

```bash
./build/auth4agent whoami
```

Expected:

```text
Mode: Agent
DID: did:agent:example.com:5e5009ae48d7e3ea
Operator: 550e8400-e29b-41d4-a716-446655440000
Server: http://localhost:8080
```

---

# 16. Assign Agent Scopes

Authorization is operator-controlled.

Agents cannot assign permissions to themselves.

---

## Windows

```powershell
./build/auth4agent.exe agent-scopes set --scopes read:payments,read:orders
```

---

## macOS/Linux

```bash
./build/auth4agent agent-scopes set --scopes read:payments,read:orders
```

Expected:

```json
{
  "updated": true,
  "did": "did:agent:example.com:5e5009ae48d7e3ea",
  "allowed_scopes": ["read:payments", "read:orders"]
}
```

---

# 17. List Agent Scopes

## Windows

```powershell
./build/auth4agent.exe agent-scopes list
```

## macOS/Linux

```bash
./build/auth4agent agent-scopes list
```

Expected:

```text
Allowed scopes for did:agent:example.com:5e5009ae48d7e3ea:
  - read:payments
  - read:orders
```

---

# 18. Request JWT Token

Agents authenticate using:
- DID
- challenge signing
- proof exchange

---

## Windows

```powershell
./build/auth4agent.exe issue --scope read:payments --aud https://api.example.com
```

---

## macOS/Linux

```bash
./build/auth4agent issue --scope read:payments --aud https://api.example.com
```

Flow:
1. Request challenge
2. Sign challenge locally
3. Exchange signed proof
4. Receive JWT

Expected:

```text
token issued

did: did:agent:example.com:5e5009ae48d7e3ea
scope: read:payments
audience: https://api.example.com
expires_at: 2024-01-15T11:30:00Z

eyJhbGciOiJFZERTQSIsImtpZCI6ImtleS0xIiwidHlwIjoiSldUIn0...
```

---

# 19. Save Token to Environment Variable

## Windows (PowerShell)

```powershell
$env:JWT_TOKEN = "eyJhbGciOiJFZERTQSIsImtpZCI6ImtleS0xIiwidHlwIjoiSldUIn0..."
```

## macOS/Linux

```bash
export JWT_TOKEN="eyJhbGciOiJFZERTQSIsImtpZCI6ImtleS0xIiwidHlwIjoiSldUIn0..."
```

---

# 20. Authorization Enforcement

If an agent requests unauthorized scopes:

## Windows

```powershell
./build/auth4agent.exe issue --scope admin:root --aud https://api.example.com
```

## macOS/Linux

```bash
./build/auth4agent issue --scope admin:root --aud https://api.example.com
```

Expected:

```text
scope not allowed
```

Authorization is deny-by-default.

**Note**: When a challenge expires (410 Gone) or is already used (409 Conflict), request a new challenge automatically.

---

# 21. Update Server Configuration

If the server URL changes:

## Windows

```powershell
./build/auth4agent.exe set-server --url http://new-server:8080
```

## macOS/Linux

```bash
./build/auth4agent set-server --url http://new-server:8080
```

---

# 22. Revoke Tokens and Agents

## Revoke a Specific Token

## Windows

```powershell
./build/auth4agent.exe revoke token --token "JWT_TOKEN" --reason "compromised"
```

## macOS/Linux

```bash
./build/auth4agent revoke token --token "JWT_TOKEN" --reason "compromised"
```

## Revoke an Entire Agent (Permanent)

## Windows

```powershell
./build/auth4agent.exe revoke agent --id AGENT_UUID --reason "decommissioned"
```

## macOS/Linux

```bash
./build/auth4agent revoke agent --id AGENT_UUID --reason "decommissioned"
```

## Suspend an Agent (Temporary)

## Windows

```powershell
./build/auth4agent.exe revoke suspend --agent-id AGENT_UUID --reason "investigating"
```

## macOS/Linux

```bash
./build/auth4agent revoke suspend --agent-id AGENT_UUID --reason "investigating"
```

## Reactivate a Suspended Agent

## Windows

```powershell
./build/auth4agent.exe revoke reactivate --agent-id AGENT_UUID
```

## macOS/Linux

```bash
./build/auth4agent revoke reactivate --agent-id AGENT_UUID
```

## List Revoked/Suspended Agents

## Windows

```powershell
./build/auth4agent.exe revoke list
```

## macOS/Linux

```bash
./build/auth4agent revoke list
```

---

# 23. Verify JWT Offline

Performs:
- signature verification
- issuer validation
- expiration validation
- JWKS validation

No online introspection required.

---

## Windows

```powershell
./build/auth4agent.exe verify --token "JWT_TOKEN"
```

---

## macOS/Linux

```bash
./build/auth4agent verify --token "JWT_TOKEN"
```

Expected:

```text
token verification

valid: true
expired: false
issuer: https://auth.example.com
subject: did:agent:example.com:5e5009ae48d7e3ea
audience: [https://api.example.com]
scope: read:payments
operator_id: 550e8400-e29b-41d4-a716-446655440000
issued_at: 2024-01-15 10:30:00 +0000 UTC
expires_at: 2024-01-15 11:30:00 +0000 UTC
```

---

# 24. Verify JWT Online

Performs:
- server-side verification
- revocation checks
- online introspection

---

## Windows

```powershell
./build/auth4agent.exe verify --token "JWT_TOKEN" --online
```

---

## macOS/Linux

```bash
./build/auth4agent verify --token "JWT_TOKEN" --online
```

Expected:

```text
online verification: success
```

---

# 25. JSON Output

Most commands support:

```text
--json
```

Example:

## Windows

```powershell
./build/auth4agent.exe verify --token "JWT_TOKEN" --json
```

## macOS/Linux

```bash
./build/auth4agent verify --token "JWT_TOKEN" --json
```

---

# 26. Configuration Directory

## Windows

```text
C:\Users\<USER>\.auth4agents\
```

## macOS/Linux

```text
~/.auth4agents/
```

---

# 27. Generated Files

## Operator

```text
operator.json
keys/operator.key
```

## Agent

```text
agent.json
keys/agent.key
```

**Private keys never leave local machine.**

---

# 28. Security Notes

- Never commit `.auth4agents` directory
- Never share `.key` files
- Backup identities securely
- Use HTTPS in production (`AUTH4AGENT_BASE_URL=https://...`)
- Rotate server signing keys regularly
- Restrict operator access
- Verify DNS ownership before production deployment
- Use `revoke token` immediately if token is compromised
- Regularly audit `revoke list` for unauthorized revocations

---

# 29. Common Commands Reference

| Action | Command |
|--------|---------|
| Check server health | `health` |
| Init operator | `init --operator` |
| Register operator | `register operator` |
| Get DNS instructions | `verify-operator instructions` |
| Verify operator | `verify-operator confirm` |
| Check verification status | `verify-operator status` |
| Init agent | `init` |
| Register agent | `register agent --operator-id <id>` |
| Show identity | `whoami` |
| Assign scopes | `agent-scopes set --scopes <scopes>` |
| List scopes | `agent-scopes list` |
| Issue token | `issue --scope <scope> --aud <audience>` |
| Verify token (offline) | `verify --token <token>` |
| Verify token (online) | `verify --token <token> --online` |
| Revoke token | `revoke token --token <token>` |
| Revoke agent | `revoke agent --id <agent-id>` |
| Suspend agent | `revoke suspend --agent-id <agent-id>` |
| Reactivate agent | `revoke reactivate --agent-id <agent-id>` |
| List revoked agents | `revoke list` |
| Update server URL | `set-server --url <url>` |
| JSON output | `--json` flag on any command |

---

# 30. Full Example Flow

```bash
# 1. Check server health
auth4agent health

# 2. Operator setup
auth4agent init --operator --domain example.com --server https://auth.example.com
auth4agent register operator
auth4agent verify-operator instructions
# ... add DNS TXT record and wait ...
auth4agent verify-operator confirm
auth4agent verify-operator status

# 3. Agent setup
auth4agent init --domain example.com --server https://auth.example.com
auth4agent register agent --operator-id <operator-id>
auth4agent whoami

# 4. Authorization
auth4agent agent-scopes set --scopes read:payments,read:orders
auth4agent agent-scopes list

# 5. Authentication
auth4agent issue --scope read:payments --aud https://api.example.com

# 6. Verification
auth4agent verify --token "<jwt>" --online

# 7. Revocation (if needed)
auth4agent revoke token --token "<jwt>" --reason "compromised"
auth4agent revoke list
```

---

# 31. Troubleshooting

## Challenge Errors

| Error | Solution |
|-------|----------|
| `410 Gone` | Challenge expired - request new challenge |
| `409 Conflict` | Challenge already used - request new challenge |
| `404 Not Found` | Challenge not found - request new challenge |
| `agent is not active` | Agent revoked or suspended - check status |

## Connection Errors

```bash
# Check server health first
auth4agent health

# Update server URL if changed
auth4agent set-server --url http://new-server:8080
```

## Scope Errors

```bash
# List current scopes
auth4agent agent-scopes list

# Update scopes if needed
auth4agent agent-scopes set --scopes read:payments,write:reports
```