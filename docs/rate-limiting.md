# Authentication Security & Resource Throttling

Flan Media Server concentrates security controls into targeted safeguards designed to protect host CPU resources, prevent brute-force authentication attempts, and eliminate unnecessary network overhead.

---

## 1. Authentication Rate Limiting (Dual-Key PIN Lockout)

### Purpose

Prevents automated brute-force attacks against 4 to 6-digit numeric account PINs while protecting host CPU resources from bcrypt exhaustion.

### Protected Endpoints

* `POST /api/login`

### Mechanics

* **Dual-Key In-Memory Tracking:** Tracked strictly in-memory using a combined key of `(Client IP, Target User ID)`. This prevents an attacker from locking out legitimate users on other household devices, and prevents database write amplification during brute-force attempts. Expired records are safely evicted on subsequent access and periodic sweeps to bound memory consumption under 1,000 active entries.
* **Proxy Support:** Respects `X-Forwarded-For` when `TRUSTED_PROXIES` is configured.
* **Lockout Progression:**
  * 1 to 3 failed attempts: Standard invalid credentials response (`HTTP 401 Unauthorized`).
  * 4 failed attempts: Artificial 2-second processing delay (**applied before acquiring the bcrypt worker slot** to prevent blocking legitimate logins).
  * 5 failed attempts: 5-minute account lockout (`HTTP 429 Too Many Requests`).
  * Subsequent failures: Exponential lockout backoff (capped at 1 hour).

### Behavior on Breach

Returns `HTTP 429 Too Many Requests`:

```json
{
  "error": "Account locked. Try again in 5 minutes"
}
```

---

## 2. CPU Starvation Protection (Bcrypt Worker Gate)

### Purpose

Protects low-power single-board computers (such as Raspberry Pi and Orange Pi) from CPU starvation caused by rapid sequential or concurrent bcrypt hashing comparisons.

### Mechanics

* **Serialized Comparison Gate:** Bcrypt hashing comparisons are serialized across the process via a buffered worker channel:
  ```go
  bcryptGate = make(chan struct{}, 1)
  ```
* **Backpressure Timeout:** If the gate is occupied and the wait exceeds 5 seconds, the request returns `HTTP 503 Service Unavailable` with `Retry-After: 1` to prevent CPU exhaustion.
* **Pre-Gate Artificial Delays:** Artificial delays on repeated failed attempts (e.g. attempt 4) execute **before** acquiring the bcrypt worker slot, ensuring that legitimate logins are never held up behind a penalized brute-force attempt.

---

## 3. Direct Streaming & Simplified Traffic Decisions

* **Direct Range Delivery Without Lease Governors:** Video streaming endpoints (`/stream/video/{file_id}`) delegate directly to Go's standard `http.ServeContent` with HTTP 206 Partial Content (Range requests) and kernel `sendfile`. In a private household of 1 to 5 members, direct 1080p streaming consumes only ~1 MB/s per client. Linux VFS readahead and OS page cache handle concurrent reads naturally without brittle userspace lease trackers, artificial grace timers, or fragile frontend unload beacons.
* **Scan Deduplication:** Library scans are protected by a standard Go `singleflight.Group`. Simultaneous scan triggers from multiple admin sessions coalesce into a single background crawl.
* **No Generalized Token Buckets:** Generalized per-IP API token buckets are omitted. On a private server for household members, in-memory token maps add runtime overhead without practical security benefit.
* **Zero External Rate Limits:** Because metadata scraping is 100% local and offline, no rate-limiting logic or API quotas exist for external services like TMDB.

---

## 4. Related Documentation

* [Master System Architecture](design.md)
* [Security Threat Model](threat-model.md)
* [Storage Architecture](storage.md)
