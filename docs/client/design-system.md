# Web Client Design System & Foundations

Visual tokens, typography scale, sidebar-only layout, and tactile interaction principles for Flan Media Server's web client.

---

## 1. Core Visual Principles

* **Sidebar-Only Navigation:** Platform navigation is strictly consolidated into the left sidebar (`Video`, `Books`, and `Manage Server`). The top header bar contains zero navigation links.
* **High-Contrast Neo-Tactile Aesthetic:** Bold 2px/3px black borders (`#000000`), solid purple header and sidebar, and high-readability surfaces.
* **No Random Hovers or Delayed Transitions:** Zero bouncy float transforms, no hover scale delays, and no blur filters (`backdrop-filter`). Controls give instant, crisp visual feedback.
* **Clicks In Place:** Physical, tactile button feel. Buttons depress slightly on `:active` (`transform: translateY(2px)`), making interaction immediately obvious and accessible on touchscreens, mice, and TV remotes.
* **Zero Emojis:** Pure SVG vector icons are used exclusively for all buttons, avatars, and indicators. No OS-dependent emoji rendering.
* **System Font Stack:** Uses the device's native system font stack with zero external font network requests (`font-family: system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif`).

---

## 2. Color Palette & Design Tokens

All styles are declared as CSS variables in `web/static/css/style.css`:

```css
:root {
    /* Brand & Structural Surfaces */
    --brand-purple:     #9c7cd8; /* Header bar and left sidebar background */
    --brand-active:     #6d52a8; /* Active/selected button background */
    --bg-main:          #f4f0fa; /* Main content area (pale lavender) */
    --surface-card:     #ffffff; /* Card background */
    --surface-button:   #eae8f2; /* Inactive button background */

    /* High-Contrast Borders */
    --border-black:     #000000; /* Crisp 2px and 3px structural borders */
    --border-subtle:    #2c2e3e; /* Inner divider borders */
    --border-focus:     #000000; /* High-visibility keyboard/remote focus outline */

    /* Typography */
    --text-primary:     #12131a; /* Pitch dark off-black on light surfaces */
    --text-header:      #ffffff; /* Pure white text on purple surfaces */
    --text-muted:       #555869; /* Secondary metadata text */

    /* Functional Accents */
    --accent-lavender:  #9c7cd8; /* Primary action color, card title bands */
    --color-danger:     #e5534b; /* Lockouts, delete actions */
    --color-success:    #57ab5a; /* Completed indicators */
}
```

---

## 3. Sidebar-Only Global Layout

All navigation across the entire platform occurs through the left sidebar. The top header is reserved exclusively for branding and utility tools.

```text
+-------------------------------------------------------------------------+
| Flan Media Server                                           ( ? )   [Avatar]|
+----------+--------------------------------------------------------------+
| ←─────── |                                                              |
| [ Video] |   +---------------------------------------------+  [ 🔍 ]    |
|          |   | Search input...                             |            |
| [ Books] |   +---------------------------------------------+            |
| ────────►|                                                              |
|          |   +---------------+  +---------------+  +---------------+    |
|          |   | [Poster Area] |  | [Poster Area] |  | [Poster Area] |    |
|          |   |---------------|  |---------------|  |---------------|    |
|          |   | [Purple Band] |  | [Purple Band] |  | [Purple Band] |    |
| [Manage] |   +---------------+  +---------------+  +---------------+    |
+----------+--------------------------------------------------------------+
```

### 1. Top Header Bar (`header.top-bar`)
* **Background:** Solid purple (`--brand-purple: #9c7cd8`). Height: 60px.
* **Left Title:** Bold text `Flan Media Server` (white, high-contrast, ~1.5rem).
* **Right Utility Controls (No navigation links):**
  * `( ? )` **Software Manual Button:** White circle, 36px diameter, thick 2px black border. Navigates to the dedicated Software Handbook view (`/manual`).
  * `[Avatar]` **User Profile Badge:** 40px rounded rectangle with thick 2px black border, displaying user's selected SVG avatar or custom image. Depresses on click to open the "My Profile" modal.

### 2. Exclusive Left Navigation Sidebar (`aside.sidebar`)
* **Background:** Solid purple (`--brand-purple: #9c7cd8`). Width: ~200px.
* **Divider:** `3px solid #000000` vertical border separating sidebar from main content.
* **Tactile Navigation Buttons:**
  * **`Video`**: Navigates to `/video`. When active, fills with solid darker purple (`#6d52a8`) and white text. When inactive, light gray (`#eae8f2`) with black text.
  * **`Books`**: Navigates to `/books`. Same active/inactive tactile states.
  * **`Manage Server`** *(bottom)*: Navigates to `/manage`.
* All buttons have a thick `2px solid #000000` border and `8px` rounded corners.

### 3. Main Content Area (`main.content`)
* **Background:** Pale lavender (`--bg-main: #f4f0fa`).
* **Search Bar:** Large white input field with `2px solid #000000` border, followed by a square `🔍` search button.
* **Media Grid:** Direct grid of high-contrast cards. Each card has a white poster area and a **solid purple footer band** (`#6d52a8`) at the bottom.

---

## 4. Split-Screen Start & Login (`login.html` - Image 1)

Unauthenticated users see a split-panel screen with zero on-screen keypad bloat:

```text
+-------------------------------------------------------------------------+
| Flan Media Server                                                       |
+------------------------------------------+------------------------------+
|                                          | User:                        |
|   ←───────────────────────────           | +--------------------------+ |
|                                          | | mike                   v | |
|        Welcome                           | +--------------------------+ |
|                                          | Pin:                         |
|   ───────────────────────────►           | +--------------------------+ |
|                                          | |                          | |
|                                          | +--------------------------+ |
|                                          |                              |
|                                          | +--------------------------+ |
|                                          | | Access                   | |
|                                          | +--------------------------+ |
+------------------------------------------+------------------------------+
```

* **Left Panel:** Light lavender background with bold slanted *"Welcome"* between arrows.
* **Right Panel:** Solid purple background with thick 3px black divider border, containing:
  * `User:` dropdown `<select>` of existing user profiles.
  * `Pin:` input `<input type="password">`.
  * `[ Access ]` purple button with white text and thick 2px black border.

---

## 5. Software Manual & System Handbook (`manual.html` - `GET /manual`)

Accessible via the `( ? )` button as a dedicated 2-column view with a sticky Table of Contents and 1-click terminal copy blocks (100% offline).

### Keyboard Shortcuts (Video Player)
| Key | Action |
| :--- | :--- |
| `Space` / `K` | Play / Pause |
| `F` | Toggle fullscreen |
| `M` | Mute / Unmute audio |
| `←` / `→` | Skip backward / forward 10 seconds |
| `↑` / `↓` | Volume up / down (5% increments) |

### Direct Play Media Guidelines
* **Video:** MP4 (H.264/AAC), WebM (VP9/Opus, AV1), and web-safe MKV. Plus 1-click `[ ⬇ VLC / Download ]` fallback for unsupported audio codecs (AC3, EAC3, DTS).
* **Books:** PDF (browser-native viewer in new tab) and EPUB (instant direct download for native reader apps like Apple Books, Moon+ Reader, Kindle).

---

## 6. Whimsical Preset & Custom Avatars

* **Bundled High-Contrast SVG Presets:** Shipped directly inside `web/static/assets/avatars/` as clean high-contrast SVGs:
  * Mascot (Flan boy)
  * Flan (Caramel pudding)
  * Cat
  * Ghost
  * Robot
  * Star
* **Custom Avatar Uploads:** Users can upload a personalized photo via `/manage` or the "My Profile" modal (JPEG, PNG, WebP up to 2MB), saved to `data/avatars/{user_id}.ext`.

---

## 7. Related Documentation

* [Tactile Components](components.md)
* [Page Templates & Wireframes](pages.md)
* [Master System Architecture](../design.md)
