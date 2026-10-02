# Mobile Responsiveness & Adaptive Navigation Architecture

Design philosophy, breakpoint specifications, and CSS layout strategies for delivering a fast, tactile mobile experience on single-board computer homelabs without native app development or JavaScript UI frameworks.

---

## 1. Core Mobile Philosophy

1. **Zero JavaScript Drawer Lag:** Avoid sliding hamburger navigation drawers that require overlay backdrops, gesture listeners, and z-index acrobatics. Navigation state is handled with pure CSS media queries.
2. **Thumb-Zone Ergonomics:** On mobile smartphones ($< 768\text{px}$), the user's thumb rests naturally at the bottom of the viewport. The navigation bar shifts from the left sidebar to a **fixed bottom navigation bar**.
3. **Touch Target Accessibility:** All buttons and interactive card elements enforce minimum $44\text{px} \times 44\text{px}$ touch targets with instantaneous tactile depression (`:active { transform: translateY(2px); }`) and zero hover delay.
4. **Memory Neutral:** Adapting layout with CSS media queries consumes zero additional Go heap memory or background CPU cycles.

---

## 2. Breakpoint Matrix

| Viewport Category | Screen Width Range | Sidebar / Navigation | Media Catalog Grid | Split Login Screen |
| :--- | :--- | :--- | :--- | :--- |
| **Desktop / Laptop** | `> 1024px` | Persistent Left Sidebar (160px) | 4–6 columns fluid | Side-by-side (50% / 50%) |
| **Tablet (Landscape)** | `769px – 1024px` | Persistent Left Sidebar (160px) | 3–4 columns fluid | Side-by-side (50% / 50%) |
| **Mobile / Tablet (Portrait)** | `481px – 768px` | **Fixed Bottom Navigation Bar (56px)** | 2–3 columns fluid | Stacked vertically |
| **Small Smartphone** | `≤ 480px` | **Fixed Bottom Navigation Bar (56px)** | 2 columns compact (scaled cards) | Stacked vertically |

---

## 3. Layout Transformations

### A. Desktop Layout (`> 768px`)
```text
+-------------------------------------------------------------------------+
| Flan Media Server                                           ( ? ) [Avatar]|
+----------+--------------------------------------------------------------+
| ←─────── |                                                              |
| [ Video] |   [ Search catalog...                          ]  [ 🔍 ]     |
| [ Books] |                                                              |
| ────────►|   +---------------+  +---------------+  +---------------+    |
|          |   | [Poster Area] |  | [Poster Area] |  | [Poster Area] |    |
|          |   |---------------|  |---------------|  |---------------|    |
|          |   | [Purple Band] |  | [Purple Band] |  | [Purple Band] |    |
| [Manage] |   +---------------+  +---------------+  +---------------+    |
+----------+--------------------------------------------------------------+
```

### B. Mobile Layout (`≤ 768px`)
```text
+-------------------------------------------------------------------------+
| Flan Media Server                                           ( ? ) [Avatar]|
+-------------------------------------------------------------------------+
|                                                                         |
|   [ Search catalog...                                       ]  [ 🔍 ]   |
|   Videos (5 titles)                                                     |
|                                                                         |
|   +-----------------------+     +-----------------------+               |
|   | [Poster Area]         |     | [Poster Area]         |               |
|   |-----------------------|     |-----------------------|               |
|   | [Purple Band: Title]  |     | [Purple Band: Title]  |               |
|   +-----------------------+     +-----------------------+               |
|                                                                         |
|   (Content scrolls freely above bottom bar)                             |
|                                                                         |
+-------------------------------------------------------------------------+
|        [ Video ]              [ Books ]              [ Manage ]         |  <-- Fixed Bottom Bar
+-------------------------------------------------------------------------+
```

---

## 4. CSS Media Query Implementation Guide

### 1. Adaptive Navigation Bar (Sidebar to Bottom Bar)

```css
/* ==========================================================================
   DESKTOP DEFAULT (Left Sidebar)
   ========================================================================== */
.left-sidebar {
  width: 160px;
  background-color: var(--sidebar-lilac);
  border-right: 2px solid #000000;
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
}

/* ==========================================================================
   MOBILE BREAKPOINT (<= 768px): Relocate to Bottom Bar
   ========================================================================== */
@media (max-width: 768px) {
  /* Convert left sidebar into flush bottom bar */
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
    border-top: 2.5px solid #000000;
    z-index: 100;
    background: #ffffff;
    box-shadow: 0 -2px 0px rgba(0, 0, 0, 0.08);
  }

  /* Distribute navigation buttons evenly across bottom width */
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
    background: var(--active-purple) !important;
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
    background: var(--active-purple) !important;
    color: #ffffff !important;
  }

  .desktop-manage-label { display: none !important; }
  .mobile-manage-label { display: inline !important; }

  /* Hide desktop decorative arrow graphics on mobile */
  .sidebar-arrow {
    display: none !important;
  }

  /* Add bottom offset so scrolling content is never covered by bottom bar */
  .main-content {
    padding-bottom: 76px !important;
  }
}
```

---

### 2. Split-Screen Start & Login Adaptation

On desktop, the start screen divides into two side-by-side panels. On mobile, it seamlessly stacks vertically:

```css
@media (max-width: 768px) {
  .split-login-container {
    flex-direction: column;
    height: auto;
    min-height: calc(100vh - 56px);
  }

  .login-left-panel {
    border-right: none;
    border-bottom: 2px solid #000000;
    padding: 30px 20px;
  }

  .welcome-banner-text {
    font-size: 2.2rem;
  }

  .login-right-panel {
    padding: 32px 20px 60px 20px;
  }

  .login-form-box {
    max-width: 100%;
  }
}
```

---

### 3. Media Catalog Grid Reflow & Card Scaling

On mobile screens, media cards are scaled down and locked to a **2-column grid** rather than collapsing into a giant 1-column stack. This allows users to see 4–6 cards per screen.

```css
@media (max-width: 768px) {
  /* Maintain 2-column tactile grid */
  .media-grid {
    grid-template-columns: repeat(2, 1fr) !important;
    gap: 10px !important;
  }

  /* Compact card aspect ratio & typography */
  .media-card {
    aspect-ratio: 2 / 2.9 !important;
    border: 2px solid #000000 !important;
    box-shadow: 2px 2px 0px #000000 !important;
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

### 4. Container Detail & Player Views

```css
@media (max-width: 768px) {
  /* Stack poster above metadata on mobile */
  .detail-header-card {
    flex-direction: column;
    align-items: center;
    text-align: center;
  }

  .detail-poster-box {
    width: 140px;
  }

  /* File item rows wrap buttons for small screens */
  .playable-item-row {
    flex-direction: column;
    align-items: flex-start;
    gap: 10px;
  }

  .playable-item-actions {
    width: 100%;
    justify-content: flex-end;
  }
}
```

---

## 5. Related Documentation

* [Master System Architecture](design.md)
* [Design System & High-Contrast Foundations](client/design-system.md)
* [Tactile Component Specifications](client/components.md)
* [Page Templates & Wireframes](client/pages.md)
