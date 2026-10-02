# Flan Accessibility Standards & Future Compliance Guide

Guidelines, rationales, and technical invariants for maintaining **WCAG 2.1 / 2.2 Level AA compliance** across Flan Media Server's web client and templates.

---

## 1. Why Accessibility is Essential for Flan

Flan Media Server is built on the philosophy of extreme simplicity, local-first operation, and low resource overhead (15–20 MB RAM) for single-board computers (Raspberry Pi, homelabs). Accessibility is not an afterthought or an external plugin—it is a core architectural pillar for several key reasons:

1. **Diverse Hardware & Input Modalities:**
   * Living room TVs and media boxes driven by infrared remotes, keyboard dongles, or game controllers rely on strict, predictable keyboard navigation and high-contrast visible focus rings.
   * Handheld smartphones and touch tablets require thumb-friendly bottom navigation and minimum $44\times44\text{px}$ hit targets.
   * Desktop monitors and laptops need fluid responsive reflow without content clipping.
2. **Diverse Household Users:**
   * Homelabs are shared across families, including kids, elderly relatives, and people with varying visual acuity, color perception deficiencies (protanopia, deuteranopia), motor challenges, or vestibular sensitivity.
3. **Zero-Overhead Accessibility:**
   * Semantic HTML5 elements (`<dialog>`, `<nav>`, `<a>`, `<button>`) and CSS variables introduce **zero additional Go heap memory** or background JavaScript CPU cycles. Good accessibility is inherently lightweight.

---

## 2. Core Architectural Invariants (Must Be Kept)

Future template authors, contributors, and maintainers must adhere to the following non-negotiable rules:

### A. Layout Containment & Static Sidebar Navigation

```text
+-------------------------------------------------------------------------+
| Top Header Bar (56px) - Sticky / Fixed                                  |
+-------------------+-----------------------------------------------------+
| Static Sidebar    | Scrollable Main Content (overflow-y: auto)          |
| (Height = 100vh   |                                                     |
|  - 56px)          | [ Search Input ]                                    |
| [ Video ]         |                                                     |
| [ Books ]         | [ Card ]  [ Card ]  [ Card ]  [ Card ]              |
|                   | [ Card ]  [ Card ]  [ Card ]  [ Card ]              |
| [Manage Server]   |                                                     |
| (Pinned bottom)   | (Only this right panel scrolls!)                    |
+-------------------+-----------------------------------------------------+
```

* **Invariant:** The left sidebar must remain **strictly static** at viewport height (`calc(100vh - 56px)` on desktop, or fixed 56px bottom bar on mobile $\le 768\text{px}$).
* **Never let the sidebar stretch with the results grid.** If `.main-content` contains 100 items, only `.main-content` scrolls internally.
* **Flexbox Rule:** `.app-container` and `.main-content` must declare `min-height: 0;` and `overflow: hidden;` / `overflow-y: auto;` respectively. This prevents the parent grid from pushing the pinned `[Manage Server]` button out of the visible screen.

---

### B. Color Contrast & Palette Invariants (WCAG 1.4.3)

All text and essential UI controls must satisfy WCAG AA contrast thresholds:
* **Normal text (< 18pt or < 14pt bold):** Minimum **`4.5:1`** ratio against its background.
* **Large text & interactive borders/icons:** Minimum **`3.0:1`** ratio.

#### Approved Color Pairings Matrix

| Component | Foreground | Background | Contrast Ratio | Compliance Status |
| :--- | :--- | :--- | :--- | :--- |
| **Top Header Title & Icons** | `#ffffff` | `#724799` (Header Purple) | **`6.83:1`** | **PASS (AA)** |
| **Active Sidebar / Card Band** | `#ffffff` | `#6d4ca6` (Active Purple) | **`6.49:1`** | **PASS (AA)** |
| **Login Access Button** | `#ffffff` | `#5b3794` (Access Purple) | **`8.68:1`** | **PASS (AAA)** |
| **VLC / Fallback Button** | **`#000000`** | `#ff6f00` (Amber Orange) | **`7.53:1`** | **PASS (AAA)** |
| **Login Error Feedback** | `#700000` (Deep Red) | `#c7aee7` (Lilac Panel) | **`6.30:1`** | **PASS (AA)** |
| **Modal / Card Error Text** | `#8b0000` (Deep Red) | `#ffffff` (White Card) | **`10.01:1`** | **PASS (AAA)** |
| **Success Indicator Text** | `#2e7d32` (Forest Green) | `#ffffff` (White Card) | **`5.13:1`** | **PASS (AA)** |
| **Primary Text on Main Canvas**| `#101015` (Pitch Dark) | `#f4effa` (Lavender) | **`16.78:1`** | **PASS (AAA)** |
| **Skip-to-Content Link** | `#000000` | `#ffeb3b` (Safety Yellow) | **`17.20:1`** | **PASS (AAA)** |

> [!WARNING]
> **Common Pitfalls to Avoid:**
> 1. **Never use white text on `#ff6f00`:** `#ffffff` on `#ff6f00` yields only **`2.79:1`** (a direct WCAG AA violation). Always use `#000000` text on orange buttons.
> 2. **Never lighten the header purple above `#794ea1`:** Tints lighter than `#724799` (such as `#9b6ebd` or `#9c7cd8`) drop below the 4.5:1 threshold for white text.
> 3. **Never use light red for error text on tinted panels:** On lilac `#c7aee7`, bright red `#b71c1c` has only 3.3:1 contrast; use `#700000` or `#5a0000`.

---

### C. Keyboard Navigation & Focus Visibility (WCAG 2.1.1 & 2.4.7)

1. **Universal Visible Focus Outline:**
   * All focusable controls must render a crisp, high-contrast focus indicator:
     ```css
     :focus-visible {
       outline: 3px solid #000000 !important;
       outline-offset: 2px !important;
     }
     ```
   * **Rule:** Never apply `outline: none;` without providing an immediate `:focus-visible` replacement.
2. **First Element is Skip Link (WCAG 2.4.1):**
   * The first interactive element in `<body>` must always be:
     ```html
     <a href="#main-content" class="skip-link">Skip to main content</a>
     ```
   * It remains off-screen until focused via `Tab`, jumping keyboard users directly over the header and sidebar.
3. **Semantic Interactive Elements (No Clickable Divs):**
   * Media catalog cards must be semantic `<a>` tags with `href` destinations, or `<button>` elements. Never bind click listeners to unadorned `<div>` or `<span>` tags.
4. **Modal Dialog Management (WCAG 2.1.2 & 2.4.3):**
   * Modals must use the HTML5 `<dialog class="modal-dialog">` element with `aria-labelledby`.
   * Open with `.showModal()` (native backdrop + focus trapping).
   * Close with `.close()` on `✕` button, outside backdrop click, or native `Escape` key.
   * **Focus Restoration:** Maintain a reference to `document.activeElement` before opening, and invoke `lastFocusedElement.focus()` upon closing.

---

### D. Screen Reader Clarity & ARIA Directives (WCAG 1.3.1, 2.4.4 & 4.1.3)

1. **Contextual Action Labels:**
   * In lists with multiple action buttons, avoid repetitive naked labels. Provide full context:
     ```html
     <!-- BAD -->
     <button>Play</button>
     <button>Download</button>

     <!-- GOOD -->
     <button aria-label="Play S01E01 - Pilot">▶ Play</button>
     <button aria-label="Download or stream in VLC: S01E01 - Pilot">⬇ VLC / Download</button>
     ```
2. **Icon-Only Buttons:**
   * Buttons displaying only glyphs or SVGs must have an explicit `aria-label`:
     * Manual button: `aria-label="Software Manual & About"`
     * Avatar button: `aria-label="My Profile & Settings" aria-haspopup="dialog"`
     * Search button: `aria-label="Search catalog"`
     * Modal close button: `aria-label="Close dialog"`
3. **Active Page State:**
   * Always reflect the active navigation button with `aria-current="page"`.
4. **Dynamic Status & Error Messages:**
   * Elements that receive dynamic feedback text without page reload (PIN failure, rescan status, profile updates) must declare:
     ```html
     <div id="login-feedback" role="alert" aria-live="assertive"></div>
     <div id="rescan-feedback" role="alert" aria-live="polite"></div>
     ```

---

### E. Motion Sensitivity & Ergonomics (WCAG 2.2.2 & 2.3.3)

1. **Vestibular Safety (`prefers-reduced-motion`):**
   * Users who configure reduced motion in their OS must not experience continuous scrolling animations or shaking elements:
     ```css
     @media (prefers-reduced-motion: reduce) {
       *, *::before, *::after {
         animation-duration: 0.01ms !important;
         animation-iteration-count: 1 !important;
         transition-duration: 0.01ms !important;
         scroll-behavior: auto !important;
       }
       .media-card:hover .card-title.is-overflowing .title-text,
       .media-card:focus-visible .card-title.is-overflowing .title-text {
         animation: none !important;
       }
       .shake {
         animation: none !important;
       }
     }
     ```
2. **Marquee Activation on Keyboard Focus:**
   * Overflowing card title marquees must trigger on `:focus-visible` as well as `:hover`, allowing keyboard and screen magnifier users to read titles that exceed the card footer width.
3. **Unrestricted Text Selection:**
   * Never set `user-select: none;` on the `body` or general text content. Users must be able to select, copy, inspect, or use text-to-speech tools on media titles, metadata, and manual instructions.

---

## 3. Pre-Commit / Pre-Release Verification Checklist

Before shipping new templates, CSS rules, or client features, run through this quick 5-step checklist:

- [ ] **Tab Order Test:** Press `Tab` repeatedly from a clean page load. Does the Skip Link appear first? Does every button, input, and card receive a thick 3px black focus ring? Can every item be opened with `Enter`?
- [ ] **Static Sidebar Test:** In a populated library view, scroll the mouse wheel or arrow keys. Does the results grid scroll freely while the left sidebar (and the bottom "Manage Server" button) remains static and visible in the viewport?
- [ ] **Dialog Trap & Escape Test:** Open the "My Profile" modal with `Enter`. Does focus jump into the profile name input? Press `Escape`. Does the modal close and return focus to the avatar button?
- [ ] **Contrast Calculation:** Check any newly introduced color tokens with the formula \((L_1 + 0.05) / (L_2 + 0.05) \ge 4.5\).
- [ ] **Screen Reader Check:** Ensure all SVG buttons and repetitive list actions have clear, human-understandable `aria-label` attributes.
