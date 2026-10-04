# Web Client Pages & Template Specifications

Layouts, wireframes, and interaction behaviors for the 9 server-rendered HTML templates in Flan Media Server.

---

## 1. Template Inventory

A dedicated top header bar (`.top-header`, 56px) provides persistent brand identity and utility controls (`?` manual button and user profile avatar) across all authenticated views. Platform navigation is cleanly consolidated into the left sidebar rail on desktop (bottom rail on mobile).

* **`login.html` (`GET /login`):** Focused tactile keypad with 1-tap household profile avatar buttons, numeric PIN input, and instant access button.
* **`video.html` (`GET /` or `GET /video`):** Video catalog grid with Continue Watching top shelf, quick category chips, and media cards.
* **`search.html` (`GET /search`):** Full-canvas interactive search page with real-time query matching across videos and books, category/status filter chips, and rich media cards.
* **`video_detail.html` (`GET /video/{id}`):** Video detail view with poster, synopsis, playable episode/cut rows, and signed VLC download links.
* **`books.html` (`GET /books`):** Book catalog grid displaying single volumes and multi-volume series.
* **`book_detail.html` (`GET /books/{id}`):** Book detail view with cover, author, overview, and download/PDF reading buttons.
* **`watch.html` (`GET /watch/{file_id}`):** Video player view with Plyr, auto-resume prompt, periodic progress sync, and VLC audio fallback banner.
* **`read.html` (`GET /read/{file_id}`):** Document viewer for PDFs (native browser view) and direct file download for EPUBs.
* **`manage.html` (`GET /manage`):** Administration console with runtime metrics, library rescan trigger, and user profile management.
* **`manual.html` (`GET /manual`):** Offline user handbook with navigation table of contents, storage guide, and CLI administrator recovery instructions.

---

## 2. Template Specifications & Wireframes

### Page 1: 2-Step Sequential Login (`login.html` - `GET /login`)

#### Step 1: Profile Selection ("Who is watching?")
```text
+-------------------------------------------------------------------------+
| Flan Media Server                                                       |
+-------------------------------------------------------------------------+
|                                                                         |
|                +---------------------------------------+                |
|                |           Flan Media Server           |                |
|                |           Who is watching?            |                |
|                |                                       |                |
|                |  +---------+  +---------+  +---------+|                |
|                |  | [Icon]  |  | [Icon]  |  | [Icon]  ||                |
|                |  |  Mike   |  | Wesley  |  |  Guest  ||                |
|                |  | [User]  |  | [Admin] |  | [User]  ||                |
|                |  +---------+  +---------+  +---------+|                |
|                +---------------------------------------+                |
|                                                                         |
+-------------------------------------------------------------------------+
```

#### Step 2: Dedicated PIN Entry
```text
+-------------------------------------------------------------------------+
| Flan Media Server                                                       |
+-------------------------------------------------------------------------+
|                                                                         |
|                +---------------------------------------+                |
|                | [Icon]  Wesley                        |                |
|                |         [Admin]                       |                |
|                |---------------------------------------|                |
|                | Enter your 4-digit PIN                |                |
|                |                                       |                |
|                | +-----------------------------------+ |                |
|                | | • • • •                           | |                |
|                | +-----------------------------------+ |                |
|                |                                       |                |
|                | +-----------------------------------+ |                |
|                | | Access Library →                  | |                |
|                | +-----------------------------------+ |                |
|                |                                       |                |
|                | +-----------------------------------+ |                |
|                | | ← Switch Profile                  | |                |
|                | +-----------------------------------+ |                |
|                +---------------------------------------+                |
|                                                                         |
+-------------------------------------------------------------------------+
```
* **Step 1 (Zero-BS Profile Picker):** Unauthenticated visitors see only the household profile cards under *"Who is watching?"*. No form inputs or extraneous buttons compete for attention.
* **Step 2 (Focused PIN Entry):** Selecting an avatar transitions to a focused screen displaying the selected user's badge, numeric PIN input (`inputmode="numeric"`), instant submit button, and a **`[ ← Switch Profile ]`** return button.
* **Keyboard Navigation:** Tab and Enter select profiles in Step 1; Escape or Switch button returns from Step 2.
* **First-Run Onboarding:** If zero users exist in the database, the server redirects to `/setup`, requiring the 6-character terminal bootstrap setup token output to host stdout/journal before creating the administrator account.

---

### Page 2: Video Catalog (`video.html` - `GET /video`)
```text
+-----------------------------------------------------------------------------------------------+
| FLAN MEDIA SERVER (70px Height, Bold 1.7rem, Title Only)                 (?) [Avatar Frame]  |
+-----------+-----------------------------------------------------------------------------------+
|           |                                                                                   |
|  [  Q  ]  |   === CONTINUE WATCHING =======================================================   |
|  SEARCH   |   +--------------------------+   +--------------------------+                     |
|           |   | [Thumb] Breaking Bad     |   | [Thumb] Spirited Away    |                     |
|  [ |> ]   |   |         S01E02           |   |         Movie            |                     |
|   VIDEO   |   |         [ In Progress ]  |   |         [ In Progress ]  |                     |
|  (Active) |   |         [ Resume > ]     |   |         [ Resume > ]     |                     |
|           |   +--------------------------+   +--------------------------+                     |
|  [ [] ]   |                                                                                   |
|   BOOKS   |   VIDEOS & MOVIES   [ All ] [ Movies ] [ Series ] [ In Progress ]      (23 items)  |
|           |   +-------+ +-------+ +-------+ +-------+ +-------+ +-------+ +-------+           |
|   -----   |   | Poster| | Poster| | Poster| | Poster| | Poster| | Poster| | Poster|           |
|           |   |-------| |-------| |-------| |-------| |-------| |-------| |-------|           |
|  [  #  ]  |   | Title | | Title | | Title | | Title | | Title | | Title | | Title |           |
|  MANAGE   |   +-------+ +-------+ +-------+ +-------+ +-------+ +-------+ +-------+           |
|           |                                                                                   |
| 108px Rail|   [ <- Prev ]   [ 1 ]   [ 2 ]   [ Next -> ]    Showing 1-21 of 23                 |
+-----------+-----------------------------------------------------------------------------------+
```
* **Persistent Top Header:** 70px top header displays brand title on the left (text only, no emblem or logo icon) and quick utility controls (`(?)` manual and user avatar badge) on the right.
* **Discord x YouTube Hybrid Rail:** 108px chunky industrial rail featuring 60px squircle tiles, 30px SVG icons, 6px left-edge indicator pills (`|`), and 12px stacked text labels (`SEARCH`, `VIDEO`, `BOOKS`, `MANAGE`).
* **Continue Watching Top Shelf:** Surfaces active resume sessions sorted by activity date descending (`last_watched_at DESC`). Features lightweight `[ In Progress ]` status tags and instant 1-tap `[ Resume ]` triggers without heavy progress bars or duration calculations.
* **Canvas Immediate Gratification:** No search bar canyon consuming vertical space. The shelf sits immediately at the top of the canvas upon loading.
* **Inline Category Chips:** Clean `.catalog-header-bar` with section title, item counter badge, and 1-click filter chips (`[ All ] [ Movies ] [ Series ] [ In Progress ]`).
* **Dedicated Search Route:** Full search is triggered via the rail `[ Search ]` button or keyboard shortcuts (`/`, `Ctrl+K`), navigating directly to the dedicated Search page (`/search`).
* **Fluid 7-Column Scaling (Low Zoom / Large Displays):** Rather than leaving empty lavender void on 4K, 1440p, or low zoom levels (80%, 67%, 50%), `--catalog-max-width` progressively expands:
  * `1680px`: `1840px` max width (cards scale to ~235px wide)
  * `1920px`: `2060px` max width (cards scale to ~270px wide)
  * `2200px`: `2280px` max width (cards scale to ~305px wide)
  * `2560px`+: `2400px` max width (cards scale to ~325px wide)
  Content caps strictly at `2400px` to prevent infinite stretching or UI distortion at 50% zoom while keeping cards bold, prominent, and readable.
* **21 Items Max / Page Pagination:** Strict chunking to 21 items max per page ($7 \times 3 = 21$) with tactile numbered controls.

---

### Page 3: Search Page (`search.html` - `GET /search`)
```text
+-----------------------------------------------------------------------------------------------+
| FLAN MEDIA SERVER (70px Height, Bold 1.7rem, Title Only)                 (?) [Avatar Frame]  |
+-----------+-----------------------------------------------------------------------------------+
|           |                                                                                   |
| | [  Q  ] |   +---------------------------------------------------------------------------+   |
|   SEARCH  |   | [Q] Search across videos, movies, series, books...                    [X] |   |
|  (Active) |   +---------------------------------------------------------------------------+   |
|           |                                                                                   |
|   [ |> ]  |   TYPE: [ All ] [ Videos ] [ Books ]    STATUS: [ All ] [ In Progress ]           |
|   VIDEO   |   Found 18 titles matching "dune"                                                 |
|           |                                                                                   |
|   [ [] ]  |   +-------+ +-------+ +-------+ +-------+ +-------+ +-------+ +-------+           |
|   BOOKS   |   | Poster| | Poster| | Cover | | Poster| | Cover | | Poster| | Poster|           |
|           |   |-------| |-------| |-------| |-------| |-------| |-------| |-------|           |
|   -----   |   | Title | | Title | | Title | | Title | | Title | | Title | | Title |           |
|           |   +-------+ +-------+ +-------+ +-------+ +-------+ +-------+ +-------+           |
|   [  #  ] |                                                                                   |
|   MANAGE  |                                                                                   |
+-----------+-----------------------------------------------------------------------------------+
```
* **First-Class Responsive Page:** Replaces cramped modal dialogs with a full-canvas, responsive search experience without text truncation or nested scrollbars.
* **Unified Navigation:** Rail button `[ Q ] Search` and mobile bottom bar tab navigate to `/search`, highlighting the search tab with the active indicator pill and squircle latch.
* **Full-Width Search Bar:** 56px neo-brutalist search input with instant clear button `[ X ]` and auto-focus on page load.
* **Query-Driven Results:** Results appear once the user types a query. When the query is empty, displays an inviting initial prompt card (`Search Your Library`) with zero card clutter.
* **Standard Media Card Reuse:** Results reuse the exact catalog `.media-card` component (upper `.card-poster` with format badge and placeholder/artwork, lower purple `.card-footer-band` with title and episode/author subtitle; zero Drive badges).
* **Filter Chips:** 1-click toggles for media type (`All`, `Videos`, `Books`) and status (`All`, `In Progress`, `Unwatched`).
* **Zero-Match State:** Clean empty-state card with a 1-tap `[ Clear Search & Filters ]` recovery button.

---

### Page 4: Video Detail View (`video_detail.html` - `GET /video/{id}`)
```text
+-------------------------------------------------------------------------------------+
| Flan Media Server                                                   ( ? )  [Avatar] |
+---------+---------------------------------------------------------------------------+
|         |  <- Back to Videos                                                        |
| | [ |> ]|                                                                           |
|   VIDEO |  +-------------+   BREAKING BAD (2008)   [ Edit Details ]                 |
|         |  | [Poster]    |   A high school chemistry teacher diagnosed              |
|   [ [] ]|  |             |   with lung cancer turns to manufacturing...             |
|   BOOKS |  +-------------+                                                          |
|         |                                                                           |
|   ---   |  PLAYABLE ITEMS / EPISODES                                                |
|         |  +----------------------------------------------------------------------+ |
|   [ # ] |  | 1. S01E01 - Pilot (48m)   [== Watched ==] [Play] [VLC]               | |
|  MANAGE |  +----------------------------------------------------------------------+ |
|         |  | 2. S01E02 - Cat's (48m)   [== 18m left =] [Play] [VLC]               | |
|         |  +----------------------------------------------------------------------+ |
+---------+---------------------------------------------------------------------------+
```
* **Admin Edit Details Button:** For administrators, displays `[ Edit Details ]` next to the container title. Opens the Video Editor modal to adjust Title, Release Year, Overview, upload/replace cover art, and rename or reorder individual episodes without touching files on disk.
* **Playable List:** Displays container files, duration, watch progress, `[ Play ]`, and `[ VLC / Download ]`.
* **Direct Streaming:** Clicking Play opens `/watch/{file_id}`.
* **VLC Signed Link & Format Awareness:** Clicking `[ VLC / Download ]` provides a 4-hour signed download URL (`/download/video/{id}?exp=...&u=...&sig=...`) to stream in external players without needing browser session cookies. For MKV container files (which lack native browser support in Chrome, Safari, and iOS), the row highlights `[ VLC / Download ]` with an informative format badge to guide users to seamless playback.

---

### Page 5: Books Catalog (`books.html` - `GET /books`)
```text
+-------------------------------------------------------------------------------------+
| Flan Media Server                                                   ( ? )  [Avatar] |
+---------+---------------------------------------------------------------------------+
|         |                                                                           |
|   [ Q ] |   === JUMP BACK IN ===================================================    |
|  SEARCH |   +--------------------------+   +--------------------------+             |
|         |   | [Cover] Dune             |   | [Cover] Snow Crash       |             |
|   [ |> ]|   |         Frank Herbert    |   |         Neal Stephenson  |             |
|   VIDEO |   |         [ Reading ]      |   |         [ Reading ]      |             |
|         |   |         [ Continue > ]   |   |         [ Continue > ]   |             |
| | [ [] ]|   +--------------------------+   +--------------------------+             |
|   BOOKS |                                                                           |
|  (Active)   BOOKS & PUBLICATIONS  [ All ] [ EPUB ] [ PDF ] [ Reading ]   (13 items) |
|   ---   |   +-------+ +-------+ +-------+ +-------+ +-------+ +-------+ +-------+   |
|         |   | Cover | | Cover | | Cover | | Cover | | Cover | | Cover | | Cover |   |
|   [ # ] |   |-------| |-------| |-------| |-------| |-------| |-------| |-------|   |
|  MANAGE |   | Title | | Title | | Title | | Title | | Title | | Title | | Title |   |
|         |   +-------+ +-------+ +-------+ +-------+ +-------+ +-------+ +-------+   |
|         |                                                                           |
|         |   [ <- Prev ]   [ 1 ]   [ Next -> ]    Showing 1-13 of 13                 |
+---------+---------------------------------------------------------------------------+
```
* **Persistent Top Header & Hybrid Rail Navigation:** 70px top header and 108px left rail sidebar with `Books` active.
* **Jump Back In Shelf:** Shows currently reading books directly at the top of the canvas with clean `[ Reading ]` status tags and 1-tap `[ Continue > ]` triggers.
* **Inline Category Chips:** Quick format switches (`[ All ] [ EPUB ] [ PDF ] [ Reading ]`) directly on the `.catalog-header-bar`.
* **Dedicated Search Route:** Full book search triggered via the `[ Search ]` rail button or `/` / `Ctrl+K`.
* **Fluid Scaled Grid:** 7-column layout scaling progressively on large screens / low zoom levels up to the 2400px boundary.
* Clicking any book card navigates to `/books/{id}`.

---

### Page 6: Book Detail View (`book_detail.html` - `GET /books/{id}`)
* Displays book cover, title, author, and description.
* Lists readable volume files:
  * For **PDFs**: `[ Open PDF ]` (opens browser viewer tab), `[ Download ]`, and status toggle (`Reading` / `Finished`).
  * For **EPUBs**: `[ Download EPUB ]` (direct file download) and status toggle (`Reading` / `Finished`).
* Reading progress is saved as discrete state (`unread`, `reading`, `finished`).

---

### Page 7: Video Player View (`watch.html` - `GET /watch/{file_id}`)
```text
+-------------------------------------------------------------------------+
| [← Back]  Breaking Bad - S01E02                                         |
+-------------------------------------------------------------------------+
|                                                                         |
|   [ RESUME PLAYBACK ]                                                   |
|   You were watching at 24:12.                                           |
|   [ Resume from 24:12 ]          [ Start from Beginning ]               |
|                                                                         |
|-------------------------------------------------------------------------|
|  [Play] [-10s] [+10s]  00:24:12 / 00:48:00       [1.0x] [Mute] [Full]   |
+-------------------------------------------------------------------------+
| Audio silent or stuttering? [ Open in External VLC Player / Download ]  |
+-------------------------------------------------------------------------+
```
* **Clean Player Canvas:** Top bar auto-hides after 3 seconds of mouse inactivity.
* **Auto-Resume:** Displays resume prompt if previous position exceeds 10 seconds.
* **Audio & Container Fallback:** Provides a 4-hour signed VLC stream link below the canvas. If the browser throws `MEDIA_ERR_SRC_NOT_SUPPORTED` (common for MKV containers or unsupported AC3/DTS audio codecs), the player renders a clear fallback overlay: *"Your browser cannot decode this video container/codec directly without transcoding. [ Launch in VLC / Stream Externally ]"* alongside a 1-tap *"Copy Stream Link"* button.

---

### Page 8: Document Reader View (`read.html` - `GET /read/{file_id}`)
* **PDF Mode:** Serves file with `Content-Type: application/pdf`, opening in the user's browser viewport.
* **EPUB Mode:** Direct file download with `Content-Disposition: attachment`.

---

### Page 9: Management Console (`manage.html` - `GET /manage`)
```text
+-------------------------------------------------------------------------+
| Flan Media Server                                           ( ? )   [Avatar]|
+----------+--------------------------------------------------------------+
|          |  MANAGE SERVER                                               |
| [ Video] |                                                              |
|          |  Media Storage Sources:         [ + Add Folder ] [Upload Files]|
| [ Books] |  • USB Movies (/mnt/usb1/movies)     [Video] [Scan] [Edit] [Remove]|
|          |  • TV Shows (/mnt/usb2/tv)           [Video] [Scan] [Edit] [Remove]|
| [Manage] |  • Books Library (./media/books)     [Books] [Scan] [Edit] [Remove]|
| (Active) |  [ ⟳ Rescan All Sources ]                                    |
|          |                                                              |
|          |  Household Profiles:                                         |
|          |  • [Icon] Wesley (Admin)  [ Change Avatar ]  [ Edit PIN ]    |
|          |  • [Icon] Mike (User)     [ Change Avatar ]  [ Reset PIN ]   |
|          |  [ + Add New User ]                                          |
|----------+--------------------------------------------------------------+
```
* **Clean Management Console:** Zero clutter. Focuses strictly on storage sources and household accounts (system performance metrics deferred to future dedicated sub-menu).
* **Scan Button:** Clicking `[ Scan ]` on any storage source launches the focused **Ingestion Pipeline Wizard** (`#ingestion-pipeline-modal`).
* **Edit Button (`#edit-source-modal`):** Clicking `[ Edit ]` allows renaming the folder's friendly display name in Flan without altering or renaming the physical directory on disk.
* **Remove Button & Safety Warning:** Clicking `[ Remove ]` triggers a safety warning dialog (`#remove-source-modal`) confirming that the storage folder will only be unregistered from Flan and no disk files will be deleted.
* **Storage Sources:** Displays all active storage mounts with media type and available disk space. Adding a source validates `.flan-keep` and automatically launches the pipeline to review discovered candidates.
* **Upload Files:** `[ Upload Files ]` allows direct in-browser upload of books or videos straight to a selected storage source with real-time direct streaming to disk.
* **Rescan Button:** Triggers a scan across all active storage sources.

---

### Edit Storage Folder Modal (`#edit-source-modal`)
```text
+-------------------------------------------------------------------------+
| Edit Storage Folder                                                 [X] |
+-------------------------------------------------------------------------+
| Folder Display Name:                                                    |
| [ USB Movies                                                          ] |
|                                                                         |
| Filesystem Path: (Read-only)                                            |
| /mnt/usb1/movies                                                        |
| (Physical directory on disk will not be changed or moved)               |
|                                                                         |
| [ Save Changes ]                                             [ Cancel ] |
+-------------------------------------------------------------------------+
```
* **Friendly Name Editing:** Safely updates `storage_sources.name` in SQLite.
* **Disk Safety Guarantee:** Path remains immutable to prevent disrupting disk mounts, torrent seeders, or permissions.

---

### Ingestion Pipeline Modal (`#ingestion-pipeline-modal`)
```text
+-------------------------------------------------------------------------+
| Ingest Media: USB Movies (/mnt/usb1/movies)                         [X] |
+-------------------------------------------------------------------------+
| Discovered: 4 containers, 22 files.                                     |
| [x] Select All (3 selected)               [ Filter items...           ] |
|-------------------------------------------------------------------------|
| [x] [Dir] the.wire.s01.720p-ctu (13 files)                            |
|     Title: [ The Wire                 ]   Year: [ 2002 ]  [View Files v]|
|-------------------------------------------------------------------------|
| [x] [Dir] alien.1979.directors.cut (2 files)                          |
|     Title: [ Alien                    ]   Year: [ 1979 ]  [View Files v]|
|-------------------------------------------------------------------------|
| [ ] [Dir] sample_clips_and_extras (5 files)                           |
|     (Unchecked - will be skipped)                                       |
|-------------------------------------------------------------------------|
| [ Ingest Selected Items (2) ]                                [ Cancel ] |
+-------------------------------------------------------------------------+
```
* **Manual Selection:** Checkboxes allow administrators to include or exclude specific containers.
* **Inline Typo Corrections:** Title and Release Year inputs allow fixing names right in the pipeline before committing to the database.
* **Pristine Database:** Unchecked or skipped items are never written to SQLite. Only user-approved items are committed with `metadata_locked = 1`.

---

### Page 10: Software Manual (`manual.html` - `GET /manual`)
```text
+-------------------------------------------------------------------------+
| [ ← Back to Library ]           User Guide                     Offline  |
+---------------------+---------------------------------------------------+
| Table of Contents   | 1. ABOUT FLAN                                     |
|                     | Self-hosted personal media server for streaming   |
| 1. About Flan       | videos and reading books across home devices.     |
| 2. Navigation       |                                                   |
| 3. Watching Videos  | 2. NAVIGATING THE INTERFACE                       |
| 4. Reading Books    | Top header tools and left rail / mobile tabs.     |
| 5. Adding Media     |                                                   |
| 6. Profiles & PINs  | 3. WATCHING VIDEOS & PLAYER CONTROLS              |
| 7. CLI Admin        | Direct play, VLC fallback, and keyboard keys.     |
|                     |                                                   |
|                     | 7. COMMAND-LINE ADMINISTRATION (CLI)              |
|                     | ./flan --reset-admin   [ Copy ]                   |
+---------------------+---------------------------------------------------+
```
* **Full-Width Canvas:** Sidebar is hidden to prioritize reading width.
* **Sticky Navigation:** Table of contents links directly to 7 numbered user-guide sections.
* **Terminal Snippets:** Commands include copy buttons.
* **User-Focused Guidance:** Practical, straightforward explanations for browsing, playback, reading, file storage, profiles, and CLI recovery without internal developer specs or buzzwords.

---

## 3. Accessibility & Keyboard Navigation (WCAG 2.1 AA)

1. **Skip-to-Content Link:** The first focusable element on authenticated pages is `<a href="#main-content" class="skip-link">Skip to main content</a>`.
2. **Tab Sequence:**
   * Header: Skip link $\to$ Brand title $\to$ Manual button (`?`) $\to$ User avatar.
   * Sidebar: `Video` $\to$ `Books` $\to$ `Manage Server` (active item flagged with `aria-current="page"`).
   * Content: Search input $\to$ Search button $\to$ Media cards in natural DOM order.
3. **Card Activation:** Cards are semantic links (`<a>`) with descriptive `aria-label` attributes, activated with the `Enter` key.
4. **Modal Dialogs:** Built using native `<dialog>` elements with focus trapping, `Escape` key dismissal, and focus restoration to the trigger button.
5. **Live Status Feedback:** Error alerts and rescan notifications use `role="alert"` and `aria-live="polite"` for screen readers.

---

## 4. Related Documentation

* [Design System](design-system.md)
* [Component Specifications](components.md)
* [Accessibility Guide](accessibility.md)
* [Mobile Responsiveness](responsiveness.md)
* [Master System Architecture](../design.md)
