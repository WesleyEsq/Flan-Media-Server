# Cross-Compilation & SBC Deployment Guide

Flan Media Server compiles with `CGO_ENABLED=0` thanks to the pure-Go SQLite driver (`modernc.org/sqlite`). Standalone, statically linked binaries can be cross-compiled for any target architecture directly from your development machine without external C toolchains.

---

## 1. Hardware Target Matrix & Build Commands

Build commands use `-ldflags="-s -w"` to strip symbol and DWARF debug tables, shrinking binary size by 25–35% for faster loading from micro-SD cards:

| Target Device Family | Target Architecture | Cross-Compilation Command |
| :--- | :--- | :--- |
| **Raspberry Pi Zero, Zero W, 1** | ARMv6 (32-bit) | `CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 go build -ldflags="-s -w" -o flan-armv6 ./cmd/flan` |
| **Raspberry Pi 2, 3 (32-bit OS), Orange Pi** | ARMv7 (32-bit) | `CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -ldflags="-s -w" -o flan-armv7 ./cmd/flan` |
| **Raspberry Pi 3, 4, 5, Zero 2 W (64-bit OS)** | ARM64 (64-bit) | `CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o flan-arm64 ./cmd/flan` |
| **x86_64 Homelab / Mini PC / Intel NUC** | AMD64 (64-bit) | `CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o flan-amd64 ./cmd/flan` |
| **Local Host Development** | Host Native | `go build -o flan ./cmd/flan` |

---

## 2. Runtime Memory Tuning

On memory-constrained boards (e.g. 512 MB Pi Zero), run Flan with runtime limits to keep memory strictly within 15–20 MB:

```bash
GOMEMLIMIT=16MiB GOGC=30 ./flan
```

* `GOMEMLIMIT=16MiB`: Soft ceiling prompting the Go runtime to trigger GC before heap allocations exceed 16 MB.
* `GOGC=30`: Aggressive GC trigger ratio (default is 100), ensuring garbage collection runs frequently while individual heaps are small.

---

## 3. Production systemd Service

For 24/7 homelab operation, manage Flan as a systemd service.

### 1. Host Preparation

```bash
sudo useradd -r -s /bin/false flan
sudo mkdir -p /var/lib/flan/data /var/lib/flan/media
sudo chown -R flan:flan /var/lib/flan
sudo cp flan-arm64 /usr/local/bin/flan
sudo chmod +x /usr/local/bin/flan
```

### 2. Service Unit (`/etc/systemd/system/flan.service`)

```ini
[Unit]
Description=Flan Media Server
After=network.target local-fs.target

[Service]
Type=simple
User=flan
Group=flan
WorkingDirectory=/var/lib/flan
ExecStart=/usr/local/bin/flan
Restart=always
RestartSec=5s

# Hard memory limits and process defenses
MemoryMax=32M
Environment="GOMEMLIMIT=16MiB"
Environment="GOGC=30"
Environment="DB_PATH=/var/lib/flan/data/flan.db"
Environment="MEDIA_DIR=/var/lib/flan/media"

# Security hardening & sandboxing
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/var/lib/flan
PrivateTmp=true
NoNewPrivileges=true

[Install]
WantedBy=multi-user.target
```

### 3. Enable & Start

```bash
sudo systemctl daemon-reload
sudo systemctl enable --now flan
sudo systemctl status flan
```

---

## 4. Related Documentation

* [Master System Architecture](design.md)
* [Storage Architecture](storage.md)
* [Database Wear-Leveling](database.md)
