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
git clone https://github.com/auth4agents/cli.git
cd cli
```

---

## macOS/Linux

```bash
git clone https://github.com/auth4agents/cli.git
cd cli
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
./build/auth4agent.exe
```

---

## macOS/Linux

```bash
go build -o build/auth4agent
```

Run:

```bash
./build/auth4agent
```

---

# 4. Operator Setup

Operator setup establishes:
- trust domain
- root identity
- authorization authority

This is typically performed once per organization/domain.

---

# 5. Initialize Operator Identity

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

---

# 6. Register Operator

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
```

Save:
```text
operator_id
```

---

# 7. Get DNS Verification Instructions

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
TXT record name
TXT record value
```

---

# 8. Add DNS TXT Record

Example:

```text
Type: TXT
Host: _auth4agents
Value: verification_token
```

Wait for DNS propagation.

---

# 9. Confirm Operator Verification

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
```

---

# 10. Agent Setup

Agents are autonomous machine identities.

Each agent:
- owns its own DID
- owns its own private key
- authenticates independently

---

# 11. Initialize Agent Identity

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

Expected DID:

```text
did:agent:example.com:...
```

---

# 12. Register Agent

Use operator ID from:
```text
register operator
```

output.

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
```

---

# 13. Assign Agent Scopes

Authorization is operator-controlled.

Agents cannot assign permissions to themselves.

---

## Windows

```powershell
./build/auth4agent.exe agent-scopes set --scopes read:payments
```

---

## macOS/Linux

```bash
./build/auth4agent agent-scopes set --scopes read:payments
```

Example scopes:

```text
read:payments
read:orders
write:reports
```

---

# 14. Request JWT Token

Agents authenticate using:
- DID
- challenge signing
- proof exchange

---

## Windows

```powershell
./build/auth4agent.exe issue --scope read:payments --aud api.example.com
```

---

## macOS/Linux

```bash
./build/auth4agent issue --scope read:payments --aud api.example.com
```

Flow:
1. Request challenge
2. Sign challenge locally
3. Exchange signed proof
4. Receive JWT

Expected:

```text
token issued
```

---

# 15. Authorization Enforcement

If an agent requests unauthorized scopes:

```powershell
./build/auth4agent.exe issue --scope admin:root --aud api.example.com
```

Expected:

```text
scope not allowed
```

Authorization is deny-by-default.

---

# 16. Verify JWT Offline

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
- issuer
- audience
- scope
- expiration
- signature verified

---

# 17. Verify JWT Online

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

# 18. JSON Output

Most commands support:

```text
--json
```

Example:

## Windows

```powershell
./build/auth4agent.exe verify --token "JWT_TOKEN" --json
```

---

## macOS/Linux

```bash
./build/auth4agent verify --token "JWT_TOKEN" --json
```

---

# 19. Configuration Directory

## Windows

```text
C:\Users\<USER>\.auth4agent\
```

---

## macOS/Linux

```text
~/.auth4agent/
```

---

# 20. Generated Files

## Operator

```text
operator.json
keys/operator.key
```

---

## Agent

```text
agent.json
keys/agent.key
```

Private keys never leave local machine.

---

# 21. Security Notes

- Never commit `.auth4agent`
- Never share `.key` files
- Backup identities securely
- Use HTTPS in production
- Rotate server signing keys
- Restrict operator access
- Verify DNS ownership before production deployment

---

# 22. Common Commands

| Action | Command |
|---|---|
| Init operator | `init --operator` |
| Register operator | `register operator` |
| Verify operator | `verify-operator confirm` |
| Init agent | `init` |
| Register agent | `register agent` |
| Assign scopes | `agent-scopes set` |
| Issue token | `issue` |
| Verify token | `verify` |

---

# 23. Full Example Flow

```text
1. init --operator
2. register operator
3. verify-operator instructions
4. verify-operator confirm
5. init
6. register agent
7. agent-scopes set
8. issue
9. verify
```