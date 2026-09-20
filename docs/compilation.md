# Cross-Compilation & SBC Deployment Guide

This document details target hardware architectures, cross-compilation procedures, binary size optimization, and production systemd configuration for deploying Flan Media Server on single-board computers (SBCs) and low-power Linux devices.

---

## 1. Overview & Toolchain Independence

Flan Media Server is designed to run on resource-constrained hardware, including 32-bit and 64-bit ARM boards. Because the database engine uses the pure-Go SQLite driver (`modernc.org/sqlite`) and the standard Go HTTP stack, the entire application compiles with `CGO_ENABLED=0`.

This eliminates external C toolchains (such as `gcc-arm-linux-gnueabihf` or `aarch64-linux-gnu-gcc`). You can cross-compile standalone, statically linked binaries for any target architecture directly from your primary development machine using the standard Go toolchain.

---

## 2. Hardware & Architecture Target Matrix

| Device Family | Architecture | Go Build Environment | Example Hardware |
| :--- | :--- | :--- | :--- |
| **ARMv6 (32-bit)** | `linux/arm` (v6) | `CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6` | Raspberry Pi Zero, Zero W, 1 Model A+/B+ |
| **ARMv7 (32-bit)** | `linux/arm` (v7) | `CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7` | Raspberry Pi 2, Pi 3 (32-bit OS), Orange Pi Zero |
| **ARM64 (64-bit)** | `linux/arm64` | `CGO_ENABLED=0 GOOS=linux GOARCH=arm64` | Raspberry Pi 3, 4, 5 (64-bit OS), Zero 2 W, Rock Pi |
| **x86_64 (64-bit)** | `linux/amd64` | `CGO_ENABLED=0 GOOS=linux GOARCH=amd64` | Intel NUC, x86 thin clients, homelab servers |

---

## 3. Compilation Commands

### Standard Local Build

To compile for the host machine:

```bash
go build -o flan ./cmd/flan
```

### Cross-Compiling for Raspberry Pi Boards

Execute these commands from your local workstation:

```bash
# 1. Raspberry Pi Zero / Zero W / 1 (ARMv6)
CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 go build -ldflags="-s -w" -o flan-armv6 ./cmd/flan

# 2. Raspberry Pi 2 / 3 running 32-bit OS (ARMv7)
CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -ldflags="-s -w" -o flan-armv7 ./cmd/flan

# 3. Raspberry Pi 3 / 4 / 5 or Zero 2 W running 64-bit OS (ARM64)
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o flan-arm64 ./cmd/flan

# 4. Standard x86_64 Homelab / Mini PC (AMD64)
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o flan-amd64 ./cmd/flan
```

---

## 4. Binary Footprint Optimization

Micro-SD cards used on single-board computers have limited sequential read speeds and high random-read latency. Reducing executable size shortens binary load times and reduces disk churn on boot:

+ **Linker Stripping (`-ldflags="-s -w"`):**
  + `-s`: Omits the symbol table and debug information.
  + `-w`: Omits DWARF debugging information.
  + **Result:** Shrinks the final executable size by approximately 25% to 35% without altering runtime execution or stability.
+ **UPX Compression (Optional Caution):**
  + Compressing the binary with UPX reduces disk usage further, but requires unpacking the binary into RAM on every execution. On 512MB RAM boards (like the Pi Zero), UPX can cause startup memory spikes. Plain uncompressed binaries with stripped debug flags are recommended.

---

## 5. Memory & Garbage Collection Tuning

On single-board computers, operating system memory must be preserved. Run Flan with disciplined garbage collection ceilings:

```bash
GOMEMLIMIT=16MiB GOGC=30 ./flan
```

+ `GOMEMLIMIT=16MiB`: Soft memory target instructing the Go runtime to trigger garbage collection before heap allocations cross 16MB.
+ `GOGC=30`: Aggressive GC trigger ratio (default is 100), ensuring garbage collection runs frequently while individual heaps are small.

---

## 6. Production systemd Deployment

For reliable 24/7 homelab operation, manage Flan as a persistent systemd service.

### 1. File Placement

Copy the cross-compiled binary and configuration to the SBC:

```bash
# On the SBC:
sudo useradd -r -s /bin/false flan
sudo mkdir -p /var/lib/flan/data /var/lib/flan/media
sudo chown -R flan:flan /var/lib/flan
sudo cp flan-arm64 /usr/local/bin/flan
sudo chmod +x /usr/local/bin/flan
```

### 2. Service Unit Configuration

Create `/etc/systemd/system/flan.service`:

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

# Sandboxing & security hardening
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/var/lib/flan
PrivateTmp=true
NoNewPrivileges=true

[Install]
WantedBy=multi-user.target
```

### 3. Enable and Start

```bash
sudo systemctl daemon-reload
sudo systemctl enable --now flan
sudo systemctl status flan
```

---

### Related Documentation

+ [Master System Specifications](docs/design.md)
+ [Storage Architecture & Drive Resiliency](docs/storage.md)
+ [Database Schema & Wear-Leveling Pragmas](docs/database.md)
+ [System Directories & Package Guide](docs/directories.md)
