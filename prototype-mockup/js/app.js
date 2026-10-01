/**
 * FLAN MEDIA SERVER - PROTOTYPE INTERACTIVE CLIENT
 */

(function () {
  "use strict";

  // State
  let state = {
    currentUser: FLAN_MOCK_DATA.users[0], // default to 'mike'
    activeTab: "video", // 'video' | 'books' | 'manage'
    searchQuery: "",
    selectedMediaId: null,
    selectedFileId: null,
    isHelpOpen: false
  };

  // DOM Elements
  const el = {
    topHeader: document.getElementById("top-header"),
    brandTitle: document.getElementById("brand-title"),
    userAvatar: document.getElementById("user-avatar"),
    helpBtn: document.getElementById("help-btn"),
    manageShortcutBtn: document.getElementById("manage-shortcut-btn"),
    appContainer: document.getElementById("app-container"),
    leftSidebar: document.getElementById("left-sidebar"),
    sidebarVideoBtn: document.getElementById("nav-video-btn"),
    sidebarBooksBtn: document.getElementById("nav-books-btn"),
    sidebarManageBtn: document.getElementById("nav-manage-btn"),
    mainContent: document.getElementById("main-content"),
    helpModal: document.getElementById("help-modal"),
    closeHelpBtn: document.getElementById("close-help-btn")
  };

  // Helper: Create SVG Arrow
  function getArrowSVG(direction) {
    if (direction === "left") {
      return `<svg viewBox="0 0 100 12" preserveAspectRatio="none">
        <line x1="100" y1="6" x2="6" y2="6" stroke="#000" stroke-width="3.5" />
        <polyline points="14,1 5,6 14,11" fill="none" stroke="#000" stroke-width="3.5" stroke-linecap="round" stroke-linejoin="round" />
      </svg>`;
    } else {
      return `<svg viewBox="0 0 100 12" preserveAspectRatio="none">
        <line x1="0" y1="6" x2="94" y2="6" stroke="#000" stroke-width="3.5" />
        <polyline points="86,1 95,6 86,11" fill="none" stroke="#000" stroke-width="3.5" stroke-linecap="round" stroke-linejoin="round" />
      </svg>`;
    }
  }

  // Router dispatcher based on window.location.hash
  function handleRoute() {
    const hash = window.location.hash || "#video";
    const parts = hash.split("/");
    const route = parts[0];
    const param = parts[1];

    // If on login, hide standard shell sidebar and top header utility tools
    if (route === "#login") {
      renderLoginScreen();
      return;
    }

    // Ensure user is logged in
    if (!state.currentUser) {
      window.location.hash = "#login";
      return;
    }

    // Show app shell
    el.topHeader.classList.remove("hidden");
    el.leftSidebar.classList.remove("hidden");
    el.appContainer.classList.remove("hidden");

    if (route === "#video") {
      state.activeTab = "video";
      updateSidebarActive();
      renderCatalogView("video");
    } else if (route === "#books") {
      state.activeTab = "books";
      updateSidebarActive();
      renderCatalogView("books");
    } else if (route === "#detail") {
      const type = parts[1];
      const id = parseInt(parts[2], 10);
      renderDetailView(type, id);
    } else if (route === "#watch") {
      const fileId = parseInt(param, 10);
      renderWatchView(fileId);
    } else if (route === "#manage") {
      state.activeTab = "manage";
      updateSidebarActive();
      renderManageView();
    } else {
      window.location.hash = "#video";
    }
  }

  function updateSidebarActive() {
    el.sidebarVideoBtn.classList.remove("active", "focused");
    el.sidebarBooksBtn.classList.remove("active", "focused");
    el.sidebarManageBtn.classList.remove("active");

    if (state.activeTab === "video") {
      el.sidebarVideoBtn.classList.add("active");
    } else if (state.activeTab === "books") {
      el.sidebarBooksBtn.classList.add("active");
    } else if (state.activeTab === "manage") {
      el.sidebarManageBtn.classList.add("active");
    }
  }

  // =========================================================================
  // VIEW 1: SPLIT START & LOGIN SCREEN (Image 1)
  // =========================================================================
  function renderLoginScreen() {
    el.leftSidebar.classList.add("hidden");
    el.mainContent.style.padding = "0";

    el.mainContent.innerHTML = `
      <div class="split-login-container">
        <!-- Left Panel: Lavender with tilted Welcome and arrows -->
        <div class="login-left-panel">
          <div class="welcome-graphic-wrap">
            <div class="welcome-arrow-top">
              ${getArrowSVG("left")}
            </div>
            <h1 class="welcome-banner-text">Welcome</h1>
            <div class="welcome-arrow-bottom">
              ${getArrowSVG("right")}
            </div>
          </div>
        </div>

        <!-- Right Panel: Lilac with User dropdown, Pin input, and Access button -->
        <div class="login-right-panel">
          <form id="login-form" class="login-form-box">
            <div class="form-field-group">
              <label class="form-field-label" for="user-select">User:</label>
              <div class="tactile-select-wrap">
                <select id="user-select">
                  ${FLAN_MOCK_DATA.users
                    .map(
                      (u) =>
                        `<option value="${u.username}" ${
                          state.currentUser && state.currentUser.username === u.username
                            ? "selected"
                            : ""
                        }>${u.username}</option>`
                    )
                    .join("")}
                </select>
                <div class="select-chevron-box">
                  <svg viewBox="0 0 10 8">
                    <polygon points="1,1 9,1 5,7" />
                  </svg>
                </div>
              </div>
            </div>

            <div class="form-field-group">
              <label class="form-field-label" for="pin-input">Pin:</label>
              <input
                id="pin-input"
                class="tactile-pin-input"
                type="password"
                maxlength="6"
                placeholder="••••"
                autofocus
                autocomplete="current-password"
              />
            </div>

            <button type="submit" id="access-btn" class="tactile-access-btn">
              Access
            </button>

            <div id="login-feedback" class="login-feedback"></div>
          </form>
        </div>
      </div>
    `;

    const form = document.getElementById("login-form");
    const userSelect = document.getElementById("user-select");
    const pinInput = document.getElementById("pin-input");
    const feedback = document.getElementById("login-feedback");
    const accessBtn = document.getElementById("access-btn");

    form.addEventListener("submit", (e) => {
      e.preventDefault();
      const selectedUsername = userSelect.value;
      const enteredPin = pinInput.value.trim();

      const user = FLAN_MOCK_DATA.users.find(
        (u) => u.username === selectedUsername
      );

      if (user && user.pin === enteredPin) {
        state.currentUser = user;
        feedback.textContent = "";
        window.location.hash = "#video";
      } else {
        feedback.className = "login-feedback error";
        feedback.textContent = "Invalid PIN. Try '1234' for mike or '0000' for wesley.";
        pinInput.classList.add("shake");
        pinInput.value = "";
        setTimeout(() => pinInput.classList.remove("shake"), 400);
      }
    });
  }

  // =========================================================================
  // VIEW 2 & 3: MEDIA CATALOG WITH SEARCH (Images 2 & 3)
  // =========================================================================
  function renderCatalogView(type) {
    el.mainContent.style.padding = "24px 32px";

    const items =
      type === "video" ? FLAN_MOCK_DATA.videos : FLAN_MOCK_DATA.books;
    const filtered = items.filter((item) =>
      item.title.toLowerCase().includes(state.searchQuery.toLowerCase())
    );

    el.mainContent.innerHTML = `
      <section class="search-section">
        <div class="search-bar-row">
          <div class="search-input-wrap">
            <input
              id="catalog-search"
              class="search-input"
              type="text"
              placeholder="${type === 'video' ? 'Search videos...' : 'Search books...'}"
              value="${state.searchQuery}"
              autocomplete="off"
            />
          </div>
          <button id="catalog-search-btn" class="search-btn" title="Search">
            <svg viewBox="0 0 24 24">
              <circle cx="11" cy="11" r="7" />
              <line x1="16.5" y1="16.5" x2="22" y2="22" />
            </svg>
          </button>
        </div>
        <div class="search-tagline">Let’s see today!</div>
      </section>

      <section class="media-grid">
        ${
          filtered.length === 0
            ? `<div style="grid-column: 1/-1; text-align:center; padding: 40px; font-weight:700; font-size:1.2rem;">No items matching "${state.searchQuery}"</div>`
            : filtered
                .map((item) => {
                  return `
            <div class="media-card" data-id="${item.id}" data-type="${type}">
              <!-- Upper Poster Area (75% height) -->
              <div class="card-poster">
                <span class="card-badge">${item.badge}</span>
                <div class="card-poster-placeholder">
                  <svg viewBox="0 0 24 24">
                    ${
                      type === "video"
                        ? '<polygon points="5,3 19,12 5,21" fill="none" stroke="currentColor" stroke-width="2"/>'
                        : '<path d="M4 19.5A2.5 2.5 0 0 1 6.5 17H20"/> <path d="M6.5 2H20v20H6.5A2.5 2.5 0 0 1 4 19.5v-15A2.5 2.5 0 0 1 6.5 2z"/>'
                    }
                  </svg>
                  <span style="font-weight:700; font-size: 0.9rem; color: #333;">${item.title}</span>
                </div>
              </div>

              <!-- Lower Purple Footer Band (25% height - Image 3) -->
              <div class="card-footer-band">
                <div class="card-title">${item.title}</div>
                <div class="card-subtitle">${
                  type === "video"
                    ? `${item.files.length} ${item.files.length === 1 ? "file" : "episodes"}`
                    : `${item.author || "Book"}`
                }</div>
              </div>
            </div>
          `;
                })
                .join("")
        }
      </section>
    `;

    // Bind Search events
    const searchInput = document.getElementById("catalog-search");
    searchInput.focus();
    // Keep cursor at end of input
    searchInput.setSelectionRange(
      searchInput.value.length,
      searchInput.value.length
    );

    searchInput.addEventListener("input", (e) => {
      state.searchQuery = e.target.value;
      renderCatalogView(type);
    });

    // Bind Card Click -> Details View
    const cards = el.mainContent.querySelectorAll(".media-card");
    cards.forEach((card) => {
      card.addEventListener("click", () => {
        const id = card.getAttribute("data-id");
        window.location.hash = `#detail/${type}/${id}`;
      });
    });
  }

  // =========================================================================
  // VIEW 4: MEDIA CONTAINER DETAIL VIEW
  // =========================================================================
  function renderDetailView(type, id) {
    el.mainContent.style.padding = "24px 32px";
    const items =
      type === "video" ? FLAN_MOCK_DATA.videos : FLAN_MOCK_DATA.books;
    const item = items.find((x) => x.id === id);

    if (!item) {
      window.location.hash = `#${type}`;
      return;
    }

    el.mainContent.innerHTML = `
      <div class="detail-container">
        <div>
          <button id="detail-back-btn" class="tactile-action-btn secondary">
            ← Back to ${type === "video" ? "Videos" : "Books"}
          </button>
        </div>

        <div class="detail-header-card">
          <div class="detail-poster-box">
            <div class="card-poster-placeholder" style="height: 100%;">
              <svg viewBox="0 0 24 24">
                ${
                  type === "video"
                    ? '<polygon points="5,3 19,12 5,21" fill="none" stroke="currentColor" stroke-width="2"/>'
                    : '<path d="M4 19.5A2.5 2.5 0 0 1 6.5 17H20"/><path d="M6.5 2H20v20H6.5A2.5 2.5 0 0 1 4 19.5v-15A2.5 2.5 0 0 1 6.5 2z"/>'
                }
              </svg>
            </div>
          </div>
          <div class="detail-meta-box">
            <h1 class="detail-title">${item.title}</h1>
            ${item.author ? `<div style="font-weight: 600; color: #555;">Author: ${item.author}</div>` : ""}
            <p class="detail-synopsis">${item.overview}</p>
          </div>
        </div>

        <h2 style="font-size: 1.3rem; font-weight: 800; margin-top: 10px;">
          ${type === "video" ? "PLAYABLE EPISODES & CUTS" : "AVAILABLE VOLUMES & EDITIONS"}
        </h2>

        <div class="playable-items-list">
          ${item.files
            .map((file, idx) => {
              const hasProgress = file.progress > 0 && !file.isFinished;
              return `
              <div class="playable-item-row">
                <div class="playable-item-title">
                  <span style="opacity: 0.5;">#${idx + 1}</span>
                  <span>${file.title}</span>
                  <span style="font-size: 0.85rem; font-weight: 600; color: #666;">
                    ${file.duration || file.size}
                  </span>
                </div>
                <div class="playable-item-actions">
                  ${
                    file.isFinished
                      ? '<span style="font-size: 0.85rem; font-weight: 700; color: #2e7d32; border: 1.5px solid #2e7d32; padding: 2px 8px; border-radius: 4px;">✓ Watched</span>'
                      : hasProgress
                      ? `<span style="font-size: 0.85rem; font-weight: 700; color: #6d4ca6; border: 1.5px solid #6d4ca6; padding: 2px 8px; border-radius: 4px;">At ${file.formattedPos || file.progress + "%"}</span>`
                      : ""
                  }
                  <button class="tactile-action-btn play-file-btn" data-file-id="${file.id}">
                    ${type === "video" ? "▶ Play" : "📖 Read"}
                  </button>
                </div>
              </div>
            `;
            })
            .join("")}
        </div>
      </div>
    `;

    document.getElementById("detail-back-btn").addEventListener("click", () => {
      window.location.hash = `#${type}`;
    });

    const playBtns = el.mainContent.querySelectorAll(".play-file-btn");
    playBtns.forEach((btn) => {
      btn.addEventListener("click", () => {
        const fileId = btn.getAttribute("data-file-id");
        window.location.hash = `#watch/${fileId}`;
      });
    });
  }

  // =========================================================================
  // VIEW 5: VIDEO PLAYER PREVIEW (With Auto-Resume prompt)
  // =========================================================================
  function renderWatchView(fileId) {
    // Look up file in all videos
    let targetFile = null;
    let parentVideo = null;
    for (const v of FLAN_MOCK_DATA.videos) {
      const f = v.files.find((x) => x.id === fileId);
      if (f) {
        targetFile = f;
        parentVideo = v;
        break;
      }
    }

    if (!targetFile) {
      window.location.hash = "#video";
      return;
    }

    el.leftSidebar.classList.add("hidden");
    el.mainContent.style.padding = "0";

    const hasResume = targetFile.positionSeconds > 10 && !targetFile.isFinished;

    el.mainContent.innerHTML = `
      <div class="player-container">
        <div class="player-top-bar">
          <button id="player-exit-btn" class="player-back-btn">← Back</button>
          <span style="font-weight: 700;">${parentVideo.title} • ${targetFile.title}</span>
        </div>

        <div class="player-screen-area">
          ${
            hasResume
              ? `
            <div id="resume-banner" class="resume-banner-overlay">
              <h3 style="font-size: 1.4rem; font-weight: 900; letter-spacing: -0.5px;">RESUME PLAYBACK</h3>
              <p style="font-size: 1.1rem; color: #ddd;">
                You were watching at <strong>${targetFile.formattedPos}</strong>.
              </p>
              <div class="resume-buttons-row">
                <button id="resume-accept-btn" class="tactile-action-btn" style="padding: 10px 24px; font-size: 1.1rem;">
                  ▶ Resume from ${targetFile.formattedPos}
                </button>
                <button id="resume-restart-btn" class="tactile-action-btn secondary" style="padding: 10px 20px;">
                  ↺ Start from Beginning
                </button>
              </div>
            </div>
          `
              : ""
          }

          <div style="text-align: center; color: #777;">
            <svg viewBox="0 0 24 24" style="width: 80px; height: 80px; stroke: #555; fill: none; stroke-width: 1.5; margin-bottom: 12px;">
              <polygon points="5,3 19,12 5,21" />
            </svg>
            <div style="font-size: 1.2rem; font-weight: 700;">Direct-Play Web Stream</div>
            <div style="font-size: 0.9rem; margin-top: 4px; color: #555;">Zero-copy kernel sendfile (HTTP 206)</div>
          </div>
        </div>
      </div>
    `;

    document.getElementById("player-exit-btn").addEventListener("click", () => {
      window.location.hash = `#detail/video/${parentVideo.id}`;
    });

    if (hasResume) {
      document.getElementById("resume-accept-btn").addEventListener("click", () => {
        document.getElementById("resume-banner").classList.add("hidden");
      });
      document.getElementById("resume-restart-btn").addEventListener("click", () => {
        document.getElementById("resume-banner").classList.add("hidden");
        targetFile.positionSeconds = 0;
        targetFile.progress = 0;
      });
    }
  }

  // =========================================================================
  // VIEW 6: MANAGE SERVER VIEW
  // =========================================================================
  function renderManageView() {
    el.mainContent.style.padding = "24px 32px";
    const m = FLAN_MOCK_DATA.serverMetrics;

    el.mainContent.innerHTML = `
      <div class="manage-container">
        <h1 style="font-size: 1.8rem; font-weight: 900; letter-spacing: -0.5px;">Manage Server</h1>

        <div class="manage-card">
          <h2>System Performance & Semaphore</h2>
          <div class="metrics-row">
            <div class="metric-box">
              <span style="font-weight: 700; color: #555;">Active RAM Footprint</span>
              <span class="metric-val">${m.memoryUsed} / ${m.memoryLimit}</span>
              <span style="font-size: 0.8rem; color: #666;">Enforced by GOMEMLIMIT=16MiB</span>
            </div>
            <div class="metric-box">
              <span style="font-weight: 700; color: #555;">Video Streams</span>
              <span class="metric-val">${m.activeStreams} / ${m.maxStreams}</span>
              <span style="font-size: 0.8rem; color: #666;">Disk Head Thrashing Guard</span>
            </div>
          </div>
        </div>

        <div class="manage-card">
          <h2>Media Libraries</h2>
          <p style="color: #444; line-height: 1.4;">
            Flan uses fixed paths under <code>./media/video</code> and <code>./media/books</code>.
          </p>
          <div style="display: flex; gap: 12px; margin-top: 8px;">
            <button id="rescan-btn" class="tactile-action-btn" style="padding: 10px 20px;">
              ⟳ Rescan All Media
            </button>
            <button id="mock-upload-btn" class="tactile-action-btn secondary" style="padding: 10px 20px;">
              ⬆ Upload Media Files
            </button>
          </div>
          <div id="rescan-feedback" style="font-weight: 700; color: #2e7d32; min-height: 20px;"></div>
        </div>

        <div class="manage-card">
          <h2>Household Profiles</h2>
          <div style="display: flex; flex-direction: column; gap: 10px;">
            ${FLAN_MOCK_DATA.users
              .map(
                (u) => `
              <div style="display: flex; justify-content: space-between; align-items: center; padding: 10px 14px; background: #fafafa; border: 2px solid #ddd; border-radius: 6px;">
                <span style="font-weight: 700;">${u.username} <span style="font-size: 0.8rem; color: #666;">(${u.role})</span></span>
                <span style="font-family: monospace; font-size: 0.9rem;">PIN: ••••</span>
              </div>
            `
              )
              .join("")}
          </div>
        </div>
      </div>
    `;

    const rescanBtn = document.getElementById("rescan-btn");
    const feedback = document.getElementById("rescan-feedback");
    rescanBtn.addEventListener("click", () => {
      rescanBtn.textContent = "Scanning...";
      feedback.textContent = "";
      setTimeout(() => {
        rescanBtn.textContent = "⟳ Rescan All Media";
        feedback.textContent = "✓ Library scan complete: 5 video containers, 3 book containers indexed.";
      }, 600);
    });
  }

  // =========================================================================
  // GLOBAL LISTENERS & NAVIGATION
  // =========================================================================
  function initListeners() {
    window.addEventListener("hashchange", handleRoute);

    // Sidebar navigation clicks
    el.sidebarVideoBtn.addEventListener("click", () => {
      state.searchQuery = "";
      window.location.hash = "#video";
    });

    el.sidebarBooksBtn.addEventListener("click", () => {
      state.searchQuery = "";
      window.location.hash = "#books";
    });

    el.sidebarManageBtn.addEventListener("click", () => {
      window.location.hash = "#manage";
    });

    el.brandTitle.addEventListener("click", () => {
      window.location.hash = "#video";
    });

    el.manageShortcutBtn.addEventListener("click", () => {
      window.location.hash = "#manage";
    });

    // Avatar click -> Prompt to logout or switch account
    el.userAvatar.addEventListener("click", () => {
      const confirmLogout = confirm(
        `Logged in as "${state.currentUser.username}". Return to Start screen?`
      );
      if (confirmLogout) {
        window.location.hash = "#login";
      }
    });

    // Help Dialog
    el.helpBtn.addEventListener("click", () => {
      el.helpModal.classList.remove("hidden");
    });

    el.closeHelpBtn.addEventListener("click", () => {
      el.helpModal.classList.add("hidden");
    });

    el.helpModal.addEventListener("click", (e) => {
      if (e.target === el.helpModal) {
        el.helpModal.classList.add("hidden");
      }
    });
  }

  // Start app
  initListeners();
  handleRoute();
})();
