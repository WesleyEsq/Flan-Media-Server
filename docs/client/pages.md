# Web Client Pages & Template Specifications

Layouts, wireframes, and interaction behaviors for the 9 server-rendered Go HTML templates in Flan Media Server.

---

## 1. Template Inventory

All platform navigation is consolidated strictly into the left sidebar. The top header contains zero navigation links.

| Template | Route | Primary Purpose | Key Features |
| :--- | :--- | :--- | :--- |
| **`login.html`** | `GET /login` | Split Start Screen | Left "Welcome" graphic + right User dropdown, PIN field, and `[ Access ]` button. |
| **`video.html`** | `GET /` or `GET /video` | Video Catalog | Unified grid of all video titles under search bar; cards feature solid purple footers. |
| **`video_detail.html`** | `GET /video/{id}` | Video Container Detail | Poster, synopsis, ordered list of playable files with `[ Play ]` and `[ ⬇ VLC / Download ]` buttons. |
| **`books.html`** | `GET /books` | Books Catalog | Grid of all book titles and multi-volume series under search bar. |
| **`book_detail.html`** | `GET /books/{id}` | Book Container Detail | Cover, author, overview, and list of volume files with `[ 📖 Open PDF ]` and `[ ⬇ Download ]` buttons. |
| **`watch.html`** | `GET /watch/{file_id}` | Video Player | Blacked-out player with Plyr, auto-resume prompt, 5s progress sync, and audio codec fallback bar. |
| **`read.html`** | `GET /read/{file_id}` | Document Reader | Native full-viewport browser PDF viewer or direct file download page. |
| **`manage.html`** | `GET /manage` | Management Console | System metrics, `[ ⟳ Rescan All Media ]` trigger, and whimsical avatar profile management. |
| **`manual.html`** | `GET /manual` | Software Handbook & About | 2-column offline manual: sticky TOC, media placement guide, VLC fallback, CLI admin recovery, and specs. |

---

## 2. Template Specifications & Wireframes

### Page 1: Split-Screen Start & Login (`login.html` - `GET /login` - Image 1)
```text
+-------------------------------------------------------------------------+
| Flan Media Server                                                       |
+------------------------------------------+------------------------------+
|                                          | User:                        |
|   ←───────────────────────────           | +--------------------------+ |
|                                          | | mike                   v | |
|        Welcome                           | +--------------------------+ |
|                                          | Pin:                         |
|   ───────────────────────────►           | +--------------------------+ |
|                                          | |                          | |
|                                          | +--------------------------+ |
|                                          |                              |
|                                          | +--------------------------+ |
|                                          | | Access                   | |
|                                          | +--------------------------+ |
+------------------------------------------+------------------------------+
```
* **No Sidebar:** The Start screen uses a full-width split view.
* **Accessible Inputs:** Browser-native `<select>` for user selection and standard password input for PIN.
* **First Run Onboarding:** If zero users exist in the database, the form prompts to create the initial admin account.

---

### Page 2: Video Catalog (`video.html` - `GET /video` - Image 2)
```text
+-------------------------------------------------------------------------+
| Flan Media Server                                           ( ? )   [Avatar]|
+----------+--------------------------------------------------------------+
| ←─────── |                                                              |
| [ Video] |   +---------------------------------------------+  [ 🔍 ]    |
| (Active) |   | Search videos...                            |            |
| [ Books] |   +---------------------------------------------+            |
| ────────►|                                                              |
|          |   +---------------+  +---------------+  +---------------+    |
|          |   | [Poster Area] |  | [Poster Area] |  | [Poster Area] |    |
|          |   |---------------|  |---------------|  |---------------|    |
|          |   | [Purple Band] |  | [Purple Band] |  | [Purple Band] |    |
| [Manage] |   +---------------+  +---------------+  +---------------+    |
+----------+--------------------------------------------------------------+
```
* **Sidebar Navigation:** The `Video` button is active (solid darker purple background `#6d52a8` with white text).
* **Direct Grid:** Media cards appear immediately below the search bar.
* **Card Anatomy:** White rectangle with thick 2px black border, image area, and solid purple footer band holding the title.
* **Click Action:** Clicking any card navigates to `/video/{id}`.

---

### Page 3: Video Detail View (`video_detail.html` - `GET /video/{id}`)
```text
+-------------------------------------------------------------------------+
| Flan Media Server                                           ( ? )   [Avatar]|
+----------+--------------------------------------------------------------+
|          |  ← Back to Videos                                            |
| [ Video] |                                                              |
| (Active) |  +-------------+   BREAKING BAD (2008)                       |
| [ Books] |  | [Poster]    |   A high school chemistry teacher diagnosed |
|          |  |             |   with lung cancer turns to manufacturing...|
|          |  +-------------+                                             |
|          |                                                              |
|          |  PLAYABLE ITEMS / EPISODES                                   |
|          |  +---------------------------------------------------------+ |
|          |  | 1. S01E01 - Pilot (48m)   [== Watched ==] [Play] [⬇ VLC]| |
|          |  +---------------------------------------------------------+ |
|          |  | 2. S01E02 - Cat's (48m)   [== 18m left =] [Play] [⬇ VLC]| |
|          |  +---------------------------------------------------------+ |
| [Manage] |                                                              |
+----------+--------------------------------------------------------------+
```
* **Playable List:** Displays all files for this container, duration, watch progress, `[ ▶ Play ]`, and `[ ⬇ VLC / Download ]`.
* **Direct Streaming:** Clicking Play opens `/watch/{file_id}`.
* **VLC Fallback:** Clicking `[ ⬇ VLC / Download ]` gives the direct stream URL to open in external players (VLC/MPV) or save locally, solving unplayable browser audio codecs (AC3/DTS).

---

### Page 4: Books Catalog (`books.html` - `GET /books`)
```text
+-------------------------------------------------------------------------+
| Flan Media Server                                           ( ? )   [Avatar]|
+----------+--------------------------------------------------------------+
|          |                                                              |
| [ Video] |   +---------------------------------------------+  [ 🔍 ]    |
|          |   | Search books...                             |            |
| [ Books] |   +---------------------------------------------+            |
| (Active) |                                                              |
|          |   +---------------+  +---------------+  +---------------+    |
|          |   | [Cover Area]  |  | [Cover Area]  |  | [Cover Area]  |    |
|          |   |---------------|  |---------------|  |---------------|    |
|          |   | [Purple Band] |  | [Purple Band] |  | [Purple Band] |    |
| [Manage] |   +---------------+  +---------------+  +---------------+    |
+----------+--------------------------------------------------------------+
```
* `Books` button is active in the sidebar.
* Clicking any book card navigates to `/books/{id}`.

---

### Page 5: Book Detail View (`book_detail.html` - `GET /books/{id}`)
* Displays book cover, title, author, and description.
* Lists readable volume files:
  * For **PDFs**: `[ 📖 Open PDF ]` (opens browser native PDF viewer in new tab) and `[ ⬇ Download ]`.
  * For **EPUBs**: `[ ⬇ Download EPUB ]` (instant direct download to open in Apple Books, Moon+ Reader, Kindle, etc.).

---

### Page 6: Video Player View (`watch.html` - `GET /watch/{file_id}`)
```text
+-------------------------------------------------------------------------+
| [← Back]  Breaking Bad - S01E02                                         |
+-------------------------------------------------------------------------+
|                                                                         |
|   [ RESUME PLAYBACK ]                                                   |
|   You were watching at 24:12.                                           |
|   [ ▶ Resume from 24:12 ]        [ ↺ Start from Beginning ]             |
|                                                                         |
|-------------------------------------------------------------------------|
|  [▶] [10s↺] [10s↻]  00:24:12 / 00:48:00          [1.0x] [🔊] [⛶]        |
+-------------------------------------------------------------------------+
| Audio silent or stuttering? [ Open in External VLC Player / Download ↗ ] |
+-------------------------------------------------------------------------+
```
* **Distraction-Free:** Header bar auto-hides after 3 seconds of mouse inactivity.
* **Auto-Resume Prompt:** Prompts to resume if `position_data > 10`.
* **Sync Loop:** `player.js` syncs position to `POST /api/progress` every 5 seconds.
* **Audio Codec Fallback Bar:** Unobtrusive banner providing an instant external VLC stream link for tracks encoded in AC3, EAC3, or DTS.

---

### Page 7: Document Reader View (`read.html` - `GET /read/{file_id}`)
* **PDF Mode:** Serves PDF directly with `Content-Type: application/pdf`, opening in the user's full browser viewport with native zoom and search (zero JS dependencies).
* **EPUB Mode:** Direct file download with `Content-Disposition: attachment` for reading in native reader apps.

---

### Page 8: Management Console (`manage.html` - `GET /manage`)
```text
+-------------------------------------------------------------------------+
| Flan Media Server                                           ( ? )   [Avatar]|
+----------+--------------------------------------------------------------+
|          |  MANAGE SERVER                                               |
| [ Video] |                                                              |
|          |  System Status: Memory: 14.2 MB / 16 MB  •  Streams: 1 / 3   |
| [ Books] |                                                              |
|          |  Media Library:                                              |
|          |  • ./media/video (populated via SMB / SCP / USB)             |
| [Manage] |  • ./media/books                                             |
| (Active) |  [ ⟳ Rescan All Media ]                                      |
|          |                                                              |
|          |  Household Profiles:                                         |
|          |  • [Icon] Wesley (Admin)  [ ✨ Change Avatar ]  [ Edit PIN ]  |
|          |  • [Icon] Mike (User)     [ ✨ Change Avatar ]  [ Reset PIN ] |
|          |  [ + Add New User ]                                          |
+----------+--------------------------------------------------------------+
```
* `Manage Server` button is active in the sidebar.
* **Rescan Button:** Walks `./media/video` and `./media/books` synchronously and updates SQLite.
* **Whimsical Avatar Picker Modal:** Triggered via `[ ✨ Change Avatar ]` or clicking the top-right header avatar:
  * Select from 6 high-contrast SVG presets (Mascot, Flan pudding, Cat, Robot, Ghost, Star).
  * Or upload a custom photo (JPEG, PNG, WebP up to 2MB) saved to `data/avatars/{user_id}.ext`.
* **User Management:** Create/delete household profiles and change PINs.

---

### Page 9: Software Manual & System Handbook (`manual.html` - `GET /manual`)
```text
+-------------------------------------------------------------------------+
| [ ← Back to Library ]       Flan Software Handbook [v0.2.0]  100% Offline|
+---------------------+---------------------------------------------------+
| Table of Contents   | 1. MEDIA STORAGE & PLACEMENT                      |
|                     | Drop video folders in ./media/video and books in  |
| 1. Media Storage    | ./media/books via SMB, USB drive, or SCP.         |
| 2. VLC Fallback     | Click [ ⟳ Rescan All Media ] for sub-200ms index. |
| 3. Profiles & PINs  |                                                   |
| 4. CLI Admin Reset  | 2. DIRECT VIDEO PLAYBACK & VLC FALLBACK           |
| 5. System Specs     | Direct play via HTTP 206 range requests.          |
|                     | If browser audio fails on AC3/DTS, click the      |
|                     | purple [ ⬇ VLC / Download ] button to stream raw. |
|                     |                                                   |
|                     | 4. SERVER CLI ADMINISTRATION & FAILSAFE RESET     |
|                     | ./flan --reset-admin   [ Copy ]                   |
+---------------------+---------------------------------------------------+
```
* **Full-Width Canvas:** Left sidebar hides automatically when viewing the manual to maximize readable technical documentation width.
* **Sticky Table of Contents:** Left TOC navigation panel anchors to sections 1 through 5.
* **1-Click Copy Blocks:** Terminal commands (e.g. `./flan --reset-admin`) feature tactile copy buttons.
* **Mobile Layout:** On $\le 768\text{px}$, TOC shifts to horizontal scrollable pills at the top, and content flows in a single column.

---

## 3. Accessibility & Keyboard Navigation (WCAG 2.1 AA)

All 9 templates are designed for seamless keyboard navigation, high-contrast readability, and screen reader clarity:

1. **Skip-to-Content Link:** The first focusable item on every authenticated page is `<a href="#main-content" class="skip-link">Skip to main content</a>`, jumping past the header and sidebar directly into catalog items or detail panels.
2. **Predictable Tab Sequence:**
   * `Top Header`: Skip Link $\to$ Brand Logo $\to$ Manual Button (`?`) $\to$ User Avatar Badge (`[Avatar]`).
   * `Sidebar Navigation`: `Video` $\to$ `Books` $\to$ `Manage Server`. Active tab is flagged with `aria-current="page"`.
   * `Main Catalog`: Search input $\to$ Search button $\to$ Media Cards in natural DOM reading order.
3. **Card Activation:** Catalog cards are semantic `<a>` tags with `role="article"` and descriptive `aria-label` attributes (e.g. `Blade Runner (Movie)`), activated instantly with `Enter`.
4. **Accessible Modals:** All dialogs use native HTML5 `<dialog class="modal-dialog">` with `aria-labelledby`. Focus is automatically trapped within the modal while open, `Escape` safely closes the dialog, and focus is restored to the triggering button on dismissal.
5. **Live Status Feedback:** Form feedback (login errors, PIN validation, media rescan status) uses `role="alert"` and `aria-live="polite"` / `aria-live="assertive"` for instantaneous screen reader announcement.

---

## 4. Related Documentation

* [Design System & Foundations](design-system.md)
* [Tactile Components](components.md)
* [Accessibility & WCAG Compliance Guide](../accessibility.md)
* [User Flows & Journeys](../diagrams/user-flows.md)

