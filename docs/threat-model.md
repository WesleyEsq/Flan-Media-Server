# Security Threat Model & Mitigations

Security boundaries, attack vectors, and defensive postures for Flan Media Server operating in private homelab networks, VPN mesh overlays (Tailscale, WireGuard), or behind reverse proxies.

---

## 1. Operating Environment & Threat Context

* **Network Environment:** Private local area networks (LANs), overlay networks (Tailscale, WireGuard), or reverse proxies (Nginx, Caddy, Cloudflare Tunnels).
* **Threat Actors:** Compromised IoT devices or guest devices on the local Wi-Fi, untrusted household members attempting privilege escalation, or automated vulnerability scanners.
* **Eliminated Attack Surfaces:** Because third-party API integrations (such as TMDB) are completely omitted, Server-Side Request Forgery (SSRF) and external credential leakage risks do not exist.

---

## 2. Threat Analysis & Mitigations

### 1. Path Traversal

* **Threat:** An attacker attempts to read arbitrary files from the host system (such as `/etc/shadow` or SSH keys) using directory traversal payloads (`../../`) in media request parameters.
* **Mitigations:**
  * **Opaque Numeric Identifiers:** Public endpoints accept only integer IDs (`/stream/video/42`), never client-supplied filesystem paths.
  * **Canonical Path Validation:** File paths resolved from the database are verified using `filepath.EvalSymlinks`, asserting that the target resides strictly within the verified media root directory.
  * **Embedded Web Assets:** Static assets and templates are served from Go's `embed.FS`, fully isolated from the host filesystem.

### 2. PIN Brute-Force Attacks

* **Threat:** An attacker automates rapid login attempts against short numeric PINs (4 to 6 digits) over the local network.
* **Mitigations:**
  * **Bcrypt Hashing:** PINs are stored as salted bcrypt hashes.
  * **Progressive Lockout:** 5 consecutive failed attempts trigger a 5-minute lockout with exponential backoff on subsequent failures.
  * **CPU Throttling:** Concurrent bcrypt hashing operations are serialized (maximum 1 concurrent comparison) to prevent CPU starvation on low-power devices.

### 3. Malicious Avatar Uploads

* **Threat:** An attacker uploads excessively large files to exhaust disk space, or uploads malicious scripts disguised as images to execute cross-site scripting (XSS).
* **Mitigations:**
  * **Size Limitation:** A strict 2 MB upload ceiling is enforced via `http.MaxBytesReader`.
  * **Format & Magic Byte Validation:** File headers are validated against allowed image formats (JPEG, PNG, WebP). SVG and HTML formats are rejected.
  * **Isolated Storage:** Files are written to `data/avatars/{user_id}.ext` with non-executable permissions (`0644`) and served with `X-Content-Type-Options: nosniff`.

### 4. Session Tampering & Privilege Escalation

* **Threat:** An attacker modifies session cookies to impersonate an administrator or bypass authorization checks.
* **Mitigations:**
  * **HMAC-SHA256 Signing:** Cookie payloads are formatted as `userID:role:tokenVersion:issuedAt:signature` and verified in constant time (`hmac.Equal`).
  * **Instant Token Revocation:** The server tracks `token_version` in memory. Modifying a PIN increments the user's version, instantly invalidating all existing cookies without requiring database lookups.
  * **Secure Cookie Attributes:** Cookies are set with `HttpOnly`, `SameSite=Lax`, and the `__Host-` prefix when TLS is active.
  * **Cross-Origin Protection:** Cross-origin request forgery is prevented using Go standard origin and fetch metadata checks (`Sec-Fetch-Site`).

### 5. Unmounted Media Hijacking

* **Threat:** An external drive fails to mount at boot, causing the server to initialize against an empty folder on root flash and corrupting the media catalog.
* **Mitigations:**
  * **Marker Verification:** Startup terminates immediately if `.flan-keep` is missing from either `./data/` or `./media/`.
  * **Safe Scanner Abort:** The scanner verifies marker files prior to crawling and halts safely without purging database records if a media directory is unavailable.

### 6. First-Run Admin Takeover

* **Threat:** An unauthorized user on the local network visits the web interface before the administrator and registers the primary administrative account.
* **Mitigations:**
  * **Terminal Bootstrap Token:** On initial boot with zero users, Flan outputs a random 6-character bootstrap setup token to the terminal/systemd journal. The setup endpoint (`/setup`) strictly requires this token to create the initial administrator account.

### 7. Denial-of-Service via Account Lockout

* **Threat:** A malicious user deliberately enters bad PINs to lock other family members out of their accounts.
* **Mitigations:**
  * **Dual-Key Isolation:** Lockout state is tracked using a combined key of Client IP address and Target User ID. A lockout triggered from one device does not affect users accessing their accounts from other household devices.
  * **Proxy Awareness:** When running behind a reverse proxy, the client IP is extracted from `X-Forwarded-For` only when requests originate from configured `TRUSTED_PROXIES`.

### 8. External Player URL Tampering

* **Threat:** An attacker attempts to forge or replay signed download URLs used for VLC/MPV streaming.
* **Mitigations:**
  * **HMAC Signed URLs:** URLs carry `exp`, `u`, and `sig` query parameters validated in constant time.
  * **Time Expiration:** Signed URLs expire after 4 hours.

---

## 3. Account Recovery & Failsafes

Because Flan runs in self-contained homelabs without external email infrastructure, recovery relies on local host access:

* **Standard Accounts:** Administrators can reset user PINs directly through the `/manage` console.
* **Administrator Recovery:** If the administrator forgets their PIN, host shell access allows resetting it directly:

  ```bash
  ./flan --reset-admin
  ```

  This command updates the PIN in SQLite, hashes it with bcrypt, and increments `token_version` to invalidate all active sessions.

---

## 4. Related Documentation

* [Master System Architecture](design.md)
* [Rate Limiting & Throttling](rate-limiting.md)
* [Storage Architecture](storage.md)
* [Database Schema & Queries](database.md)
