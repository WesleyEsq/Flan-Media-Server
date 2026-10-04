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

* **Wide / 4K Desktop (`>= 1680px`):** Persistent Top Header (56px), Streamlined Left Rail Sidebar (160px), strict 7-column max hard limit on media grid, 21 items max per page.
* **Standard Desktop / Laptop (`1025px - 1679px`):** Persistent Top Header (56px), Streamlined Left Rail Sidebar (160px), 4–6 column fluid media grid.
* **Tablet Landscape (`769px - 1024px`):** Persistent Top Header (56px), Streamlined Left Rail Sidebar (160px), 3–4 column fluid media grid.
* **Tablet Portrait / Large Mobile (`481px - 768px`):** Persistent Top Header (56px), fixed bottom navigation bar (56px), 2–3 column media grid.
* **Small Mobile (`<= 480px`):** Persistent Top Header (56px), fixed bottom navigation bar (56px), 2-column compact media grid.

---

## 3. Layout Diagrams

### Desktop Viewport (`> 768px: Persistent Top Header + Streamlined Rail Sidebar`)
```text
+-------------------------------------------------------------------------------------+
| Flan Media Server                                                   ( ? )  [Avatar] |
+----------+--------------------------------------------------------------------------+
|          |                                                                          |
| [ Video] |   [ Search catalog...                          ] [Search] [ Filter [v] ] |
| (Active) |                                                                          |
|          |   === CONTINUE WATCHING ==============================================   |
| [ Books] |   +--------------------------+   +--------------------------+            |
|          |   | [Thumb] In-Progress Item |   | [Thumb] In-Progress Item |            |
|          |   |         [== Progress ==] |   |         [== Progress ==] |            |
|          |   |         [ Resume > ]     |   |         [ Resume > ]     |            |
|          |   +--------------------------+   +--------------------------+            |
|          |                                                                          |
|          |   === MEDIA GRID (7 Columns Max Hard Limit) ==========================   |
|          |   +-------+ +-------+ +-------+ +-------+ +-------+ +-------+ +-------+  |
|          |   | Poster| | Poster| | Poster| | Poster| | Poster| | Poster| | Poster|  |
|          |   |-------| |-------| |-------| |-------| |-------| |-------| |-------|  |
|          |   | Title | | Title | | Title | | Title | | Title | | Title | | Title |  |
|          |   +-------+ +-------+ +-------+ +-------+ +-------+ +-------+ +-------+  |
|          |                                                                          |
|          |   [ <- Prev ]   [ 1 ]   [ 2 ]   [ Next -> ]    Showing 1-21 of 23        |
| [Manage] |                                                                          |
+----------+--------------------------------------------------------------------------+
```

### Mobile Viewport (`<= 768px: Top Header + Bottom Rail`)
```text
+-------------------------------------------------------------------------+
| Flan Media Server                                           ( ? ) [Avatar]|  <-- Persistent Top Header
+-------------------------------------------------------------------------+
|                                                                         |
|   [ Search catalog...                        ] [Search] [ Filter [v] ]  |
|   ALL VIDEOS - 23 TITLES                                                |
|                                                                         |
|   === CONTINUE WATCHING =============================================   |
|   +-----------------------------------------------------------------+   |
|   | [Thumb] In-Progress Item     [== Progress ==]      [ Resume > ] |   |
|   +-----------------------------------------------------------------+   |
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
|        [ Video ]              [ Books ]              [ Manage ]         |  <-- Fixed Bottom Bar
+-------------------------------------------------------------------------+
```

---

## 4. Key CSS Implementations

### Bottom Navigation Bar (`<= 768px`)

```css
@media (max-width: 768px) {
  /* Convert left sidebar into fixed bottom bar */
  .left-sidebar {
    position: fixed;
    bottom: 0;
    left: 0;
    right: 0;
    width: 100%;
    height: 56px;
    flex-direction: row;
    padding: 0 !important;
    gap: 0 !important;
    align-items: stretch !important;
    border-right: none;
    border-top: 2px solid #000000;
    z-index: 100;
    background: #ffffff;
  }

  /* Distribute navigation buttons evenly across bottom bar */
  .sidebar-nav {
    display: flex !important;
    flex-direction: row !important;
    flex: 2 !important;
    height: 100% !important;
    margin: 0 !important;
    padding: 0 !important;
    gap: 0 !important;
  }

  .sidebar-nav-btn {
    flex: 1;
    height: 100%;
    margin: 0;
    padding: 0 4px;
    border: none;
    border-radius: 0;
    border-right: 2px solid #000000;
    background: #ffffff;
    color: #111111;
    font-size: 0.95rem;
    font-weight: 800;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
  }

  .sidebar-nav-btn.active {
    background: var(--brand-active) !important;
    color: #ffffff !important;
  }

  .sidebar-bottom {
    flex: 1;
    height: 100%;
    margin: 0;
    padding: 0;
    display: flex;
  }

  .sidebar-manage-btn {
    width: 100%;
    height: 100%;
    margin: 0;
    padding: 0 4px;
    border: none;
    border-radius: 0;
    background: #ffffff;
    color: #111111;
    font-size: 0.95rem;
    font-weight: 800;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
  }

  .sidebar-manage-btn.active {
    background: var(--brand-active) !important;
    color: #ffffff !important;
  }

  .desktop-manage-label { display: none !important; }
  .mobile-manage-label { display: inline !important; }

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

---

## 5. Related Documentation

* [Design System](design-system.md)
* [Component Specifications](components.md)
* [Page Templates & Wireframes](pages.md)
* [Accessibility Guide](accessibility.md)
* [Master System Architecture](../design.md)
