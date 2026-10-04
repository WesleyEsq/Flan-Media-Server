# Accessibility Standards & Verification Guide

Standards and technical invariants for maintaining WCAG 2.1 AA compliance across Flan Media Server's web client and HTML templates.

---

## 1. Core Requirements

* **Input Modality Support:** Interfaces must be navigable using keyboard navigation, TV remotes, touch screens, and assistive screen readers.
* **Semantic HTML:** Semantic HTML5 elements (`<dialog>`, `<nav>`, `<a>`, `<button>`) must be used instead of styled generic `<div>` tags.
* **Zero Script Overhead:** Accessibility features rely on native browser primitives and CSS variables without adding background JavaScript processing.

---

## 2. Technical Invariants

### A. Layout Containment
The left sidebar (or bottom navigation bar on mobile) must remain static within the viewport.
* The sidebar is constrained to viewport height (`calc(100vh - 56px)` on desktop).
* In populated library views, only the main content area (`.main-content`) scrolls internally.
* `.app-container` and `.main-content` declare `min-height: 0` and `overflow-y: auto` to prevent content from pushing navigation buttons off-screen.

### B. Color Contrast Tokens
All text and essential controls must satisfy WCAG AA contrast ratios:
* Standard text (< 18pt or < 14pt bold): minimum 4.5:1 ratio.
* Large text and interactive borders/icons: minimum 3.0:1 ratio.

**Approved Pairings:**
* White text (`#ffffff`) on Header Purple (`#724799`): Contrast 6.8:1 (Pass AA)
* White text (`#ffffff`) on Active Button Purple (`#6d4ca6`): Contrast 6.5:1 (Pass AA)
* White text (`#ffffff`) on Access Button (`#5b3794`): Contrast 8.7:1 (Pass AAA)
* Black text (`#000000`) on VLC Button (`#ff6f00`): Contrast 7.5:1 (Pass AAA)
* Dark Red text (`#700000`) on Lilac Panel (`#c7aee7`): Contrast 6.3:1 (Pass AA)
* Dark Red text (`#8b0000`) on White Card (`#ffffff`): Contrast 10.0:1 (Pass AAA)
* Green text (`#2e7d32`) on White Card (`#ffffff`): Contrast 5.1:1 (Pass AA)
* Off-black text (`#101015`) on Pale Lavender Canvas (`#f4effa`): Contrast 16.8:1 (Pass AAA)

**Contrast Rules:**
* Never use white text on orange background (`#ff6f00`); always use black text (`#000000`).
* Never use bright red text on tinted backgrounds; use dark red (`#700000`).

### C. Keyboard Navigation & Focus Management
* **Focus Indicator:** All interactive elements must declare a visible outline:
  ```css
  :focus-visible {
    outline: 3px solid #000000 !important;
    outline-offset: 2px !important;
  }
  ```
  `outline: none` is forbidden unless paired with `:focus-visible`.
* **Skip Link:** The first interactive element in `<body>` must be:
  ```html
  <a href="#main-content" class="skip-link">Skip to main content</a>
  ```
  It remains hidden off-screen until focused with `Tab`.
* **Modal Dialogs:** Modals must use the HTML5 `<dialog>` element. Focus is automatically trapped when opened via `.showModal()`, `Escape` closes the dialog, and focus is restored to the initiating button upon dismissal.

### D. Screen Reader Attributes
* **Descriptive Action Labels:** Avoid ambiguous repeated button labels:
  ```html
  <!-- Preferred -->
  <button aria-label="Play S01E01 - Pilot">Play</button>
  <button aria-label="Download or stream in VLC: S01E01 - Pilot">VLC / Download</button>
  ```
* **Icon Buttons:** SVG-only buttons must declare an `aria-label`:
  * Software manual: `aria-label="Software Manual & About"`
  * User avatar: `aria-label="My Profile & Settings" aria-haspopup="dialog"`
  * Search button: `aria-label="Search catalog"`
  * Modal close: `aria-label="Close dialog"`
* **Active Navigation State:** The current navigation link must include `aria-current="page"`.
* **Dynamic Feedback:** Asynchronous status messages must use live regions:
  ```html
  <div id="login-feedback" role="alert" aria-live="assertive"></div>
  <div id="rescan-feedback" role="alert" aria-live="polite"></div>
  ```

### E. Motion & Text Selection
* **Reduced Motion:** When `prefers-reduced-motion: reduce` is detected, animations and smooth scrolling must be disabled:
  ```css
  @media (prefers-reduced-motion: reduce) {
    *, *::before, *::after {
      animation-duration: 0.01ms !important;
      transition-duration: 0.01ms !important;
      scroll-behavior: auto !important;
    }
    .media-card:hover .card-title.is-overflowing .title-text,
    .media-card:focus-visible .card-title.is-overflowing .title-text {
      animation: none !important;
    }
  }
  ```
* **Text Selection:** `user-select: none` is prohibited on readable text content.

---

## 3. Verification Checklist

* [ ] **Tab Navigation:** Pressing `Tab` from a fresh page load focuses the Skip Link first, followed by header controls, sidebar items, and catalog cards.
* [ ] **Focus Rings:** Every interactive control displays a 3px black focus ring when focused via keyboard.
* [ ] **Dialog Trap & Escape:** Opening a dialog moves focus to the first input. Pressing `Escape` closes the dialog and returns focus to the triggering element.
* [ ] **Contrast Compliance:** All text elements meet or exceed 4.5:1 contrast against their backgrounds.
* [ ] **ARIA Labels:** Icon buttons and action rows contain descriptive `aria-label` attributes.

---

## 4. Related Documentation

* [Design System](design-system.md)
* [Component Specifications](components.md)
* [Mobile Responsiveness](responsiveness.md)
* [Page Templates & Wireframes](pages.md)
