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

## 4. Circular Header Action Buttons (`.header-icon-btn`)

Positioned on the right side of the top purple bar for utility controls (no navigation links).

### HTML Structure
```html
<!-- Help Button -->
<button class="header-icon-btn" id="open-help" aria-label="Help">?</button>

<!-- Manage Shortcut -->
<a href="/manage" class="header-icon-btn" aria-label="Manage Server">
    <svg class="icon" viewBox="0 0 24 24"><path d="M...gear"/></svg>
</a>

<!-- Profile Avatar Badge -->
<a href="/login" class="header-avatar-badge" aria-label="Profile">
    <img src="/static/assets/avatars/flan.svg" alt="Avatar">
</a>
```

### Styling Rules
* Circular white background, 36px diameter, thick `2px solid #000000` border, centered content.
* Avatar badge is a 36px rounded square with `2px solid #000000` border.

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

## 6. Playable File Row (`.file-row`)

Displayed inside the Video or Book detail view to list playable episodes, cuts, or book volumes.

### HTML Structure
```html
<div class="file-row">
    <div class="file-row__info">
        <span class="file-row__title">S01E01 - Pilot</span>
        <span class="file-row__duration">48m</span>
    </div>
    <!-- Progress Indicator (if started) -->
    <div class="file-row__progress">
        <div class="file-row__progress-fill" style="width: 45%;"></div>
    </div>
    <a href="/watch/101" class="btn btn--primary">Play</a>
</div>
```

---

## 7. Modal Dialog Component (`<dialog class="modal">`)

Used for the e-manual, upload dialog, and user editing.

```html
<dialog class="modal" id="modal-id">
    <div class="modal__container">
        <header class="modal__header">
            <h2 class="modal__title">Title</h2>
            <button class="modal__close" aria-label="Close">✕</button>
        </header>
        <div class="modal__body"><!-- Content --></div>
        <footer class="modal__footer">
            <button class="btn btn--secondary">Cancel</button>
            <button class="btn btn--primary">Confirm</button>
        </footer>
    </div>
</dialog>
```

* High-contrast `2px solid #000000` border, white surface, native `<dialog>` with focus trapping.

---

## 8. Related Documentation

* [Design System & Tokens](design-system.md)
* [Page Templates & Wireframes](pages.md)
* [Master System Architecture](../design.md)
