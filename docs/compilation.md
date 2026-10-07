# Compilation & Homelab Deployment Guide

Flan Media Server compiles with `CGO_ENABLED=0` using the pure-Go SQLite driver (`modernc.org/sqlite`). Statically linked binaries can be cross-compiled for any target architecture directly from your development machine without external C toolchains.

---

## 1. Prerequisites & Toolchain

* **Go Version:** Go 1.24 or newer.
* **Target Architectures:** AMD64, ARM64, ARMv7 (32-bit).

---

## 2. Cross-Compilation Commands

Build commands use `-ldflags="-s -w"` to strip symbol and debug information, reducing binary size by 25–35% for faster loading from disk or flash storage.

### Local Development
```bash
go build -o flan ./cmd/flan
```

### x86_64 / AMD64 (Standard PCs, Mini PCs, NUCs)
```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o flan-amd64 ./cmd/flan
```

### ARM64 / AArch64 (Raspberry Pi 3/4/5 64-bit, Orange Pi, Rock Pi)
```bash
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o flan-arm64 ./cmd/flan
```

### ARMv7 32-bit (Raspberry Pi 2/3 32-bit OS, Older SBCs)
```bash
CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -ldflags="-s -w" -o flan-armv7 ./cmd/flan
```

### ARMv6 32-bit (Raspberry Pi 1, Raspberry Pi Zero)
```bash
CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 go build -ldflags="-s -w" -o flan-armv6 ./cmd/flan
```

---

## 2. Resource Management & Runtime Tuning

Flan is designed for low memory and CPU overhead:
1. **Zero-Copy Streaming:** File delivery uses Go's `http.ServeContent`, delegating to Linux `sendfile` to stream bytes directly from filesystem cache to socket without copying into heap buffers.
2. **Direct Passthrough:** Omits CPU-intensive real-time transcoding.
3. **Bounded Database Cache:** SQLite runs in WAL mode with a bounded ~2 MB page cache per connection.
4. **Embedded Assets:** Templates and static files are compiled directly into the binary via `embed.FS`.

Process resource metrics (allocated heap, system memory, goroutines, and uptime) can be viewed directly on the `/manage` console.

### Low-Memory Tuning (512 MB SBCs)

When running on devices with 512 MB of RAM, the Go runtime can be constrained via environment variables:

```bash
GOMEMLIMIT=32MiB GOGC=50 ./flan
```

On devices with 1 GB or more RAM, default Go runtime settings are recommended.

---

## 3. Production systemd Service Setup

To run Flan as a persistent service on a Linux host:

### 1. Create Dedicated User and Directories

```bash
sudo useradd -r -s /bin/false flan
sudo mkdir -p /var/lib/flan/data /var/lib/flan/media
sudo touch /var/lib/flan/data/.flan-keep /var/lib/flan/media/.flan-keep
sudo chown -R flan:flan /var/lib/flan

# Copy binary for target architecture
sudo cp flan-arm64 /usr/local/bin/flan
sudo chmod +x /usr/local/bin/flan
```

### 2. Service Unit (`/etc/systemd/system/flan.service`)

```ini
[Unit]
Description=Flan Media Server
After=network.target local-fs.target
RequiresMountsFor=/var/lib/flan/media /var/lib/flan/data

[Service]
Type=simple
User=flan
Group=flan
WorkingDirectory=/var/lib/flan
ExecStart=/usr/local/bin/flan
Restart=always
RestartSec=5s

# Environment configuration
Environment="PORT=4907"
Environment="DATA_DIR=/var/lib/flan/data"
Environment="DB_PATH=/var/lib/flan/data/flan.db"
Environment="MEDIA_DIR=/var/lib/flan/media"

# Security hardening
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/var/lib/flan
PrivateTmp=true
NoNewPrivileges=true

[Install]
WantedBy=multi-user.target
```

### 3. Enable and Start Service

```bash
sudo systemctl daemon-reload
sudo systemctl enable --now flan
sudo systemctl status flan
```

---

## 4. Related Documentation

* [Master System Architecture](design.md)
* [Storage Architecture](storage.md)
* [Authentication Security & Rate Limiting](rate-limiting.md)
