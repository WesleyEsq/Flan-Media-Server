# Web Client Design System & Foundations

Visual tokens, typography scale, global navigation bar, and built-in e-manual for Flan Media Server's web client.

---

## 1. Core Visual Principles

* **Zero Animation Bloat:** No bouncy transforms, hover scaling, or heavy CSS gradients. Surfaces are flat and matte, ensuring instant browser painting on low-power Smart TVs and older mobile devices.
* **No Backdrop Blur:** CSS `backdrop-filter: blur(...)` is strictly avoided to prevent GPU and CPU rendering lag on constrained hardware.
* **High Contrast (WCAG AAA):** Text and surfaces exceed 14:1 contrast ratios for effortless reading across a living room on low-brightness TV displays.
* **System Font Stack:** Uses the device's native system font stack with zero web fonts, eliminating network latency and layout shifts (`font-family: system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif`).

---

## 2. Color Palette & Design Tokens

All styles are consolidated in `web/static/css/style.css`:

```css
:root {
    /* Surfaces */
    --bg-base:        #12131a; /* Deep matte slate background */
    --bg-surface:     #1a1b24; /* Flat card, shelf, and nav background */
    --bg-elevated:    #242633; /* Modals, popups, and hover states */
    --bg-active:      #2d2640; /* Selected item with lavender tint */

    /* Borders */
    --border-subtle:  #2c2e3e; /* Visible card and divider borders */
    --border-focus:   #bb9af7; /* High-visibility keyboard/remote focus ring */

    /* Typography */
    --text-primary:   #f0f2fc; /* High-contrast off-white (>14:1 ratio) */
    --text-secondary: #9ba1b8; /* Muted metadata (year, runtime, author) */
    --text-muted:     #6b728d; /* Inactive labels and placeholders */

    /* Accents & Functional */
    --accent-lavender:#bb9af7; /* Primary buttons, active tabs, progress bars */
    --accent-hover:   #cbb0f9; /* Button hover state */
    --color-danger:   #f7768e; /* Lockout warnings, delete buttons */
    --color-success:  #9ece6a; /* Completed checkmarks, health badges */
}
```

---

## 3. Typography Scale

| Level | Size | Weight | Line Height | Usage |
| :--- | :--- | :--- | :--- | :--- |
| **Display / Header** | 1.5rem (24px) | Bold (700) | 1.2 | Page titles, TV series detail heading |
| **Shelf Title** | 1.25rem (20px)| 600 | 1.3 | Dashboard section headings ("Continue Watching") |
| **Card Title** | 0.9375rem (15px)| 500 | 1.4 | Movie, episode, and book titles (truncated 2 lines) |
| **Metadata / Subtitle**| 0.8125rem (13px)| 400 | 1.4 | Year, duration, episode numbers, authors (`--text-secondary`) |
| **Badges & Pills** | 0.75rem (12px)| 600 | 1.0 | Format pills (`4K`, `EPUB`, `PDF`), genre filter tags |

---

## 4. Global Navigation Bar

Fixed at the top of the viewport (`sticky; top: 0; z-index: 50`) on all pages except fullscreen playback:

```text
+-------------------------------------------------------------------------+
| [Flan Logo]   Home   Videos   Books                 [? Help] [⚙] [Avatar] |
+-------------------------------------------------------------------------+
```

* **Brand (Left):** Text logo `Flan` linking to `/`.
* **Nav Links (Center/Left):** `Home` (`/`), `Videos` (`/videos`), `Books` (`/books`). Active page indicated by a solid lavender bottom border (`border-bottom: 2px solid var(--accent-lavender)`).
* **Utility Controls (Right):**
  * `? Help`: Opens the built-in HTML `<dialog>` e-manual.
  * `⚙ Settings`: Opens administrative settings (visible to admin).
  * `Profile Avatar`: Chosen icon on user's accent color tile. Clicking opens dropdown to switch profiles or log out.
* **Mobile Handling:** Links wrap cleanly or scroll horizontally without a heavy JavaScript hamburger menu.

---

## 5. Built-in E-Manual (`<dialog id="manual-dialog">`)

Accessible on any page via `? Help`, built with native HTML `<dialog>` (zero JS libraries; traps focus, closes on `Escape`):

### Keyboard Shortcuts (Video Player)
| Key | Action |
| :--- | :--- |
| `Space` / `K` | Play / Pause |
| `F` | Toggle fullscreen |
| `M` | Mute / Unmute audio |
| `←` / `→` | Skip backward / forward 10 seconds |
| `↑` / `↓` | Volume up / down (5% increments) |
| `0`–`9` | Jump to 0% through 90% duration |

### Direct Play Codec Guidelines
* **Native Web Video:** MP4 (H.264/AAC), WebM (VP9/Opus, AV1), and web-safe MKV.
* **Unsupported Formats:** Non-web codecs (e.g. DivX in `.avi`, DTS/AC3 in `.mkv`) display a fallback prompt allowing direct file download for VLC playback.
* **Reading Material:** EPUB (ePub.js reader) and PDF (native browser viewer).

---

## 6. Embedded Third-Party Vendor Assets

All dependencies are vendored into `web/static/vendor/` and embedded into the binary via `embed.FS`:

| Library | Directory | Version Target | License | Purpose |
| :--- | :--- | :--- | :--- | :--- |
| **Plyr** | `web/static/vendor/plyr/` | v3.7+ | MIT | Lightweight HTML5 video player styled with lavender variables (`--plyr-color-main: #bb9af7`). |
| **ePub.js** | `web/static/vendor/epubjs/` | v0.3+ | BSD-2-Clause | In-browser EPUB unpacker and paginated reader. |
| **JSZip** | `web/static/vendor/epubjs/` | v3.10+ | MIT | Client-side ZIP decompression dependency for ePub.js. |

---

## 7. Related Documentation

* [Web Client Components](components.md)
* [Page Templates & Wireframes](pages.md)
* [Master System Architecture](../design.md)
