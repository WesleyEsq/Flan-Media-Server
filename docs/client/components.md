# Web Client Component Specifications

Functional specifications, markup structures, and interaction behaviors for core UI components in Flan Media Server.

---

## 1. Media Card Component (`.card`)

The primary catalog unit representing a video or book container.

### Visual Rules
* **Borders:** `2px solid #000000` border, `8px` rounded corners.
* **Poster Area:** White surface (`#ffffff`) holding the media cover image.
* **Footer Band:** Solid purple strip (`--brand-active: #6d4ca6`) across the bottom of the card with white title text (`#ffffff`).
* **Interactions:** Depresses immediately into active state on click (`transform: translateY(2px)`).
* **Aspect Ratios:**
  * Video Poster (`.card--video`): 2:3 aspect ratio.
  * Book Cover (`.card--book`): 1:1.4 aspect ratio.

### HTML Structure
```html
<article class="card card--video">
    <a href="/video/42" class="card__link">
        <div class="card__poster-wrap">
            <img src="/covers/video/42" alt="Title" class="card__poster" loading="lazy">
        </div>
        <div class="card__footer-band">
            <h3 class="card__title">Blade Runner</h3>
        </div>
    </a>
</article>
```

---

## 2. Sidebar Navigation Buttons (`.sidebar-btn`)

Used in the left navigation sidebar for `Video`, `Books`, and `Manage Server`.

### HTML Structure
```html
<nav class="sidebar-nav" aria-label="Main Navigation">
    <a href="/video" class="sidebar-btn sidebar-btn--active" aria-current="page">Video</a>
    <a href="/books" class="sidebar-btn">Books</a>
    <a href="/manage" class="sidebar-btn sidebar-btn--bottom">Manage Server</a>
</nav>
```

### Styling Rules
* **Border:** `2px solid #000000`, `8px` rounded corners.
* **Active State (`.sidebar-btn--active`):** Solid darker purple (`--brand-active: #6d4ca6`) background with white bold text.
* **Inactive State:** Light gray surface (`#eae8f2`) with dark bold text (`#12131a`).
* **Active Click:** Depresses on click (`:active { transform: translateY(2px); }`).
* **Focus Outline:** Visible 3px black outline (`:focus-visible { outline: 3px solid #000; outline-offset: 2px; }`).

---

## 3. Search Bar Component (`.search-bar`)

Positioned at the top of the main content area.

### HTML Structure
```html
<form class="search-bar" action="/video" method="GET">
    <input type="text" name="q" class="search-bar__input" placeholder="Search..." autocomplete="off">
    <button type="submit" class="search-bar__btn" aria-label="Search">
        <svg class="icon" viewBox="0 0 24 24"><path d="M...magnifying-glass"/></svg>
    </button>
</form>
```

### Styling Rules
* Input field has a white background with a `2px solid #000000` border and `8px` rounded corners.
* Search button is a high-contrast square with `2px solid #000000` border containing an inline SVG icon.

---

## 4. Top Header & Utility Controls (`.top-header`)

The persistent top header bar (`.top-header`, 56px height) sits at the top of the interface across both desktop and mobile viewports. It carries the brand title on the left and utility controls (`?` manual button and user avatar badge) on the right. Scoped search remains within the catalog canvas.

### HTML Structure
```html
<header id="top-header" class="top-header">
    <div id="brand-title" class="brand-title brand-trigger" title="Go to Catalog" tabindex="0" role="link" aria-label="Flan Media Server - Go to Catalog">
        <span>Flan Media Server</span>
    </div>

    <div class="header-tools">
        <!-- Software Manual Trigger -->
        <button id="help-btn" class="tool-circle-btn help-trigger-btn" aria-label="Software Manual & About" title="Software Manual & About">?</button>

        <!-- Profile Avatar Badge -->
        <button id="user-avatar" type="button" class="avatar-badge user-avatar-trigger" aria-label="My Profile & Settings" aria-haspopup="dialog" title="My Profile & Settings">
            <!-- Injected user SVG or custom photo -->
        </button>
    </div>
</header>
```

### Styling Rules
* **Header Bar:** 56px height, solid purple background (`var(--header-purple): #6d4ca6`), `2.5px solid #000000` bottom border, sticky at top.
* **Brand Title:** Bold white text (`font-size: 1.35rem; font-weight: 700`), interactive cursor, navigates to `#video` on click/Enter.
* **Manual Button (`#help-btn`):** Circular white button, 38px diameter, `2.5px solid #000000` border, centered `?` symbol, navigates to `#manual`.
* **Avatar Badge (`#user-avatar`):** 40px rounded rectangle button with `3px solid #3ea6ff` border and black ring shadow (`box-shadow: 0 0 0 2px #000`), opens Unified My Profile modal on click.

---

## 5. Continue Watching / Jump Back In Shelf (`.continue-shelf`)

Surfaces active in-progress media sessions at the top of the catalog:

### Visual & Functional Rules
* **Video Resume Cards:** Shows item title, current episode / cut, formatted progress (`24:12 / 48m`), a tactile purple progress bar, and a 1-tap `[ Resume ]` button that deep-links directly into `/watch/{file_id}`.
* **Book Reading Cards:** Shows title, author, and discrete status badges (e.g. `READING - Page 142 of 680`). No arbitrary percentage numbers are calculated for books.
* **Activity Sorting:** Sorted by last interaction timestamp descending (`last_watched_at` / `last_read_at DESC`).
* **Zero-Waste Collapse:** If 0 items are in progress or if the user filters by `Unwatched`, the shelf collapses completely from the DOM.

---

## 6. Collapsible Filter & Sort Drawer (`.filter-drawer`)

Hides granular catalog filters behind a tactile toggle button to preserve a clean browsing surface:

### Visual & Functional Rules
* **Toggle Button (`.filter-toggle-btn`):** Positioned beside the search button (`[ Filter [v] ]` / `[ Filter [^] ]`). Turns active purple when the drawer is open.
* **Pill Groups:**
  * **Format Pills:** `All`, `Movies`, `Series` (for video); `All`, `EPUB`, `PDF` (for books).
  * **Status Pills:** `All`, `In Progress`, `Unwatched` / `Unread`.
  * **Source Drive Pills:** Dynamic pills for each configured storage folder.
  * **Sort Order:** `Shuffle / Random` (default order), `Recently Added`, `Title (A-Z)`, `Release Year`.
* **Tactile Styling:** High-contrast pills with `2px solid #000000` borders and purple active state.

---

## 7. Pagination Component (`.pagination-bar`)

Enforces strict, predictable catalog chunking:

### Functional Rules
* **Capacity Limit:** Capped at **21 items maximum per page** ($7 \text{ columns} \times 3 \text{ rows} = 21$).
* **Touch Targets:** Minimum 44px hit targets on all navigation buttons for WCAG 2.5.5 compliance.
* **Controls:** Includes `[ <- Prev ]`, numbered page buttons `[ 1 ] [ 2 ] ...`, `[ Next -> ]`, and clear range indicators (`Showing 1-21 of 23`).
* **Instant Keyboard Access:** Fully navigable via Tab and Enter keys.

---

## 5. 2-Step Sequential Login (`.login-pad-viewport`)

The centered tactile authentication pad with 2-step state transitions:

### Step 1: Profile Selection (`.login-picker-card`)
```html
<div class="login-pad-viewport">
    <div class="login-pad-card login-picker-card">
        <div class="login-pad-header">
            <h1 class="login-pad-title">Flan Media Server</h1>
            <p class="login-pad-subtitle">Who is watching?</p>
        </div>

        <!-- Household Profile Cards -->
        <div class="login-profiles-grid" role="group" aria-label="Household Profiles">
            <button type="button" class="login-profile-card" data-username="wesley">
                <div class="login-profile-avatar-frame"><!-- 72px User Avatar SVG --></div>
                <span class="login-profile-username">wesley</span>
                <span class="login-profile-role-pill">Admin</span>
            </button>
            <button type="button" class="login-profile-card" data-username="mike">
                <div class="login-profile-avatar-frame"><!-- 72px User Avatar SVG --></div>
                <span class="login-profile-username">mike</span>
                <span class="login-profile-role-pill">User</span>
            </button>
        </div>
    </div>
</div>
```

### Step 2: Dedicated PIN Entry (`.login-pin-card`)
```html
<div class="login-pad-viewport">
    <div class="login-pad-card login-pin-card">
        <div class="login-user-banner">
            <div class="login-user-avatar-badge"><!-- 60px User Avatar SVG --></div>
            <div class="login-user-info">
                <h2 class="login-user-name">wesley</h2>
                <span class="login-profile-role-pill">Admin</span>
            </div>
        </div>

        <p class="login-pad-subtitle">Enter your 4-digit PIN</p>

        <!-- Tactile PIN Access Form -->
        <form id="login-form" class="login-pin-form" action="/api/login" method="POST">
            <div class="form-field-group">
                <input class="tactile-pin-input" id="pin-input" type="password" name="pin" maxlength="6" inputmode="numeric" placeholder="••••" autofocus required>
            </div>

            <button type="submit" class="tactile-access-btn" id="access-btn">Access Library →</button>
            <button type="button" class="tactile-switch-btn" id="switch-profile-btn">← Switch Profile</button>
            <div id="login-feedback" class="login-feedback" role="alert" aria-live="assertive"></div>
        </form>
    </div>
</div>
```

---

## 6. Playable File Row (`.file-row`)

Lists playable video files, cuts, or book volumes in detail views.

### HTML Structure
```html
<!-- Video File Row with Direct Play & VLC Fallback -->
<div class="file-row">
    <div class="file-row__info">
        <span class="file-row__title">S01E01 - Pilot</span>
        <span class="file-row__duration">48m</span>
    </div>
    <div class="file-row__progress">
        <div class="file-row__progress-fill" style="width: 45%;"></div>
    </div>
    <div class="file-row__actions">
        <a href="/watch/101" class="btn btn--primary">Play</a>
        <!-- Signed URL allows external players (VLC, MPV) to stream without session cookies -->
        <a href="/download/video/101?exp=1700000000&u=1&sig=a1b2c3d4..." class="btn btn--vlc" title="Direct download or stream in VLC">VLC / Download</a>
    </div>
</div>

<!-- Book File Row -->
<div class="file-row">
    <div class="file-row__info">
        <span class="file-row__title">Volume 1 (PDF)</span>
        <span class="file-row__duration">18.4 MB</span>
        <span class="badge badge--status">Reading</span>
    </div>
    <div class="file-row__actions">
        <a href="/stream/book/201" target="_blank" class="btn btn--primary">Open PDF</a>
        <a href="/download/book/201" class="btn btn--secondary">Download</a>
    </div>
</div>
```

---

## 7. Avatar Picker Component (`.avatar-picker`)

Used in `/manage` and profile settings:

* **Preset Grid:** 3-column grid displaying 6 SVG options (Mascot, Flan, Cat, Ghost, Robot, Star).
* **Selection State:** `3px solid --brand-active: #6d4ca6` border.
* **Custom Photo Upload:** Standard image file input (`accept="image/png, image/jpeg, image/webp"`, maximum 2 MB).

---

## 8. Modal Dialog Component (`<dialog class="modal-dialog">`)

Uses native HTML5 `<dialog>` for focus trapping, backdrop support, and native `Escape` key dismissal:

```html
<dialog id="profile-modal" class="modal-dialog" aria-labelledby="profile-title">
    <header class="modal-header">
        <h2 id="profile-title" class="modal-title">My Profile</h2>
        <button class="modal-close-btn" type="button" aria-label="Close dialog">X</button>
    </header>
    <div class="modal-body"><!-- Content --></div>
</dialog>
```

* **Focus Restoration:** When opened with `.showModal()`, focus moves to the first actionable input. When closed, focus returns to the triggering button.
* **Sticky Header:** The header is sticky so the title and close button remain visible while scrolling modal content.

---

## 9. Catalog Status Strip

Provides a compact status banner above media card grids:

* **Browsing Mode:** Displays current section and count:
  `ALL VIDEOS & SERIES - 5 TITLES`
* **Filter Mode:** Displays active query and a clear button:
  `[FILTER] "gundam" (1 match)   [ Clear ]`

---

## 10. Horizontal Title Marquee

For container titles that exceed card footer width:

* **Default:** Truncates with ellipsis (`text-overflow: ellipsis; white-space: nowrap;`).
* **Hover / Focus:** On mouseover or `:focus-visible`, if `scrollWidth > clientWidth`, CSS translation scrolls the title so the entire text is readable.
* **Reduced Motion:** If `prefers-reduced-motion: reduce` is active, the animation is disabled and static ellipsis truncation is maintained.

---

## 11. Video Detail Editor Modal (`#media-edit-modal`)

Used on `/video/{id}` by administrators to modify display metadata without modifying disk files:

```html
<dialog id="media-edit-modal" class="modal-dialog" aria-labelledby="edit-media-title">
    <header class="modal-header">
        <h2 id="edit-media-title" class="modal-title">Edit Video Metadata</h2>
        <button class="modal-close-btn" type="button" aria-label="Close dialog">X</button>
    </header>
    <form class="modal-form" id="media-edit-form">
        <label for="edit-title">Title:</label>
        <input type="text" id="edit-title" name="title" class="tactile-text-input" required>

        <label for="edit-year">Release Year:</label>
        <input type="number" id="edit-year" name="release_year" class="tactile-text-input" min="1900" max="2100">

        <label for="edit-overview">Synopsis / Overview:</label>
        <textarea id="edit-overview" name="overview" class="tactile-textarea" rows="3"></textarea>

        <label for="edit-cover-file">Replace Poster Art (Max 2MB):</label>
        <input type="file" id="edit-cover-file" accept="image/png, image/jpeg, image/webp" class="tactile-file-input">

        <h3>Playable Files / Episodes</h3>
        <div class="edit-files-list">
            <!-- Rows with title text input, up/down order buttons, and hide toggle -->
        </div>

        <button type="submit" class="tactile-action-btn">Save Changes</button>
    </form>
</dialog>
```

---

## 12. Ingestion Pipeline Checklist Component (`.ingest-pipeline`)

Used in the on-demand Ingestion Pipeline modal (`#ingestion-pipeline-modal`) to inspect, select, and edit gathered media before committing to the database:

* **Selection Toolbar:**
  * Master `Select All` checkbox with counter (`N selected`).
  * Instant search filter input to quickly find specific gathered items.
* **Candidate Item Card:**
  * **Checkbox:** Toggle whether to import or skip this container.
  * **Raw Folder Indicator:** Displays the raw directory name on disk (e.g. `[Dir] the.wire.s01.720p-ctu`).
  * **Inline Metadata Inputs:** Clean text input for `Title` and number input for `Release Year`, allowing immediate typo fixes without extra modal layers.
  * **File List Accordion:** Expandable button (`View Files ▾`) revealing all detected episodes/cuts with individual file sizes and formats.
* **Action Footer:**
  * High-contrast `[ Ingest Selected Items (N) ]` tactile button. Only checked items are imported into SQLite.

---

## 13. Storage Source Manager (`.storage-source-row`)

Lists active media directories in `/manage`:

* Displays friendly name, absolute path, media type tag (`Video` / `Books`), and free disk space.
* **Scan Button:** Directly launches the Ingestion Pipeline wizard (`#ingestion-pipeline-modal`) for this specific storage folder.
* **Edit Button (`#edit-source-modal`):** Opens a modal allowing administrators to rename the folder's friendly display name (stored in SQLite `storage_sources.name`). Preserves underlying disk paths, preventing breaking symlinks, mounts, or torrent seeds.
* **Remove Button with Warning Modal:** Opens `#remove-source-modal` to confirm unregistering the folder from Flan without deleting any disk files.
* **Upload Files (`#upload-media-modal`):** Direct in-browser upload streaming straight to destination storage.
* **Add Source Modal (`#add-source-modal`):** Provides path input, one-click `.flan-keep` creation, and an immediate prompt to run the ingestion pipeline.

---

## 14. Related Documentation

* [Design System](design-system.md)
* [Page Templates & Wireframes](pages.md)
* [Accessibility Guide](accessibility.md)
* [Mobile Responsiveness](responsiveness.md)
* [Master System Architecture](../design.md)
* [Database Schema & Queries](../database.md)
* [Storage Architecture](../storage.md)
* [Local-First Metadata Engine](../scraper.md)
