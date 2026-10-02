# Web Client Component Specifications

Functional specifications, markup structures, and tactile interaction behaviors for core UI components in Flan Media Server.

---

## 1. Media Card Component (`.card`)

The primary catalog unit representing a video or book container.

### Visual & Tactile Rules
* **Borders:** Crisp `2px solid #000000` border, `8px` rounded corners.
* **Poster Area:** Pure white surface (`#ffffff`) holding the media cover image.
* **Footer Band:** Solid purple strip (`#6d52a8`) across the bottom of the card with pure white title text (`#ffffff`), as shown in the mockup.
* **Interactions:** Zero hover scaling or float delays. Clicks immediately into active state (`transform: translateY(2px)`).
* **Aspect Ratios:**
  * **Video Poster (`.card--video`):** Standard 2:3 aspect ratio.
  * **Book Cover (`.card--book`):** Standard 1:1.4 aspect ratio.

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

## 2. Exclusive Sidebar Navigation Buttons (`.sidebar-btn`)

Used in the left navigation sidebar for **`Video`**, **`Books`**, and **`Manage Server`**.

### HTML Structure
```html
<nav class="sidebar-nav">
    <a href="/video" class="sidebar-btn sidebar-btn--active">Video</a>
    <a href="/books" class="sidebar-btn">Books</a>
    <a href="/manage" class="sidebar-btn sidebar-btn--bottom">Manage Server</a>
</nav>
```

### Styling Rules
* **Border:** Thick `2px solid #000000`, `8px` rounded corners.
* **Active State (`.sidebar-btn--active`):** Solid darker purple (`#6d52a8`) background with pure white bold text.
* **Inactive State:** Light gray/white surface (`#eae8f2`) with dark bold text (`#12131a`).
* **Tactile Click:** Depresses slightly on click (`:active { transform: translateY(2px); }`) with zero hover float delays.

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
* Input field has a pure white background with a bold `2px solid #000000` border and `8px` rounded corners.
* Search button is a high-contrast square with `2px solid #000000` border and pure SVG magnifying glass.

---

## 4. Header Utility Controls (`.header-tools`)

Positioned on the right side of the top purple bar strictly for utility controls (no navigation links).

### HTML Structure
```html
<!-- Software Manual Trigger (Navigates to /manual or #view-manual) -->
<button class="header-icon-btn" id="open-help" aria-label="Software Manual & About">?</button>

<!-- Profile Avatar Badge (Opens My Profile Modal) -->
<button class="header-avatar-badge" id="open-profile" aria-label="My Profile">
    <img src="/static/assets/avatars/mascot.svg" alt="Avatar">
</button>
```

### Styling Rules
* Manual button is a circular white surface, 36px diameter, thick `2px solid #000000` border, centered `?` symbol.
* Avatar badge is a 40px rounded rectangle with `3px solid #3ea6ff` accent border and `2px solid #000000` outline, which depresses on click to open the "My Profile" modal.

---

## 5. Split-Screen Start & Login Form (`.start-screen`)

Replaces the complex 3x4 on-screen keypad and avatar tiles with the clean split design from Image 1.

### HTML Structure
```html
<div class="start-screen">
    <!-- Left Panel: Graphic Welcome -->
    <div class="start-screen__welcome">
        <div class="welcome-arrow welcome-arrow--left">←───────────</div>
        <h1 class="welcome-title">Welcome</h1>
        <div class="welcome-arrow welcome-arrow--right">───────────→</div>
    </div>

    <!-- Right Panel: Accessible Access Form -->
    <div class="start-screen__form-panel">
        <form class="access-form" action="/api/login" method="POST">
            <label class="access-form__label" for="user-select">User:</label>
            <select class="access-form__select" id="user-select" name="user_id">
                <option value="1">mike</option>
            </select>

            <label class="access-form__label" for="pin-input">Pin:</label>
            <input class="access-form__input" id="pin-input" type="password" name="pin" maxlength="6" autofocus>

            <button type="submit" class="access-form__btn">Access</button>
        </form>
    </div>
</div>
```

### Styling Rules
* `User:` select and `Pin:` input fields have white/light-gray backgrounds, thick `2px solid #000000` borders, and rounded corners.
* `[ Access ]` button is solid purple with white text, thick `2px solid #000000` border, and instant tactile depression on click.

---

## 6. Playable File Row & Fallback Action (`.file-row`)

Displayed inside the Video or Book detail view to list playable episodes, cuts, or book volumes.

### HTML Structure
```html
<!-- Video File Row with Direct Play & VLC Fallback -->
<div class="file-row">
    <div class="file-row__info">
        <span class="file-row__title">S01E01 - Pilot</span>
        <span class="file-row__duration">48m</span>
    </div>
    <!-- Progress Indicator (if started) -->
    <div class="file-row__progress">
        <div class="file-row__progress-fill" style="width: 45%;"></div>
    </div>
    <div class="file-row__actions">
        <a href="/watch/101" class="btn btn--primary">Play</a>
        <a href="/download/video/101" class="btn btn--vlc" title="Direct download or stream in VLC">⬇ VLC / Download</a>
    </div>
</div>

<!-- Book File Row with Read / Download -->
<div class="file-row">
    <div class="file-row__info">
        <span class="file-row__title">Volume 1 (PDF)</span>
        <span class="file-row__duration">18.4 MB</span>
    </div>
    <div class="file-row__actions">
        <a href="/stream/book/201" target="_blank" class="btn btn--primary">📖 Open PDF</a>
        <a href="/download/book/201" class="btn btn--secondary">⬇ Download</a>
    </div>
</div>
```

---

## 7. Whimsical Avatar Picker Component (`.avatar-picker`)

Used in `/manage` and accessible by clicking the header avatar badge to customize household profiles.

### Visual & Tactile Rules
* **Preset Grid:** 3-column tactile grid displaying 6 high-contrast SVG companions (Mascot, Flan, Cat, Ghost, Robot, Star).
* **Active Indicator:** 3px solid `#6d52a8` border with 4px offset shadow.
* **Custom Photo Upload:** Direct image file input (`accept="image/png, image/jpeg, image/webp"`, enforced max 2 MB).
* **Instant Tactile Feedback:** Selecting an avatar immediately updates the header avatar badge without a page reload.

---

## 8. Modal Dialog Component (`<div class="modal-card">` or `<dialog>`)

Used for the unified **My Profile** modal (name, avatar picker/upload, PIN, log out) and **Add New User** modal in `/manage`.

```html
<div class="modal-card">
    <header class="modal-header">
        <h2 class="modal-title">My Profile</h2>
        <button class="modal-close-btn" aria-label="Close">✕</button>
    </header>
    <div class="modal-body"><!-- Content --></div>
</div>
```

* **Sticky Pinned Header:** `modal-header` is sticky (`position: sticky; top: -24px; background: #fff;`) so the title and `✕` close button remain visible at all times while scrolling.
* **Tactile Scrollbar:** Bounded to `max-height: 85vh` with `overflow-y: auto` and a custom high-contrast purple thumb scrollbar with 2px borders.

---

## 9. Compact Catalog Status & Filter Strip

Replaces heavy section headings with an unobtrusive, high-density status strip (~24px height, `0.78rem` uppercase monospace/bold tactile font).

* **Default Browsing Mode:** Displays section identity and total count without pushing content down:
  `ALL VIDEOS & SERIES • 5 TITLES` or `BOOKS & PUBLICATIONS • 3 TITLES`
* **Search / Filter Mode:** Smoothly mutates into an active filter tag with an interactive 1-click reset button:
  `[FILTER] "gundam" (1 match)   [ ✕ Clear ]`
* **Touch & Click:** Clicking `[ ✕ Clear ]` resets search input and immediately restores full library view.

---

## 10. Horizontal Title Marquee on Hover

When media container titles exceed the width of the card's solid purple footer band:

* **Idle State:** Text truncates with a clean ellipsis (`text-overflow: ellipsis; white-space: nowrap;`).
* **Hover State:** If `scrollWidth > clientWidth`, `.is-overflowing` is dynamically added, activating a CSS marquee animation:
  * 0%–20%: Pauses at origin so the opening words are readable.
  * 80%–100%: Smoothly translates left by `calc(scrollWidth - clientWidth + 8px)` and pauses so the end of the title is readable.
  * Resets smoothly on mouse leave. Titles that fit without overflow remain completely static.

---

## 11. Related Documentation

* [Design System & Tokens](design-system.md)
* [Page Templates & Wireframes](pages.md)
* [Master System Architecture](../design.md)

