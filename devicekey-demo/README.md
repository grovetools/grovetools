# Device-key identity walkthrough

This small, standard-library-only Go module contrasts a copyable bearer secret with a device identity proved by an Ed25519 private key. Run one lesson at a time:

```sh
go run . -step 1
go run . -step 2
```

Steps 3–5 are reserved for enrollment, signed handshakes, and enforcement/revocation in the next slice.

## Mental model

| Question | Token world (today) | Device world (redesign) |
|---|---|---|
| What names a machine? | A ULID in JSON, but the server discards it | The ULID remains the stable device name |
| What proves identity? | Possession of one long-lived, human-copied token | A signature made by that device's private key |
| What does the server trust? | `sha256(token)` | An approved device ID and public key |
| Can a copied credential impersonate it? | Yes, indefinitely and from any machine | Copying the ID and public key is insufficient |
| Can one machine be revoked? | Not without replacing the shared token everywhere | Yes, by revoking that device's approved key and sessions |
| Are bearer tokens gone? | No; the long-lived token is the identity | No; short-lived sessions are derived after key proof and refreshed silently |

The key distinction is **assertion versus proof**. A device ID says “I am Alice's laptop.” A valid signature proves that the speaker possesses the private key enrolled for Alice's laptop. The private key itself never needs to leave that machine.

## Step 1 — Assertion vs proof

The program mints a random, 26-character Crockford-base32 ID (ULID-shaped for teaching purposes) and stores the ID and friendly name as JSON. An attacker can copy those bytes exactly, so a server that trusts only the claimed ID accepts a perfect impersonation.

Alice also mints an Ed25519 keypair. She signs a message containing her device ID and requested action. Her enrolled public key verifies that signature. The attacker can copy the ID, message, and public key, but a signature made with the attacker's private key does not verify as Alice. The unexported private key is the proof.

```sh
go run . -step 1
```

## Step 2 — The bearer status quo

The toy server stores only `sha256(token)`, like a conventional bearer-token server. Alice's token authenticates successfully—but pasting the exact same token into a different machine also succeeds. The API receives no independently proved machine identity, so both callers collapse into one principal. Revoking only the copied credential is impossible because it is the same credential.

```sh
go run . -step 2
```

## What this toy simplifies

Even when all five lessons are implemented, a real system must also add persistent and transactional storage, TLS, request rate limiting, secure private-key storage, key rotation and recovery, careful clock handling, audit logs, and production-grade secret/session lifecycle controls. This first slice uses in-memory state, prints demonstrations to stdout, and models a ULID-shaped random ID rather than implementing the full ULID timestamp format.
