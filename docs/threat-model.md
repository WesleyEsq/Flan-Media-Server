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
* **Threat:** The server streams files from the host filesystem. If user input directly determines file paths, an attacker could request sensitive host files such as /etc/shadow or SSH keys.
* **Impact:** Critical. Complete exposure of the host operating system.
* **Mitigations:**
  + **Opaque Numeric Identifiers:** The web client and streaming endpoints never accept file paths. All media requests use integer IDs (e.g. /stream/42).
  + **Canonical Path Verification:** When looking up a media file from the database, the server resolves all symbolic links using filepath.EvalSymlinks and verifies that the resulting absolute path starts with one of the configured media directory root paths.
  + **Static Asset Isolation:** Static files (HTML, CSS, JS) are embedded into the Go binary using embed.FS. The file server never serves from the operating system root.

---

## Vector 2: PIN Brute-Forcing
* **Threat:** Profiles are authenticated using 4 or 6-digit numeric PINs for usability on TVs and mobile devices. A 4-digit PIN has only 10,000 combinations (0000 to 9999), which an automated script could test within seconds over a local network.
* **Impact:** High. Unauthorized profile takeover.
* **Mitigations:**
  + **Progressive Rate Limiting:** Track failed PIN attempts per profile and per client IP address.
  + **Lockout Schedule:** After 5 consecutive failed attempts, enforce a 5-minute lockout. Each subsequent failure doubles the delay.
  + **Bcrypt Storage:** PINs are never stored in plaintext. They are salted and hashed using bcrypt before being saved to SQLite.

---

## Vector 3: Disk Space Exhaustion (Denial of Service)
* **Threat:** SBCs typically run on small micro-SD cards or flash drives (16gb to 64gb). An attacker or runaway upload could fill the storage partition, causing the Linux kernel to panic or system services to crash.
* **Impact:** High. System unavailability.
* **Mitigations:**
  + **Admin-Only Upload Permissions:** File and folder upload endpoints are strictly restricted to the admin profile. Standard profiles cannot upload files.
  + **Pre-Upload Disk Check:** Before accepting an upload, query available disk space via statfs. If free space is below a safety threshold (e.g. 2gb), the upload is rejected immediately with HTTP 507 Insufficient Storage.
  + **Streaming Directly to Disk:** File uploads are streamed straight from the network socket to disk in 32kb chunks via r.MultipartReader and io.Copy. This ensures memory usage remains near zero and prevents Out-Of-Memory (OOM) crashes during large file transfers.

---

## Vector 4: File Type Confusion and Stored Script Execution
* **Threat:** An attacker uploads malicious HTML, SVG, or executable scripts disguised as media files, attempting to execute cross-site scripting (XSS) in an admin's browser session.
* **Impact:** High. Session hijacking and administrative takeover.
* **Mitigations:**
  + **Strict Extension Whitelist:** The server only accepts approved extensions:
    + Video: .mp4, .webm, .mkv (if direct play compatible), .avi
    + Books: .epub, .pdf
    + Images: .jpg, .jpeg, .png, .webp
  + **Explicit MIME Headers:** When serving media, the server explicitly sets the Content-Type header (such as video/mp4) and adds X-Content-Type-Options: nosniff. This prevents the browser from interpreting video files as executable HTML or script content.
  + **File Permissions:** Uploaded files are written with 0644 permissions (read and write only, no execution flag).

---

## Vector 5: Server-Side Request Forgery (SSRF) in Scraper
* **Threat:** If the cover and metadata scraper accepts arbitrary URLs from users, an attacker could instruct the server to make requests to internal network services (for example, hitting a router admin page at http://192.168.1.1).
* **Impact:** Medium to High. Internal network reconnaissance.
* **Mitigations:**
  + **Hardcoded External Domains:** The scraper only queries trusted, hardcoded public API domains (The Movie Database and OpenLibrary).
  + **Private IP Blocking:** The HTTP client used for scraping rejects any redirects or targets resolving to private, loopback, or link-local IP ranges (127.0.0.0/8, 10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16).

---

## Vector 6: Privilege Escalation
* **Threat:** A standard household user attempts to call administrative endpoints to initiate scans or alter other user accounts.
* **Impact:** Medium. Unauthorized configuration changes.
* **Mitigations:**
  + **Role Verification Middleware:** Endpoints that initiate scans, upload files, or manage users verify that the active session belongs to a user with the admin role.
  + **HMAC-Signed Session Security:** Session cookies are cryptographically signed using HMAC-SHA256 with a server secret. Cookies are configured with HttpOnly, SameSite=Lax, and Secure flags when running over HTTPS.

---

## 3. Account Recovery and Failsafes

Because the server runs locally without external email dependencies, account recovery uses two straightforward mechanisms:

1. **Standard Profile Recovery:** The admin can reset or change any standard user's PIN directly from the settings page.
2. **Command-Line Failsafe:** If the admin forgets their PIN, physical or SSH access to the host allows running:
   ```bash
   ./flan --reset-admin
   ```
   This interactive command resets the admin account directly in the local SQLite database.
