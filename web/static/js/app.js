/**
 * Flan Media Server - Client Application Script
 */
document.addEventListener("DOMContentLoaded", () => {
  initModals();
  initLoginFlow();
  initSetupFlow();
  initProfileModal();
  initManageView();
  initMediaEditModal();
  initSearchPage();
  initVideoPlayer();
  initBookActions();
  initKeyboardShortcuts();
});

/* =========================================================================
   1. ACCESSIBLE MODAL DIALOGS
   ========================================================================= */
function initModals() {
  document.querySelectorAll("dialog.modal-dialog").forEach((dialog) => {
    // Close on click outside (backdrop)
    dialog.addEventListener("click", (e) => {
      const rect = dialog.getBoundingClientRect();
      const isInDialog = (
        rect.top <= e.clientY && e.clientY <= rect.top + rect.height &&
        rect.left <= e.clientX && e.clientX <= rect.left + rect.width
      );
      if (!isInDialog) {
        dialog.close();
      }
    });

    // Close buttons inside dialog
    dialog.querySelectorAll(".modal-close-btn, [data-close-modal]").forEach((btn) => {
      btn.addEventListener("click", () => dialog.close());
    });
  });
}

function openModal(id) {
  const dialog = document.getElementById(id);
  if (dialog && typeof dialog.showModal === "function") {
    dialog.showModal();
    const firstInput = dialog.querySelector("input:not([type=hidden]), button:not(.modal-close-btn)");
    if (firstInput) firstInput.focus();
  }
}

function closeModal(id) {
  const dialog = document.getElementById(id);
  if (dialog && typeof dialog.close === "function") {
    dialog.close();
  }
}

/* =========================================================================
   2. 2-STEP SEQUENTIAL LOGIN
   ========================================================================= */
function initLoginFlow() {
  const pickerCard = document.getElementById("login-picker-card");
  const pinCard = document.getElementById("login-pin-card");
  if (!pickerCard || !pinCard) return;

  let selectedUserID = null;

  document.querySelectorAll(".login-profile-card").forEach((card) => {
    card.addEventListener("click", () => {
      selectedUserID = card.dataset.userId;
      const username = card.dataset.username;
      const role = card.dataset.role;
      const avatarHtml = card.querySelector(".login-profile-avatar-frame").innerHTML;

      // Populate Step 2 banner
      document.getElementById("login-banner-avatar").innerHTML = avatarHtml;
      document.getElementById("login-banner-username").textContent = username;
      const roleTag = document.getElementById("login-banner-role");
      if (roleTag) {
        roleTag.textContent = role === "admin" ? "Admin" : "User";
      }

      pickerCard.style.display = "none";
      pinCard.style.display = "block";

      const pinInput = document.getElementById("pin-input");
      if (pinInput) {
        pinInput.value = "";
        pinInput.focus();
      }
    });
  });

  const switchBtn = document.getElementById("switch-profile-btn");
  if (switchBtn) {
    switchBtn.addEventListener("click", () => {
      selectedUserID = null;
      pinCard.style.display = "none";
      pickerCard.style.display = "block";
      const feedback = document.getElementById("login-feedback");
      if (feedback) feedback.textContent = "";
    });
  }

  const loginForm = document.getElementById("login-form");
  if (loginForm) {
    loginForm.addEventListener("submit", async (e) => {
      e.preventDefault();
      const pinInput = document.getElementById("pin-input");
      const feedback = document.getElementById("login-feedback");
      const pin = pinInput.value.trim();

      if (!pin) return;
      feedback.textContent = "";
      feedback.style.color = "var(--color-danger)";

      try {
        const res = await fetch("/api/login", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ user_id: parseInt(selectedUserID, 10), pin }),
        });

        const data = await res.json();
        if (res.ok && data.success) {
          window.location.href = "/video";
        } else {
          feedback.textContent = data.error || "Incorrect PIN";
          pinInput.classList.add("shake");
          setTimeout(() => pinInput.classList.remove("shake"), 500);
          pinInput.value = "";
          pinInput.focus();
        }
      } catch (err) {
        feedback.textContent = "Network error. Please try again.";
      }
    });
  }
}

/* =========================================================================
   3. FIRST-RUN SETUP
   ========================================================================= */
function initSetupFlow() {
  const setupForm = document.getElementById("setup-form");
  if (!setupForm) return;

  setupForm.addEventListener("submit", async (e) => {
    e.preventDefault();
    const token = document.getElementById("setup-token-input").value.trim();
    const username = document.getElementById("setup-username-input").value.trim();
    const pin = document.getElementById("setup-pin-input").value.trim();
    const feedback = document.getElementById("setup-feedback");

    feedback.textContent = "";
    try {
      const res = await fetch("/api/setup", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ token, username, pin }),
      });

      const data = await res.json();
      if (res.ok && data.success) {
        window.location.href = "/video";
      } else {
        feedback.textContent = data.error || "Setup failed";
      }
    } catch (err) {
      feedback.textContent = "Failed to connect to server";
    }
  });
}

/* =========================================================================
   4. USER PROFILE & SETTINGS MODAL
   ========================================================================= */
function initProfileModal() {
  const avatarBtn = document.getElementById("user-avatar");
  if (avatarBtn) {
    avatarBtn.addEventListener("click", () => openModal("profile-modal"));
  }

  const profileForm = document.getElementById("profile-form");
  if (profileForm) {
    profileForm.addEventListener("submit", async (e) => {
      e.preventDefault();
      const displayName = document.getElementById("profile-displayname-input")?.value.trim();
      const pin = document.getElementById("profile-pin-input")?.value.trim();
      const feedback = document.getElementById("profile-feedback");
      const selectedAvatar = document.querySelector('input[name="preset_avatar"]:checked')?.value;

      feedback.textContent = "";

      try {
        const res = await fetch("/api/users/self", {
          method: "PUT",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({
            display_name: displayName,
            avatar_icon: selectedAvatar || "flan",
            pin: pin || undefined,
          }),
        });
        const data = await res.json();
        if (res.ok) {
          feedback.style.color = "var(--color-success)";
          feedback.textContent = "Profile updated successfully";
          setTimeout(() => window.location.reload(), 800);
        } else {
          feedback.style.color = "var(--color-danger)";
          feedback.textContent = data.error || "Update failed";
        }
      } catch (err) {
        feedback.textContent = "Network error";
      }
    });
  }

  const logoutBtn = document.getElementById("profile-logout-btn");
  if (logoutBtn) {
    logoutBtn.addEventListener("click", async () => {
      try {
        await fetch("/api/logout", { method: "POST" });
      } finally {
        window.location.href = "/login";
      }
    });
  }
}

/* =========================================================================
   5. MANAGEMENT CONSOLE (SOURCES, INGESTION, USERS)
   ========================================================================= */
function initManageView() {
  // Add Source trigger
  const addSourceBtn = document.getElementById("open-add-source-btn");
  if (addSourceBtn) {
    addSourceBtn.addEventListener("click", () => openModal("add-source-modal"));
  }

  // Upload Files trigger
  const uploadFilesBtn = document.getElementById("open-upload-btn");
  if (uploadFilesBtn) {
    uploadFilesBtn.addEventListener("click", () => openModal("upload-media-modal"));
  }

  // Add User trigger
  const addUserBtn = document.getElementById("open-add-user-btn");
  if (addUserBtn) {
    addUserBtn.addEventListener("click", () => openModal("add-user-modal"));
  }

  // Add Source Form
  const addSourceForm = document.getElementById("add-source-form");
  if (addSourceForm) {
    addSourceForm.addEventListener("submit", async (e) => {
      e.preventDefault();
      const name = document.getElementById("source-name-input").value.trim();
      const media_type = document.getElementById("source-type-select").value;
      const folder_path = document.getElementById("source-path-input").value.trim();
      const feedback = document.getElementById("add-source-feedback");

      try {
        const res = await fetch("/api/sources", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ name, media_type, folder_path }),
        });
        const data = await res.json();
        if (res.ok) {
          closeModal("add-source-modal");
          window.location.reload();
        } else {
          feedback.textContent = data.error || "Failed to add source";
        }
      } catch (err) {
        feedback.textContent = "Network error";
      }
    });
  }

  // Scan Source trigger (Ingestion Pipeline Wizard)
  document.querySelectorAll(".btn-scan-source").forEach((btn) => {
    btn.addEventListener("click", () => {
      const sourceId = btn.dataset.sourceId;
      const sourceName = btn.dataset.sourceName;
      launchIngestionPipeline(sourceId, sourceName);
    });
  });

  // Edit Source trigger
  document.querySelectorAll(".btn-edit-source").forEach((btn) => {
    btn.addEventListener("click", () => {
      const id = btn.dataset.sourceId;
      const name = btn.dataset.sourceName;
      const path = btn.dataset.sourcePath;
      document.getElementById("edit-source-id").value = id;
      document.getElementById("edit-source-name-input").value = name;
      document.getElementById("edit-source-path-display").textContent = path;
      openModal("edit-source-modal");
    });
  });

  const editSourceForm = document.getElementById("edit-source-form");
  if (editSourceForm) {
    editSourceForm.addEventListener("submit", async (e) => {
      e.preventDefault();
      const id = document.getElementById("edit-source-id").value;
      const name = document.getElementById("edit-source-name-input").value.trim();
      try {
        const res = await fetch(`/api/sources/${id}`, {
          method: "PUT",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ name }),
        });
        if (res.ok) {
          closeModal("edit-source-modal");
          window.location.reload();
        }
      } catch (err) {}
    });
  }

  // Remove Source trigger
  document.querySelectorAll(".btn-remove-source").forEach((btn) => {
    btn.addEventListener("click", () => {
      const id = btn.dataset.sourceId;
      const name = btn.dataset.sourceName;
      const path = btn.dataset.sourcePath;
      document.getElementById("remove-source-name").textContent = name;
      document.getElementById("remove-source-path").textContent = path;
      const confirmBtn = document.getElementById("confirm-remove-source-btn");
      confirmBtn.onclick = async () => {
        await fetch(`/api/sources/${id}`, { method: "DELETE" });
        closeModal("remove-source-modal");
        window.location.reload();
      };
      openModal("remove-source-modal");
    });
  });

  // Rescan all sources
  const rescanAllBtn = document.getElementById("rescan-all-btn");
  if (rescanAllBtn) {
    rescanAllBtn.addEventListener("click", async () => {
      rescanAllBtn.disabled = true;
      rescanAllBtn.textContent = "Scanning...";
      try {
        const res = await fetch("/api/scan", { method: "POST" });
        const data = await res.json();
        alert(`Maintenance scan complete: ${data.items_scanned || 0} items checked.`);
        window.location.reload();
      } catch (err) {
        alert("Scan failed");
      } finally {
        rescanAllBtn.disabled = false;
        rescanAllBtn.textContent = "⟳ Rescan All Sources";
      }
    });
  }

  // Add User Form
  const addUserForm = document.getElementById("add-user-form");
  if (addUserForm) {
    addUserForm.addEventListener("submit", async (e) => {
      e.preventDefault();
      const username = document.getElementById("new-user-name").value.trim();
      const pin = document.getElementById("new-user-pin").value.trim();
      const role = document.getElementById("new-user-role").value;
      const avatarIcon = document.querySelector('input[name="new_preset_avatar"]:checked')?.value || "flan";
      const feedback = document.getElementById("add-user-feedback");

      try {
        const res = await fetch("/api/users", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ username, pin, role, avatar_icon: avatarIcon }),
        });
        const data = await res.json();
        if (res.ok) {
          closeModal("add-user-modal");
          window.location.reload();
        } else {
          feedback.textContent = data.error || "Failed to create user";
        }
      } catch (err) {
        feedback.textContent = "Network error";
      }
    });
  }

  // Upload Form with Progress
  const uploadForm = document.getElementById("upload-media-form");
  if (uploadForm) {
    uploadForm.addEventListener("submit", (e) => {
      e.preventDefault();
      const sourceId = document.getElementById("upload-target-source").value;
      const title = document.getElementById("upload-title-input").value.trim();
      const fileInput = document.getElementById("upload-file-input");
      const file = fileInput.files[0];
      if (!file) return;

      const progressWrap = document.getElementById("upload-progress-wrap");
      const progressFill = document.getElementById("upload-progress-fill");
      const percentText = document.getElementById("upload-percent-text");
      const startBtn = document.getElementById("start-upload-btn");
      const feedback = document.getElementById("upload-feedback");

      progressWrap.style.display = "flex";
      startBtn.disabled = true;

      const formData = new FormData();
      formData.append("source_id", sourceId);
      formData.append("title", title);
      formData.append("file", file);

      const xhr = new XMLHttpRequest();
      xhr.open("POST", "/api/upload", true);

      xhr.upload.onprogress = (event) => {
        if (event.lengthComputable) {
          const percent = Math.round((event.loaded / event.total) * 100);
          progressFill.style.width = percent + "%";
          percentText.textContent = percent + "%";
        }
      };

      xhr.onload = () => {
        if (xhr.status === 200) {
          feedback.style.color = "var(--color-success)";
          feedback.textContent = "Upload complete!";
          setTimeout(() => {
            closeModal("upload-media-modal");
            window.location.reload();
          }, 800);
        } else {
          feedback.style.color = "var(--color-danger)";
          feedback.textContent = "Upload failed (" + xhr.status + ")";
          startBtn.disabled = false;
        }
      };

      xhr.onerror = () => {
        feedback.textContent = "Network error during upload";
        startBtn.disabled = false;
      };

      xhr.send(formData);
    });
  }
}

/* =========================================================================
   6. INGESTION PIPELINE WIZARD
   ========================================================================= */
async function launchIngestionPipeline(sourceId, sourceName) {
  const modal = document.getElementById("ingestion-pipeline-modal");
  const subtitle = document.getElementById("pipeline-source-subtitle");
  const list = document.getElementById("pipeline-items-list");
  const commitBtn = document.getElementById("commit-pipeline-btn");
  const feedback = document.getElementById("pipeline-feedback");

  if (!modal || !list) return;

  subtitle.textContent = `Gathering items from ${sourceName}...`;
  list.innerHTML = `<div style="padding: 20px; text-align: center; font-weight: 700;">Scanning files on disk...</div>`;
  feedback.textContent = "";
  openModal("ingestion-pipeline-modal");

  try {
    const res = await fetch(`/api/sources/${sourceId}/scan`, { method: "POST" });
    if (!res.ok) {
      list.innerHTML = `<div style="color: var(--color-danger); padding: 20px;">Scan failed. Verify drive is mounted.</div>`;
      return;
    }

    const items = await res.json();
    if (!items || items.length === 0) {
      list.innerHTML = `<div style="padding: 20px; text-align: center;">No new media items found in this directory.</div>`;
      commitBtn.disabled = true;
      return;
    }

    renderPipelineItems(items, sourceId);
  } catch (err) {
    list.innerHTML = `<div style="color: var(--color-danger); padding: 20px;">Failed to communicate with scanner.</div>`;
  }
}

function renderPipelineItems(items, sourceId) {
  const list = document.getElementById("pipeline-items-list");
  const commitBtn = document.getElementById("commit-pipeline-btn");
  const counter = document.getElementById("pipeline-selected-counter");

  list.innerHTML = "";
  commitBtn.disabled = false;

  items.forEach((item, idx) => {
    const card = document.createElement("div");
    card.className = "pipeline-candidate-card";
    card.style.border = "2px solid #000";
    card.style.borderRadius = "8px";
    card.style.padding = "10px";
    card.style.background = "#fff";

    card.innerHTML = `
      <div style="display: flex; align-items: center; gap: 10px;">
        <input type="checkbox" class="pipeline-item-check" data-idx="${idx}" checked style="width: 20px; height: 20px;" />
        <div style="flex: 1;">
          <div style="font-size: 0.8rem; color: #555; font-family: monospace;">📁 ${item.folder_path} (${(item.files || []).length} files)</div>
          <div style="display: flex; gap: 8px; margin-top: 6px;">
            <input type="text" class="tactile-text-input pipeline-item-title" value="${item.title}" style="flex: 2; padding: 4px 8px; font-size: 0.9rem;" placeholder="Title" />
            <input type="number" class="tactile-text-input pipeline-item-year" value="${item.release_year || ''}" style="width: 80px; padding: 4px 8px; font-size: 0.9rem;" placeholder="Year" />
          </div>
        </div>
      </div>
    `;
    list.appendChild(card);
  });

  const updateCounter = () => {
    const checked = list.querySelectorAll(".pipeline-item-check:checked").length;
    counter.textContent = `${checked} selected`;
    commitBtn.textContent = `Ingest Selected Items (${checked})`;
    commitBtn.disabled = checked === 0;
  };
  updateCounter();

  list.querySelectorAll(".pipeline-item-check").forEach((cb) => {
    cb.addEventListener("change", updateCounter);
  });

  // Master Select All
  const masterSelect = document.getElementById("pipeline-select-all");
  if (masterSelect) {
    masterSelect.checked = true;
    masterSelect.addEventListener("change", () => {
      list.querySelectorAll(".pipeline-item-check").forEach((cb) => {
        cb.checked = masterSelect.checked;
      });
      updateCounter();
    });
  }

  // Filter input
  const filterInput = document.getElementById("pipeline-filter-input");
  if (filterInput) {
    filterInput.value = "";
    filterInput.addEventListener("input", () => {
      const q = filterInput.value.toLowerCase();
      list.querySelectorAll(".pipeline-candidate-card").forEach((card) => {
        const title = card.querySelector(".pipeline-item-title").value.toLowerCase();
        card.style.display = title.includes(q) ? "block" : "none";
      });
    });
  }

  commitBtn.onclick = async () => {
    const approved = [];
    list.querySelectorAll(".pipeline-item-check:checked").forEach((cb) => {
      const idx = parseInt(cb.dataset.idx, 10);
      const card = cb.closest(".pipeline-candidate-card");
      const title = card.querySelector(".pipeline-item-title").value.trim();
      const year = parseInt(card.querySelector(".pipeline-item-year").value, 10) || 0;
      const baseItem = items[idx];
      baseItem.title = title;
      baseItem.release_year = year;
      approved.push(baseItem);
    });

    if (approved.length === 0) return;

    commitBtn.disabled = true;
    commitBtn.textContent = "Committing to catalog...";

    try {
      const res = await fetch(`/api/sources/${sourceId}/ingest`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ approved }),
      });
      if (res.ok) {
        closeModal("ingestion-pipeline-modal");
        window.location.reload();
      } else {
        alert("Failed to commit items");
        commitBtn.disabled = false;
      }
    } catch (err) {
      alert("Network error");
      commitBtn.disabled = false;
    }
  };
}

/* =========================================================================
   7. VIDEO DETAIL EDITOR MODAL
   ========================================================================= */
function initMediaEditModal() {
  const editBtn = document.getElementById("open-media-edit-btn");
  if (editBtn) {
    editBtn.addEventListener("click", () => openModal("media-edit-modal"));
  }

  const editForm = document.getElementById("media-edit-form");
  if (editForm) {
    editForm.addEventListener("submit", async (e) => {
      e.preventDefault();
      const videoId = document.getElementById("edit-video-id").value;
      const title = document.getElementById("edit-video-title").value.trim();
      const year = parseInt(document.getElementById("edit-video-year").value, 10) || 0;
      const videoType = document.getElementById("edit-video-type").value;
      const overview = document.getElementById("edit-video-overview").value.trim();
      const feedback = document.getElementById("media-edit-feedback");

      // Playable files
      const files = [];
      document.querySelectorAll(".edit-file-row").forEach((row) => {
        files.push({
          file_id: parseInt(row.dataset.fileId, 10),
          custom_title: row.querySelector(".edit-file-title").value.trim(),
          order_index: parseInt(row.querySelector(".edit-file-order").value, 10) || 0,
          is_hidden: row.querySelector(".edit-file-hidden").checked,
        });
      });

      try {
        const res = await fetch(`/api/media/video/${videoId}`, {
          method: "PUT",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ title, release_year: year, video_type: videoType, overview, files }),
        });
        if (res.ok) {
          closeModal("media-edit-modal");
          window.location.reload();
        } else {
          feedback.textContent = "Failed to save changes";
        }
      } catch (err) {
        feedback.textContent = "Network error";
      }
    });
  }
}

/* =========================================================================
   8. SEARCH PAGE (DEBOUNCED LIVE QUERY)
   ========================================================================= */
function initSearchPage() {
  const searchInput = document.getElementById("search-page-input");
  const clearBtn = document.getElementById("search-page-clear-btn");
  const resultsGrid = document.getElementById("search-media-grid");
  const promptCard = document.getElementById("search-prompt-card");
  const metaStatus = document.getElementById("search-meta-status");
  if (!searchInput || !resultsGrid) return;

  let debounceTimer = null;

  const performSearch = async () => {
    const q = searchInput.value.trim();
    if (!q) {
      if (promptCard) promptCard.style.display = "block";
      resultsGrid.innerHTML = "";
      if (metaStatus) metaStatus.textContent = "Enter search terms above";
      return;
    }

    if (promptCard) promptCard.style.display = "none";
    if (metaStatus) metaStatus.textContent = `Searching for "${q}"...`;

    try {
      const res = await fetch(`/api/search?q=${encodeURIComponent(q)}`);
      const data = await res.json();
      const videos = data.videos || [];
      const books = data.books || [];
      const total = videos.length + books.length;

      if (metaStatus) metaStatus.textContent = `Found ${total} items matching "${q}"`;

      resultsGrid.innerHTML = "";
      if (total === 0) {
        resultsGrid.innerHTML = `<div style="grid-column: 1 / -1; padding: 40px; text-align: center; font-weight: 700;">No titles matching your search.</div>`;
        return;
      }

      videos.forEach((v) => {
        const card = createCatalogCard("video", v.video_id, v.title, v.release_year ? String(v.release_year) : "Video", v.cover_path);
        resultsGrid.appendChild(card);
      });

      books.forEach((b) => {
        const card = createCatalogCard("books", b.book_id, b.title, b.author || "Book", b.cover_path);
        resultsGrid.appendChild(card);
      });
    } catch (err) {
      if (metaStatus) metaStatus.textContent = "Search error";
    }
  };

  searchInput.addEventListener("input", () => {
    clearTimeout(debounceTimer);
    debounceTimer = setTimeout(performSearch, 250);
  });

  if (clearBtn) {
    clearBtn.addEventListener("click", () => {
      searchInput.value = "";
      performSearch();
      searchInput.focus();
    });
  }
}

function createCatalogCard(type, id, title, subtitle, coverPath) {
  const article = document.createElement("article");
  article.className = "card card--" + (type === "video" ? "video" : "book");
  const coverUrl = coverPath ? `/covers/${type === "video" ? "video" : "book"}/${id}` : "";

  article.innerHTML = `
    <a href="/${type}/${id}" class="card__link">
      <div class="card__poster-wrap">
        ${coverUrl ? `<img src="${coverUrl}" alt="${title}" class="card__poster" loading="lazy">` : `<div style="width:100%;height:100%;display:flex;align-items:center;justify-content:center;background:#eae8f2;font-size:2rem;font-weight:900;">🎬</div>`}
      </div>
      <div class="card__footer-band">
        <h3 class="card__title">${title}</h3>
        <span class="card__subtitle">${subtitle}</span>
      </div>
    </a>
  `;
  return article;
}

/* =========================================================================
   9. VIDEO PLAYER (PLYR + PROGRESS SYNC)
   ========================================================================= */
function initVideoPlayer() {
  const playerEl = document.getElementById("player");
  if (!playerEl || typeof Plyr === "undefined") return;

  const fileId = playerEl.dataset.fileId;
  const initialPos = parseFloat(playerEl.dataset.initialPos) || 0;
  const resumePrompt = document.getElementById("resume-prompt-card");

  const player = new Plyr(playerEl, {
    controls: ["play-large", "play", "progress", "current-time", "duration", "mute", "volume", "captions", "settings", "fullscreen"],
    keyboard: { focused: true, global: true },
  });

  // Resume prompt
  if (resumePrompt && initialPos > 10) {
    resumePrompt.style.display = "block";
    document.getElementById("btn-resume-pos")?.addEventListener("click", () => {
      player.currentTime = initialPos;
      resumePrompt.style.display = "none";
      player.play();
    });
    document.getElementById("btn-start-over")?.addEventListener("click", () => {
      player.currentTime = 0;
      resumePrompt.style.display = "none";
      player.play();
    });
  }

  // Periodic progress sync loop (every 15s)
  let lastReportedTime = 0;
  const syncProgress = (isFinished = false) => {
    const curTime = player.currentTime || 0;
    const dur = player.duration || 1;
    if (Math.abs(curTime - lastReportedTime) < 5 && !isFinished) return;
    lastReportedTime = curTime;

    const percentage = Math.min(100, (curTime / dur) * 100);
    fetch("/api/progress", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        media_type: "video",
        file_id: parseInt(fileId, 10),
        position_data: curTime.toFixed(1),
        percentage: percentage,
        is_finished: isFinished,
      }),
    }).catch(() => {});
  };

  setInterval(() => {
    if (!player.paused) syncProgress(false);
  }, 15000);

  player.on("pause", () => syncProgress(false));
  player.on("ended", () => syncProgress(true));

  // Handle format errors
  player.on("error", () => {
    const fallbackBanner = document.getElementById("vlc-fallback-banner");
    if (fallbackBanner) fallbackBanner.style.display = "block";
  });
}

/* =========================================================================
   10. BOOK PROGRESS ACTIONS
   ========================================================================= */
function initBookActions() {
  document.querySelectorAll(".btn-book-status").forEach((btn) => {
    btn.addEventListener("click", async () => {
      const fileId = btn.dataset.fileId;
      const currentStatus = btn.dataset.status;
      const nextStatus = currentStatus === "reading" ? "finished" : "reading";

      try {
        await fetch("/api/progress", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({
            media_type: "book",
            file_id: parseInt(fileId, 10),
            position_data: nextStatus,
            is_finished: nextStatus === "finished",
          }),
        });
        window.location.reload();
      } catch (err) {}
    });
  });
}

/* =========================================================================
   11. GLOBAL KEYBOARD SHORTCUTS
   ========================================================================= */
function initKeyboardShortcuts() {
  window.addEventListener("keydown", (e) => {
    if (e.target.tagName === "INPUT" || e.target.tagName === "TEXTAREA") {
      return;
    }
    // "/" or "Ctrl+K" jumps to /search
    if (e.key === "/" || (e.ctrlKey && e.key === "k")) {
      e.preventDefault();
      window.location.href = "/search";
    }
  });
}
