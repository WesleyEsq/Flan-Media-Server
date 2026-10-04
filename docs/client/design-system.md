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

## 3. Dedicated Top Header & Discord x YouTube Hybrid Rail Layout

A dedicated top header bar (`header.top-header`, 70px height) sits at the top of the interface across both desktop and mobile viewports, providing persistent brand identity (text-only) and utility controls (`?` manual button and user avatar). Platform navigation is housed in a chunky 108px Discord x YouTube hybrid rail on desktop (shifting to a 4-tab 64px bottom bar on mobile). Search is accessed as a dedicated full-canvas page from the rail, allowing the "Continue Watching" shelf to sit directly at the top of the video catalog.

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

### 1. Top Header Bar (`header.top-header`)
* **Height:** 70px (`var(--header-height: 70px)`, sticky top across desktop, 64px on mobile).
* **Background:** Deep purple (`var(--header-purple): #724799`) with `3px solid #000000` bottom border.
* **Left:** Interactive brand title `Flan Media Server` (1.7rem, weight 800) linking back to the `#video` catalog on click/Enter. Text-only (no logo icon).
* **Right:** Utility tools group containing squircle `(?)` software manual button (52px) and user avatar profile badge (52px, zero blue borders, tactile shadow).

### 2. Discord x YouTube Hybrid Rail (`aside.left-sidebar`)
* **Width:** 108px (`var(--sidebar-width: 108px)`, height `calc(100vh - 70px)` on desktop).
* **Background:** Soft lilac surface (`var(--sidebar-lilac): #d1b8ee`) with `3px solid #000000` right border.
* **Nav Items:** `[ Q ] Search` $\rightarrow$ `[ |> ] Video` $\rightarrow$ `[ [] ] Books` $\rightarrow$ `[Divider]` $\rightarrow$ `[ # ] Manage`.
* **Discord Squircle Tiles:** 60px by 60px tactile containers with `2.5px solid #000000` borders and 16px border radius. Active tile fills with `--active-purple` and white icon.
* **Icons:** 30px by 30px SVGs with 2.4–2.5px bold stroke weight.
* **Discord Left Indicator Pill:** 6px wide bar on the rail's left edge. Expands to 46px on active items and 18px on hover.
* **YouTube Stacked Labels:** Crisp 12px uppercase bold labels centered directly below each squircle tile (`SEARCH`, `VIDEO`, `BOOKS`, `MANAGE`).
* **Tactile Divider:** 50px wide by 2.5px tall divider line separating media libraries from pinned `Manage Server`.

### 3. Main Content Canvas & Media Grid
* **Background:** Pale lavender canvas (`var(--canvas-lavender): #f4effa`).
* **Instant Media Immersion:** No in-canvas search bar canyon. The **Continue Watching top shelf sits right at the top of the canvas**, offering immediate access to active sessions.
* **Catalog Header & Category Chips:** Clean `.catalog-header-bar` featuring the section title, item count badge, and tactile 1-click filter chips (`[ All ] [ Movies ] [ Series ]` for videos, `[ All ] [ EPUB ] [ PDF ]` for books).
* **Dedicated Search Page:** Search is accessed via the `[ Search ]` rail button or keyboard shortcuts (`/` or `Ctrl+K`), navigating to a full-canvas, responsive search page (`#search`). Results appear once the user types a query, reusing the standard catalog `.media-card` components (with zero Drive badges) to maintain 100% visual consistency and prevent title truncation. An initial prompt card is displayed when no query is typed.
* **Fluid Tiered Scaling (7-Column Max Hard Limit):** `.media-grid` maintains a hard cap of 7 columns on desktop viewports. To prevent cavernous empty lavender space on large monitors and low zoom settings (80%, 67%, 50%), `--catalog-max-width` progressively scales:
  * Baseline: `1680px`
  * Displays $\ge 1680\text{px}$: `1840px` (cards $\approx 235\text{px}$ wide)
  * Displays $\ge 1920\text{px}$: `2060px` (cards $\approx 270\text{px}$ wide)
  * Displays $\ge 2200\text{px}$: `2280px` (cards $\approx 305\text{px}$ wide)
  * Displays $\ge 2560\text{px}$: `2400px` (cards $\approx 325\text{px}$ wide, hard boundary cap)
  Card heights scale naturally via `aspect-ratio: 3 / 4.4`, and inner SVGs / typography scale up proportionally.
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

## 5. Software Manual & User Guide (`/manual`)

The offline user guide provides direct, practical guidance organized into 7 sections: About Flan, Interface Navigation, Watching Videos & Player Controls, Reading Books, Adding Media & Storage Layout, Household Profiles & PINs, and Host CLI Administration.

### Video Player Shortcuts
* `Space` / `K`: Play or Pause
* `F`: Toggle fullscreen
* `M`: Mute or unmute audio
* `Left` / `Right Arrow`: Skip backward or forward 10 seconds
* `Up` / `Down Arrow`: Adjust volume

### Supported Media Formats
* **Video:** Direct-play MP4 (H.264/AAC) and WebM (VP9/Opus, AV1) in browser HTML5 `<video>`. Matroska (MKV) containers and multi-channel audio tracks (AC3, E-AC3, DTS) stream directly via short-lived signed URLs in external players (VLC, MPV) or direct download.
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
