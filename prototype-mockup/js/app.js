/**
 * FLAN MEDIA SERVER - PROTOTYPE INTERACTIVE CLIENT
 */

(function () {
  "use strict";

  // State
  let state = {
    currentUser: FLAN_MOCK_DATA.users[0], // default to 'mike'
    activeTab: "video", // 'video' | 'books' | 'manage' | 'manual'
    previousTab: "video",
    searchQuery: "",
    selectedMediaId: null,
    selectedFileId: null
  };

  // DOM Elements
  const el = {
    topHeader: document.getElementById("top-header"),
    brandTitle: document.getElementById("brand-title"),
    userAvatar: document.getElementById("user-avatar"),
    helpBtn: document.getElementById("help-btn"),
    appContainer: document.getElementById("app-container"),
    leftSidebar: document.getElementById("left-sidebar"),
    sidebarVideoBtn: document.getElementById("nav-video-btn"),
    sidebarBooksBtn: document.getElementById("nav-books-btn"),
    sidebarManageBtn: document.getElementById("nav-manage-btn"),
    mainContent: document.getElementById("main-content"),

    // Unified My Profile Modal
    profileModal: document.getElementById("profile-modal"),
    closeProfileBtn: document.getElementById("close-profile-btn"),
    profileForm: document.getElementById("profile-form"),
    profileUsernameInput: document.getElementById("profile-username-input"),
    profileAvatarsGrid: document.getElementById("profile-avatars-grid"),
    profileAvatarFile: document.getElementById("profile-avatar-file"),
    profilePinInput: document.getElementById("profile-pin-input"),
    profileFeedback: document.getElementById("profile-feedback"),
    profileLogoutBtn: document.getElementById("profile-logout-btn"),

    // Add User Modal (Manage Server)
    addUserModal: document.getElementById("add-user-modal"),
    closeAddUserBtn: document.getElementById("close-add-user-btn"),
    addUserForm: document.getElementById("add-user-form"),
    newUserName: document.getElementById("new-user-name"),
    newUserPin: document.getElementById("new-user-pin"),
    newUserRole: document.getElementById("new-user-role"),
    newUserAvatarsGrid: document.getElementById("new-user-avatars-grid"),
    addUserFeedback: document.getElementById("add-user-feedback")
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

  // Render User Avatar Badge (Top Header)
  function renderUserAvatarBadge(user) {
    if (!user || !el.userAvatar) return;
    if (user.avatarType === "custom" && user.customAvatarData) {
      el.userAvatar.innerHTML = `<img src="${user.customAvatarData}" alt="${user.username}" style="width:100%; height:100%; object-fit:cover;" />`;
    } else {
      const key = user.avatarKey || "mascot";
      const preset = FLAN_MOCK_DATA.presetAvatars[key] || FLAN_MOCK_DATA.presetAvatars.mascot;
      el.userAvatar.innerHTML = preset.svg;
    }
  }

  // =========================================================================
  // UNIFIED MY PROFILE MODAL MANAGER
  // =========================================================================
  let profileTargetUser = null;
  let tempAvatarType = null;
  let tempAvatarKey = null;
  let tempCustomAvatarData = null;

  function openProfileModal(user) {
    profileTargetUser = user || state.currentUser;
    tempAvatarType = profileTargetUser.avatarType || "preset";
    tempAvatarKey = profileTargetUser.avatarKey || "mascot";
    tempCustomAvatarData = profileTargetUser.customAvatarData || null;

    el.profileUsernameInput.value = profileTargetUser.username;
    el.profilePinInput.value = "";
    el.profileFeedback.textContent = "";
    el.profileFeedback.className = "avatar-feedback";
    el.profileAvatarFile.value = "";

    renderProfileAvatarsGrid();
    el.profileModal.classList.remove("hidden");
  }

  function renderProfileAvatarsGrid() {
    el.profileAvatarsGrid.innerHTML = Object.entries(FLAN_MOCK_DATA.presetAvatars)
      .map(([key, item]) => {
        const isActive = tempAvatarType === "preset" && tempAvatarKey === key && !tempCustomAvatarData;
        return `
          <button type="button" class="preset-avatar-btn ${isActive ? "active" : ""}" data-key="${key}" title="${item.name}">
            ${item.svg}
            <span class="preset-avatar-label">${item.name}</span>
          </button>
        `;
      })
      .join("");

    const btns = el.profileAvatarsGrid.querySelectorAll(".preset-avatar-btn");
    btns.forEach((btn) => {
      btn.addEventListener("click", () => {
        const key = btn.getAttribute("data-key");
        tempAvatarType = "preset";
        tempAvatarKey = key;
        tempCustomAvatarData = null;

        btns.forEach((b) => b.classList.remove("active"));
        btn.classList.add("active");

        el.profileFeedback.className = "avatar-feedback success";
        el.profileFeedback.textContent = `Selected "${FLAN_MOCK_DATA.presetAvatars[key].name}"! Click 'Save Changes' to apply.`;
      });
    });
  }

  function closeProfileModal() {
    el.profileModal.classList.add("hidden");
    profileTargetUser = null;
  }

  // =========================================================================
  // ADD USER MODAL MANAGER (Manage Server)
  // =========================================================================
  let newUserSelectedAvatarKey = "mascot";

  function openAddUserModal() {
    el.newUserName.value = "";
    el.newUserPin.value = "";
    el.newUserRole.value = "user";
    newUserSelectedAvatarKey = "mascot";
    el.addUserFeedback.textContent = "";
    el.addUserFeedback.className = "avatar-feedback";

    // Render avatar choices
    el.newUserAvatarsGrid.innerHTML = Object.entries(FLAN_MOCK_DATA.presetAvatars)
      .map(([key, item]) => {
        const isActive = key === newUserSelectedAvatarKey;
        return `
          <button type="button" class="preset-avatar-btn ${isActive ? "active" : ""}" data-key="${key}" title="${item.name}">
            ${item.svg}
            <span class="preset-avatar-label">${item.name}</span>
          </button>
        `;
      })
      .join("");

    const btns = el.newUserAvatarsGrid.querySelectorAll(".preset-avatar-btn");
    btns.forEach((btn) => {
      btn.addEventListener("click", () => {
        newUserSelectedAvatarKey = btn.getAttribute("data-key");
        btns.forEach((b) => b.classList.remove("active"));
        btn.classList.add("active");
      });
    });

    el.addUserModal.classList.remove("hidden");
  }

  function closeAddUserModal() {
    el.addUserModal.classList.add("hidden");
  }

  // =========================================================================
  // ROUTER DISPATCHER
  // =========================================================================
  function handleRoute() {
    const hash = window.location.hash || "#video";
    const parts = hash.split("/");
    const route = parts[0];
    const param = parts[1];

    if (route === "#login") {
      renderLoginScreen();
      return;
    }

    if (!state.currentUser) {
      window.location.hash = "#login";
      return;
    }

    // Show app shell
    el.topHeader.classList.remove("hidden");
    el.leftSidebar.classList.remove("hidden");
    el.appContainer.classList.remove("hidden");
    renderUserAvatarBadge(state.currentUser);

    if (route !== "#manual") {
      el.mainContent.classList.remove("manual-mode");
    }

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
      el.mainContent.classList.remove("manual-mode");
      state.activeTab = "manage";
      updateSidebarActive();
      renderManageView();
    } else if (route === "#manual") {
      if (state.activeTab !== "manual") {
        state.previousTab = state.activeTab;
      }
      state.activeTab = "manual";
      updateSidebarActive();
      renderManualView();
    } else {
      el.mainContent.classList.remove("manual-mode");
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
  // VIEW 1: SPLIT START & LOGIN SCREEN
  // =========================================================================
  function renderLoginScreen() {
    el.leftSidebar.classList.add("hidden");
    el.mainContent.style.padding = "0";

    el.mainContent.innerHTML = `
      <div class="split-login-container">
        <!-- Left Panel: Welcome Graphic -->
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

        <!-- Right Panel: Access Form -->
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

    form.addEventListener("submit", (e) => {
      e.preventDefault();
      const selectedUsername = userSelect.value;
      const enteredPin = pinInput.value.trim();

      const user = FLAN_MOCK_DATA.users.find(
        (u) => u.username === selectedUsername
      );

      if (user && user.pin === enteredPin) {
        state.currentUser = user;
        renderUserAvatarBadge(user);
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
  // VIEW 2 & 3: MEDIA CATALOG WITH DYNAMIC HEADER
  // =========================================================================
  function renderCatalogView(type) {
    el.mainContent.style.padding = "24px 32px";

    const items =
      type === "video" ? FLAN_MOCK_DATA.videos : FLAN_MOCK_DATA.books;
    const filtered = items.filter((item) =>
      item.title.toLowerCase().includes(state.searchQuery.toLowerCase())
    );

    // Compact Tactile Status & Filter Strip
    const isSearching = state.searchQuery.trim().length > 0;
    let filterStripHtml = "";
    if (isSearching) {
      filterStripHtml = `
        <div class="catalog-status-strip filtered">
          <span class="status-indicator-tag">FILTER</span>
          <span class="status-query-text">"${state.searchQuery}"</span>
          <span class="status-match-count">(${filtered.length} match${filtered.length === 1 ? "" : "es"})</span>
          <button id="clear-search-btn" class="tactile-clear-btn" title="Clear Search">
            ✕ Clear
          </button>
        </div>
      `;
    } else {
      const typeLabel =
        type === "video" ? "ALL VIDEOS & SERIES" : "BOOKS & PUBLICATIONS";
      filterStripHtml = `
        <div class="catalog-status-strip">
          <span class="status-section-name">${typeLabel}</span>
          <span class="status-divider">•</span>
          <span class="status-total-count">${items.length} TITLES</span>
        </div>
      `;
    }

    el.mainContent.innerHTML = `
      <section class="search-section">
        <div class="search-bar-row">
          <div class="search-input-wrap">
            <input
              id="catalog-search"
              class="search-input"
              type="text"
              placeholder="${type === 'video' ? 'Search videos & movies...' : 'Search books & documents...'}"
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

        <!-- Compact Tactile Status & Filter Strip -->
        ${filterStripHtml}
      </section>

      <section class="media-grid">
        ${
          filtered.length === 0
            ? `<div style="grid-column: 1/-1; text-align:center; padding: 40px; font-weight:700; font-size:1.1rem; color: #555;">No media matching "${state.searchQuery}"</div>`
            : filtered
                .map((item) => {
                  return `
            <div class="media-card" data-id="${item.id}" data-type="${type}">
              <!-- Upper Poster Area -->
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

              <!-- Lower Purple Footer Band -->
              <div class="card-footer-band">
                <div class="card-title" title="${item.title}">
                  <span class="title-text">${item.title}</span>
                </div>
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
    searchInput.setSelectionRange(
      searchInput.value.length,
      searchInput.value.length
    );

    searchInput.addEventListener("input", (e) => {
      state.searchQuery = e.target.value;
      renderCatalogView(type);
    });

    const clearSearchBtn = document.getElementById("clear-search-btn");
    if (clearSearchBtn) {
      clearSearchBtn.addEventListener("click", () => {
        state.searchQuery = "";
        renderCatalogView(type);
      });
    }

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
                  ${
                    type === "video"
                      ? `
                      <button class="tactile-action-btn play-file-btn" data-file-id="${file.id}">
                        ▶ Play
                      </button>
                      <button class="tactile-action-btn vlc download-vlc-btn" data-url="${file.downloadUrl}" data-title="${file.title}" title="Direct download or stream in VLC">
                        ⬇ VLC / Download
                      </button>
                    `
                      : `
                      ${
                        file.format === "pdf"
                          ? `<button class="tactile-action-btn read-pdf-btn" data-url="${file.readUrl}" data-title="${file.title}">📖 Open PDF</button>`
                          : ""
                      }
                      <button class="tactile-action-btn download download-book-btn" data-url="${file.downloadUrl}" data-title="${file.title}">
                        ⬇ Download (${file.format.toUpperCase()})
                      </button>
                    `
                  }
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

    const vlcBtns = el.mainContent.querySelectorAll(".download-vlc-btn");
    vlcBtns.forEach((btn) => {
      btn.addEventListener("click", () => {
        const url = btn.getAttribute("data-url");
        const title = btn.getAttribute("data-title");
        const streamEndpoint = `http://${window.location.hostname || "localhost"}:4907${url}`;
        alert(
          `VLC & DIRECT DOWNLOAD LINK:\n\nTitle: ${title}\nURL: ${streamEndpoint}\n\n• For VLC: Open VLC → Media → Open Network Stream → Paste URL.\n• For Download: Direct file streaming bypasses browser audio codec limits!`
        );
      });
    });

    const bookDlBtns = el.mainContent.querySelectorAll(".download-book-btn");
    bookDlBtns.forEach((btn) => {
      btn.addEventListener("click", () => {
        const title = btn.getAttribute("data-title");
        alert(`Downloading "${title}"...\nSaved to your device for reading in Apple Books, Moon+ Reader, or Kindle!`);
      });
    });

    const pdfBtns = el.mainContent.querySelectorAll(".read-pdf-btn");
    pdfBtns.forEach((btn) => {
      btn.addEventListener("click", () => {
        const title = btn.getAttribute("data-title");
        alert(`Opening "${title}" in native browser PDF tab!`);
      });
    });
  }

  // =========================================================================
  // VIEW 5: VIDEO PLAYER PREVIEW
  // =========================================================================
  function renderWatchView(fileId) {
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

        <!-- Audio Codec Fallback Bar -->
        <div class="player-codec-fallback-bar">
          <span>Audio silent or video stuttering? (Browser lacks AC3 / DTS / HEVC support)</span>
          <button id="player-vlc-btn">Open in External VLC Player / Download ↗</button>
        </div>
      </div>
    `;

    document.getElementById("player-exit-btn").addEventListener("click", () => {
      window.location.hash = `#detail/video/${parentVideo.id}`;
    });

    document.getElementById("player-vlc-btn").addEventListener("click", () => {
      alert(
        `External Stream URL for VLC:\nhttp://${window.location.hostname || "localhost"}:4907${targetFile.downloadUrl}\n\nPaste into VLC (Media > Open Network Stream) or download directly for 100% audio compatibility!`
      );
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
            Flan uses fixed paths under <code>./media/video</code> and <code>./media/books</code> populated directly via Samba, SCP, or external USB drive.
          </p>
          <div style="display: flex; gap: 12px; margin-top: 8px;">
            <button id="rescan-btn" class="tactile-action-btn" style="padding: 10px 20px;">
              ⟳ Rescan All Media
            </button>
          </div>
          <div id="rescan-feedback" style="font-weight: 700; color: #2e7d32; min-height: 20px;"></div>
        </div>

        <div class="manage-card">
          <div style="display: flex; justify-content: space-between; align-items: center; border-bottom: 2px solid #eee; padding-bottom: 8px;">
            <h2 style="border: none; padding: 0; margin: 0;">Household Profiles</h2>
            <button id="add-user-btn" class="tactile-action-btn" style="padding: 6px 14px; font-size: 0.85rem;">
              + Add New User
            </button>
          </div>

          <div style="display: flex; flex-direction: column; gap: 12px; margin-top: 8px;">
            ${FLAN_MOCK_DATA.users
              .map((u) => {
                const avatarPreview =
                  u.avatarType === "custom" && u.customAvatarData
                    ? `<img src="${u.customAvatarData}" style="width:36px; height:36px; border-radius:4px; border:2px solid #3ea6ff; object-fit:cover;" />`
                    : `<div style="width:36px; height:36px; border-radius:4px; border:2px solid #3ea6ff; overflow:hidden;">${
                        (FLAN_MOCK_DATA.presetAvatars[u.avatarKey] || FLAN_MOCK_DATA.presetAvatars.mascot).svg
                      }</div>`;

                return `
              <div style="display: flex; justify-content: space-between; align-items: center; padding: 10px 14px; background: #fafafa; border: 2px solid #000; border-radius: 6px;">
                <div style="display: flex; align-items: center; gap: 12px;">
                  ${avatarPreview}
                  <div>
                    <div style="font-weight: 700; font-size: 1.05rem;">
                      ${u.username} <span style="font-size: 0.8rem; color: #666;">(${u.role})</span>
                    </div>
                    <span style="font-family: monospace; font-size: 0.85rem; color: #555;">PIN: ••••</span>
                  </div>
                </div>
                <div style="display: flex; gap: 8px;">
                  <button class="tactile-action-btn secondary edit-user-btn" data-username="${u.username}" style="padding: 6px 12px; font-size: 0.85rem;">
                    ✏️ Edit Profile
                  </button>
                </div>
              </div>
            `;
              })
              .join("")}
          </div>
        </div>
      </div>
    `;

    // Rescan button
    const rescanBtn = document.getElementById("rescan-btn");
    const feedback = document.getElementById("rescan-feedback");
    rescanBtn.addEventListener("click", () => {
      rescanBtn.textContent = "Scanning...";
      feedback.textContent = "";
      setTimeout(() => {
        rescanBtn.textContent = "⟳ Rescan All Media";
        feedback.textContent = "✓ Library scan complete: 5 video containers, 3 book containers indexed.";
      }, 500);
    });

    // Add user button
    document.getElementById("add-user-btn").addEventListener("click", openAddUserModal);

    // Edit user buttons
    const editBtns = el.mainContent.querySelectorAll(".edit-user-btn");
    editBtns.forEach((btn) => {
      btn.addEventListener("click", () => {
        const username = btn.getAttribute("data-username");
        const user = FLAN_MOCK_DATA.users.find((u) => u.username === username);
        if (user) {
          openProfileModal(user);
        }
      });
    });
  }

  // =========================================================================
  // VIEW 5: DEDICATED SOFTWARE MANUAL & ABOUT VIEW (#view-manual)
  // =========================================================================
  function renderManualView() {
    el.mainContent.classList.add("manual-mode");
    el.mainContent.style.padding = "0";

    el.mainContent.innerHTML = `
      <div class="manual-view-container">
        <!-- Top Navigation Bar (Sticky on Mobile) -->
        <header class="manual-top-bar">
          <button id="manual-back-btn" class="manual-back-btn" title="Back to Library">
            <span class="back-text-desktop">← Back to Library</span>
            <span class="back-text-mobile">← Back</span>
          </button>
          <div class="manual-title-cluster">
            <span class="manual-main-title">Software Handbook</span>
            <span class="manual-version-pill">v0.2.0</span>
          </div>
          <div class="manual-offline-tag">
            100% Offline Reference
          </div>
        </header>

        <!-- Two-Column Layout -->
        <div class="manual-layout">
          <!-- Left Table of Contents (Sticky on Mobile) -->
          <nav class="manual-toc-sidebar">
            <div class="toc-heading">Table of Contents</div>
            <a class="toc-nav-link active" href="#sec-storage">1. Storage & SMB</a>
            <a class="toc-nav-link" href="#sec-playback">2. Playback & VLC</a>
            <a class="toc-nav-link" href="#sec-profiles">3. Accounts & PINs</a>
            <a class="toc-nav-link" href="#sec-cli">4. CLI Admin</a>
            <a class="toc-nav-link" href="#sec-about">5. About & Specs</a>
          </nav>

          <!-- Right Reading Content Pane -->
          <div class="manual-content-pane" id="manual-content-pane">
            <!-- Section 1 -->
            <section id="sec-storage" class="manual-section-card">
              <div class="manual-section-header">
                <span class="manual-section-num">1</span>
                <h2 class="manual-section-title">Media Storage & Placement</h2>
              </div>
              <p class="manual-body-text">
                Flan eliminates complex in-app file uploaders in favor of standard homelab storage management. 
                Place media files directly into the server directory using SMB shares, NFS mounts, SSH/rsync, or an external USB hard drive:
              </p>
              <div class="manual-code-box">
                <div class="manual-code-header">
                  <span class="manual-code-lang">MEDIA FILE PLACEMENT</span>
                  <button class="manual-copy-btn" data-copy="./media/video/
./media/books/">Copy</button>
                </div>
                <pre class="manual-code-pre"><code>./media/video/&lt;Show or Movie Title&gt;/ep01.mp4
./media/video/&lt;Show or Movie Title&gt;/poster.jpg
./media/books/&lt;Book Title&gt;.epub
./media/books/&lt;Document Title&gt;.pdf</code></pre>
              </div>
              <p class="manual-body-text">
                <strong>Sub-second Rescanning:</strong> Whenever you add or organize files, go to <strong>Manage Server</strong> and click 
                <code>[ ⟳ Rescan All Media ]</code>. Flan performs a synchronous file walk and updates its local SQLite database in under 200ms.
              </p>
            </section>

            <!-- Section 2 -->
            <section id="sec-playback" class="manual-section-card">
              <div class="manual-section-header">
                <span class="manual-section-num">2</span>
                <h2 class="manual-section-title">Direct Video Playback & VLC Fallback</h2>
              </div>
              <p class="manual-body-text">
                Flan is built for low-power devices like Raspberry Pi single-board computers (15–20 MB RAM budget). 
                To ensure maximum battery life and zero CPU strain, <strong>Flan never performs server-side video transcoding</strong>.
              </p>
              <ul class="manual-body-text" style="padding-left: 20px; display: flex; flex-direction: column; gap: 8px;">
                <li><strong>Native Browser Play:</strong> Videos formatted with H.264 / AAC or WebM play instantly in any web browser using native zero-copy HTTP 206 range requests.</li>
                <li><strong>Codec Fallback (AC3, EAC3, DTS, 10-bit HEVC):</strong> If a video lacks audio in your browser due to proprietary codec licensing, simply click the purple <strong><code>[ ⬇ VLC / Download ]</code></strong> button. This streams the raw container directly into external media players (VLC, MPV, IINA) or downloads it locally.</li>
                <li><strong>E-Books:</strong> PDFs open directly in a clean browser viewing tab, while EPUBs download with one click to your favorite e-reader application.</li>
              </ul>
            </section>

            <!-- Section 3 -->
            <section id="sec-profiles" class="manual-section-card">
              <div class="manual-section-header">
                <span class="manual-section-num">3</span>
                <h2 class="manual-section-title">Household Profiles & Whimsical Avatars</h2>
              </div>
              <p class="manual-body-text">
                Flan supports independent household profiles so family members maintain separate watch histories and progress bars:
              </p>
              <ul class="manual-body-text" style="padding-left: 20px; display: flex; flex-direction: column; gap: 8px;">
                <li><strong>Fast Switching:</strong> Click your avatar in the top-right header to change your display name, switch companion avatars, or log out.</li>
                <li><strong>Tactile Avatars:</strong> Choose from 6 bundled high-contrast SVG companion presets (Mascot, Flan, Cat, Ghost, Robot, Star) or upload a custom image (JPEG, PNG, WebP up to 2MB).</li>
                <li><strong>PIN Security:</strong> Accounts are safeguarded by 4-to-6 digit numeric PINs salted and hashed via bcrypt.</li>
              </ul>
            </section>

            <!-- Section 4 -->
            <section id="sec-cli" class="manual-section-card">
              <div class="manual-section-header">
                <span class="manual-section-num">4</span>
                <h2 class="manual-section-title">Server CLI Administration & Failsafe Reset</h2>
              </div>
              <p class="manual-body-text">
                Because Flan runs strictly offline with zero external cloud dependencies, administrative recovery is performed directly on the host machine:
              </p>
              <div class="manual-code-box">
                <div class="manual-code-header">
                  <span class="manual-code-lang">CLI ADMIN COMMANDS</span>
                  <button class="manual-copy-btn" data-copy="./flan --reset-admin">Copy</button>
                </div>
                <pre class="manual-code-pre"><code># Reset administrator PIN to default '0000'
./flan --reset-admin

# Run on custom port with bounded memory
PORT=8080 GOMEMLIMIT=16MiB ./flan</code></pre>
              </div>
              <p class="manual-body-text">
                <strong>Storage Portability:</strong> All database state is preserved in <code>./data/flan.db</code> with SQLite WAL mode enabled. To migrate or back up your server, simply copy the <code>data/</code> folder.
              </p>
            </section>

            <!-- Section 5 -->
            <section id="sec-about" class="manual-section-card">
              <div class="manual-section-header">
                <span class="manual-section-num">5</span>
                <h2 class="manual-section-title">About Flan & System Architecture</h2>
              </div>
              <table class="manual-specs-table">
                <tr>
                  <td class="spec-label">Project</td>
                  <td class="spec-val"><strong>Flan Media Server</strong></td>
                </tr>
                <tr>
                  <td class="spec-label">Version</td>
                  <td class="spec-val"><strong>v0.2.0-prototype</strong></td>
                </tr>
                <tr>
                  <td class="spec-label">Author</td>
                  <td class="spec-val">Wesley Esquivel (<a href="https://github.com/WesleyEsq" target="_blank" style="color: #6d52a8; font-weight: 800;">MechanicalSpeak</a>)</td>
                </tr>
                <tr>
                  <td class="spec-label">Memory Ceiling</td>
                  <td class="spec-val">15–20 MB Resident RAM (<code>GOMEMLIMIT=16MiB</code>)</td>
                </tr>
                <tr>
                  <td class="spec-label">Database</td>
                  <td class="spec-val">Pure-Go SQLite WAL & Wear Leveling</td>
                </tr>
                <tr>
                  <td class="spec-label">License</td>
                  <td class="spec-val">Apache 2.0 Open Source</td>
                </tr>
              </table>
              <div style="margin-top: 12px; padding: 12px 14px; background: #faf7fd; border: 1.5px solid #000; border-radius: 4px; font-size: 0.88rem; font-style: italic; color: #444;">
                "A simple server for people that think they want a media server, but in reality just want to host movies and books for themselves and their kids."
              </div>
            </section>
          </div>
        </div>
      </div>
    `;

    // Bind Back button
    document.getElementById("manual-back-btn").addEventListener("click", () => {
      el.mainContent.classList.remove("manual-mode");
      const target = state.previousTab || "video";
      window.location.hash = "#" + target;
    });

    // Bind TOC scroll links
    const tocLinks = document.querySelectorAll(".toc-nav-link");
    tocLinks.forEach((link) => {
      link.addEventListener("click", (e) => {
        e.preventDefault();
        tocLinks.forEach((l) => l.classList.remove("active"));
        link.classList.add("active");
        const targetId = link.getAttribute("href").substring(1);
        const targetEl = document.getElementById(targetId);
        if (targetEl) {
          targetEl.scrollIntoView({ behavior: "smooth", block: "start" });
        }
      });
    });

    // Bind Copy buttons
    document.querySelectorAll(".manual-copy-btn").forEach((btn) => {
      btn.addEventListener("click", () => {
        const text = btn.getAttribute("data-copy");
        navigator.clipboard.writeText(text).then(() => {
          const original = btn.textContent;
          btn.textContent = "✓ Copied";
          setTimeout(() => {
            btn.textContent = original;
          }, 1500);
        });
      });
    });
  }

  // =========================================================================
  // GLOBAL LISTENERS & INITIALIZATION
  // =========================================================================
  function initListeners() {
    window.addEventListener("hashchange", handleRoute);

    // Title Horizontal Scrolling on Hover (Marquee)
    document.addEventListener("mouseover", (e) => {
      const card = e.target.closest(".media-card");
      if (!card) return;
      const titleEl = card.querySelector(".card-title");
      const spanEl = titleEl ? titleEl.querySelector(".title-text") : null;
      if (!titleEl || !spanEl) return;

      if (spanEl.scrollWidth > titleEl.clientWidth) {
        const overflow = spanEl.scrollWidth - titleEl.clientWidth + 8;
        titleEl.classList.add("is-overflowing");
        titleEl.style.setProperty("--scroll-offset", `-${overflow}px`);
        const duration = Math.max(3, overflow / 22);
        titleEl.style.setProperty("--marquee-duration", `${duration}s`);
      } else {
        titleEl.classList.remove("is-overflowing");
      }
    });

    document.addEventListener("mouseout", (e) => {
      const card = e.target.closest(".media-card");
      if (!card) return;
      if (!card.contains(e.relatedTarget)) {
        const spanEl = card.querySelector(".card-title .title-text");
        if (spanEl) {
          spanEl.style.transform = "translateX(0)";
        }
      }
    });

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

    // Header avatar click -> Opens Unified My Profile Modal
    el.userAvatar.addEventListener("click", () => {
      openProfileModal(state.currentUser);
    });

    // Profile Modal Listeners
    el.closeProfileBtn.addEventListener("click", closeProfileModal);
    el.profileModal.addEventListener("click", (e) => {
      if (e.target === el.profileModal) {
        closeProfileModal();
      }
    });

    // Profile photo upload
    el.profileAvatarFile.addEventListener("change", (e) => {
      const file = e.target.files[0];
      if (!file) return;

      if (!file.type.match(/^image\/(png|jpeg|jpg|webp)$/)) {
        el.profileFeedback.className = "avatar-feedback error";
        el.profileFeedback.textContent = "Error: Please upload a PNG, JPEG, or WebP photo.";
        return;
      }

      if (file.size > 2 * 1024 * 1024) {
        el.profileFeedback.className = "avatar-feedback error";
        el.profileFeedback.textContent = "Error: File size exceeds 2 MB limit.";
        return;
      }

      const reader = new FileReader();
      reader.onload = (event) => {
        tempAvatarType = "custom";
        tempCustomAvatarData = event.target.result;
        el.profileFeedback.className = "avatar-feedback success";
        el.profileFeedback.textContent = "✓ Photo loaded! Click 'Save Changes' to update.";
        el.profileAvatarsGrid.querySelectorAll(".preset-avatar-btn").forEach((b) => b.classList.remove("active"));
      };
      reader.readAsDataURL(file);
    });

    // Profile Form Submit (Save Changes)
    el.profileForm.addEventListener("submit", (e) => {
      e.preventDefault();
      if (!profileTargetUser) return;

      const newName = el.profileUsernameInput.value.trim();
      const newPin = el.profilePinInput.value.trim();

      if (!newName) {
        el.profileFeedback.className = "avatar-feedback error";
        el.profileFeedback.textContent = "Username cannot be empty.";
        return;
      }

      profileTargetUser.username = newName;
      profileTargetUser.avatarType = tempAvatarType;
      profileTargetUser.avatarKey = tempAvatarKey;
      profileTargetUser.customAvatarData = tempCustomAvatarData;

      if (newPin) {
        profileTargetUser.pin = newPin;
      }

      renderUserAvatarBadge(state.currentUser);
      if (state.activeTab === "manage") {
        renderManageView();
      }

      el.profileFeedback.className = "avatar-feedback success";
      el.profileFeedback.textContent = "✓ Profile updated successfully!";
      setTimeout(closeProfileModal, 600);
    });

    // Profile Log Out Button
    el.profileLogoutBtn.addEventListener("click", () => {
      closeProfileModal();
      state.currentUser = null;
      window.location.hash = "#login";
    });

    // Add User Form Submit
    el.closeAddUserBtn.addEventListener("click", closeAddUserModal);
    el.addUserModal.addEventListener("click", (e) => {
      if (e.target === el.addUserModal) {
        closeAddUserModal();
      }
    });

    el.addUserForm.addEventListener("submit", (e) => {
      e.preventDefault();
      const name = el.newUserName.value.trim();
      const pin = el.newUserPin.value.trim();
      const role = el.newUserRole.value;

      if (!name || !pin) {
        el.addUserFeedback.className = "avatar-feedback error";
        el.addUserFeedback.textContent = "Please fill in all fields.";
        return;
      }

      const existing = FLAN_MOCK_DATA.users.find(
        (u) => u.username.toLowerCase() === name.toLowerCase()
      );
      if (existing) {
        el.addUserFeedback.className = "avatar-feedback error";
        el.addUserFeedback.textContent = "A user with this username already exists.";
        return;
      }

      const newUser = {
        id: Date.now(),
        username: name,
        pin: pin,
        role: role,
        avatarType: "preset",
        avatarKey: newUserSelectedAvatarKey,
        customAvatarData: null
      };

      FLAN_MOCK_DATA.users.push(newUser);
      el.addUserFeedback.className = "avatar-feedback success";
      el.addUserFeedback.textContent = `✓ Created account for "${name}"!`;

      if (state.activeTab === "manage") {
        renderManageView();
      }

      setTimeout(closeAddUserModal, 600);
    });

    // Help Button -> Dedicated Software Manual View
    el.helpBtn.addEventListener("click", () => {
      window.location.hash = "#manual";
    });
  }

  // Start app
  initListeners();
  handleRoute();
})();
