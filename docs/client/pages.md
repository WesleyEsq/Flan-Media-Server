# Web Client Pages & Template Specifications

Layouts, wireframes, and interaction behaviors for the 9 server-rendered HTML templates in Flan Media Server.

---

## 1. Template Inventory

A dedicated top header bar (`.top-header`, 56px) provides persistent brand identity and utility controls (`?` manual button and user profile avatar) across all authenticated views. Platform navigation is cleanly consolidated into the left sidebar rail on desktop (bottom rail on mobile).

* **`login.html` (`GET /login`):** Focused tactile keypad with 1-tap household profile avatar buttons, numeric PIN input, and instant access button.
* **`video.html` (`GET /` or `GET /video`):** Video catalog grid with search bar and cards displaying solid purple title footers.
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
+-------------------------------------------------------------------------------------+
| Flan Media Server                                                   ( ? )  [Avatar] |
+----------+--------------------------------------------------------------------------+
|          |                                                                          |
| [ Video] |   +---------------------------------------+  [Search]  [ Filter [v] ]    |
| (Active) |   | Search videos...                      |                              |
|          |   +---------------------------------------+                              |
| [ Books] |                                                                          |
|          |   === CONTINUE WATCHING ==============================================   |
|          |   +--------------------------+   +--------------------------+            |
|          |   | [Thumb] Breaking Bad     |   | [Thumb] Spirited Away    |            |
|          |   |         S01E02 (24:12)   |   |         1:33:45 (75%)    |            |
|          |   |         [== Prog ==]     |   |         [==== Prog ==]   |            |
|          |   |         [ Resume > ]     |   |         [ Resume > ]     |            |
|          |   +--------------------------+   +--------------------------+            |
|          |                                                                          |
|          |   === ALL VIDEOS (7 Columns Max Hard Limit) ==========================   |
|          |   +-------+ +-------+ +-------+ +-------+ +-------+ +-------+ +-------+  |
|          |   | Poster| | Poster| | Poster| | Poster| | Poster| | Poster| | Poster|  |
|          |   |-------| |-------| |-------| |-------| |-------| |-------| |-------|  |
|          |   | Title | | Title | | Title | | Title | | Title | | Title | | Title |  |
|          |   +-------+ +-------+ +-------+ +-------+ +-------+ +-------+ +-------+  |
|          |                                                                          |
|          |   [ <- Prev ]   [ 1 ]   [ 2 ]   [ Next -> ]    Showing 1-21 of 23        |
| [Manage] |                                                                          |
+----------+--------------------------------------------------------------------------+
```
* **Persistent Top Header:** Restored 56px top header displays brand title on the left and quick utility controls (`(?)` manual and user avatar badge) on the right. Search remains scoped in the catalog canvas.
* **Rail Navigation:** Pure industrial buttons for `Video`, `Books`, and pinned `Manage Server` with zero decorative arrows and no distracting count badges.
* **Continue Watching Top Shelf:** Surfaces active resume sessions sorted by activity date descending (`last_watched_at DESC`). Complete with episode indicators, tactile progress bars, and instant 1-tap `[ Resume ]` triggers. Completely collapses when 0 items are in progress.
* **Collapsible Filter & Sort Drawer:** Toggled via `[ Filter [v] ]` beside the search bar. Provides format pills (`All`, `Movies`, `Series`), status pills (`All`, `In Progress`, `Unwatched`), storage drive filters, and sort options (including `Shuffle / Random` by default).
* **7-Column Max Hard Limit:** `.media-grid` is capped at a strict maximum of 7 columns on wide displays (`max-width: 1680px`), guaranteeing optimal readability.
* **21 Items Max / Page Pagination:** Strict chunking to 21 items max per page ($7 \times 3 = 21$) with tactile numbered controls.

---

### Page 3: Video Detail View (`video_detail.html` - `GET /video/{id}`)
```text
+-------------------------------------------------------------------------------------+
| Flan Media Server                                                   ( ? )  [Avatar] |
+----------+--------------------------------------------------------------------------+
|          |  <- Back to Videos                                                       |
| [ Video] |                                                                          |
| (Active) |  +-------------+   BREAKING BAD (2008)   [ Edit Details ]                |
|          |  | [Poster]    |   A high school chemistry teacher diagnosed             |
| [ Books] |  |             |   with lung cancer turns to manufacturing...            |
|          |  +-------------+                                                         |
|          |                                                                          |
|          |  PLAYABLE ITEMS / EPISODES                                               |
|          |  +---------------------------------------------------------------------+ |
|          |  | 1. S01E01 - Pilot (48m)   [== Watched ==] [Play] [VLC]              | |
|          |  +---------------------------------------------------------------------+ |
|          |  | 2. S01E02 - Cat's (48m)   [== 18m left =] [Play] [VLC]              | |
| [Manage] |  +---------------------------------------------------------------------+ |
+----------+--------------------------------------------------------------------------+
```
* **Admin Edit Details Button:** For administrators, displays `[ Edit Details ]` next to the container title. Opens the Video Editor modal to adjust Title, Release Year, Overview, upload/replace cover art, and rename or reorder individual episodes without touching files on disk.
* **Playable List:** Displays container files, duration, watch progress, `[ Play ]`, and `[ VLC / Download ]`.
* **Direct Streaming:** Clicking Play opens `/watch/{file_id}`.
* **VLC Signed Link:** Clicking `[ VLC / Download ]` provides a 4-hour signed download URL (`/download/video/{id}?exp=...&u=...&sig=...`) to stream in external players without needing browser session cookies.

---

### Page 4: Books Catalog (`books.html` - `GET /books`)
```text
+-------------------------------------------------------------------------------------+
| Flan Media Server                                                   ( ? )  [Avatar] |
+----------+--------------------------------------------------------------------------+
|          |                                                                          |
| [ Video] |   +---------------------------------------+  [Search]  [ Filter [v] ]    |
|          |   | Search books...                       |                              |
|          |   +---------------------------------------+                              |
| [ Books] |                                                                          |
| (Active) |   === JUMP BACK IN ==================================================    |
|          |   +--------------------------+   +--------------------------+            |
|          |   | [Cover] Dune             |   | [Cover] Snow Crash       |            |
|          |   |         Frank Herbert    |   |         Neal Stephenson  |            |
|          |   |         [ READING - p.142]   |         [ READING - p.88 ]|           |
|          |   |         [ Continue > ]   |   |         [ Continue > ]   |            |
|          |   +--------------------------+   +--------------------------+            |
|          |                                                                          |
|          |   === ALL BOOKS & PUBLICATIONS ======================================    |
|          |   +-------+ +-------+ +-------+ +-------+ +-------+ +-------+ +-------+  |
|          |   | Cover | | Cover | | Cover | | Cover | | Cover | | Cover | | Cover |  |
|          |   |-------| |-------| |-------| |-------| |-------| |-------| |-------|  |
|          |   | Title | | Title | | Title | | Title | | Title | | Title | | Title |  |
|          |   +-------+ +-------+ +-------+ +-------+ +-------+ +-------+ +-------+  |
|          |                                                                          |
|          |   [ <- Prev ]   [ 1 ]   [ Next -> ]    Showing 1-13 of 13                |
| [Manage] |                                                                          |
+----------+--------------------------------------------------------------------------+
```
* **Persistent Top Header & Rail Navigation:** `Books` button is active in the left rail sidebar. Top header provides instant access to help manual and user profile modal.
* **Jump Back In Shelf:** Shows currently reading books with discrete status tags (e.g., `READING - Page 142 of 680`). No arbitrary percentages.
* **Format Filters:** Allows toggling between `EPUB` and `PDF` files.
* Clicking any book card navigates to `/books/{id}`.

---

### Page 5: Book Detail View (`book_detail.html` - `GET /books/{id}`)
* Displays book cover, title, author, and description.
* Lists readable volume files:
  * For **PDFs**: `[ Open PDF ]` (opens browser viewer tab), `[ Download ]`, and status toggle (`Reading` / `Finished`).
  * For **EPUBs**: `[ Download EPUB ]` (direct file download) and status toggle (`Reading` / `Finished`).
* Reading progress is saved as discrete state (`unread`, `reading`, `finished`).

---

### Page 6: Video Player View (`watch.html` - `GET /watch/{file_id}`)
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
* **Periodic Sync:** Client script posts current position to `/api/progress` every 15 seconds (and upon pause or navigating away).
* **Audio Fallback Bar:** Provides a 4-hour signed VLC stream link if the browser lacks support for the file's audio track.

---

### Page 7: Document Reader View (`read.html` - `GET /read/{file_id}`)
* **PDF Mode:** Serves file with `Content-Type: application/pdf`, opening in the user's browser viewport.
* **EPUB Mode:** Direct file download with `Content-Disposition: attachment`.

---

### Page 8: Management Console (`manage.html` - `GET /manage`)
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

### Page 9: Software Manual (`manual.html` - `GET /manual`)
```text
+-------------------------------------------------------------------------+
| [ ← Back to Library ]       Flan Software Handbook             Offline  |
+---------------------+---------------------------------------------------+
| Table of Contents   | 1. MEDIA STORAGE & PLACEMENT                      |
|                     | Drop video folders in ./media/video and books in  |
| 1. Media Storage    | ./media/books via SMB, USB drive, or SCP.         |
| 2. VLC Fallback     | Click [ Rescan All Media ] to index new items.    |
| 3. Profiles & PINs  |                                                   |
| 4. CLI Admin Reset  | 2. DIRECT VIDEO PLAYBACK & VLC FALLBACK           |
| 5. System Specs     | Direct play via HTTP 206 range requests.          |
|                     | If browser audio fails on AC3/DTS, click the      |
|                     | [ VLC / Download ] button to stream raw.          |
|                     |                                                   |
|                     | 4. SERVER CLI ADMINISTRATION & FAILSAFE RESET     |
|                     | ./flan --reset-admin   [ Copy ]                   |
+---------------------+---------------------------------------------------+
```
* **Full-Width Canvas:** Sidebar is hidden to prioritize reading width.
* **Sticky Navigation:** Table of contents links directly to numbered sections.
* **Terminal Snippets:** Commands include copy buttons.

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
