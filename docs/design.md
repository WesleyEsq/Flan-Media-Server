# Master System Architecture

High-level architectural design and core technical decisions for Flan Media Server.

---

## 1. System Constraints & Operating Targets

Flan is built specifically for low-power Linux devices (Raspberry Pi Zero, 1–5, low-spec ARM/x86 boards) serving 1 to 5 household users with extreme simplicity.

| Parameter | Specification | Engineering Rationale |
| :--- | :--- | :--- |
| **Memory Budget** | 15–20 MB RAM | Soft ceiling via `GOMEMLIMIT=16MiB` and aggressive GC via `GOGC=30`. |
| **CPU Budget** | Low-power ARM cores | Real-time transcoding is omitted; all media streams directly. |
| **Navigation System** | **Sidebar-Only** | All platform navigation is consolidated strictly into the left sidebar (`Video`, `Books`, and `Manage Server`). The top header contains zero navigation links. |
| **Media Architecture** | Container + File List | Unifies Movies and TV Series into `videos` + `video_files`. Unifies Books into `books` + `book_files`. |
| **Network Footprint** | **100% Offline** | Zero external API calls (no TMDB). Folder name = title; `poster.jpg` = cover. |
| **Storage Paths** | Fixed Directories | `./media/video` and `./media/books`. Eliminates `libraries` table. |
| **Concurrent Streams** | Max 3 active | Enforced by counting semaphore to prevent USB mechanical drive thrashing. |
| **UI Aesthetics** | High-contrast tactile | Thick 2px black borders, purple header & sidebar, split Start screen, zero emojis, zero hover bloat. |
| **Deployment** | Single static binary | Pure-Go SQLite (`modernc.org/sqlite`, `CGO_ENABLED=0`), assets embedded via `embed.FS`. |

---

## 2. Software Architecture (MVC)

The codebase implements a clean Model-View-Controller pattern using Go 1.22+ standard library features:

```text
HTTP Request → [Middleware: Auth & Stream Semaphore] → Controller (RegisterRoutes) → Model (Store) → SQLite WAL
                                                      ↓
                                ViewModel Assembly → View (html/template) → Client Response
```

* **Model Layer (`internal/model`):** Domain structs (`Video`, `VideoFile`, `Book`, `BookFile`, `User`, `Progress`), sentinel errors (`ErrNotFound`, `ErrDuplicate`), and pure SQL queries for the 6-table schema.
* **View Layer (`web/templates`, `web/static`):** 9 server-rendered Go HTML templates and vanilla CSS/JS packaged directly into the binary via `embed.FS`.
* **Controller Layer (`internal/controller`):** HTTP transport adapters, request decoders, status code mapping, ViewModel assembly, and route registration onto `http.ServeMux`.
* **Middleware Layer (`internal/middleware`):** Minimalist cross-cutting filters: HMAC cookie authentication and the concurrent stream governor semaphore.

---

## 3. Key Technical Decisions

### Exclusive Navigation & Clean Header Utilities
* **No Header Navigation Links:** The top purple bar strictly contains the platform brand title (*"Flan Media Server"*) on the left, and utility controls on the right: `(?)` Software Manual & About view (`/manual`), and the User Profile Avatar. The redundant `⚙` settings cog has been removed to preserve a clutter-free header and prevent unauthorized configuration access for non-admin household users.
* **Left Sidebar (Desktop) & Bottom Navigation Bar (Mobile):** Serves as the single navigation backbone across the entire app with three tactile buttons:
  * **`Video`**: Navigates to `/video`. When active, fills with solid purple and white text.
  * **`Books`**: Navigates to `/books`. When active, fills with solid purple and white text.
  * **`Manage Server`**: Navigates to `/manage`.
  * **Mobile Adaptation:** On screens $\le 768\text{px}$, the sidebar seamlessly transitions via pure CSS into a fixed 56px bottom navigation bar for thumb ergonomics.
* **Unified "My Profile" Modal:** Clicking the user avatar opens an interactive profile dialog allowing household members to edit their display name, choose a whimsical SVG companion or upload a custom photo (max 2MB), update their 4-digit PIN, or log out via a prominent `[ Log Out ]` button.

### Split-Screen Start & Login (Image 1)
* Unauthenticated visits route to a clean split Start screen:
  * **Left:** Light lavender panel with a bold stylized *"Welcome"* graphic between arrows (stacks on mobile).
  * **Right:** Solid purple panel with `User:` dropdown (`<select>`), `Pin:` input field, and `[ Access ]` button.
* Replaces complex 3x4 on-screen keypad components with standard, accessible browser inputs.

### 100% Offline & Fixed Directory Layout
* Media is placed directly into:
  * `./media/video/<Container Title>/<filename>.<ext>` (with optional `poster.jpg`)
  * `./media/books/<Container Title>/<filename>.<ext>`
* No dynamic multi-library database tables, no TMDB API keys, and no external network timeouts.

### Direct Play, Zero-Copy Streaming & VLC Fallback
* **No Transcoding:** Videos stream in original direct-play formats to preserve SBC CPU and SD card life.
* **Linux `sendfile`:** Video delivery delegates to Go's standard library `http.ServeContent` with HTTP 206 Partial Content (Range requests), transferring bytes directly from filesystem cache to socket.
* **1-Click VLC / External Player Fallback:** Browsers lack built-in decoders for multi-channel AC3, EAC3, or DTS audio. To bypass this without CPU-heavy transcoding, every media card and the video player provide a `[ ⬇ Download / VLC ]` link. Users can paste the network stream directly into VLC/MPV or download the file instantly.
* **Whimsical Avatars:** Household users can personalize their accounts with 6 curated high-contrast SVG mascot presets or upload a custom photo avatar (max 2MB).

### Two-Safeguard Rate Limiting
1. **Stream Governor Semaphore (Safeguard 1):** Caps active streams at 3 (`HTTP 429` on exceed) to protect mechanical USB hard drive read heads.
2. **PIN Lockout Protection (Safeguard 2):** 5 failed attempts trigger a 5-minute lockout with exponential backoff on subsequent failures.

---

## 4. Server Endpoints Summary

### Web Pages (Server-Rendered HTML - 9 Templates)
| Method | Route | Description |
| :--- | :--- | :--- |
| `GET` | `/login` | Split Start screen (Welcome banner + User dropdown & PIN input) |
| `GET` | `/` or `/video` | Video catalog (uniform card grid with purple footers and search bar) |
| `GET` | `/video/{id}` | Video detail view showing cover, synopsis, playable files, and VLC/download links |
| `GET` | `/books` | Books catalog (uniform card grid) |
| `GET` | `/books/{id}` | Book detail view showing cover, author, overview, and direct download/PDF read buttons |
| `GET` | `/watch/{file_id}` | Video player page (Plyr) with auto-resume prompt & VLC audio fallback link |
| `GET` | `/read/{file_id}` | Document reader / download page (native browser PDF tab or direct EPUB download) |
| `GET` | `/manage` | Server status, rescan trigger for `./media`, and user profile management |
| `GET` | `/manual` | Dedicated Software Handbook & System Manual (offline docs, storage guide, CLI admin reset, specs) |

### Media Delivery & Streaming
| Method | Route | Description |
| :--- | :--- | :--- |
| `GET` | `/stream/video/{file_id}` | Zero-copy byte range delivery via `sendfile` (HTTP 206) |
| `GET` | `/stream/book/{file_id}` | Serves book file (PDF/EPUB) with Range support |
| `GET` | `/download/{type}/{file_id}` | Direct file download with `Content-Disposition: attachment` for VLC and offline reading |
| `GET` | `/covers/{type}/{id}` | Serves locally cached cover images with long-lived browser caching |
| `GET` | `/static/*` | Serves embedded CSS, JS, player assets, and bundled SVG avatars |

### JSON Management & State APIs
| Category | Method | Route | Description |
| :--- | :--- | :--- | :--- |
| **Auth** | `POST` | `/api/login` | Validates PIN, checks lockout, and issues signed cookie |
| | `POST` | `/api/logout` | Clears session cookie |
| **Users** | `GET` | `/api/users` | Lists profile names for Start screen dropdown |
| | `POST` | `/api/users` | Admin creates a new profile |
| | `PUT` | `/api/users/{id}` | Updates profile avatar preset icon (`mascot`, `flan`, `cat`, etc.) |
| | `POST` | `/api/users/{id}/avatar` | Uploads custom photo avatar (PNG/JPEG/WebP, max 2MB) |
| | `PUT` | `/api/users/{id}/pin`| Changes PIN and increments token version |
| | `DELETE`| `/api/users/{id}`| Admin deletes a user profile |
| **Media** | `POST` | `/api/scan` | Scans `./media/video` and `./media/books` |
| | `PUT` | `/api/media/{type}/{id}` | Edits title and overview synopsis |
| | `DELETE`| `/api/media/{type}/{id}`| Admin removes an item from catalog |
| **Progress**| `GET` | `/api/progress/{type}/{file_id}` | Retrieves saved playback or reading position |
| | `POST` | `/api/progress` | Syncs current playback or reading progress |

---

## 5. Architectural References

* **Data & Storage:** [Database Schema (6 Tables)](database.md) • [Storage Architecture](storage.md) • [Local-First Metadata Engine](scraper.md)
* **Security & Traffic:** [Two-Safeguard Rate Limiting](rate-limiting.md) • [Threat Model](threat-model.md)
* **Client & UI:** [Design System](client/design-system.md) • [Tactile Components](client/components.md) • [Page Templates (9 Views)](client/pages.md)
* **Operations:** [Compilation & Deployment](compilation.md) • [Testing Strategy](testing.md)
* **Diagrams:** [Data Flow](diagrams/data-flow.md) • [User Flows](diagrams/user-flows.md)
