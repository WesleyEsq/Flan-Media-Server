# Web Client Design System & Foundations

Visual tokens, typography scale, sidebar layout, and interaction principles for Flan Media Server's web client.

---

## 1. Core Visual Principles

* **Sidebar Navigation:** Platform navigation is consolidated into the left sidebar (`Video`, `Books`, and `Manage Server`). The top header bar contains no page navigation links.
* **High-Contrast Styling:** 2px and 3px solid black borders (`#000000`), solid purple header and sidebar, and clean, high-readability surfaces.
* **Instant Transitions:** No float transforms, hover scale delays, or blur filters (`backdrop-filter`). Controls provide immediate visual feedback.
* **Tactile Buttons:** Buttons depress slightly on `:active` (`transform: translateY(2px)`), making interaction evident on touchscreens, mice, and TV remotes.
* **Zero Emojis:** Standard SVG vector icons are used for buttons, indicators, and controls.
* **System Font Stack:** Uses native system fonts with zero external font network requests:
  `font-family: system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif`
* **Accessibility Compliance:** WCAG 2.1 AA compliant contrast ratios ($\ge 4.5:1$) on text surfaces, visible keyboard focus rings (`:focus-visible`), support for `prefers-reduced-motion: reduce`, unrestricted user text selection, and skip-to-content links.

---

## 2. Color Palette & Design Tokens

Styles are declared as CSS variables in `web/static/css/style.css`:

```css
:root {
    /* Brand & Structural Surfaces (WCAG AA Compliant) */
    --brand-purple:     #724799; /* Header bar background (contrast >= 6.8:1 with pure white text) */
    --brand-active:     #6d4ca6; /* Active button background (contrast 6.5:1 with white text) */
    --brand-access:     #5b3794; /* Login Access button (contrast 8.7:1 with white text) */
    --bg-main:          #f4f0fa; /* Main content area (pale lavender) */
    --surface-card:     #ffffff; /* Card background */
    --surface-button:   #eae8f2; /* Inactive button background */

    /* High-Contrast Borders & Focus */
    --border-black:     #000000; /* Crisp 2px and 3px structural borders */
    --border-subtle:    #2c2e3e; /* Inner divider borders */
    --border-focus:     #000000; /* High-visibility keyboard/remote focus outline (:focus-visible 3px) */

    /* Typography */
    --text-primary:     #12131a; /* Pitch dark off-black on light surfaces (contrast 16.8:1) */
    --text-header:      #ffffff; /* Pure white text on purple surfaces */
    --text-muted:       #555869; /* Secondary metadata text (contrast 6.3:1) */

    /* Functional Accents & Alerts */
    --accent-lavender:  #6d4ca6; /* Primary action color, card title bands */
    --color-danger:     #700000; /* Lockouts, login errors (contrast 6.3:1 on lilac, 11:1 on white) */
    --color-success:    #2e7d32; /* Completed indicators (contrast 5.1:1 on white) */
    --btn-vlc-bg:       #ff6f00; /* VLC action button (pair strictly with #000000 text for 7.5:1 contrast) */
    --btn-vlc-text:     #000000;
}
```

---

## 3. Dedicated Top Header & Streamlined Rail Sidebar Layout

A dedicated top header bar (`header.top-header`, 56px height) sits at the top of the interface across both desktop and mobile viewports, providing persistent brand identity and utility controls (`?` manual button and user avatar). Platform navigation is streamlined into a clean left rail sidebar on desktop (shifting to a fixed bottom rail on mobile). Scoped search remains within the catalog canvas.

```text
+-------------------------------------------------------------------------------------+
| Flan Media Server                                                   ( ? )  [Avatar] |
+----------+--------------------------------------------------------------------------+
|          |                                                                          |
| [ Video] |   +---------------------------------------+  [Search]  [ Filter [v] ]    |
| (Active) |   | Search input...                       |                              |
|          |   +---------------------------------------+                              |
| [ Books] |                                                                          |
|          |   === CONTINUE WATCHING ==============================================   |
|          |   +--------------------------+   +--------------------------+            |
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

### 1. Top Header Bar (`header.top-header`)
* **Height:** 56px (sticky top across desktop and mobile).
* **Background:** Deep purple (`var(--header-purple): #6d4ca6`) with `2.5px solid #000000` bottom border.
* **Left:** Interactive brand title `Flan Media Server` linking back to the `#video` catalog on click/Enter.
* **Right:** Utility tools group containing circular `(?)` software manual button (38px) and user avatar profile badge (40px).

### 2. Streamlined Left Rail Sidebar (`aside.left-sidebar`)
* **Width:** 160px (height `calc(100vh - 56px)` on desktop).
* **Background:** Soft lilac surface (`var(--sidebar-lilac): #e2d9f3`) with `2.5px solid #000000` right border.
* **Navigation Rail:** Minimalist industrial buttons for `Video` and `Books` at the top, and pinned `Manage Server` at the bottom.
* **Zero Distractions:** Zero decorative arrows (`.sidebar-arrow`) and zero counter badges.

### 3. Main Content Canvas & Media Grid
* **Background:** Pale lavender canvas (`var(--canvas-lavender): #f4f0fa`).
* **Search & Filter Drawer:** In-canvas search bar paired with tactile `[ Filter [v] ]` toggle button revealing format, status, source drive, and sort order options.
* **Continue Watching Top Shelf:** Surfaces in-progress items sorted by interaction date (`last_watched_at DESC`). Collapses cleanly if 0 items are in progress.
* **7-Column Max Hard Limit:** `.media-grid` is capped at a strict maximum of 7 columns (`max-width: 1680px; @media (min-width: 1680px) { grid-template-columns: repeat(7, 1fr); }`) to ensure cards remain comfortably readable on 4K and ultra-wide displays.
* **21 Items Max / Page Pagination:** Strict chunking to 21 items max per page ($7 \times 3 = 21$) with tactile numbered pagination controls.

---

## 4. 2-Step Sequential Login (`login.html`)

Unauthenticated visitors experience a clean, sequential 2-step authentication flow:

1. **Step 1 (Profile Selection):**
   * Displays the brand title and generous household profile cards with 72px avatar frames and role pills.
   * Selecting an avatar triggers a smooth transition into Step 2.
2. **Step 2 (Dedicated PIN Prompt):**
   * Displays the chosen user's identity banner (`.login-user-banner`), framed avatar, and role tag.
   * Auto-focuses the numeric PIN input (`inputmode="numeric"`).
   * Provides a primary **`[ Access Library → ]`** action button and a secondary **`[ ← Switch Profile ]`** tactile button.
   * Wrong PIN entries shake the input field (`.shake`) with high-contrast error messaging.

Both states render inside a centered console card (`.login-pad-card`) with solid 2px black borders and a 4px shadow offset on pale lavender canvas (`--bg-main: #f4f0fa`).

---

## 5. Software Manual & Keyboard Shortcuts (`/manual`)

The manual view provides direct guidance for video playback and file formats:

### Video Player Shortcuts
* `Space` / `K`: Play or Pause
* `F`: Toggle fullscreen
* `M`: Mute or unmute audio
* `Left` / `Right Arrow`: Skip backward or forward 10 seconds
* `Up` / `Down Arrow`: Adjust volume in 5% increments

### Supported Media Formats
* **Video:** Direct-play MP4 (H.264/AAC), WebM (VP9/Opus, AV1), and web-safe MKV. A signed `[ VLC / Download ]` link is available for external playback of multi-channel AC3 or DTS audio.
* **Books:** PDF (browser-native viewer tab) and EPUB (direct file download).

---

## 6. Avatar Options

* **Preset SVG Icons:** Included in `web/static/assets/avatars/`:
  * Mascot (Flan boy)
  * Flan (Custard)
  * Cat
  * Ghost
  * Robot
  * Star
* **Custom Photo Avatars:** Users can upload a personalized photo via `/manage` or the profile modal (JPEG, PNG, WebP up to 2MB).

---

## 7. Related Documentation

* [Component Specifications](components.md)
* [Page Templates & Wireframes](pages.md)
* [Accessibility Guide](accessibility.md)
* [Mobile Responsiveness](responsiveness.md)
* [Master System Architecture](../design.md)
