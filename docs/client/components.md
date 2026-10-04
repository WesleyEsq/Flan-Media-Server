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
  * Video Poster (`.card--video`): 2:3 aspect ratio (rendered as `aspect-ratio: 3 / 4.4` in catalog grids).
  * Book Cover (`.card--book`): 1:1.4 aspect ratio.
* **Fluid Tiered Scaling (Low Zoom / Large Displays):** Rather than remaining static 200px tiles on high-res displays, card dimensions scale up dynamically across progressive breakpoints as `--catalog-max-width` expands from 1680px to 2400px:
  * Cards scale from ~200px wide up to ~325px wide and ~475px tall.
  * Poster play SVGs scale from 48px up to 76px.
  * Title typography scales from 0.92rem up to 1.2rem.
  * Preserves bold, prominent readability similar to modern YouTube layouts until capping strictly at 2400px.

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

## 2. Discord x YouTube Hybrid Rail (`.sidebar-nav-item`)

Used in the 108px chunky industrial rail sidebar for `Search`, `Video`, `Books`, and `Manage Server`. Blends Discord's tactile squircle containers and left indicator pill with YouTube's icon-above-label hierarchy.

### HTML Structure
```html
<nav class="sidebar-nav" aria-label="Media Library">
    <!-- Search Item (Navigates to /search) -->
    <a href="/search" id="nav-search-btn" class="sidebar-nav-item" title="Search Library (/)" aria-label="Search Library">
        <div class="nav-indicator-pill" aria-hidden="true"></div>
        <div class="nav-squircle">
            <svg viewBox="0 0 24 24" width="30" height="30" fill="none" stroke="currentColor" stroke-width="2.5"><circle cx="11" cy="11" r="7"/><line x1="16.5" y1="16.5" x2="22" y2="22"/></svg>
        </div>
        <span class="nav-rail-label">Search</span>
    </a>

    <!-- Video Item -->
    <button id="nav-video-btn" class="sidebar-nav-item active" aria-current="page" title="Videos & Movies">
        <div class="nav-indicator-pill" aria-hidden="true"></div>
        <div class="nav-squircle">
            <svg viewBox="0 0 24 24" width="30" height="30" fill="none" stroke="currentColor" stroke-width="2.4"><rect x="2" y="4" width="20" height="16" rx="4"/><polygon points="10,8 16,12 10,16"/></svg>
        </div>
        <span class="nav-rail-label">Video</span>
    </button>

    <!-- Books Item -->
    <button id="nav-books-btn" class="sidebar-nav-item" title="Books & Publications">
        <div class="nav-indicator-pill" aria-hidden="true"></div>
        <div class="nav-squircle">
            <svg viewBox="0 0 24 24" width="30" height="30" fill="none" stroke="currentColor" stroke-width="2.4"><path d="M2 3h6a4 4 0 0 1 4 4v14a3 3 0 0 0-3-3H2z"/><path d="M22 3h-6a4 4 0 0 0-4 4v14a3 3 0 0 1 3-3h7z"/></svg>
        </div>
        <span class="nav-rail-label">Books</span>
    </button>
</nav>

<div class="sidebar-bottom">
    <div class="sidebar-divider" aria-hidden="true"></div>
    <button id="nav-manage-btn" class="sidebar-nav-item sidebar-manage-item" title="Manage Server">
        <div class="nav-indicator-pill" aria-hidden="true"></div>
        <div class="nav-squircle">
            <svg viewBox="0 0 24 24" width="30" height="30" fill="none" stroke="currentColor" stroke-width="2.4"><rect x="2" y="2" width="20" height="8" rx="2"/><rect x="2" y="14" width="20" height="8" rx="2"/><line x1="6" y1="6" x2="6.01" y2="6"/><line x1="6" y1="18" x2="6.01" y2="18"/></svg>
        </div>
        <span class="nav-rail-label">Manage</span>
    </button>
</div>
```

### Styling Rules
* **Rail Width:** 108px chunky industrial rail (`var(--sidebar-width: 108px)`).
* **Discord Indicator Pill (`.nav-indicator-pill`):** 6px wide bar on left rail edge. Expands to 46px height on active item (`--header-purple`), 18px preview on hover, and 0 height when inactive.
* **Squircle Tile (`.nav-squircle`):** 60px by 60px container, `2.5px solid #000000` border, `16px` border radius, `2.5px 2.5px 0px #000000` shadow. Turns solid purple (`#6d4ca6`) with white icon when active.
* **SVG Icons:** 30px by 30px with 2.4–2.5px stroke weight; punchy visibility from across the room.
* **Stacked Label (`.nav-rail-label`):** 12px (0.75rem) uppercase bold text directly below the squircle with 0.6px letter spacing.
* **Divider Line (`.sidebar-divider`):** 50px wide by 2.5px tall separator before the pinned Manage Server button.
* **Active State Click:** Smooth depression (`:active { transform: translateY(2px); }`).

---

## 3. Dedicated Search Page & Inline Category Chips

Search is implemented as a first-class responsive page (`/search`) rather than a modal dialog to eliminate text truncation, nested scrollbars, and mobile clipping. Full search is triggered via the rail `[ Search ]` button or keyboard shortcuts (`/` or `Ctrl+K`), while quick category switching on catalog pages is supported via inline chips.

### Dedicated Search Page HTML Structure
```html
<div class="search-page-container">
    <header class="search-page-header">
        <div class="search-page-input-wrap">
            <svg class="search-page-input-icon" viewBox="0 0 24 24" width="26" height="26"><circle cx="11" cy="11" r="7"/><line x1="16.5" y1="16.5" x2="22" y2="22"/></svg>
            <input id="search-page-input" class="search-page-input" type="text" placeholder="Search across videos, movies, series, books..." autocomplete="off">
            <button id="search-page-clear-btn" class="search-page-clear-btn" title="Clear Search">✕</button>
        </div>

        <div class="search-page-controls-row">
            <div class="search-page-filter-groups">
                <div class="search-page-filter-group">
                    <span class="search-page-filter-label">Type:</span>
                    <button class="search-page-chip active" data-filter="type" data-val="all">All</button>
                    <button class="search-page-chip" data-filter="type" data-val="video">Videos</button>
                    <button class="search-page-chip" data-filter="type" data-val="books">Books</button>
                </div>
                <div class="search-page-filter-group">
                    <span class="search-page-filter-label">Status:</span>
                    <button class="search-page-chip active" data-filter="status" data-val="all">All</button>
                    <button class="search-page-chip" data-filter="status" data-val="in-progress">In Progress</button>
                    <button class="search-page-chip" data-filter="status" data-val="unwatched">Unwatched</button>
                </div>
            </div>
            <div class="search-page-meta-status">Enter search terms above</div>
        </div>
    </header>

    <section class="search-page-results-section">
        <!-- Initial prompt state shown before typing -->
        <div class="search-prompt-card">
            <div class="search-prompt-icon">...</div>
            <div class="search-prompt-title">Search Your Library</div>
            <p class="search-prompt-desc">Start typing a title, director, author, or keyword to find videos and books.</p>
        </div>

        <!-- Or media-grid once user types, reusing standard catalog media-card markup -->
        <!-- <section class="media-grid" id="search-media-grid"> ... </section> -->
    </section>
</div>
```

### Inline Category Chips HTML Structure (`.catalog-header-bar`)
```html
<div class="catalog-header-bar">
    <div class="catalog-title-wrap">
        <h2 class="catalog-title">Videos & Movies</h2>
        <span class="catalog-count-badge">23 items</span>
    </div>
    <div class="catalog-category-chips">
        <button class="cat-chip active">All</button>
        <button class="cat-chip">Movies</button>
        <button class="cat-chip">Series</button>
        <button class="cat-chip">In Progress</button>
    </div>
</div>
```

---

## 4. Top Header & Utility Controls (`.top-header`)

The persistent top header bar (`.top-header`, 70px height) anchors the application with generous architectural proportions across desktop and mobile viewports. It carries the brand title on the left (text only, no logo or emblems) and utility controls on the right.

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
* **Header Bar:** 70px height (`var(--header-height: 70px)`), solid purple background (`var(--header-purple): #724799`), `3px solid #000000` bottom border, sticky at top.
* **Brand Title:** Bold white text (`font-size: 1.7rem; font-weight: 800; letter-spacing: -0.4px`), interactive cursor, navigates to `#video` on click/Enter. Text-only (no logo icon).
* **Manual Button (`#help-btn`):** 52px squircle white button (`border-radius: 14px`), `2.5px solid #000000` border, `2px 2px 0px #000000` tactile shadow, bold centered `?` symbol (`font-size: 1.55rem; font-weight: 900`), navigates to `#manual`.
* **Avatar Badge (`#user-avatar`):** 52px squircle button (`border-radius: 14px`), `2.5px solid #000000` border (zero blue borders), `2px 2px 0px #000000` tactile shadow, clean `#fdedd0` cream background, opens Unified My Profile modal on click.

---

## 5. Continue Watching / Jump Back In Shelf (`.continue-shelf`)

Surfaces active in-progress media sessions at the top of the catalog:

### Visual & Functional Rules
* **Video Resume Cards:** Shows item title, episode name or subtitle, high-contrast `[ In Progress ]` status tag, and a 1-tap `[ Resume ]` button linking to `/watch/{file_id}`. No progress bars, time percentages, or duration math.
* **Book Reading Cards:** Shows book title, author, `[ Reading ]` status tag, and a 1-tap `[ Continue ]` button. No page count strings or percentages.
* **Activity Sorting:** Sorted by last interaction timestamp descending (`last_watched_at` / `last_read_at DESC`).
* **Zero-Waste Collapse:** If 0 items are in progress or if the user filters by `Unwatched`, the shelf collapses completely from the DOM.
* **Zero-Overhead Queries:** Eliminates complex database timecode or EPUB page parsing queries.

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
