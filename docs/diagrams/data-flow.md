# Data Flow & Architecture Diagrams

This document visualizes the internal and external data flows using standard UML sequence diagrams and Data Flow Diagrams (DFD) rendered with Mermaid.

---

## 1. Context Data Flow Architecture (Level 0 DFD)

Illustrates the high-level boundaries between external actors, upstream providers, host storage tiers, and the core Flan Media Server daemon process.

```mermaid
flowchart LR
    subgraph Clients["Clients & Users"]
        direction TB
        Browser["Client Browsers<br/>(TV, Mobile, Desktop)"]
        Admin["Administrator<br/>(Setup & Intake)"]
    end

    subgraph Core["Flan Media Server Process"]
        Daemon["Flan Daemon (:4907)<br/>• net/http & HTML Engine<br/>• Stream Governor Semaphore<br/>• SQLite WAL Persistence"]
    end

    subgraph External["External Cloud Services"]
        TMDB["TMDB / OpenLibrary<br/>(Metadata & Covers)"]
    end

    subgraph Storage["Host Storage Tiers"]
        direction TB
        AppData["App Storage (NVMe / SD)<br/>• flan.db (WAL Mode)<br/>• data/covers/"]
        BulkMedia["Bulk Storage (USB / SATA)<br/>• Video & Document Libraries<br/>• Sidecar Subtitles"]
    end

    %% Client Interactions
    Browser -->|"1. HTTP Range & Progress"| Daemon
    Daemon -->|"2. 206 Partial (sendfile) & UI"| Browser

    %% Admin Interactions
    Admin -->|"3. Setup & Multipart Uploads"| Daemon
    Daemon -->|"4. Admin Views & Status"| Admin

    %% Upstream Metadata
    Daemon <-->|"5. Throttled HTTPS Queries & Posters"| TMDB

    %% Storage Interactions
    Daemon <-->|"6. Read/Write SQLite & Covers"| AppData
    Daemon <-->|"7. Zero-Copy Reads & Chunked Writes"| BulkMedia
```

### Context Data Flow Specification

| Ref | Channel | Direction | Protocol / Mechanism | Description & Payload |
| :--- | :--- | :--- | :--- | :--- |
| **1** | Client Inbound | Browser → Daemon | HTTP/1.1 (TCP :4907) | Page navigation, user PIN authentication, HTTP Range requests (`bytes=start-end`), and throttled playback progress (`POST /api/progress`). |
| **2** | Client Outbound | Daemon → Browser | HTTP 200 / 206 Partial | Server-rendered HTML templates with embedded CSS/JS, WebVTT subtitles, and zero-copy byte ranges streamed via Linux `sendfile`. |
| **3** | Admin Inbound | Admin → Daemon | HTTP/1.1 | Initial onboarding setup (`/setup`), library scan triggers (`/api/scan`), and multipart file uploads streamed in 32kb chunks (`/api/upload`). |
| **4** | Admin Outbound | Daemon → Admin | HTTP 200 / JSON | Server settings interface, hardware storage stats, scan task progress, and upload status responses. |
| **5** | Upstream Scraper | Daemon ↔ TMDB | HTTPS (Throttled ~2.8 req/s) | Outbound cleaned title queries to TMDB/OpenLibrary; inbound JSON metadata (ratings, overviews, genres) and streamed cover images. |
| **6** | App Persistence | Daemon ↔ App Storage | POSIX File I/O & SQLite WAL | Fast random I/O: SQLite transactions (`flan.db`) with 2mb page cache, `.flan-keep` mount verification, and local artwork storage (`data/covers/`). |
| **7** | Media Storage | Daemon ↔ Bulk Storage | Linux VFS / Kernel `sendfile` | Sequential media reads via kernel zero-copy transfer to network sockets, and direct socket-to-disk 32kb writes during admin uploads. |

---

## 2. Component Data Flow Diagram (Level 1 DFD)

Breaks down the server process into distinct internal modules, illustrating how data moves between network listeners, middlewares, handlers, and persistence stores.

```mermaid
flowchart LR
    subgraph Input["Network Input"]
        Request["Incoming HTTP Request"]
    end

    subgraph Core["Internal Processing Modules"]
        Governor["1.0 Stream Governor (Semaphore)"]
        Router["2.0 Router & Auth Middleware"]
        PageEngine["3.0 Template Engine (html/template)"]
        StreamEngine["4.0 Streaming Engine (http.ServeContent)"]
        UploadEngine["5.0 Upload Handler (MultipartReader)"]
        Scanner["6.0 Background Scanner & Scraper"]
    end

    subgraph Stores["Persistence & Hardware Stores"]
        SQLite[("SQLite Database: flan.db")]
        Covers[("Local Cover Cache: data/covers/")]
        HostDrives[("Host Filesystem: /mnt/...")]
    end

    subgraph Output["Network Output"]
        Response["HTTP Response (HTML, JSON, Video Bytes)"]
    end

    Request --> Governor
    Governor --> Router

    Router -- "GET /stream/{id}" --> StreamEngine
    Router -- "GET / (Pages)" --> PageEngine
    Router -- "POST /api/upload" --> UploadEngine

    StreamEngine --> SQLite
    StreamEngine --> HostDrives
    StreamEngine --> Response

    PageEngine --> SQLite
    PageEngine --> Covers
    PageEngine --> Response

    UploadEngine --> HostDrives
    UploadEngine --> SQLite
    UploadEngine --> Response

    Scanner --> HostDrives
    Scanner --> SQLite
    Scanner --> Covers
```

---

## 3. Zero-Copy Video Streaming Data Flow (UML Sequence)

Shows the exact call path for a media stream request, highlighting why video data bypasses Go's garbage-collected application heap.

```mermaid
sequenceDiagram
    autonumber
    actor Client as Browser (HTML5 Video)
    participant Gov as Stream Governor (Semaphore)
    participant Handler as Stream Handler (Go)
    participant DB as SQLite (flan.db)
    participant VFS as Linux Kernel VFS
    participant Socket as Network TCP Socket

    Client->>Gov: GET /stream/42 with Range: bytes=1048576-
    Gov->>Gov: Acquire stream token (active <= 3)
    Gov->>Handler: Forward request

    Handler->>DB: Query file path for media_id = 42
    DB-->>Handler: Return /mnt/hdd1/movies/Dune.mp4

    Handler->>VFS: os.Open("/mnt/hdd1/movies/Dune.mp4")
    VFS-->>Handler: Return file descriptor (fd_in)

    Handler->>Handler: Calculate Content-Range and byte count
    Handler->>VFS: Invoke sendfile(fd_out, fd_in, offset, count)

    Note over VFS,Socket: Kernel transfers pages directly from<br/>filesystem cache to network socket.<br/>Zero bytes enter Go user-space RAM!

    VFS-->>Socket: Stream raw bytes to client
    Socket-->>Client: HTTP 206 Partial Content

    Client->>Handler: Connection closed or stream finished
    Handler->>Gov: Release stream token
```

---

## 4. Ingestion & Scraping Data Flow (UML Sequence)

Illustrates how newly added files are discovered, checked for local artwork, matched online, and stored locally.

```mermaid
sequenceDiagram
    autonumber
    participant Scanner as Library Scanner (Go)
    participant Disk as Media Storage (/mnt/...)
    participant TMDB as TMDB API (HTTPS)
    participant Covers as Local Cover Dir (data/covers/)
    participant DB as SQLite (flan.db)

    Scanner->>Disk: fs.WalkDir(mediaRoot)
    Disk-->>Scanner: File discovered: Dune.Part.Two.2024.1080p.mp4

    Scanner->>DB: Check if file_path exists in media_items
    DB-->>Scanner: Not found (New file)

    Scanner->>Disk: Check for local poster.jpg in same folder
    alt Local poster.jpg exists
        Disk-->>Scanner: Found local poster.jpg
        Scanner->>DB: Insert media_items with local cover path
    else No local poster found
        Scanner->>Scanner: Regex clean: Title="Dune Part Two", Year=2024
        Scanner->>TMDB: Search query via HTTPS (throttled at 2.8 req/s)
        TMDB-->>Scanner: Return canonical metadata, rating, genres, poster URL
        Scanner->>Covers: Stream image bytes directly to disk (io.Copy)
        Scanner->>DB: Insert media_items with cached cover, rating, and genres
    end
```

---

## 5. Client Progress Tracking Data Flow

Illustrates the state sync loop between client playback and the SQLite database.

```mermaid
sequenceDiagram
    autonumber
    actor User as User Playing Video
    participant Player as Plyr Video Player (JS)
    participant API as Progress API Handler (Go)
    participant DB as SQLite (flan.db)

    User->>Player: Watching video
    loop Every 5 seconds (throttled)
        Player->>API: POST /api/progress {media_id: 42, position_seconds: 1450, total: 9600}
        API->>API: Extract user_id from HMAC signed cookie
        API->>DB: Execute atomic UPSERT into playback_progress
        DB-->>API: Row updated
        API-->>Player: HTTP 200 OK
    end

    User->>Player: Reaches end of video (ended event)
    Player->>API: POST /api/progress {media_id: 42, is_finished: 1}
    API->>DB: Mark is_finished = 1 in playback_progress
    DB-->>API: Row updated
    API-->>Player: HTTP 200 OK
```

---

## 6. Admin Zero-Memory File Upload Data Flow

Illustrates how large video files are uploaded through the browser and streamed directly to disk without bloating RAM.

```mermaid
sequenceDiagram
    autonumber
    actor Admin as Admin Browser
    participant UploadHandler as Upload API (Go)
    participant OS as Linux System (statfs)
    participant Disk as Target Drive (/mnt/hdd1/...)
    participant DB as SQLite (flan.db)

    Admin->>UploadHandler: POST /api/upload (multipart form stream)
    UploadHandler->>UploadHandler: Verify admin role from signed cookie
    UploadHandler->>OS: Query free space on destination filesystem via statfs

    alt Free space < 2 GB
        UploadHandler-->>Admin: HTTP 507 Insufficient Storage
    else Free space adequate
        loop For each file part in multipart stream
            UploadHandler->>Disk: os.OpenFile(destinationPath, O_CREATE|O_WRONLY, 0644)
            UploadHandler->>Disk: io.Copy in 32kb buffers (Socket -> Disk)
            Disk-->>UploadHandler: File write complete
            UploadHandler->>DB: Index newly created file in media_items
        end
        UploadHandler-->>Admin: HTTP 200 OK (Upload Successful)
    end
```

---

### Related Documentation

+ [User Flows & Journey Diagrams](user-flows.md)
+ [Master System Specifications](../design.md)
+ [Five-Zone Rate Limiting Architecture](../rate-limiting.md)
+ [Storage Architecture & Drive Resiliency](../storage.md)
