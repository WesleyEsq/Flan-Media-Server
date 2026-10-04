# Mobile Responsiveness & Adaptive Navigation

Breakpoint specifications and CSS layout strategies for mobile and small-screen devices.

---

## 1. Principles

1. **CSS-Only Transitions:** Navigation shifts between sidebar and bottom bar using pure CSS media queries without JavaScript drawer toggles or overlay backdrops.
2. **Bottom Navigation for Touch:** On screens 768px wide or smaller, the navigation moves to a fixed 56px bottom bar for thumb accessibility.
3. **Minimum Hit Targets:** All buttons maintain minimum 44px by 44px hit targets with instant active depression (`:active { transform: translateY(2px); }`).
4. **Zero Overhead:** Responsive adaptations rely entirely on standard CSS with no client-side runtime overhead.

---

## 2. Screen Breakpoints

* **Ultra-Wide / 4K Desktop & 50% Zoom Boundary (`>= 2560px`):** `--catalog-max-width: 2400px` (strict hard cap preventing infinite stretch on 32:9 monitors or 50% crazy zoom). Discord x YouTube Rail (108px), Top Header (70px), 7 columns, cards $\approx 325\text{px}$ wide, 21 items max per page.
* **Large Desktop & 67% Zoom (`2200px - 2559px`):** `--catalog-max-width: 2280px`, Discord x YouTube Rail (108px), Top Header (70px), 7 columns, cards $\approx 305\text{px}$ wide.
* **Standard 1080p Desktop & 80% Zoom (`1920px - 2199px`):** `--catalog-max-width: 2060px`, Discord x YouTube Rail (108px), Top Header (70px), 7 columns, cards $\approx 270\text{px}$ wide.
* **Base Wide Desktop / High Res (`1680px - 1919px`):** `--catalog-max-width: 1840px`, Discord x YouTube Rail (108px), Top Header (70px), 7 columns, cards $\approx 235\text{px}$ wide.
* **Standard Desktop / Laptop (`1025px - 1679px`):** Persistent Top Header (70px), Discord x YouTube Hybrid Rail (108px), 4–6 column fluid media grid (`--catalog-max-width: 1680px`).
* **Tablet Landscape (`769px - 1024px`):** Persistent Top Header (70px), Discord x YouTube Hybrid Rail (108px), 3–4 column fluid media grid.
* **Tablet Portrait / Large Mobile (`481px - 768px`):** Top Header (64px), fixed 4-tab bottom navigation bar (64px), 2–3 column media grid.
* **Small Mobile (`<= 480px`):** Top Header (64px), fixed 4-tab bottom navigation bar (64px), 2-column compact media grid.

---

## 3. Layout Diagrams

### Desktop Viewport (`> 768px: Persistent 70px Header + 108px Hybrid Rail`)
```text
+-----------------------------------------------------------------------------------------------+
| FLAN MEDIA SERVER (70px Height, Bold 1.7rem, Title Only)                 (?) [Avatar Frame]  |
+-----------+-----------------------------------------------------------------------------------+
|           |                                                                                   |
|  [  Q  ]  |   === CONTINUE WATCHING =======================================================   |
|  SEARCH   |   +--------------------------+   +--------------------------+                     |
|           |   | [Thumb] In-Progress Item |   | [Thumb] In-Progress Item |                     |
|  [ |> ]   |   |         [ In Progress ]  |   |         [ In Progress ]  |                     |
|   VIDEO   |   |         [ Resume > ]     |   |         [ Resume > ]     |                     |
|  (Active) |   +--------------------------+   +--------------------------+                     |
|           |                                                                                   |
|  [ [] ]   |   VIDEOS & MOVIES   [ All ] [ Movies ] [ Series ] [ In Progress ]      (23 items)  |
|   BOOKS   |   === MEDIA GRID (7 Columns Max Hard Limit) ===================================   |
|           |   +-------+ +-------+ +-------+ +-------+ +-------+ +-------+ +-------+           |
|   -----   |   | Poster| | Poster| | Poster| | Poster| | Poster| | Poster| | Poster|           |
|           |   |-------| |-------| |-------| |-------| |-------| |-------| |-------|           |
|  [  #  ]  |   | Title | | Title | | Title | | Title | | Title | | Title | | Title |           |
|  MANAGE   |   +-------+ +-------+ +-------+ +-------+ +-------+ +-------+ +-------+           |
|           |                                                                                   |
| 108px Rail|   [ <- Prev ]   [ 1 ]   [ 2 ]   [ Next -> ]    Showing 1-21 of 23                 |
+-----------+-----------------------------------------------------------------------------------+
```

### Mobile Viewport (`<= 768px: 64px Top Header + 64px 4-Tab Bottom Bar`)
```text
+-------------------------------------------------------------------------+
| Flan Media Server                                           ( ? ) [Avatar]|  <-- 64px Top Header
+-------------------------------------------------------------------------+
|                                                                         |
|   === CONTINUE WATCHING =============================================   |
|   +-----------------------------------------------------------------+   |
|   | [Thumb] In-Progress Item     [ In Progress ]       [ Resume > ] |   |
|   +-----------------------------------------------------------------+   |
|                                                                         |
|   ALL VIDEOS   [ All ] [ Movies ] [ Series ]                            |
|                                                                         |
|   +-----------------------+     +-----------------------+               |
|   | [Poster Area]         |     | [Poster Area]         |               |
|   |-----------------------|     |-----------------------|               |
|   | [Purple Band: Title]  |     | [Purple Band: Title]  |               |
|   +-----------------------+     +-----------------------+               |
|                                                                         |
|   [ <- Prev ]   [ 1 ]   [ 2 ]   [ Next -> ]                             |
|                                                                         |
|   (Content scrolls freely above bottom bar)                             |
|                                                                         |
+-------------------------------------------------------------------------+
|  [ Q ] Search   [ |> ] Video   [ [] ] Books   [ # ] Manage              |  <-- 64px Bottom Bar
+-------------------------------------------------------------------------+
```

---

## 4. Key CSS Implementations

### Bottom Navigation Bar (`<= 768px`)

```css
@media (max-width: 768px) {
  /* Top header remains active at top */
  .top-header {
    height: 64px;
    padding: 0 16px;
  }

  /* Convert left sidebar into fixed 64px 4-tab bottom bar */
  .left-sidebar {
    position: fixed;
    bottom: 0;
    left: 0;
    right: 0;
    width: 100% !important;
    height: 64px !important;
    flex-direction: row !important;
    padding: 0 !important;
    gap: 0 !important;
    align-items: stretch !important;
    border-right: none !important;
    border-top: 2.5px solid #000000 !important;
    z-index: 100;
    background: #ffffff !important;
    box-shadow: 0 -2px 0px rgba(0, 0, 0, 0.08);
  }

  /* Distribute navigation buttons (Search, Video, Books) */
  .sidebar-nav {
    display: flex !important;
    flex-direction: row !important;
    flex: 3 !important;
    height: 100% !important;
    margin: 0 !important;
    padding: 0 !important;
    gap: 0 !important;
    box-sizing: border-box !important;
  }

  /* Pinned Manage Server item */
  .sidebar-bottom {
    flex: 1 !important;
    height: 100% !important;
    margin: 0 !important;
    padding: 0 !important;
    display: flex !important;
    flex-direction: row !important;
    align-items: stretch !important;
    border-top: none !important;
    box-sizing: border-box !important;
  }

  .sidebar-divider {
    display: none !important;
  }

  .sidebar-nav-item {
    flex: 1 !important;
    height: 100% !important;
    margin: 0 !important;
    padding: 4px 2px !important;
    border: none !important;
    border-radius: 0 !important;
    border-right: 2px solid #000000 !important;
    background: #ffffff !important;
    display: flex !important;
    flex-direction: column !important;
    align-items: center !important;
    justify-content: center !important;
    cursor: pointer !important;
    text-align: center !important;
  }

  .sidebar-manage-item {
    border-right: none !important;
  }

  .sidebar-nav-item .nav-indicator-pill {
    display: none !important;
  }

  .sidebar-nav-item .nav-squircle {
    width: 32px !important;
    height: 32px !important;
    border: 2px solid #000000 !important;
    border-radius: 9px !important;
  }

  .sidebar-nav-item .nav-squircle svg {
    width: 18px !important;
    height: 18px !important;
  }

  .sidebar-nav-item .nav-rail-label {
    font-size: 0.65rem !important;
    margin-top: 2px !important;
  }

  .sidebar-nav-item.active {
    background: #f3ecfb !important;
  }

  /* Add bottom offset so scrolling content clears the bottom bar */
  .main-content {
    padding-bottom: 76px !important;
  }
}
```

### Management Console Card Reflow (`<= 768px`)

```css
@media (max-width: 768px) {
  /* Management card headers stack vertically to prevent text truncation */
  .manage-card-header {
    flex-direction: column !important;
    align-items: stretch !important;
    gap: 10px !important;
  }

  .manage-header-actions {
    display: flex !important;
    flex-wrap: wrap !important;
    gap: 8px !important;
    width: 100% !important;
  }

  .manage-header-actions .btn {
    flex: 1 1 auto !important;
  }

  /* Storage source item layout reflows into stacked rows */
  .storage-source-item {
    flex-direction: column !important;
    align-items: stretch !important;
    gap: 12px !important;
  }

  .storage-source-actions {
    display: flex !important;
    flex-wrap: wrap !important;
    gap: 8px !important;
    width: 100% !important;
    justify-content: flex-end !important;
  }

  .storage-source-actions .btn {
    flex: 1 1 auto !important;
    min-height: 40px !important;
  }
}
```

### Media Grid Reflow

```css
@media (max-width: 768px) {
  /* Maintain 2-column grid on mobile */
  .media-grid {
    grid-template-columns: repeat(2, 1fr) !important;
    gap: 10px !important;
  }

  .media-card {
    aspect-ratio: 2 / 2.9 !important;
    border: 2px solid #000000 !important;
  }

  .card-footer-band {
    height: 40px !important;
    padding: 3px 6px !important;
  }

  .card-title {
    font-size: 0.78rem !important;
  }
}

@media (max-width: 480px) {
  .media-grid {
    grid-template-columns: repeat(2, 1fr) !important;
    gap: 8px !important;
  }

  .main-content {
    padding: 10px 8px 72px 8px !important;
  }
}
```

### Progressive Large-Screen & Low-Zoom Scaling (Approach A Toolbar + 2400px Cap)

```css
/* Tier 1: Base Wide / 1080p Standard (min-width: 1680px) */
@media (min-width: 1680px) {
  :root {
    --catalog-max-width: 1840px;
  }
  .media-grid {
    grid-template-columns: repeat(7, 1fr) !important;
    gap: 20px !important;
  }
  .poster-play-btn svg { width: 54px; height: 54px; }
  .card-title { font-size: 0.98rem; }
}

/* Tier 2: 1440p / 80% Zoom (min-width: 1920px) */
@media (min-width: 1920px) {
  :root {
    --catalog-max-width: 2060px;
  }
  .media-grid {
    gap: 22px !important;
  }
  .poster-play-btn svg { width: 62px; height: 62px; }
  .card-title { font-size: 1.05rem; }
  .search-input { height: 52px; font-size: 1.05rem; }
  .search-btn, .filter-btn { height: 52px; font-size: 1.05rem; }
}

/* Tier 3: 4K / 67% Zoom (min-width: 2200px) */
@media (min-width: 2200px) {
  :root {
    --catalog-max-width: 2280px;
  }
  .top-header { height: 64px; }
  .left-sidebar { width: 96px; }
  .nav-squircle { width: 56px; height: 56px; }
  .media-grid { gap: 24px !important; }
  .poster-play-btn svg { width: 68px; height: 68px; }
  .card-title { font-size: 1.12rem; }
}

/* Tier 4: Ultra-Wide / 50% Zoom Hard Cap (min-width: 2560px) */
@media (min-width: 2560px) {
  :root {
    --catalog-max-width: 2400px; /* Hard ceiling prevents runaway stretching */
  }
  .media-grid { gap: 26px !important; }
  .poster-play-btn svg { width: 76px; height: 76px; }
  .card-title { font-size: 1.2rem; }
  .card-subtitle { font-size: 0.96rem; }
}
```

---

## 5. Related Documentation

* [Design System](design-system.md)
* [Component Specifications](components.md)
* [Page Templates & Wireframes](pages.md)
* [Accessibility Guide](accessibility.md)
* [Master System Architecture](../design.md)
