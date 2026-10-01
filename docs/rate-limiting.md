# Rate Limiting & Resource Throttling

To remain radically simple and lightweight, Flan Media Server eliminates complex leaky-bucket algorithms and concentrates traffic controls into **two essential hardware and security safeguards**.

---

## 1. The Two Essential Safeguards

```text
Incoming Traffic ───────────────────────────────────────────────┐
  │                                                             │
  ├──► [Safeguard 1: Stream Governor] Max 3 active streams (protects USB HDD from thrashing)
  └──► [Safeguard 2: PIN Auth]        5 attempts → 5 min lockout (prevents brute-force)
```

| Safeguard | Protected Resource | Target Endpoints | Mechanism & Implementation | Limit / Threshold | Action on Breach |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **1. Stream Governor** | USB HDD read heads & memory | `GET /stream/video/{file_id}` | In-memory counting semaphore | Max 3 concurrent active streams (`MAX_CONCURRENT_STREAMS`) | `HTTP 429 Too Many Requests`<br>`{"error": "Stream capacity reached (3 active streams)"}` |
| **2. PIN Auth Lockout** | Profile authentication | `POST /api/login` | Dual-key tracker (IP + User ID) with progressive delay | • 1–3 fails: standard<br>• 4 fails: +2s artificial delay<br>• 5 fails: 5 min lockout<br>• Subsequent: exponential backoff | `HTTP 429 Too Many Requests`<br>`{"error": "Account locked. Try again in 5 minutes"}` |

---

## 2. Why Other Limits Were Cut

1. **Inbound API Token Buckets (Cut):** On a personal homelab server for 1 to 5 household members, per-IP token bucket maps and periodic eviction cleaners add code overhead and heap memory tracking with zero real-world benefit.
2. **Library Scan Cooldown Mutex (Cut):** Scans are triggered manually by the admin via `/manage` or automatically on boot. A standard single-flight check prevents simultaneous runs without needing cooldown timers.
3. **Automated Scraper Tickers (Cut):** Background scraper queues and outbound request pacing are eliminated; metadata is local-first by default, with an optional on-demand button.

---

## 3. Related Documentation

* [Master System Architecture](design.md)
* [Security Threat Model](threat-model.md)
* [Storage Architecture](storage.md)
