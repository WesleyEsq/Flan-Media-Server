# Rate Limiting & Resource Throttling

On low-power single-board computers, rate limiting is an essential **stability and hardware protection defense**: it prevents CPU exhaustion, eliminates mechanical disk seek thrashing, thwarts PIN brute-force attacks, and prevents upstream API bans.

---

## 1. The Five Rate-Limiting Zones

```text
Incoming Traffic ───────────────────────────────────────────────┐
  │                                                             │
  ├──► [Zone A: Stream Governor]  Max 3 active streams (protects USB HDD)
  ├──► [Zone B: PIN Auth]         5 attempts → 5 min lockout (prevents brute-force)
  ├──► [Zone C: API Bucket]       15 req/s per IP (protects SBC CPU & SQLite)
  └──► [Zone D: Scan Cooldown]    1 scan at a time + 30s cooldown (protects I/O)

Outbound Traffic ───────────────────────────────────────────────┐
  └──► [Zone E: Scraper Pacing]   Max 2.8 req/s to TMDB (prevents IP bans)
```

| Zone | Protected Resource | Target Endpoints | Mechanism & Implementation | Limit / Threshold | Action on Breach |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **Zone A: Stream Governor** | USB HDD read heads & memory | `GET /stream/{type}/{id}` | In-memory counting semaphore | Max 3 concurrent active streams (`MAX_CONCURRENT_STREAMS`) | `HTTP 429 Too Many Requests` |
| **Zone B: PIN Auth** | Profile authentication | `POST /api/login` | Dual-key tracker (IP + User ID) with progressive delay | • 1–3 fails: standard<br>• 4 fails: +2s artificial delay<br>• 5 fails: 5 min lockout<br>• Subsequent: exponential backoff | `HTTP 429 Too Many Requests` (Lockout duration returned) |
| **Zone C: API Bucket** | SBC CPU & SQLite connection | All `/api/*` endpoints | Token Bucket per IP (`golang.org/x/time/rate`) | • Sustained: 15 req/s<br>• Burst: 30 req<br>(Static/covers burst: 60) | `HTTP 429 Too Many Requests` |
| **Zone D: Scan Cooldown** | Filesystem I/O & DB writes | `POST /api/scan`<br>`POST /api/libraries/{id}/scan` | Single-flight mutex with trailing timestamp | • 1 scan execution at a time<br>• 30-second cooldown post-scan | • In-progress: `HTTP 409 Conflict`<br>• Cooldown: `HTTP 429 Too Many Requests` |
| **Zone E: Scraper Pacing** | TMDB & OpenLibrary IP standing | Outbound HTTPS calls in `internal/scraper` | Sequential worker queue with Go `time.Ticker` | 350 ms tick rate (~2.8 req/s) | Requests queued sequentially in a single goroutine |

---

## 2. In-Memory Maintenance & Low-Footprint Rules

* **Zone C Memory Cleaner:** An in-memory cleaner runs every 5 minutes and evicts IP rate-limit records that have been idle for >5 minutes. This ensures the rate-limiter map never leaks memory or exceeds 50 KB of RAM.
* **Cooperative Lock Yielding:** During heavy directory scans or scraper runs, workers call `runtime.Gosched()` and short delays between items so pending HTTP requests can acquire SQLite locks within the `busy_timeout` window without latency spikes.

---

## 3. Related Documentation

* [Master System Architecture](design.md)
* [Security Threat Model](threat-model.md)
* [Storage Architecture](storage.md)
* [Scraping Engine](scraper.md)
