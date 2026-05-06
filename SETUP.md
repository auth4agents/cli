# Auth4Agent CLI Guide

## Requirements

- Auth4Agent server running
- Go 1.23+
- Verified operator domain
- Network access to Auth4Agent server

---

# 1. Build CLI

## Windows

```powershell
go build -o build/auth4agent.exe
````

Run:

```powershell
./build/auth4agent.exe
```

---

## macOS/Linux

```bash id="m8k4wp"
go build -o build/auth4agent
```

Run:

```bash id="h2f7qn"
./build/auth4agent
```

---

# 2. Initialize Operator Identity

Creates:

* operator keypair
* operator config
* operator identity

---

## Windows

```powershell id="c5t1ma"
./build/auth4agent.exe init --operator --domain example.com --server http://localhost:8080
```

---

## macOS/Linux

```bash id="y3w7vb"
./build/auth4agent init --operator --domain example.com --server http://localhost:8080
```

---

# 3. Operator Files

Generated automatically.

## Windows

```text id="s6q2rf"
C:\Users\<USER>\.auth4agent\
```

## macOS/Linux

```text id="n9m1xk"
~/.auth4agent/
```

Files:

```text id="q4v7ep"
operator.json
keys/operator.key
```

---

# 4. Register Operator

Registers operator with Auth4Agent server.

---

## Windows

```powershell id="w7k5fd"
./build/auth4agent.exe register operator
```

---

## macOS/Linux

```bash id="d1z9qu"
./build/auth4agent register operator
```

Expected:

```text id="j0f8ny"
operator registered
```

---

# 5. Get DNS Verification Instructions

---

## Windows

```powershell id="x8u4cp"
./build/auth4agent.exe verify-operator instructions
```

---

## macOS/Linux

```bash id="g7r3lt"
./build/auth4agent verify-operator instructions
```

Expected:

```text id="e5y1mw"
TXT record name
TXT record value
```

---

# 6. Add DNS TXT Record

Add TXT record in DNS provider.

Example:

```text id="r2p9vk"
Type: TXT
Host: _auth4agents
Value: verification_token
```

Wait for DNS propagation.

---

# 7. Confirm Operator Verification

---

## Windows

```powershell id="n6w2zx"
./build/auth4agent.exe verify-operator confirm
```

---

## macOS/Linux

```bash id="b3m8ha"
./build/auth4agent verify-operator confirm
```

Expected:

```text id="q9k5uv"
operator verified
```

---

# 8. Initialize Agent Identity

Creates:

* agent DID
* agent keypair
* DID document

---

## Windows

```powershell id="f4v1pb"
./build/auth4agent.exe init --domain example.com --server http://localhost:8080
```

---

## macOS/Linux

```bash id="u8x7tr"
./build/auth4agent init --domain example.com --server http://localhost:8080
```

Expected DID:

```text id="y5c0jd"
did:agent:example.com:...
```

Generated files:

```text id="v1n8kg"
agent.json
keys/agent.key
```

---

# 9. Register Agent

Use operator ID from:

```text id="t3q6me"
register operator
```

output.

---

## Windows

```powershell id="j2h4ys"
./build/auth4agent.exe register agent --operator-id YOUR_OPERATOR_ID
```

---

## macOS/Linux

```bash id="r8p1cw"
./build/auth4agent register agent --operator-id YOUR_OPERATOR_ID
```

Expected:

```text id="m4z9qt"
agent registered
```

---

# 10. Issue JWT Token

Requests scoped JWT using DID proof authentication.

---

## Windows

```powershell id="h7k3xn"
./build/auth4agent.exe issue --scope read:payments --aud api.example.com
```

---

## macOS/Linux

```bash id="e1u9lp"
./build/auth4agent issue --scope read:payments --aud api.example.com
```

Flow:

1. Request challenge
2. Sign challenge locally
3. Exchange proof
4. Receive JWT

Expected:

```text id="k6x2vd"
token issued
```

---

# 11. Verify JWT Offline

Performs:

* JWT decode
* expiration check
* claim inspection

No server call.

---

## Windows

```powershell id="w0m8rc"
./build/auth4agent.exe verify --token "JWT_TOKEN"
```

---

## macOS/Linux

```bash id="s4q5yb"
./build/auth4agent verify --token "JWT_TOKEN"
```

Expected:

* issuer
* audience
* scope
* expiration
* subject DID

---

# 12. Verify JWT Online

Performs:

* signature verification
* token validation
* revocation checks

Requires server access.

---

## Windows

```powershell id="p8t1fx"
./build/auth4agent.exe verify --token "JWT_TOKEN" --online
```

---

## macOS/Linux

```bash id="a6n4kv"
./build/auth4agent verify --token "JWT_TOKEN" --online
```

Expected:

```text id="u2y7mw"
verified_online: true
```

---

# 13. JSON Output

Most commands support:

```text id="f7v0qd"
--json
```

Example:

## Windows

```powershell id="b5r2nh"
./build/auth4agent.exe issue --scope read:* --aud api.example.com --json
```

## macOS/Linux

```bash id="m9x8tc"
./build/auth4agent issue --scope read:* --aud api.example.com --json
```

---

# 14. CLI Configuration Locations

## Windows

```text id="r3m7wb"
C:\Users\<USER>\.auth4agent\
```

---

## macOS/Linux

```text id="d8u1vp"
~/.auth4agent/
```

---

# 15. Configuration Files

## Operator

```text id="f1y5qn"
operator.json
```

Contains:

* operator ID
* domain
* public key
* server URL

---

## Agent

```text id="t6z2mx"
agent.json
```

Contains:

* DID
* DID document
* public key
* operator relationship

---

# 16. Key Files

Stored separately.

```text id="e0p4ka"
keys/operator.key
keys/agent.key
```

Private keys never leave local machine.

---

# 17. Security Notes

* Never commit `.auth4agent`
* Never share `.key` files
* Backup identity files securely
* Use HTTPS in production
* Rotate signing keys periodically
* Verify DNS ownership before production use

---

# 18. Common Commands

| Action            | Command                   |
| ----------------- | ------------------------- |
| Init operator     | `init --operator`         |
| Register operator | `register operator`       |
| Verify operator   | `verify-operator confirm` |
| Init agent        | `init`                    |
| Register agent    | `register agent`          |
| Issue token       | `issue`                   |
| Verify token      | `verify`                  |

---

# 19. Full Example Flow

```text id="c9k3rf"
1. init --operator
2. register operator
3. verify-operator instructions
4. verify-operator confirm
5. init
6. register agent
7. issue
8. verify
```
