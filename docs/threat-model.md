# Threat Model and Security Specification

This document details the threat landscape, security boundaries, and defensive mitigations for Flan Media Server.

---

## 1. Operating Context and Security Posture

Flan Media Server is designed to run on low-power Linux computers and single-board computers (SBCs). In a typical homelab setup, the server operates in one of three environments:

1. **Isolated Local Network (LAN):** Accessible only to household devices and Wi-Fi guests.
2. **Overlay Mesh VPN (Tailscale, WireGuard):** Accessible remotely by authenticated personal devices.
3. **Public Port Forward or Reverse Proxy:** Accessible over the public internet, usually behind Cloudflare or Nginx.

### Potential Threat Actors

+ Untrusted devices on the local Wi-Fi network (compromised IoT devices, guests).
+ Internal household users attempting unauthorized access to administrative settings or other user profiles.
+ Remote automated scanners if the port is exposed directly to the internet.

---

## 2. Attack Vectors and Mitigations

### Vector 1: Path Traversal and Arbitrary File Access

+ **Threat:** The server streams files from the host filesystem. If user input directly determines file paths, an attacker could request sensitive host files such as /etc/shadow or SSH keys.

+ **Impact:** Critical. Complete exposure of the host operating system.
+ **Mitigations:**
  + **Opaque Numeric Identifiers:** The web client and streaming endpoints never accept file paths. All media requests use typed endpoints and integer IDs (e.g. `/stream/movie/42` or `/stream/episode/42`).
  + **Canonical Path Verification:** When looking up a media file from the database, the server resolves all symbolic links using `filepath.EvalSymlinks` and verifies that the resulting absolute path starts with the verified root path of the library it belongs to (from the `libraries` table).
  + **Static Asset Isolation:** Static files (HTML, CSS, JS) are embedded into the Go binary using embed.FS. The file server never serves from the operating system root.

---

## Vector 2: PIN Brute-Forcing

+ **Threat:** Profiles are authenticated using 4 or 6-digit numeric PINs for usability on TVs and mobile devices. A 4-digit PIN has only 10,000 combinations (0000 to 9999), which an automated script could test within seconds over a local network.

+ **Impact:** High. Unauthorized profile takeover.
+ **Mitigations:**
  + **Progressive Rate Limiting:** Track failed PIN attempts per profile and per client IP address. Detailed rules are in [docs/rate-limiting.md](file:///home/wess/Documents/MechanicalSpeak/Flan-Media-Server/docs/rate-limiting.md).
  + **Lockout Schedule:** After 5 consecutive failed attempts, enforce a 5-minute lockout. Each subsequent failure doubles the delay.
  + **Bcrypt Storage:** PINs are never stored in plaintext. They are salted and hashed using bcrypt before being saved to SQLite.

---

## Vector 3: Disk Space Exhaustion (Denial of Service)

+ **Threat:** SBCs typically run on small micro-SD cards or flash drives (16gb to 64gb). An attacker or runaway upload could fill the storage partition, causing the Linux kernel to panic or system services to crash.

+ **Impact:** High. System unavailability.
+ **Mitigations:**
  + **Admin-Only Upload Permissions:** File and folder upload endpoints are strictly restricted to the admin profile. Standard profiles cannot upload files.
  + **Pre-Upload Disk Check:** Before accepting an upload, query available disk space via statfs. If free space is below a safety threshold (e.g. 2gb), the upload is rejected immediately with HTTP 507 Insufficient Storage.
  + **Streaming Directly to Disk:** File uploads are streamed straight from the network socket to disk in 32kb chunks via r.MultipartReader and io.Copy. This ensures memory usage remains near zero and prevents Out-Of-Memory (OOM) crashes during large file transfers.

---

## Vector 4: File Type Confusion and Stored Script Execution

+ **Threat:** An attacker uploads malicious HTML, SVG, or executable scripts disguised as media files, attempting to execute cross-site scripting (XSS) in an admin's browser session.

+ **Impact:** High. Session hijacking and administrative takeover.
+ **Mitigations:**
  + **Strict Extension Whitelist:** The server only accepts approved extensions:
    + Video: .mp4, .webm, .mkv (if direct play compatible), .avi
    + Books: .epub, .pdf
    + Images: .jpg, .jpeg, .png, .webp
  + **Explicit MIME Headers:** When serving media, the server explicitly sets the Content-Type header (such as video/mp4) and adds X-Content-Type-Options: nosniff. This prevents the browser from interpreting video files as executable HTML or script content.
  + **File Permissions:** Uploaded files are written with 0644 permissions (read and write only, no execution flag).

---

## Vector 5: Server-Side Request Forgery (SSRF) in Scraper

+ **Threat:** If the cover and metadata scraper accepts arbitrary URLs from users, an attacker could instruct the server to make requests to internal network services (for example, hitting a router admin page at <http://192.168.1.1>).

+ **Impact:** Medium to High. Internal network reconnaissance.
+ **Mitigations:**
  + **Hardcoded External Domains:** The scraper only queries trusted, hardcoded public API domains (The Movie Database and OpenLibrary).
  + **Private IP Blocking:** The HTTP client used for scraping rejects any redirects or targets resolving to private, loopback, or link-local IP ranges (127.0.0.0/8, 10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16).

---

## Vector 6: Privilege Escalation & Session Tampering

+ **Threat:** A standard household user attempts to call administrative endpoints to initiate scans or alter other user accounts, or attempts to forge/tamper with session cookies to elevate their role from `user` to `admin`.

+ **Impact:** High. Unauthorized administrative takeover.
+ **Mitigations:**
  + **Role Verification Middleware:** Endpoints that initiate scans, upload files, or manage users strictly verify that the active session's authenticated user record carries the `admin` role.
  + **HMAC-SHA256 Cryptographic Signing:** Session cookies use a tamper-proof payload formatted as `userID:role:issuedAt:signature`. Changing any token component invalidates the signature immediately.
  + **Constant-Time Verification:** Cookie signatures are verified using Go's standard `crypto/hmac.Equal` to eliminate side-channel timing attacks during authentication verification.
  + **Hardened Key Storage:** The 32-byte signing secret is loaded from `SESSION_SECRET` or read from a local `.session_secret` keyfile stored in the database folder. The keyfile is created with restrictive `0600` permissions (readable only by the daemon process user), preventing unauthorized local disclosure.
  + **Cookie Security Flags:** Session cookies are strictly configured with `HttpOnly` (blocking JavaScript access), `SameSite=Lax` (preventing CSRF during cross-origin navigation), and `Secure` when TLS/HTTPS is active.

---

## Vector 7: Unmounted Storage Path Hijacking & Ghost Database Initialization

+ **Threat:** When external NVMe or USB drives fail to mount on boot, the target directory remains as an empty folder on the root micro-SD card. A server that blindly initializes databases or accepts uploads will write to the root filesystem, exhausting flash storage, causing split-brain database corruption, or exposing the setup wizard to users.
+ **Impact:** High. Data corruption and system crash.
+ **Mitigations:**
  + **Marker File Verification (.flan-keep):** Before opening an external database path, verify the presence of the hidden marker file. If absent, immediately abort startup with a fatal log error.
  + **Scanner Mount Liveness:** Check that the media root exists and is not empty before pruning any catalog items.
  + **Destination statfs Validation:** Ensure upload destination filesystems match the expected external mount points and possess adequate free space before writing bytes.

---

## 3. Account Recovery and Failsafes

Because the server runs locally without external email dependencies, account recovery uses two straightforward mechanisms:

1. **Standard Profile Recovery:** The admin can reset or change any standard user's PIN directly from the settings page.
2. **Command-Line Failsafe:** If the admin forgets their PIN, physical or SSH access to the host allows running:

   ```bash
   ./flan --reset-admin
   ```

   This interactive command resets the admin account directly in the local SQLite database.

---

### Related Documentation

+ [Master System Specifications](docs/design.md)
+ [Five-Zone Rate Limiting Architecture](docs/rate-limiting.md)
+ [Storage Architecture & Mount Resiliency](docs/storage.md)
+ [Database Schema & Hardened Session Storage](docs/database.md)
+ [Admin File Upload Sequence Diagram](docs/diagrams/data-flow.md#6-admin-zero-memory-file-upload-data-flow)
