# Web Client Pages & Template Specifications

Layouts, wireframes, and interaction behaviors for the 8 server-rendered Go HTML templates in Flan Media Server.

---

## 1. Template Inventory

All platform navigation is consolidated strictly into the left sidebar. The top header contains zero navigation links.

| Template | Route | Primary Purpose | Key Features |
| :--- | :--- | :--- | :--- |
| **`login.html`** | `GET /login` | Split Start Screen | Left "Welcome" graphic + right User dropdown, PIN field, and `[ Access ]` button. |
| **`video.html`** | `GET /` or `GET /video` | Video Catalog | Unified grid of all video titles under search bar; cards feature solid purple footers. |
| **`video_detail.html`** | `GET /video/{id}` | Video Container Detail | Poster, synopsis, and ordered list of playable files (episodes or cuts) with Play buttons. |
| **`books.html`** | `GET /books` | Books Catalog | Grid of all book titles and multi-volume series under search bar. |
| **`book_detail.html`** | `GET /books/{id}` | Book Container Detail | Cover, author, overview, and list of readable files (volumes, editions) with Read/Download buttons. |
| **`watch.html`** | `GET /watch/{file_id}` | Video Player | Blacked-out player with Plyr, auto-resume prompt, and 5s progress sync. |
| **`read.html`** | `GET /read/{file_id}` | Document Reader | Dual-mode viewer: native PDF iframe and paginated ePub.js reader with download button. |
| **`manage.html`** | `GET /manage` | Management Console | System status, rescan button for `./media`, uploads, and user profile management. |

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
| Flan Media Server                                   ( ? )   ( ⚙ )   [Avatar]|
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
| Flan Media Server                                   ( ? )   ( ⚙ )   [Avatar]|
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
|          |  | 1. S01E01 - Pilot (48m)             [== Watched ==] [Play| |
|          |  +---------------------------------------------------------+ |
|          |  | 2. S01E02 - Cat's in the Bag (48m)  [== 18m left =] [Play| |
|          |  +---------------------------------------------------------+ |
| [Manage] |                                                              |
+----------+--------------------------------------------------------------+
```
* **Playable List:** Displays all files for this container, duration, watch progress, and Play button.
* **Direct Streaming:** Clicking Play opens `/watch/{file_id}`.

---

### Page 4: Books Catalog (`books.html` - `GET /books`)
```text
+-------------------------------------------------------------------------+
| Flan Media Server                                   ( ? )   ( ⚙ )   [Avatar]|
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
* Lists readable volume files with `[ Read ]` and `[ ⬇ Download ]` buttons.

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
```
* **Distraction-Free:** Header bar auto-hides after 3 seconds of mouse inactivity.
* **Auto-Resume Prompt:** Prompts to resume if `position_data > 10`.
* **Sync Loop:** `player.js` syncs position to `POST /api/progress` every 5 seconds.

---

### Page 7: Document Reader View (`read.html` - `GET /read/{file_id}`)
* **PDF Mode:** Embedded native browser viewer via `<iframe src="/stream/book/{file_id}">` (zero JS bundle bloat).
* **EPUB Mode:** Paginated in-browser ePub.js reader with font scaling controls (`[ A- ] [ A+ ]`).
* **Download Button:** Always available to open in native tablet reading apps.

---

### Page 8: Management Console (`manage.html` - `GET /manage`)
```text
+-------------------------------------------------------------------------+
| Flan Media Server                                   ( ? )   ( ⚙ )   [Avatar]|
+----------+--------------------------------------------------------------+
|          |  MANAGE SERVER                                               |
| [ Video] |                                                              |
|          |  System Status: Memory: 14.2 MB / 16 MB  •  Streams: 1 / 3   |
| [ Books] |                                                              |
|          |  Media Library:                                              |
|          |  • ./media/video                                             |
| [Manage] |  • ./media/books                                             |
| (Active) |  [ ⟳ Rescan All Media ]   [ ⬆ Upload Files ]                 |
|          |                                                              |
|          |  User Profiles:                                              |
|          |  • Wesley (Admin)  [ Edit PIN ]  [ Change Icon ]             |
|          |  • Mike (User)     [ Reset PIN ] [ Delete ]                  |
|          |  [ + Add New User ]                                          |
+----------+--------------------------------------------------------------+
```
* `Manage Server` button is active in the sidebar.
* **Rescan Button:** Walks `./media/video` and `./media/books` synchronously and updates SQLite.
* **Upload Modal:** Direct browser uploads into `./media/video/<Container>/` or `./media/books/<Container>/`.
* **User Management:** Create/delete user profiles and change PINs.

---

## 3. Related Documentation

* [Design System & Foundations](design-system.md)
* [Tactile Components](components.md)
* [User Flows & Journeys](../diagrams/user-flows.md)
