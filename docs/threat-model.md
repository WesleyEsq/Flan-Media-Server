# Security Threat Model & Defensive Mitigations

Security boundaries, attack vectors, and defensive postures for Flan Media Server running in homelabs, VPN overlays, and local networks.

---

## 1. Operating Environment & Threat Actors

* **Operating Context:** Isolated LANs, overlay mesh networks (Tailscale, WireGuard), or behind reverse proxies (Nginx, Cloudflare).
* **Threat Actors:** Compromised IoT devices or guests on local Wi-Fi, untrusted household users seeking elevated permissions, or automated internet scanners.
* **Zero Outbound Attack Surface:** Because external scraper integrations (TMDB) are eliminated, Server-Side Request Forgery (SSRF) risks are non-existent.

---

## 2. Threat Vector Matrix

| # | Attack Vector | Threat & Impact | Defensive Mitigation |
| :- | :--- | :--- | :--- |
| **1** | **Path Traversal** | Attacker accesses sensitive host files (`/etc/shadow`, SSH keys) via file parameters. *(Critical)* | • **Opaque IDs:** Endpoints accept typed numeric IDs (`/stream/video/42`), never file paths.<br>• **Canonical Path Check:** Resolves symlinks via `filepath.EvalSymlinks` and asserts that the target resides inside the verified media directory.<br>• **Embedded Static Files:** Web assets are served from `embed.FS`, completely isolated from the host root filesystem. |
| **2** | **PIN Brute-Forcing** | 4–6 digit numeric PINs are automated over the local network. *(High)* | • **Bcrypt Hashing:** PINs are salted and hashed with bcrypt.<br>• **Progressive Lockout:** 5 failed attempts trigger a 5-minute lockout with exponential backoff on subsequent failures (Safeguard 2). |
| **3** | **Malicious Avatar Uploads** | Attacker attempts to upload large files to exhaust storage or upload HTML/scripts disguised as images. *(Medium)* | • **Size Ceiling:** Hard 2 MB ceiling enforced via `http.MaxBytesReader`.<br>• **Magic Byte Validation:** Validates image magic headers (JPEG, PNG, WebP) and rejects executable/SVG content.<br>• **Sandboxed Storage:** Saved to `data/avatars/{user_id}.ext` with non-executable `0644` permissions and served with `X-Content-Type-Options: nosniff`. |
| **4** | **Session Tampering & Privilege Escalation** | Attacker tampers with cookies to elevate role or access deleted accounts. *(High)* | • **HMAC-SHA256 Signing:** Cookie payload formatted as `userID:role:tokenVersion:issuedAt:signature`, verified in constant time (`hmac.Equal`).<br>• **Instant In-Memory Revocation:** Server caches `token_version` in RAM (`map[int]int`). Resetting a PIN increments the version, revoking old sessions instantly without DB queries.<br>• **Key Security:** 32-byte secret loaded from `SESSION_SECRET` or `data/.session_secret` (`0600` permissions).<br>• **Cookie Flags:** `HttpOnly`, `SameSite=Lax`, and `Secure` when TLS is active. |
| **5** | **Ghost Database Hijacking** | External media/DB mount fails on boot; server initializes on root flash card, causing split-brain. *(High)* | • **Marker File (`.flan-keep`):** Startup halts if the DB path lacks `.flan-keep`.<br>• **Scanner Mount Check:** Scans abort safely without purging database records if a media folder is empty or absent. |

---

## 3. Account Recovery & Failsafes

Because Flan runs in self-contained homelabs without email infrastructure, recovery uses a two-tier model:

1. **Standard Users:** The administrator can change or reset any user's PIN from `/manage`.
2. **Administrator CLI Failsafe:** If the admin forgets their PIN, host shell access allows resetting it directly:
   ```bash
   ./flan --reset-admin
   ```
   This interactive CLI command prompts for a new PIN, hashes it with bcrypt, updates SQLite, and increments `token_version` to invalidate any compromised sessions.

---

## 4. Related Documentation

* [Rate Limiting Architecture](rate-limiting.md)
* [Storage Architecture](storage.md)
* [Simplified Database Schema (6 Tables)](database.md)
