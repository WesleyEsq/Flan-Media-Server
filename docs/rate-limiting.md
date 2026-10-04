# Rate Limiting & Resource Throttling

Flan Media Server concentrates traffic controls into two essential safeguards designed to protect attached storage hardware and prevent brute-force authentication attempts.

---

## 1. Concurrency Controls (Playback Lease Governor)

### Purpose

Protects attached external mechanical storage (USB HDDs) from head thrashing caused by simultaneous random read requests.

### Protected Endpoints

* `GET /stream/video/{file_id}`
* `GET /download/video/{file_id}`

### Mechanics

Counting raw HTTP Range connections using a semaphore causes false positives because modern web video players issue frequent overlapping range requests to probe containers, inspect metadata, and seek.

Flan uses an activity-based playback lease governor:

1. When a client requests video bytes or an external player accesses a signed download URL, the governor checks for an existing lease identified by `(session_id, file_id)`.
2. If a lease exists, its `last_active` timestamp updates without consuming additional capacity. Seeks and range chunks from the same viewer share the same lease.
3. If no lease exists and active leases are fewer than 3 (`MAX_CONCURRENT_STREAMS`), a new lease is issued.
4. If 3 active leases already exist for other sessions, the request is rejected immediately.
5. Inactive leases expire after 30 seconds of inactivity or when the client sends a playback-ended beacon.

### Behavior on Breach

Returns `HTTP 503 Service Unavailable` with a `Retry-After: 30` header:

```json
{
  "error": "Stream capacity reached (3 active viewers)"
}
```

---

## 2. Authentication Rate Limiting (PIN Lockout)

### Purpose

Prevents automated brute-force attacks against the 4 to 6-digit numeric account PINs while protecting host CPU resources from bcrypt exhaustion.

### Protected Endpoints

* `POST /api/login`

### Mechanics

* **Dual-Key Tracking:** Tracked using a combined key of Client IP address and Target User ID. This prevents an attacker from locking out legitimate users on other household devices.
* **Proxy Support:** Respects `X-Forwarded-For` when `TRUSTED_PROXIES` is configured.
* **CPU Protection:** Concurrent bcrypt hashing operations are serialized (maximum 1 concurrent comparison) to prevent CPU starvation on low-power devices.
* **Lockout Progression:**
  * 1 to 3 failed attempts: Normal error response.
  * 4 failed attempts: Artificial 2-second processing delay.
  * 5 failed attempts: 5-minute account lockout.
  * Subsequent failures: Exponential lockout backoff.

### Behavior on Breach

Returns `HTTP 429 Too Many Requests`:

```json
{
  "error": "Account locked. Try again in 5 minutes"
}
```

---

## 3. Simplified Design Decisions

* **No Token Buckets:** Generalized per-IP API token buckets are omitted. On a private server for 1 to 5 household members, in-memory token maps add overhead without practical security benefit.
* **Scan Deduplication:** Library scans are protected by a standard Go `singleflight.Group`. Simultaneous scan triggers from multiple admin sessions coalesce into a single background scan.
* **Zero External Rate Limits:** Because metadata scraping is 100% local, no rate-limiting logic is needed for external services like TMDB.

---

## 4. Related Documentation

* [Master System Architecture](design.md)
* [Security Threat Model](threat-model.md)
* [Storage Architecture](storage.md)
