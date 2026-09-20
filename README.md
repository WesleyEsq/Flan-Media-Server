# Flan Media Server

Server for organizing and streaming videos and books. It is meant for very underpowered Linux computers and single-board computers, for 1 to 5 people at once. Can run on extremely constrained systems and targets around 15 to 20mb RAM.

The goal is cramming a functional streaming server with a multi-page web client into an **extremely** constrained memory and CPU footprint.

## How to use

To access your catalog, visit the client in your browser by entering the IP of your device alongside the port. By default the port is 4907, but you can change that in the .env file.

+ **First run:** On first boot, the server opens a setup wizard at /setup where you create your admin profile, set a numeric PIN, and point to your media folder.
+ **Returning visits:** Shows a "Who is watching?" profile selector. Enter your PIN to access the catalog.
+ **Account recovery:** If you forget your admin PIN, you can reset it directly from the host terminal using `./flan --reset-admin`. Regular profiles can be reset by the admin from the settings page.

## How to install

To compile the server locally, build with Go in your terminal:

```bash
go build -o flan ./cmd/flan
```

To run with memory limits tuned for low-memory devices:

```bash
GOMEMLIMIT=16MiB GOGC=30 ./flan
```

You can also cross-compile directly for ARM boards like Raspberry Pi:

```bash
# Raspberry Pi Zero / 1 (ARMv6)
GOOS=linux GOARCH=arm GOARM=6 go build -o flan-armv6 ./cmd/flan

# Raspberry Pi 2 / 3 (32-bit ARMv7)
GOOS=linux GOARCH=arm GOARM=7 go build -o flan-armv7 ./cmd/flan

# Raspberry Pi 3 / 4 / 5 (64-bit ARM64)
GOOS=linux GOARCH=arm64 go build -o flan-arm64 ./cmd/flan
```

Eventually, I will also configure an Alpine Linux Docker image for containerized deployment.

## Resource use and compromises

The server does the bare minimum to deliver video directly to the client without on-the-fly transcoding. This means files must already be encoded in web-compatible formats such as MP4 (H.264/AAC), WebM (VP9/Opus), or AV1.

+ Video delivery relies on the Linux sendfile system call via Go's standard library to transfer byte ranges directly from the filesystem to the network socket. This keeps media data off the application heap.
+ The server handles Partial Content requests so media players can seek and stream individual parts of files.
+ Tailored media players: Plyr for video with subtitle support, native browser viewing for PDF documents, and an ePub.js web reader for EPUB books with direct download options.
+ Subtitles (.srt) are discovered directly on disk alongside videos and converted to WebVTT format on the fly with zero memory overhead.
+ TV shows are organized hierarchically (Series → Seasons → Episodes) so multi-episode seasons do not flood the main catalog grid.
+ Media intake supports both scanning existing folders in-place without moving files, and direct browser uploads for administrators that stream straight to disk without memory spikes.
+ Requests are handled using lightweight goroutines. Each idle connection consumes only about 2kb of memory, easily handling 10 idle connections and 2 to 3 active streams.
+ Go runtime limits like GOMEMLIMIT and GOGC keep the garbage collector disciplined so total memory stays within 15 to 20mb.

## Database and other dependencies

+ **SQLite3:** Uses a streamlined 3-table schema (users, media_items, playback_progress) running in WAL mode with a small 2mb page cache. Sessions use HMAC-signed cookies to avoid database hits on page visits.
+ **Web Client:** Built with Go html/template for multi-page rendering, along with modular vanilla JavaScript and CSS styled in a soft dark slate theme with lavender accents. The templates and static assets are embedded directly into the single binary with embed.FS, so no external assets need to be deployed.
+ **Metadata & covers:** Uses a local-first approach (checking for poster.jpg or embedded EPUB art first), falling back to TMDB for missing covers, ratings, and genre tags. Covers are stored locally for offline resilience. Administrators can also fix matches or upload custom covers manually.
+ **Libraries:** This project uses Apache 2.0 and all external libraries should use a similar license like MIT or BSD.

## Project background

I'm making this for myself. I own an SBC that I currently use as a home server. The board is already running close to its memory and CPU limits, leaving very little room for traditional, proper media servers like Jellyfin or Plex.

My home setup will rarely see more than 2 or 3 concurrent streams, with at most 10 idle connections at a time. This is explicitly designed for lightweight homelabs and personal single-board computers, not commercial or production-scale workloads.

The name comes from my hamster, Flan (custard in Spanish).

## MVP

1. Basic HTTP server with port and directory configuration loaded from .env.
2. First-time onboarding wizard (/setup) and PIN profile authentication with brute-force rate limiting.
3. Embedded multi-page web client using Go templates and vanilla JavaScript with soft dark and lavender styling.
4. Media streaming endpoint supporting HTTP Range and Partial Content transfers.
5. In-place library scanning and zero-memory admin file uploads.
6. Local-first scraping engine with TMDB fallback and admin manual override.
7. Tailored media players (Plyr for video, native PDF iframe, ePub.js for books) and subtitle delivery.
8. Streamlined 3-table SQLite catalog tracking files, metadata, and watch progress.
9. Unit and integration tests.

---

Last updated: 2026, September 19th  
Author: Wesley Esquivel.
