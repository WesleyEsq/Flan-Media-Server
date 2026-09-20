# Web Client Pages & Template Specifications

This document defines the layout, wireframe hierarchy, and interaction behaviors for all server-rendered Go HTML templates in Flan Media Server.

---

## 1. Page Template Inventory

To prevent scope creep, the web client is strictly bounded to nine server-rendered templates:

1. `index.html`: Dashboard with in-progress continue shelves and recently added media.
2. `videos.html`: Video catalog with Movies and TV Series tabs and genre filtering.
3. `show.html`: TV Series detail view with season tabs and episode lists.
4. `books.html`: Document and book catalog with author and genre filters.
5. `watch.html`: Distraction-free video player embedding Plyr with auto-resume.
6. `read.html`: Dual-mode document viewer (native PDF embed and ePub.js reader).
7. `login.html`: Profile selector ("Who is watching?") and numeric PIN keypad.
8. `setup.html`: First-time onboarding wizard for initial server bootstrapping.
9. `settings.html`: Server status, storage health, library scans, and user management.

---

## 2. Page Specifications & Wireframes

### Page 1: Dashboard (`index.html` - Route: `GET /`)

The home screen presents currently active media before general catalog items, allowing users to immediately resume watching or reading.

#### Wireframe Layout

``` text
+-------------------------------------------------------------------------+
| [  Flan]   Home   Videos   Books                 [? Help] [⚙] [Avatar] |
+-------------------------------------------------------------------------+
|                                                                         |
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
|  +--------+  +--------+  +--------+  +--------+  +--------+             |
|  | Card 1 |  | Card 2 |  | Card 3 |  | Card 4 |  | Card 5 |             |
|  +--------+  +--------+  +--------+  +--------+  +--------+             |
|                                                                         |
+-------------------------------------------------------------------------+
```

#### Behavior & Interactions

+ **Conditional Shelves:** The Continue Watching and Continue Reading shelves only render if the active user has items with `position_seconds > 10` and `is_finished = 0`. If none exist, these shelves are hidden.
+ **Horizontal Scroll with Snap:** Shelves scroll horizontally via native CSS (`scroll-snap-type: x mandatory`). On desktop, chevron buttons allow advancing by one viewport width.
+ **Empty Library State:** If the database contains zero media items, the dashboard displays an approachable empty state:
  *"No media found. Go to Settings to scan your media folder or upload files."* with a button linking to `/settings`.

---

### Page 2: Videos Catalog (`videos.html` - Route: `GET /videos`)

The primary video browsing view supporting tabbed filtering between standalone movies and TV series.

#### Wireframe

``` text
+-------------------------------------------------------------------------+
| [  Flan]   Home   Videos   Books                 [? Help] [⚙] [Avatar] |
+-------------------------------------------------------------------------+
|                                                                         |
|  VIDEOS                                                                 |
|  [ ALL ]  [ MOVIES ]  [ TV SERIES ]                                     |
|                                                                         |
|  Genres: ( All ) ( Animation ) ( Comedy ) ( Sci-Fi ) ( Drama ) ( Action )|
|                                                                         |
|  +----------------+  +----------------+  +----------------+  +--------+ |
|  | [2:3 Poster]   |  | [2:3 Poster]   |  | [2:3 Poster]   |  | ...    | |
|  | Movie Title    |  | Series Title   |  | Movie Title    |  |        | |
|  | 2024 • 1h 45m  |  | 2022 • 3 Seasons| 1999 • 2h 16m  |  |        | |
|  | ★ 8.2          |  | ★ 8.9          |  | ★ 8.7          |  |        | |
|  +----------------+  +----------------+  +----------------+  +--------+ |
|                                                                         |
+-------------------------------------------------------------------------+
```

#### Behavior & Interactions

+ **Sub-Tabs:** Clicking Movies filters the query to standalone films (`series_id IS NULL`). Clicking TV Series groups results by series title, displaying one card per show.
+ **Genre Pills:** Clicking a genre pill appends `?genre=Animation` to the URL. The active genre pill displays with a solid lavender background (`--accent-lavender`).
+ **Responsive CSS Grid:** Cards wrap automatically across rows (`grid-template-columns: repeat(auto-fill, minmax(160px, 1fr))`), accommodating screens from 320px smartphones to 4K TVs.
+ **Click Targets:** Clicking a movie card navigates to `/watch/{id}`. Clicking a TV series card navigates to `/show/{title}`.

---

### Page 3: TV Series Detail View (`show.html` - Route: `GET /show/{title}`)

A dedicated view for television shows that organizes multiple seasons and episodes without cluttering the main catalog.

#### Wireframe Layout

```text
+-------------------------------------------------------------------------+
| [  Flan]   Home   Videos   Books                 [? Help] [⚙] [Avatar] |
+-------------------------------------------------------------------------+
|  ← Back to Videos                                                       |
|                                                                         |
|  +--------------+   BREAKING BAD (2008)                                 |
|  | [Poster]     |   ★ 9.5 • 5 Seasons • Crime, Drama                    |
|  |              |   A high school chemistry teacher diagnosed with      |
|  |              |   inoperable lung cancer turns to manufacturing...    |
|  +--------------+                                                       |
|                                                                         |
|  [ Season 1 ]  [ Season 2 ]  [ Season 3 ]  [ Season 4 ]  [ Season 5 ]   |
|                                                                         |
|  EPISODES                                                               |
|  +--------------------------------------------------------------------+ |
|  | [Thumb]  1. Pilot                                            48m   | |
|  |          Diagnosed with terminal cancer, Walter White teams up...   | |
|  |          [================== Watched ====================]  [ ✓ ]  | |
|  +--------------------------------------------------------------------+ |
|  | [Thumb]  2. Cat's in the Bag...                              48m   | |
|  |          Walt and Jesse attempt to dispose of two bodies...         | |
|  |          [====== 18m left =======]                           [ ► ]  | |
|  +--------------------------------------------------------------------+ |
|  | [Thumb]  3. ...And the Bag's in the River                    48m   | |
|  |          Walt grapples with a life-or-death decision...             | |
|  |                                                              [ ► ]  | |
|  +--------------------------------------------------------------------+ |
+-------------------------------------------------------------------------+
```

#### Behavior & Interactions

+ **Season Switcher:** Clicking a season tab filters the episode list in place. The active season tab is highlighted in lavender.
+ **Progress Badges:** Each episode displays an inline watch progress bar. Completed episodes show a subtle green checkmark badge (`[ ✓ ]`).
+ **Direct Play:** Clicking anywhere on an episode row immediately launches `/watch/{media_id}`.

---

### Page 4: Books Catalog (`books.html` - Route: `GET /books`)

The catalog view for PDF documents and EPUB books.

#### Wireframe Layout

```text
+-------------------------------------------------------------------------+
| [  Flan]   Home   Videos   Books                 [? Help] [⚙] [Avatar]  |
+-------------------------------------------------------------------------+
|                                                                         |
|  BOOKS & DOCUMENTS                                                      |
|  Genres: ( All ) ( Sci-Fi ) ( Fantasy ) ( Technology ) ( Documentation )|
|                                                                         |
|  +----------------+  +----------------+  +----------------+             |
|  | [1:1.4 Cover]  |  | [1:1.4 Cover]  |  | [1:1.4 Cover]  |             |
|  | Dune           |  | Go in Action   |  | Linux Kernel   |             |
|  | Frank Herbert  |  | William Kennedy|  | Robert Love    |             |
|  | [ EPUB ]       |  | [ PDF ]        |  | [ PDF ]        |             |
|  +----------------+  +----------------+  +----------------+             |
|                                                                         |
+-------------------------------------------------------------------------+
```

#### Behavior & Interactions

+ **Card Ratio:** Book cards use a 1:1.4 aspect ratio with a subtle vertical faux spine shadow along the left border.
+ **Format Badges:** Each card displays a small format pill (`[ EPUB ]` or `[ PDF ]`).
+ **Click Actions:** Clicking a card navigates to `/read/{id}`.

---

### Page 5: Video Player View (`watch.html` - Route: `GET /watch/{id}`)

A distraction-free, blacked-out viewing screen where media playback takes center stage.

#### Wireframe Layout

```
+-------------------------------------------------------------------------+
| [← Back]  Dune: Part Two (2024)                                         |
+-------------------------------------------------------------------------+
|                                                                         |
|  +-------------------------------------------------------------------+  |
|  |                                                                   |  |
|  |   [ RESUME PLAYBACK ]                                             |  |
|  |   You were watching at 24:12.                                     |  |
|  |   [ ▶ Resume from 24:12 ]      [ ↺ Start from Beginning ]       |  |
|  |                                                                   |  |
|  |-------------------------------------------------------------------|  |
|  |  [▶] [10s↺] [10s↻]  00:24:12 / 02:46:00    [CC] [1.0x] [🔊] [⛶] |  |
|  +-------------------------------------------------------------------+  |
|                                                                         |
+-------------------------------------------------------------------------+
```

#### Behavior & Interactions

+ **Header Overlay:** A minimalist top bar containing a `[← Back]` button and the media title. Auto-hides after 3 seconds of mouse inactivity during playback.
+ **Auto-Resume Prompt:** If saved progress exists (`position_seconds > 10`), a flat modal overlay appears before playback starts asking whether to resume or start over.
+ **Player Controls (Plyr):** Custom lavender accent (`--plyr-color-main: #bb9af7`), skip 10s buttons, playback speed options (0.5x, 0.75x, 1x, 1.25x, 1.5x, 2x), and WebVTT subtitle track switcher.
+ **Progress Sync:** `player.js` emits a throttled `POST /api/progress` every 5 seconds. On video completion (`ended` event), marks `is_finished = 1`.

---

### Page 6: Document Reader View (`read.html` - Route: `GET /read/{id}`)

Provides dedicated reading interfaces tailored to the file format (PDF vs. EPUB).

#### A. PDF Mode Layout (Native Embed)

```
+-------------------------------------------------------------------------+
| [← Back]  Go Programming Language.pdf                 [ ⬇ Download PDF ]|
+-------------------------------------------------------------------------+
|                                                                         |
|  +-------------------------------------------------------------------+  |
|  |                                                                   |  |
|  |               [ Native Browser PDF Viewing Engine ]               |  |
|  |            (Embedded via <iframe src="/stream/{id}">)             |  |
|  |                                                                   |  |
|  +-------------------------------------------------------------------+  |
|                                                                         |
+-------------------------------------------------------------------------+
```

#### B. EPUB Mode Layout (ePub.js Reader)

```
+-------------------------------------------------------------------------+
| [← Back]  Dune - Frank Herbert       [ A- ] [ A+ ]   [ ⬇ Download EPUB ]|
+-------------------------------------------------------------------------+
|                                                                         |
|  [ < Prev ]                                                  [ Next > ] |
|                                                                         |
|             CHAPTER 1                                                   |
|                                                                         |
|             A beginning is the time for taking the                      |
|             most delicate care that the balances are correct.           |
|             This every sister of the Bene Gesserit knows...             |
|                                                                         |
|-------------------------------------------------------------------------|
|  Chapter 1: Arrakis                                        Page 12 / 620|
+-------------------------------------------------------------------------+
```

#### Behavior & Interactions

+ **PDF:** Uses native browser PDF rendering via an iframe, eliminating heavy JavaScript reader bundles.
+ **EPUB:** ePub.js renders chapters with font scaling controls (`A-` and `A+`) and left/right click or keyboard arrow page turning.
+ **Download Action:** Both modes provide a prominent `[ ⬇ Download ]` button for users who prefer opening files in native e-readers or PDF apps on tablets.

---

### Page 7: Profile Selector & Login (`login.html` - Route: `GET /login`)

The landing screen for returning users. Focuses on rapid profile switching via friendly avatar tiles.

#### Wireframe Layout

```
+-------------------------------------------------------------------------+
|                                                                         |
|                                  Flan                                  |
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
|                     |                           |                       |
|                     |    [ 1 ]   [ 2 ]   [ 3 ]  |                       |
|                     |    [ 4 ]   [ 5 ]   [ 6 ]  |                       |
|                     |    [ 7 ]   [ 8 ]   [ 9 ]  |                       |
|                     |    [ ⌫ ]   [ 0 ]   [ ↵ ]  |                       |
|                     +---------------------------+                       |
|                                                                         |
+-------------------------------------------------------------------------+
```

#### Behavior & Interactions

+ **Profile Selection:** Clicking a profile tile activates the numeric PIN keypad modal for that user.
+ **Keypad Input:** Supports both physical keyboard number keys and on-screen touch keypad clicks.
+ **Lockout Feedback:** If an account is rate-limited, the keypad disables and displays:
  *"Account locked due to failed attempts. Try again in 4:32."*

---

### Page 8: First-Time Setup Wizard (`setup.html` - Route: `GET /setup`)

Displayed only when the database contains zero users. Automatically locked once completed.

#### Wireframe Layout

```
+-------------------------------------------------------------------------+
|                                                                         |
|                                Welcome to Flan                         |
|                         Initial Server Setup                            |
|                                                                         |
|         +-----------------------------------------------------+         |
|         | 1. Administrator Profile                            |         |
|         |    Username: [ Wesley                            ]  |         |
|         |                                                     |         |
|         | 2. Security PIN                                     |         |
|         |    Create 4 to 6-Digit PIN: [ ****               ]  |         |
|         |                                                     |         |
|         | 3. Primary Media Folder                             |         |
|         |    Path on host: [ /mnt/storage/media            ]  |         |
|         |                                                     |         |
|         |    [ Complete Setup & Launch Catalog ]              |         |
|         +-----------------------------------------------------+         |
|                                                                         |
+-------------------------------------------------------------------------+
```

#### Behavior & Interactions

+ **Path Validation:** Submitting validates that the host directory exists and has read permissions.
+ **Automatic Lockout:** Upon creating the admin account, `/setup` is permanently disabled.

---

### Page 9: Server & Library Settings (`settings.html` - Route: `GET /settings`)

The administrative control center for storage management, scanning, and user profiles.

#### Wireframe Layout

```
+-------------------------------------------------------------------------+
| [  Flan]   Home   Videos   Books                 [? Help] [⚙] [Avatar] |
+-------------------------------------------------------------------------+
|                                                                         |
|  SETTINGS                                                               |
|                                                                         |
|  System Status & Storage                                                |
|  +--------------------------------------------------------------------+ |
|  | Host OS: Linux (ARM64)               Memory: 14.2 MB / 16 MB limit | |
|  | Active Streams: 1 / 3 max            Database: data/flan.db (WAL)  | |
|  | Root SD Free: 42.1 GB               Media USB Free: 1.4 TB         | |
|  +--------------------------------------------------------------------+ |
|                                                                         |
|  Library Folders                                                        |
|  +--------------------------------------------------------------------+ |
|  | /mnt/hdd1/movies (Videos)  •  342 files                            | |
|  | /mnt/hdd2/tv (TV Shows)    •  18 shows, 214 episodes               | |
|  | /mnt/nvme/books (Books)    •  48 books                             | |
|  |                                                                    | |
|  | [ ⟳ Scan All Libraries ]      [ ⬆ Upload Files/Folder ]            | |
|  +--------------------------------------------------------------------+ |
|                                                                         |
|  User Profiles                                                          |
|  +--------------------------------------------------------------------+ |
|  | [Avatar] Wesley (Admin)                 [ Edit PIN ] [ Change Icon]| |
|  | [Avatar] Guest (User)                   [ Reset PIN ] [ Delete ]   | |
|  |                                                                    | |
|  | [ + Add New Profile ]                                              | |
|  +--------------------------------------------------------------------+ |
|                                                                         |
+-------------------------------------------------------------------------+
```

#### Behavior & Interactions

+ **Scan Trigger:** Clicking `[ ⟳ Scan All Libraries ]` initiates a background scan and changes the button state to a disabled spinner with a 30-second cooldown timer.
+ **Upload Trigger:** Clicking `[ ⬆ Upload Files/Folder ]` opens the streaming upload modal with target drive selection.
+ **Profile Management:** Admin can add new household profiles, change avatars, and reset user PINs without terminal access.

---

### Related Documentation

+ [Web Client Design System & Foundations](design-system.md)
+ [Web Client Component Specifications](components.md)
+ [User Flows & Journey Diagrams](../diagrams/user-flows.md)
+ [Master System Specifications](../design.md)
