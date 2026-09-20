# User Design & System Flow Diagrams

This document illustrates the key user flows and technical workflows for Flan Media Server.

---

## 1. First-Time Setup Wizard

When the server runs for the first time with an empty database, it guides the administrator through a single onboarding screen before locking the setup route.

```mermaid
flowchart TD
    Start["User visits server at http://<ip>:4907"] --> CheckDB{"Are there any users in database?"}
    CheckDB -- "No (Initial Boot)" --> RedirectSetup["Redirect to /setup"]
    CheckDB -- "Yes" --> ShowLogin["Redirect to Profile Selector /login"]

    RedirectSetup --> Form["Admin fills in profile name, PIN, and media library path"]
    Form --> Submit["Submit Setup Form"]

    Submit --> ValidatePath{"Does media path exist and is readable?"}
    ValidatePath -- "No" --> Error["Show error message on setup page"]
    Error --> Form

    ValidatePath -- "Yes" --> CreateAdmin["1. Hash PIN with bcrypt<br/>2. Create Admin user in SQLite<br/>3. Save primary library path in .env"]
    CreateAdmin --> StartScan["Start background library scan"]
    StartScan --> IssueSession["Issue HMAC-signed session cookie"]
    IssueSession --> Catalog["Redirect to Catalog /"]
```

---

## 2. Profile Selection and PIN Authentication

Returning users see a friendly profile selection screen ("Who is watching?") and log in using their numeric PIN with brute-force protection and signed session cookies.

```mermaid
sequenceDiagram
    autonumber
    actor User
    participant Browser as Client Browser
    participant Server as Flan Media Server
    participant DB as SQLite Database

    User->>Browser: Opens app
    Browser->>Server: GET /login
    Server->>DB: Query profile list (names and avatars)
    DB-->>Server: Return profiles
    Server-->>Browser: Render profile selector page

    User->>Browser: Selects profile and enters PIN
    Browser->>Server: POST /api/login (user_id, pin)

    Server->>Server: Check rate limiting for this IP and profile
    alt Rate limit exceeded (5+ failed attempts)
        Server-->>Browser: HTTP 429 Too Many Requests (Lockout for 5 mins)
    else Attempt allowed
        Server->>DB: Fetch user password hash
        DB-->>Server: Return bcrypt hash
        Server->>Server: Verify bcrypt hash against entered PIN

        alt PIN is incorrect
            Server->>Server: Increment failed attempts counter in DB
            Server-->>Browser: HTTP 401 Unauthorized (Invalid PIN)
        else PIN is correct
            Server->>Server: Reset failed attempts counter in DB
            Server->>Server: Generate HMAC-signed cookie (user_id + timestamp)
            Server-->>Browser: Set HttpOnly signed session cookie & redirect to /
        end
    end
```

---

## 3. Media Streaming and Progress Tracking

Demonstrates zero-copy streaming using the Linux sendfile system call and how playback progress is automatically saved to allow resuming later.

```mermaid
sequenceDiagram
    autonumber
    actor User
    participant Browser as HTML5 Video Player (Plyr)
    participant Server as Flan Media Server
    participant DB as SQLite Database
    participant Kernel as Linux Kernel / VFS

    User->>Browser: Clicks on a video card
    Browser->>Server: GET /watch/{id}
    Server->>DB: Query media details and saved position
    DB-->>Server: Return title, duration, last position
    Server-->>Browser: Render watch.html with Plyr player

    Note over Browser,Server: Browser requests initial video chunk
    Browser->>Server: GET /stream/{id} with Range: bytes=0-
    Server->>DB: Lookup file path for media ID
    DB-->>Server: Return /media/movies/flan.mp4
    Server->>Kernel: Call sendfile from file descriptor to socket
    Kernel-->>Browser: HTTP 206 Partial Content (streams video bytes)

    Note over Browser,Server: Playback progress syncs automatically
    loop Every 5 seconds during playback
        Browser->>Server: POST /api/progress (media_id, position_seconds)
        Server->>DB: UPSERT playback_progress
        DB-->>Server: Updated
        Server-->>Browser: HTTP 200 OK
    end

    User->>Browser: Stops video or closes tab
    Note over User,Browser: Next time user visits, video resumes from saved position
```

---

## 4. File Intake: Local In-Place Scan vs Admin Web Upload

Flan Media Server supports both referencing existing media collections on the host and uploading new folders through the browser without memory bloat.

```mermaid
flowchart TD
    subgraph Local["Method A: Local Directory Scanning (Existing Collections)"]
        A1["Admin adds folder path in .env e.g. /media"] --> A2["Server validates directory existence and read permissions"]
        A2 --> A3["Background worker walks directory tree recursively"]
        A3 --> A4["Check for local poster.jpg or cover.jpg first"]
        A4 --> A5["Extract file size, format, duration, and title"]
        A5 --> A6["Insert into media_items table (files stay in place on disk)"]
    end

    subgraph Remote["Method B: Admin Web Upload (Streaming Intake)"]
        B1["Admin drags folder into Web UI upload modal"] --> B2["Browser sends multipart stream via POST /api/upload"]
        B2 --> B3["Server checks available disk space via statfs"]
        B3 -- "Disk low (<2gb)" --> B4["Reject upload with HTTP 507 Insufficient Storage"]
        B3 -- "Space OK" --> B5["r.MultipartReader reads incoming file parts"]
        B5 --> B6["io.Copy streams directly from socket to destination file in 32kb chunks"]
        B6 --> B7["Save file with 0644 permissions (non-executable)"]
        B7 --> B8["Index new file into SQLite database"]
    end
```

---

## 5. Admin Account Recovery Flow

Demonstrates the straightforward host command-line failsafe for recovering the admin account if a PIN is forgotten.

```mermaid
flowchart TD
    Forgot["Admin forgot PIN"] --> Terminal["Access server shell via SSH or local terminal"]
    Terminal --> RunCLI["Run command: ./flan --reset-admin"]
    RunCLI --> Interactive["CLI prompts for new Admin PIN"]
    Interactive --> DirectDB["Updates admin record directly in flan.db"]
    DirectDB --> Done["Admin logs in with the new PIN"]
```
