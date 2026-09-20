# Data Flow & Architecture Diagrams

This document visualizes the internal and external data flows of Flan Media Server using standard UML sequence diagrams and Data Flow Diagrams (DFD) rendered with Mermaid.

---

## 1. Context Data Flow Diagram (Level 0 DFD)

Illustrates the high-level boundary between external entities (browsers, filesystem drives, and metadata providers) and the core Flan Media Server process.

```mermaid
flowchart TD
    subgraph Entities["External Entities"]
        Client["Client Browser (Phone, TV, Laptop)"]
        Admin["Administrator"]
        TMDB["The Movie Database (TMDB API)"]
        MediaDrives["Bulk Media Storage (USB HDD / SATA)"]
        AppData["Fast Storage (NVMe / SD Card)"]
    end

    subgraph Flan["Flan Media Server Process"]
        ServerCore["Flan Media Server Daemon"]
    end

    Client -- "1. HTTP Request (Page Views, Range Headers)" --> ServerCore
    ServerCore -- "2. Streamed Media Bytes (sendfile 206 Partial)" --> Client
    ServerCore -- "3. Server-Rendered HTML & Embedded Assets" --> Client
    Client -- "4. Playback Progress Sync (POST /api/progress)" --> ServerCore

    Admin -- "5. First-Time Setup & Settings Configuration" --> ServerCore
    Admin -- "6. Streaming File Uploads (POST /api/upload)" --> ServerCore

    ServerCore -- "7. Cleaned Title/Year Search (HTTPS)" --> TMDB
    TMDB -- "8. Metadata, Ratings & Cover Image URLs" --> ServerCore

    ServerCore -- "9. Direct Zero-Copy Read (Kernel VFS)" --> MediaDrives
    ServerCore -- "10. Stream Uploaded Files (32kb chunks)" --> MediaDrives

    ServerCore -- "11. Read/Write SQLite Queries (WAL Mode)" --> AppData
    ServerCore -- "12. Save Downloaded Artwork (data/covers/)" --> AppData
```

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
