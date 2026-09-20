# Web Client Design System & Foundations

This document details the visual foundations, color palette, typography scale, global navigation bar, and built-in e-manual for Flan Media Server's web client.

---

## 1. Core Design Philosophy

The web interface is engineered for fast rendering and high readability across diverse display hardware, from low-brightness mobile screens in dark living rooms to low-powered smart TV web browsers:

+ **Zero Animation & Zero Gradient Bloat:** No bouncy transforms, no hover scaling, and no CSS gradients. Surfaces are flat and matte, ensuring immediate browser painting without GPU rendering lag.
+ **No Blur Filters:** CSS backdrop-filter (blur) is strictly avoided. It is notoriously CPU-heavy on low-power devices and older tablets.
+ **High Contrast First:** Text and icons use high-contrast color values that exceed standard accessibility guidelines (WCAG AAA), avoiding the washed-out gray-on-gray look common in dark themes.
+ **Zero Font Network Overhead:** Uses the device's native system font stack. The browser downloads zero font files, eliminating network delays and layout shifts.

---

## 2. Color Palette & Design Tokens

All colors are declared as CSS custom properties in the root stylesheet:

```css
:root {
    /* Backgrounds & Surfaces */
    --bg-base: #12131a;       /* Deep matte slate background */
    --bg-surface: #1a1b24;    /* Flat card and shelf background */
    --bg-elevated: #242633;   /* Modals, dropdowns, and hover states */
    --bg-active: #2d2640;     /* Selected item background with lavender tint */

    /* Borders & Dividers */
    --border-subtle: #2c2e3e; /* Visible component and card borders */
    --border-focus: #bb9af7;  /* High-visibility focus and active ring */

    /* High-Contrast Typography */
    --text-primary: #f0f2fc;   /* High-contrast off-white (14:1 contrast ratio) */
    --text-secondary: #9ba1b8; /* Muted metadata (year, duration, author) */
    --text-muted: #6b728d;     /* Inactive labels and placeholders */

    /* Accents & Functional Colors */
    --accent-lavender: #bb9af7; /* Primary buttons, active tabs, progress bars */
    --accent-hover: #cbb0f9;    /* Button interaction state */
    --color-danger: #f7768e;    /* Error states and delete confirmations */
    --color-success: #9ece6a;   /* Completed indicators and status badges */
}
```

### Contrast and Readability

+ **Primary Text (--text-primary: #f0f2fc) on Base (--bg-base: #12131a):** Delivers a contrast ratio of over 14:1, making text effortlessly readable from across a room on a television screen.
+ **Lavender Accent (--accent-lavender: #bb9af7):** Used deliberately for actionable controls—active navigation tabs, playback scrub bars, and primary buttons.

---

## 3. Typography Scale

The font family relies entirely on the system font stack to achieve instant rendering with zero network requests:

```css
body {
    font-family: system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
    font-size: 15px;
    line-height: 1.5;
    color: var(--text-primary);
    background-color: var(--bg-base);
}
```

### Type Scale (rem-based)

+ **Display / Page Header:** 1.5rem (24px), bold, letter-spacing: -0.02em. Used for page titles and TV series headings.
+ **Shelf / Section Title:** 1.25rem (20px), 600 weight. Used for row titles (e.g. "Continue Watching", "Movies").
+ **Card Title:** 0.9375rem (15px), 500 weight. Used for movie, episode, and book titles. Truncates cleanly after 2 lines with text-overflow.
+ **Metadata / Subtitle:** 0.8125rem (13px), 400 weight, color: var(--text-secondary). Used for release year, runtime, and author names.
+ **Badges / Tags:** 0.75rem (12px), 600 weight, uppercase, letter-spacing: 0.05em. Used for format labels (4K, MP4, EPUB) and genre pills.

---

## 4. Global Navigation Bar Specification

The top navigation bar is fixed at the top of the viewport (`position: sticky; top: 0; z-index: 50`) across all pages except fullscreen video playback.

```text
+-------------------------------------------------------------------------+
| [Flan logo]   Home   Videos   Books                 [? Help] [⚙] [Avatar] |
+-------------------------------------------------------------------------+
```

### Layout Breakdown

1. **Brand (Left):**
   + Text logo: `Flan` linking directly to the dashboard (`/`).
   + Distinct lavender highlight on hover/focus.
2. **Primary Navigation Links (Center / Left):**
   + **Home (`/`):** Dashboard with in-progress shelves and recent additions.
   + **Videos (`/videos`):** Movie and TV catalog with genre filter bar.
   + **Books (`/books`):** Document and book catalog.
   + Active link is highlighted with an underline border (`border-bottom: 2px solid var(--accent-lavender)`) and bold text.
3. **Utility Controls (Right):**
   + **E-Manual Button (`? Help`):** Opens the built-in user guide modal.
   + **Settings Button (`⚙`):** Opens library and server settings (visible to admin).
   + **Profile Avatar Badge:** Displays the active user's chosen icon on their accent color tile. Clicking opens a dropdown with options to switch profiles or log out.

### Responsive Mobile Behavior

+ On screens narrower than 640px, navigation links remain compact in the header or sit cleanly below the brand in a horizontal scroll row, ensuring full accessibility on phones without requiring a heavy hamburger menu script.

---

## 5. The Built-in E-Manual (User Guide)

The e-manual is a self-contained, accessible user guide embedded directly in every page shell. It opens when the user clicks the `? Help` button in the navigation bar.

### Implementation: Native HTML `<dialog>`

The manual uses the standard HTML `<dialog id="manual-dialog">` element:

+ **Zero External Libraries:** Standard browser dialog with native modal backdrop.
+ **Accessible by Default:** Automatically traps keyboard focus, closes on pressing `Escape`, and restores focus to the trigger button when closed.

### Content Structure of the E-Manual

#### Section 1: Keyboard Shortcuts (Video Player)

| Key | Action |
| :--- | :--- |
| `Space` or `K` | Play / Pause video |
| `F` | Toggle fullscreen mode |
| `M` | Mute / Unmute audio |
| `←` / `→` | Skip backward / forward 10 seconds |
| `↑` / `↓` | Volume up / down (5% increments) |
| `0` to `9` | Jump to 0% through 90% of duration |

#### Section 2: Supported Media Formats (Direct Play)

The server delivers files directly without transcoding. The manual explains required formats:

+ **Video:** MP4 container with H.264 video and AAC audio, or WebM container with VP9 video and Opus audio.
+ **Documents:** EPUB and PDF files.

#### Section 3: Adding New Files

+ Explains how to drop files into a configured library directory and click "Scan Libraries" in Settings, or use the Admin upload option.
+ Notes standard file naming conventions: `Movie Title (Year).mp4` and `<Series>/Season <NN>/<Series> - S<NN>E<NN> - <Title>.mp4`.

---

## 6. CSS Architecture & Coding Conventions

All client styling is consolidated into `web/static/css/style.css`:

1. **No CSS Frameworks:** No Tailwind, Bootstrap, or utility bloat. Plain, highly optimized CSS.
2. **Semantic Class Naming:** Components use simple, self-explanatory class names:
   + `.navbar`, `.navbar__brand`, `.navbar__link`
   + `.shelf`, `.shelf__title`, `.shelf__grid`, `.shelf__scroll`
   + `.card`, `.card__poster`, `.card__title`, `.card__meta`, `.card__progress`
   + `.btn`, `.btn--primary`, `.btn--secondary`
3. **Box Sizing & Reset:**

   ```css
   *, *::before, *::after {
       box-sizing: border-box;
       margin: 0;
       padding: 0;
   }
   ```

4. **Accessible Focus Rings:**
   Interactive elements (buttons, cards, links, inputs) have a clear focus ring for keyboard and remote control navigation:

   ```css
   :focus-visible {
       outline: 2px solid var(--border-focus);
       outline-offset: 2px;
   }
   ```

---

## 7. Embedded Third-Party Vendor Assets

To guarantee complete offline autonomy and eliminate external CDN dependencies, third-party frontend libraries are vendored into `web/static/vendor/` and compiled directly into the executable via `embed.FS`.

### Vendor Inventory & Responsibilities

| Package | Files | Version Target | License | Purpose |
| :--- | :--- | :--- | :--- | :--- |
| **Plyr** | `plyr/plyr.min.js`<br>`plyr/plyr.css`<br>`plyr/plyr.svg` | v3.7+ | MIT | Lightweight HTML5 video player with custom CSS variable overrides (`--plyr-color-main`), responsive touch scrub bar, and keyboard accessibility. |
| **ePub.js** | `epubjs/epub.min.js` | v0.3+ | BSD-2-Clause | Client-side EPUB unpacker and paginated reader for browser viewing. |
| **JSZip** | `epubjs/jszip.min.js` | v3.10+ | MIT | In-browser zip archive decompression dependency required by ePub.js. |

### Architectural Rules for Vendor Assets

1. **Zero External CDN Links:** Templates must never reference public CDNs (e.g. `cdn.jsdelivr.net`, `cdnjs`, Google Fonts). All resources resolve locally through `/static/vendor/...`.
2. **Deterministic Offline Operation:** The client renders identically on air-gapped homelab local networks without internet connectivity.
3. **Non-Invasive Theme Integration:** Vendor components are adapted using CSS custom property overrides (e.g. `--plyr-color-main: var(--accent-lavender)`), preserving unmodified upstream minified sources.

---

### Related Documentation

+ [Web Client Pages & Template Layouts](pages.md)
+ [Web Client Component Specifications](components.md)
+ [Master System Specifications](../design.md)
+ [User Flows & Client Navigation Diagrams](../diagrams/user-flows.md)
