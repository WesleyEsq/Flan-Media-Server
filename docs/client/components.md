# Web Client Component Specifications

Functional specifications, markup structures, aspect ratios, and interaction states for core UI components in Flan Media Server.

---

## 1. Media Card Component (`.card`)

The primary catalog unit representing movies, TV series, episodes, or books.

### Variants & Aspect Ratios
* **Video Poster (`.card--video`, `.card--series`):** Standard 2:3 aspect ratio (`aspect-ratio: 2 / 3`). Min width: 150px.
* **Book Cover (`.card--book`):** Standard 1:1.4 aspect ratio (`aspect-ratio: 1 / 1.4`). Includes a subtle 3px faux spine shadow on the left edge (`box-shadow: inset 4px 0 0 rgba(0, 0, 0, 0.5)`).
* **Episode Row (`.card--episode`):** Horizontal row with a 16:9 thumbnail, episode title, runtime, summary, and progress bar.

### HTML Structure
```html
<article class="card card--video">
    <a href="/watch/movie/42" class="card__link">
        <div class="card__poster-wrap">
            <img src="/covers/movie/42" alt="Title" class="card__poster" loading="lazy">
            <span class="card__badge">4K</span>
            <!-- Progress Bar (rendered only if position_seconds > 10) -->
            <div class="card__progress"><div class="card__progress-fill" style="width: 65%;"></div></div>
        </div>
        <div class="card__info">
            <h3 class="card__title">Dune: Part Two</h3>
            <p class="card__meta">2024 • 2h 46m <span class="card__rating">★ 8.6</span></p>
        </div>
    </a>
</article>
```

### Styling Rules
* Background is `--bg-surface`, border is `--border-subtle`, radius is `8px`.
* Title truncates with ellipsis after 1–2 lines (`overflow: hidden; text-overflow: ellipsis`).
* Hover / focus states apply a crisp outline (`border-color: var(--border-focus)`) with zero layout transform scaling.

---

## 2. Horizontal Shelf Component (`.shelf`)

Used on the dashboard for "Continue Watching", "Continue Reading", and "Recently Added".

### HTML Structure
```html
<section class="shelf">
    <div class="shelf__header">
        <h2 class="shelf__title">Continue Watching</h2>
        <div class="shelf__controls">
            <button class="shelf__arrow" data-dir="prev" aria-label="Scroll left">‹</button>
            <button class="shelf__arrow" data-dir="next" aria-label="Scroll right">›</button>
        </div>
    </div>
    <div class="shelf__track">
        <!-- Media Cards -->
        <article class="card card--video"> ... </article>
    </div>
</section>
```

### Styling Rules
* `.shelf__track` uses native CSS scroll-snap (`overflow-x: auto; scroll-snap-type: x mandatory`).
* Cards inside `.shelf__track` are flex items (`flex: 0 0 160px; scroll-snap-align: start`).
* Arrow buttons advance the track by one viewport width via vanilla JavaScript.

---

## 3. Profile Avatar Tile Component (`.avatar-tile`)

Displayed on the "Who is watching?" screen (`/login`) and profile selector.

### HTML Structure
```html
<button class="avatar-tile" data-user-id="1" data-username="Wesley">
    <div class="avatar-tile__badge" style="background-color: #bb9af7;">
        <svg class="avatar-tile__icon" viewBox="0 0 24 24"><!-- SVG Icon --></svg>
    </div>
    <span class="avatar-tile__name">Wesley</span>
</button>
```

### Styling Rules
* Badge tile: 110px by 110px, rounded corners (`border-radius: 16px`), customizable background color.
* Embedded SVG icon (hamster, cat, popcorn, dog, robot, book) filled with dark slate (`#12131a`).
* Active/hover states display a high-contrast focus ring (`outline: 2px solid var(--accent-lavender)`).

---

## 4. Numeric PIN Keypad Component (`.pin-pad`)

Touch and keyboard-friendly numeric input for entering 4 to 6-digit PINs.

### HTML Structure
```html
<div class="pin-pad">
    <div class="pin-pad__dots">
        <span class="pin-pad__dot pin-pad__dot--filled"></span>
        <span class="pin-pad__dot"></span>
        <span class="pin-pad__dot"></span>
        <span class="pin-pad__dot"></span>
    </div>
    <div class="pin-pad__grid">
        <button class="pin-pad__key" data-key="1">1</button>
        <!-- Keys 2 through 9 -->
        <button class="pin-pad__key pin-pad__key--action" data-key="backspace">⌫</button>
        <button class="pin-pad__key" data-key="0">0</button>
        <button class="pin-pad__key pin-pad__key--action pin-pad__key--submit" data-key="enter">↵</button>
    </div>
</div>
```

### Styling Rules
* Dot indicators: 14px circles; filled dots take `--accent-lavender`.
* Key grid: 3-column CSS Grid (`grid-template-columns: repeat(3, 1fr)`), 52px key height.
* Key clicks trigger haptic/visual feedback; keyboard number keys (`0`–`9`, `Backspace`, `Enter`) hook directly into the same handler.

---

## 5. Modal Dialog Component (`<dialog class="modal">`)

Used for the e-manual, admin uploads, and "Fix Match" editing.

```html
<dialog class="modal" id="modal-id">
    <div class="modal__container">
        <header class="modal__header">
            <h2 class="modal__title">Dialog Title</h2>
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

### Styling Rules
* Native `<dialog>` element with `::backdrop` styled with `rgba(10, 11, 15, 0.8)`.
* Container background is `--bg-elevated`, bordered with `--border-subtle`, max width 600px.
* Automatically traps keyboard focus and closes on `Escape`.

---

## 6. Buttons & Interactive Controls (`.btn`)

| Variant | Class | Styling & Visual Identity | Usage |
| :--- | :--- | :--- | :--- |
| **Primary** | `.btn--primary` | Solid lavender background (`--accent-lavender`), dark slate text (`#12131a`), bold (600). | Primary calls to action: "Resume", "Upload", "Save", "Submit PIN". |
| **Secondary** | `.btn--secondary` | Transparent background, subtle border (`--border-subtle`), off-white text. | "Cancel", "Back", "Close", secondary filters. |
| **Danger** | `.btn--danger` | Transparent background, red border and text (`--color-danger`). Hover fills red. | "Delete Profile", "Remove Library", "Delete Media". |
| **Filter Pill** | `.pill` | Compact badge (`padding: 4px 12px; border-radius: 16px`). Active state fills lavender. | Genre filtering pills, format tabs (`[ MOVIES ]`, `[ TV ]`). |

---

## 7. Related Documentation

* [Design System & Tokens](design-system.md)
* [Page Templates & Wireframes](pages.md)
* [Master System Architecture](../design.md)
