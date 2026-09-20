# Rate Limiting & Resource Throttling Specification

This document details the rate-limiting architecture, resource governor, and traffic controls designed to protect low-power single-board computers running Flan Media Server.

---

## 1. Why Rate Limiting is Critical on SBCs

On high-end servers, rate limiting is primarily a security tool to block denial-of-service attacks. On low-power single-board computers (such as Raspberry Pi or Orange Pi), rate limiting is an essential **stability and throughput defense**:

+ **CPU Protection:** Low-power ARM cores will peg at 100% CPU utilization if flooded with JSON API queries, starving active video streams.
+ **Mechanical Disk Thrashing:** If multiple clients stream simultaneously from a mechanical USB hard drive, head seek latency causes throughput to drop from 100 MB/s down to 3 MB/s, causing constant buffering for all viewers.
+ **PIN Security:** 4-digit PINs (10,000 combinations) can be brute-forced over a local network in seconds without rate limits.
+ **Upstream Ban Prevention:** Outbound metadata scraping against public APIs (like TMDB) will trigger IP bans if requests are not strictly metered.

To balance security and system stability without degrading the viewing experience, Flan enforces **five distinct rate-limiting zones**.

---

## 2. The Five Rate-Limiting Zones

``` text
+-------------------------------------------------------------------+
|                        Incoming Traffic                           |
+---------------------------------+---------------------------------+
                                  |
         +------------------------+------------------------+
         |                        |                        |
         v                        v                        v
  [Zone A: Stream Gov]    [Zone B: PIN Auth]      [Zone C: API Bucket]
  Max 3 active streams    5 attempts = lockout    15 req/s per client IP
  (Protects USB HDD)      (Prevents brute-force)  (Protects SBC CPU)
                                  |
                                  v
                        [Zone D: Heavy Jobs]
                        1 scan at a time + 30s cooldown
                                  |
                                  v
+-------------------------------------------------------------------+
|                        Outbound Traffic                           |
|  [Zone E: Outbound Scraper Ticker]                                |
|  Max 2.8 requests/second to TMDB / OpenLibrary (Avoids IP bans)   |
+-------------------------------------------------------------------+
```

---

### Zone A: The Concurrent Stream Governor (Mechanical Disk Defense)

+ **Protected Resource:** External USB 3.0 / SATA mechanical hard drive read heads.
+ **Target Endpoint:** GET /stream/{id}
+ **Mechanism:** In-memory counting semaphore.
+ **Configuration:** MAX_CONCURRENT_STREAMS in the .env file (default: 3).

#### How It Works

When a client requests a video stream, the handler acquires a token from the stream governor semaphore:

+ **Under Capacity (1 to 3 active streams):** The token is acquired and the Linux kernel streams byte ranges via sendfile. The drive performs smooth, high-throughput sequential reads.
+ **Over Capacity (4+ active streams):** The handler returns HTTP 429 Too Many Requests:

  ```json
  {"error": "Stream capacity reached (3 active streams). Please try again shortly."}
  ```

+ When streaming completes or the socket closes (e.g. user pauses or navigates away), the semaphore token is released immediately.

---

### Zone B: Authentication & PIN Brute-Force Protection

+ **Protected Resource:** User profile PIN authentication.
+ **Target Endpoint:** POST /api/login
+ **Mechanism:** Dual-key in-memory progressive lockout (keyed by IP address and Profile ID).

#### Schedule

+ **Attempts 1 to 3:** Standard response time.
+ **Attempt 4:** Artificial 2-second delay before responding.
+ **Attempt 5:** **5-minute account lockout**.
+ **Subsequent Failures:** Each failure after 5 doubles the lockout duration (10 minutes, 20 minutes, 40 minutes).
+ **Success:** Resets the failed attempts counter to 0 upon entering the correct PIN.

This turns a 15-second brute-force attack on 10,000 combinations into an operation that takes months.

---

### Zone C: Inbound API Rate Limiting (CPU & Database Protection)

+ **Protected Resource:** SBC CPU cores and SQLite query processing.
+ **Target Endpoints:** All JSON API routes (/api/media, /api/progress, /api/libraries).
+ **Mechanism:** Token Bucket algorithm per client IP address (using Go's standard golang.org/x/time/rate).

#### Parameters

+ **Sustained Rate:** 15 requests per second per IP.
+ **Burst Capacity:** 30 requests per IP.
+ **Memory Cleanup:** An in-memory cleaner runs every 5 minutes and evicts IP rate-limit records that have been idle for more than 5 minutes. This prevents the rate-limiter map from growing and consuming RAM.
+ **Exemptions:** Static assets (CSS, JS, icons) and cached cover images (/covers/{id}) have a higher burst allowance (60 requests) so grid browsing never stutters.

---

### Zone D: Heavy Operations Cooldown (Library Scans)

+ **Protected Resource:** Filesystem I/O and SQLite write transactions.
+ **Target Endpoint:** POST /api/scan
+ **Mechanism:** Single-worker mutex with a trailing timestamp cooldown.

#### Rules

1. **Single-Flight Only:** Only one recursive directory scan can run at any given moment.
2. **Conflict Response:** If a scan is already running, clicking "Scan" returns HTTP 409 Conflict:

   ```json
   {"error": "A library scan is already in progress."}
   ```

3. **Cooldown Window:** Once a scan finishes, a 30-second cooldown is enforced. Any scan requests during the cooldown return HTTP 429 Too Many Requests.

---

### Zone E: Outbound Web Scraper Throttling (API Ban Defense)

+ **Protected Resource:** Upstream API reputation and IP address standing with The Movie Database (TMDB) and OpenLibrary.
+ **Target:** Outbound HTTPS calls in internal/scraper.
+ **Mechanism:** Sequential worker queue governed by a Go time.Ticker.

#### Parameters for the scraper

+ **Tick Rate:** 350 milliseconds between outbound requests.
+ **Effective Rate:** ~2.8 requests per second.
+ **Queue:** Scraping jobs are processed sequentially in a single background goroutine. Even if 100 new movies are added simultaneously, the scraper queues them and meters requests safely below the standard 40-requests-per-10-seconds threshold of public APIs.
