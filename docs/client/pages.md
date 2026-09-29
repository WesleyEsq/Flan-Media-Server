# Web Client Pages & Template Specifications

Layouts, wireframes, and interaction behaviors for the 9 server-rendered Go HTML templates in Flan Media Server.

---

## 1. Template Inventory

| Template | Route | Primary Purpose | Key Features |
| :--- | :--- | :--- | :--- |
| **`index.html`** | `GET /` | Dashboard | "Continue Watching" & "Continue Reading" shelves, "Recently Added" grid. |
| **`videos.html`** | `GET /videos` | Video Catalog | Sub-tabs for Movies vs. TV Series, genre filter pills. |
| **`show.html`** | `GET /show/{id}` | TV Series View | Season selector tabs, episode rows with watch progress and checkmarks. |
| **`books.html`** | `GET /books` | Reading Catalog | Book and document grid with format badges (`EPUB`, `PDF`) and genre pills. |
| **`watch.html`** | `GET /watch/{type}/{id}`| Video Player | Blacked-out player with Plyr, auto-resume prompt, and 5s progress sync. |
| **`read.html`** | `GET /read/{id}` | Document Reader | Dual-mode viewer: native PDF iframe and paginated in-browser ePub.js. |
| **`login.html`** | `GET /login` | Profile Selector | "Who is watching?" avatar tiles and touch/keyboard PIN keypad. |
| **`setup.html`** | `GET /setup` | First-Time Wizard | Admin setup and initial library path (permanently locked once complete). |
| **`settings.html`**| `GET /settings` | Admin Console | Hardware stats, library roots, manual scan triggers, uploads, and users. |

---

## 2. Template Specifications & Wireframes

### Page 1: Dashboard (`index.html` - `GET /`)
```text
+-------------------------------------------------------------------------+
| [Flan]   Home   Videos   Books                   [? Help] [⚙] [Avatar]  |
+-------------------------------------------------------------------------+
|  CONTINUE WATCHING                                                      |
|  +----------------+  +----------------+  +----------------+             |
|  | [Poster]       |  | [Poster]       |  | [Poster]       |   [ > ]     |
|  | Dune Part Two  |  | Breaking Bad   |  | Spirited Away  |             |
|  | [== 35m left =]|  | S01E02 [=== 12m|  | [==== 1h left =]             |
|  +----------------+  +----------------+  +----------------+             |
|                                                                         |
|  CONTINUE READING                                                       |
|  +----------------+  +----------------+                                 |
|  | [Book Cover]   |  | [Book Cover]   |                       [ > ]     |
|  | Dune           |  | Go in Action   |                                 |
|  | [= Page 120 ==]|  | [== Page 45 ==]|                                 |
|  +----------------+  +----------------+                                 |
|                                                                         |
|  RECENTLY ADDED                                                         |
|  [ Card 1 ] [ Card 2 ] [ Card 3 ] [ Card 4 ] [ Card 5 ]                 |
+-------------------------------------------------------------------------+
```
* **Conditional Shelves:** Continue shelves render only if the user has items with `position_seconds > 10` and `is_finished = 0`.
* **CSS Snapping:** Shelves scroll horizontally via `scroll-snap-type: x mandatory`.
* **Empty State:** Shows a helpful prompt linking to `/settings` if the catalog contains zero items.

---

### Page 2: Videos Catalog (`videos.html` - `GET /videos`)
```text
+-------------------------------------------------------------------------+
| [Flan]   Home   Videos   Books                   [? Help] [⚙] [Avatar]  |
+-------------------------------------------------------------------------+
|  VIDEOS                                                                 |
|  [ ALL ]  [ MOVIES ]  [ TV SERIES ]                                     |
|  Genres: ( All ) ( Animation ) ( Comedy ) ( Sci-Fi ) ( Drama ) ( Action)|
|                                                                         |
|  +----------------+  +----------------+  +----------------+             |
|  | [2:3 Poster]   |  | [2:3 Poster]   |  | [2:3 Poster]   |             |
|  | Movie Title    |  | Series Title   |  | Movie Title    |             |
|  | 2024 • 1h 45m  |  | 2022 • 3 Seasons| 1999 • 2h 16m  |             |
|  | ★ 8.2          |  | ★ 8.9          |  | ★ 8.7          |             |
|  +----------------+  +----------------+  +----------------+             |
+-------------------------------------------------------------------------+
```
* **Catalog Grid:** Auto-wrapping grid (`repeat(auto-fill, minmax(160px, 1fr))`).
* **Sub-Tabs & Filtering:** Filters between movies and TV series; clicking genre pills appends `?genre=Name`.
* **Navigation:** Movie cards link to `/watch/movie/{id}`; series cards link to `/show/{id}`.

---

### Page 3: TV Series Detail View (`show.html` - `GET /show/{id}`)
```text
+-------------------------------------------------------------------------+
| [Flan]   Home   Videos   Books                   [? Help] [⚙] [Avatar]  |
+-------------------------------------------------------------------------+
|  ← Back to Videos                                                       |
|  +--------------+   BREAKING BAD (2008)                                 |
|  | [Poster]     |   ★ 9.5 • 5 Seasons • Crime, Drama                    |
|  |              |   Synopsis description...                             |
|  +--------------+                                                       |
|                                                                         |
|  [ Season 1 ]  [ Season 2 ]  [ Season 3 ]  [ Season 4 ]  [ Season 5 ]   |
|                                                                         |
|  EPISODES                                                               |
|  +--------------------------------------------------------------------+ |
|  | [Thumb] 1. Pilot (48m)                  [== Watched ==]  [ ✓ ]     | |
|  +--------------------------------------------------------------------+ |
|  | [Thumb] 2. Cat's in the Bag (48m)       [== 18m left ==] [ ► ]     | |
|  +--------------------------------------------------------------------+ |
+-------------------------------------------------------------------------+
```
* **Season Tabs:** Switches visible episodes in-place without full page reloads.
* **Progress Badges:** Displays watch progress bar per episode; completed items show checkmarks (`[ ✓ ]`).
* **Playback:** Clicking any episode row opens `/watch/episode/{id}`.

---

### Page 4: Books Catalog (`books.html` - `GET /books`)
```text
+-------------------------------------------------------------------------+
| [Flan]   Home   Videos   Books                   [? Help] [⚙] [Avatar]  |
+-------------------------------------------------------------------------+
|  BOOKS & DOCUMENTS                                                      |
|  Genres: ( All ) ( Sci-Fi ) ( Technology ) ( Documentation )            |
|                                                                         |
|  +----------------+  +----------------+  +----------------+             |
|  | [1:1.4 Cover]  |  | [1:1.4 Cover]  |  | [1:1.4 Cover]  |             |
|  | Dune           |  | Go in Action   |  | Linux Kernel   |             |
|  | Frank Herbert  |  | William Kennedy|  | Robert Love    |             |
|  | [ EPUB ]       |  | [ PDF ]        |  | [ PDF ]        |             |
|  +----------------+  +----------------+  +----------------+             |
+-------------------------------------------------------------------------+
```
* **Cover Styling:** 1:1.4 aspect ratio with faux book spine shadow on the left edge.
* **Format Badges:** Clear format indicators (`EPUB` / `PDF`). Cards link to `/read/{id}`.

---

### Page 5: Video Player View (`watch.html` - `GET /watch/{type}/{id}`)
```text
+-------------------------------------------------------------------------+
| [← Back]  Dune: Part Two (2024)                                         |
+-------------------------------------------------------------------------+
|                                                                         |
|   [ RESUME PLAYBACK ]                                                   |
|   You were watching at 24:12.                                           |
|   [ ▶ Resume from 24:12 ]        [ ↺ Start from Beginning ]             |
|                                                                         |
|-------------------------------------------------------------------------|
|  [▶] [10s↺] [10s↻]  00:24:12 / 02:46:00          [1.0x] [🔊] [⛶]        |
+-------------------------------------------------------------------------+
```
* **Distraction-Free:** Header bar auto-hides after 3 seconds of mouse inactivity.
* **Auto-Resume Prompt:** Prompts to resume if `position_seconds > 10`.
* **Sync Loop:** `player.js` syncs position to `POST /api/progress` every 5 seconds.
* **Codec Fallback:** Displays a download banner if the browser cannot decode stream streams directly.

---

### Page 6: Document Reader (`read.html` - `GET /read/{id}`)
* **PDF Mode:** Embedded directly using the native browser PDF engine via `<iframe src="/stream/book/{id}">` (zero JS bundle bloat).
* **EPUB Mode:** In-browser pagination via ePub.js with font scaling (`[ A- ] [ A+ ]`) and chapter navigation.
* **Download Button:** Prominent `[ ⬇ Download ]` button in the header across both modes for reading on external apps.

---

### Page 7: Profile Selector & Login (`login.html` - `GET /login`)
```text
+-------------------------------------------------------------------------+
|                                  Flan                                   |
|                           Who is watching?                              |
|                                                                         |
|         +-----------+          +-----------+          +-----------+     |
|         |  ( =^.^=) |          |  ( o.o )  |          |  [ 🍿 ]   |     |
|         |   Cat     |          |  Hamster  |          |  Popcorn  |     |
|         +-----------+          +-----------+          +-----------+     |
|            Wesley                 Guest                  Living Room    |
|                                                                         |
|                     +---------------------------+                       |
|                     | Enter PIN for Wesley      |                       |
|                     |        ●  ●  ○  ○         |                       |
|                     |    [ 1 ]   [ 2 ]   [ 3 ]  |                       |
|                     |    [ 4 ]   [ 5 ]   [ 6 ]  |                       |
|                     |    [ 7 ]   [ 8 ]   [ 9 ]  |                       |
|                     |    [ ⌫ ]   [ 0 ]   [ ↵ ]  |                       |
|                     +---------------------------+                       |
+-------------------------------------------------------------------------+
```
* **Avatar Grid:** Renders household profile tiles. Clicking a tile activates the numeric PIN keypad modal.
* **Keypad:** Accepts on-screen clicks or keyboard number keys.
* **Lockout:** Shows countdown timer if account is locked by Zone B rate limiting.

---

### Page 8: First-Time Setup Wizard (`setup.html` - `GET /setup`)
```text
+-------------------------------------------------------------------------+
|                                Welcome to Flan                          |
|                             Initial Server Setup                        |
|                                                                         |
|         +-----------------------------------------------------+         |
|         | 1. Admin Username:   [ Wesley                     ] |         |
|         | 2. Create PIN:       [ ****                       ] |         |
|         | 3. Media Directory:  [ /mnt/storage/media         ] |         |
|         |                                                     |         |
|         |    [ Complete Setup & Launch Catalog ]              |         |
|         +-----------------------------------------------------+         |
+-------------------------------------------------------------------------+
```
* **Bootstrapping:** Active only when `users` table is empty; permanently locked once complete.
* **Validation:** Verifies host path accessibility, hashes PIN with bcrypt, initializes admin, and starts background scan.

---

### Page 9: Server & Library Settings (`settings.html` - `GET /settings`)
```text
+-------------------------------------------------------------------------+
| [Flan]   Home   Videos   Books                   [? Help] [⚙] [Avatar]  |
+-------------------------------------------------------------------------+
|  SETTINGS                                                               |
|                                                                         |
|  System Status: Memory: 14.2 MB / 16 MB  •  Streams: 1 / 3 max          |
|  Storage:       Root SD: 42.1 GB Free    •  Media USB: 1.4 TB Free      |
|                                                                         |
|  Libraries:                                                             |
|  • /mnt/hdd1/movies (Movies) - 342 films                                |
|  • /mnt/hdd2/tv (TV Shows)   - 18 shows, 214 episodes                   |
|  • /mnt/nvme/books (Books)   - 48 books                                 |
|  [ ⟳ Scan All ]   [ + Add Library ]   [ ⬆ Upload Files/Folder ]         |
|                                                                         |
|  Metadata: TMDB API Key: [ ***************** ]  [ Save Key ]            |
|                                                                         |
|  User Profiles:                                                         |
|  • Wesley (Admin)  [ Edit PIN ]  [ Change Icon ]                        |
|  • Guest (User)    [ Reset PIN ] [ Delete ]                             |
|  [ + Add New Profile ]                                                  |
+-------------------------------------------------------------------------+
```
* **System Metrics:** Displays live memory usage, stream semaphore status, and mount free space via `statfs`.
* **Library Management:** Add/remove library roots and trigger scans (`POST /api/scan`).
* **Upload Modal:** Direct browser uploads streamed to disk in 32 KB chunks (<1 MB RAM).
* **Profiles:** Create users, change avatars, and reset PINs (which increments `token_version` to invalidate cookies).

---

## 3. Related Documentation

* [Design System & Foundations](design-system.md)
* [Component Specifications](components.md)
* [User Flows & Journeys](../diagrams/user-flows.md)
