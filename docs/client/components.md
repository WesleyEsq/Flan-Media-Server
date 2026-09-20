# Web Client Component Specifications

This document defines the HTML structure, CSS styling, and interaction behaviors for the core UI components in Flan Media Server's web client.

---

## 1. Media Card Component (`.card`)

The media card is the central building block of the interface. It represents a standalone film, a TV series, an episode, or a book.

### Card Types & Aspect Ratios

+ **Video Poster Card (`.card--video`):** Uses the standard 2:3 poster ratio (`aspect-ratio: 2 / 3`). Minimum width: 150px.
+ **TV Series Card (`.card--series`):** Same 2:3 ratio, with an episode count badge in the metadata line.
+ **Book Card (`.card--book`):** Uses a 1:1.4 book ratio (`aspect-ratio: 1 / 1.4`). Styled with a subtle 3px vertical border on the left edge (`box-shadow: inset 3px 0 0 rgba(0, 0, 0, 0.4)`) to visually emulate a physical book spine.
+ **Episode Row (`.card--episode`):** Horizontal row with a 16:9 thumbnail, episode summary, duration, and progress bar.

### HTML Structure (Movie / TV Series / Book)

```html
<article class="card card--video">
    <a href="/watch/movie/42" class="card__link">
        <div class="card__poster-wrap">
            <img src="/covers/movie/42" alt="Dune: Part Two" class="card__poster" loading="lazy">
            <span class="card__badge">4K</span>
            <!-- Progress Bar (Only rendered if position_seconds > 0) -->
            <div class="card__progress">
                <div class="card__progress-fill" style="width: 65%;"></div>
            </div>
        </div>
        <div class="card__info">
            <h3 class="card__title">Dune: Part Two</h3>
            <p class="card__meta">2024 • 2h 46m <span class="card__rating">★ 8.6</span></p>
        </div>
    </a>
</article>
```

### CSS Styling

```css
.card {
    display: flex;
    flex-direction: column;
    background-color: var(--bg-surface);
    border: 1px solid var(--border-subtle);
    border-radius: 8px;
    overflow: hidden;
}

.card__link {
    display: flex;
    flex-direction: column;
    text-decoration: none;
    color: inherit;
}

.card__poster-wrap {
    position: relative;
    width: 100%;
    aspect-ratio: 2 / 3;
    background-color: #1a1b24;
}

.card--book .card__poster-wrap {
    aspect-ratio: 1 / 1.4;
    box-shadow: inset 4px 0 0 rgba(0, 0, 0, 0.5);
}

.card__poster {
    width: 100%;
    height: 100%;
    object-fit: cover;
    display: block;
}

.card__badge {
    position: absolute;
    top: 8px;
    right: 8px;
    background-color: rgba(18, 19, 26, 0.85);
    color: var(--text-primary);
    font-size: 0.75rem;
    font-weight: 600;
    padding: 2px 6px;
    border-radius: 4px;
    border: 1px solid var(--border-subtle);
}

.card__progress {
    position: absolute;
    bottom: 0;
    left: 0;
    right: 0;
    height: 4px;
    background-color: rgba(0, 0, 0, 0.6);
}

.card__progress-fill {
    height: 100%;
    background-color: var(--accent-lavender);
}

.card__info {
    padding: 10px;
}

.card__title {
    font-size: 0.9375rem;
    font-weight: 500;
    color: var(--text-primary);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}

.card__meta {
    font-size: 0.8125rem;
    color: var(--text-secondary);
    margin-top: 4px;
}

/* Hover and Focus - No animation, high-contrast outline */
.card__link:focus-visible,
.card:hover {
    border-color: var(--border-focus);
}
```

---

## 2. Horizontal Shelf Component (`.shelf`)

Used on the dashboard for Continue Watching and Continue Reading. Content scrolls horizontally with native CSS snapping.

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
        <article class="card card--video"> ... </article>
        <article class="card card--video"> ... </article>
    </div>
</section>
```

### CSS Styling

```css
.shelf {
    margin-bottom: 2rem;
}

.shelf__header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 0.75rem;
}

.shelf__title {
    font-size: 1.25rem;
    font-weight: 600;
    color: var(--text-primary);
}

.shelf__controls {
    display: flex;
    gap: 0.5rem;
}

.shelf__arrow {
    background-color: var(--bg-surface);
    border: 1px solid var(--border-subtle);
    color: var(--text-primary);
    width: 32px;
    height: 32px;
    border-radius: 4px;
    cursor: pointer;
    font-size: 1.25rem;
    display: flex;
    align-items: center;
    justify-content: center;
}

.shelf__arrow:hover {
    border-color: var(--accent-lavender);
}

.shelf__track {
    display: flex;
    gap: 1rem;
    overflow-x: auto;
    scroll-snap-type: x mandatory;
    scrollbar-width: thin;
    scrollbar-color: var(--border-subtle) transparent;
    padding-bottom: 8px;
}

.shelf__track > .card {
    flex: 0 0 160px;
    scroll-snap-align: start;
}
```

---

## 3. Profile Avatar Tile Component (`.avatar-tile`)

Displayed on the "Who is watching?" screen (`/login`).

### HTML Structure

```html
<button class="avatar-tile" data-user-id="1" data-username="Wesley">
    <div class="avatar-tile__badge" style="background-color: #bb9af7;">
        <!-- Embedded SVG icon -->
        <svg class="avatar-tile__icon" viewBox="0 0 24 24">
            <!-- Icon path for hamster mascot, cat, popcorn, etc. -->
        </svg>
    </div>
    <span class="avatar-tile__name">Wesley</span>
</button>
```

### CSS Styling

```css
.avatar-tile {
    display: flex;
    flex-direction: column;
    align-items: center;
    background: none;
    border: none;
    cursor: pointer;
    gap: 0.75rem;
}

.avatar-tile__badge {
    width: 110px;
    height: 110px;
    border-radius: 16px;
    display: flex;
    align-items: center;
    justify-content: center;
    border: 2px solid transparent;
}

.avatar-tile__icon {
    width: 56px;
    height: 56px;
    fill: #12131a;
}

.avatar-tile__name {
    font-size: 1rem;
    font-weight: 500;
    color: var(--text-primary);
}

.avatar-tile:focus-visible .avatar-tile__badge,
.avatar-tile:hover .avatar-tile__badge {
    border-color: #ffffff;
    outline: 2px solid var(--accent-lavender);
}
```

---

## 4. Numeric PIN Keypad Component (`.pin-pad`)

Provides touch and mouse input for entering 4 to 6-digit PINs, while supporting physical keyboard numbers.

### HTML Structure

```html
<div class="pin-pad">
    <div class="pin-pad__dots">
        <span class="pin-pad__dot pin-pad__dot--filled"></span>
        <span class="pin-pad__dot pin-pad__dot--filled"></span>
        <span class="pin-pad__dot"></span>
        <span class="pin-pad__dot"></span>
    </div>
    <div class="pin-pad__grid">
        <button class="pin-pad__key" data-key="1">1</button>
        <button class="pin-pad__key" data-key="2">2</button>
        <button class="pin-pad__key" data-key="3">3</button>
        <button class="pin-pad__key" data-key="4">4</button>
        <button class="pin-pad__key" data-key="5">5</button>
        <button class="pin-pad__key" data-key="6">6</button>
        <button class="pin-pad__key" data-key="7">7</button>
        <button class="pin-pad__key" data-key="8">8</button>
        <button class="pin-pad__key" data-key="9">9</button>
        <button class="pin-pad__key pin-pad__key--action" data-key="backspace">⌫</button>
        <button class="pin-pad__key" data-key="0">0</button>
        <button class="pin-pad__key pin-pad__key--action pin-pad__key--submit" data-key="enter">↵</button>
    </div>
</div>
```

### CSS Styling

```css
.pin-pad {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 1.5rem;
    width: 240px;
}

.pin-pad__dots {
    display: flex;
    gap: 1rem;
}

.pin-pad__dot {
    width: 14px;
    height: 14px;
    border-radius: 50%;
    border: 2px solid var(--border-subtle);
    background-color: transparent;
}

.pin-pad__dot--filled {
    background-color: var(--accent-lavender);
    border-color: var(--accent-lavender);
}

.pin-pad__grid {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 0.75rem;
    width: 100%;
}

.pin-pad__key {
    height: 52px;
    background-color: var(--bg-surface);
    border: 1px solid var(--border-subtle);
    border-radius: 8px;
    color: var(--text-primary);
    font-size: 1.25rem;
    font-weight: 500;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
}

.pin-pad__key:hover,
.pin-pad__key:focus-visible {
    border-color: var(--accent-lavender);
    background-color: var(--bg-elevated);
}

.pin-pad__key--submit {
    background-color: var(--accent-lavender);
    color: #12131a;
    font-weight: 700;
}
```

---

## 5. Modal Dialog Component (`<dialog class="modal">`)

Used for the built-in e-manual, the admin upload dialog, and metadata editing.

### HTML Structure

```html
<dialog class="modal" id="manual-dialog">
    <div class="modal__container">
        <header class="modal__header">
            <h2 class="modal__title">User Guide & Shortcuts</h2>
            <button class="modal__close" aria-label="Close dialog">✕</button>
        </header>
        <div class="modal__body">
            <!-- Modal content -->
        </div>
        <footer class="modal__footer">
            <button class="btn btn--secondary modal__close-btn">Close</button>
        </footer>
    </div>
</dialog>
```

### CSS Styling

```css
.modal {
    border: none;
    background: transparent;
    padding: 0;
    max-width: 600px;
    width: 90%;
    margin: auto;
}

.modal::backdrop {
    background-color: rgba(10, 11, 15, 0.8);
}

.modal__container {
    background-color: var(--bg-elevated);
    border: 1px solid var(--border-subtle);
    border-radius: 10px;
    display: flex;
    flex-direction: column;
    overflow: hidden;
}

.modal__header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1rem 1.25rem;
    border-bottom: 1px solid var(--border-subtle);
}

.modal__title {
    font-size: 1.125rem;
    font-weight: 600;
    color: var(--text-primary);
}

.modal__close {
    background: none;
    border: none;
    color: var(--text-secondary);
    font-size: 1.25rem;
    cursor: pointer;
}

.modal__close:hover {
    color: var(--text-primary);
}

.modal__body {
    padding: 1.25rem;
    max-height: 70vh;
    overflow-y: auto;
}

.modal__footer {
    padding: 1rem 1.25rem;
    border-top: 1px solid var(--border-subtle);
    display: flex;
    justify-content: flex-end;
    gap: 0.75rem;
}
```

---

## 6. Buttons & Interactive Controls (`.btn`)

### Button Variants

```css
.btn {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    padding: 0.5rem 1rem;
    border-radius: 6px;
    font-size: 0.9375rem;
    font-weight: 500;
    cursor: pointer;
    text-decoration: none;
    border: 1px solid transparent;
}

/* Primary: Lavender filled */
.btn--primary {
    background-color: var(--accent-lavender);
    color: #12131a;
    font-weight: 600;
}

.btn--primary:hover {
    background-color: var(--accent-hover);
}

/* Secondary: Subtle border */
.btn--secondary {
    background-color: transparent;
    border-color: var(--border-subtle);
    color: var(--text-primary);
}

.btn--secondary:hover {
    border-color: var(--text-secondary);
}

/* Danger: Red warning */
.btn--danger {
    background-color: transparent;
    border-color: var(--color-danger);
    color: var(--color-danger);
}

.btn--danger:hover {
    background-color: var(--color-danger);
    color: #12131a;
}
```

---

### Related Documentation

+ [Web Client Design System & Tokens](design-system.md)
+ [Web Client Pages & Template Specifications](pages.md)
+ [Master System Specifications](../design.md)
