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
    filterExpanded: false,
    selectedFormat: "all", // 'all' | 'movies' | 'series' (or 'epub' | 'pdf')
    selectedStatus: "all", // 'all' | 'in-progress' | 'unwatched'
    selectedSource: "all", // 'all' | sourceId
    sortOrder: "random",   // 'random' (default) | 'recent' | 'az' | 'year'
    currentPage: 1,
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
    addUserFeedback: document.getElementById("add-user-feedback"),

    // Video Detail Edit Modal
    mediaEditModal: document.getElementById("media-edit-modal"),
    closeMediaEditBtn: document.getElementById("close-media-edit-btn"),
    cancelMediaEditBtn: document.getElementById("cancel-media-edit-btn"),
    mediaEditForm: document.getElementById("media-edit-form"),
    editVideoId: document.getElementById("edit-video-id"),
    editVideoTitle: document.getElementById("edit-video-title"),
    editVideoYear: document.getElementById("edit-video-year"),
    editVideoType: document.getElementById("edit-video-type"),
    editVideoOverview: document.getElementById("edit-video-overview"),
    editFilesTableWrap: document.getElementById("edit-files-table-wrap"),
    mediaEditFeedback: document.getElementById("media-edit-feedback"),

    // Add Storage Source Modal
    addSourceModal: document.getElementById("add-source-modal"),
    closeAddSourceBtn: document.getElementById("close-add-source-btn"),
    addSourceForm: document.getElementById("add-source-form"),
    sourceNameInput: document.getElementById("source-name-input"),
    sourceTypeSelect: document.getElementById("source-type-select"),
    sourcePathInput: document.getElementById("source-path-input"),
    addSourceFeedback: document.getElementById("add-source-feedback"),

    // Direct Web Upload Modal
    uploadMediaModal: document.getElementById("upload-media-modal"),
    closeUploadBtn: document.getElementById("close-upload-btn"),
    uploadMediaForm: document.getElementById("upload-media-form"),
    uploadTargetSource: document.getElementById("upload-target-source"),
    uploadTitleInput: document.getElementById("upload-title-input"),
    uploadFileInput: document.getElementById("upload-file-input"),
    uploadProgressWrap: document.getElementById("upload-progress-wrap"),
    uploadStatusText: document.getElementById("upload-status-text"),
    uploadPercentText: document.getElementById("upload-percent-text"),
    uploadProgressFill: document.getElementById("upload-progress-fill"),
    uploadFeedback: document.getElementById("upload-feedback"),
    startUploadBtn: document.getElementById("start-upload-btn"),

    // Ingestion Pipeline Modal
    ingestionPipelineModal: document.getElementById("ingestion-pipeline-modal"),
    pipelineModalTitle: document.getElementById("pipeline-modal-title"),
    pipelineSourceSubtitle: document.getElementById("pipeline-source-subtitle"),
    closePipelineBtn: document.getElementById("close-pipeline-btn"),
    cancelPipelineBtn: document.getElementById("cancel-pipeline-btn"),
    pipelineSelectAll: document.getElementById("pipeline-select-all"),
    pipelineSelectedCounter: document.getElementById("pipeline-selected-counter"),
    pipelineFilterInput: document.getElementById("pipeline-filter-input"),
    pipelineItemsList: document.getElementById("pipeline-items-list"),
    pipelineFeedback: document.getElementById("pipeline-feedback"),
    commitPipelineBtn: document.getElementById("commit-pipeline-btn"),

    // Remove Storage Source Warning Modal
    removeSourceModal: document.getElementById("remove-source-modal"),
    removeSourceName: document.getElementById("remove-source-name"),
    removeSourcePath: document.getElementById("remove-source-path"),
    closeRemoveSourceBtn: document.getElementById("close-remove-source-btn"),
    cancelRemoveSourceBtn: document.getElementById("cancel-remove-source-btn"),
    confirmRemoveSourceBtn: document.getElementById("confirm-remove-source-btn"),

    // Edit Storage Source Modal
    editSourceModal: document.getElementById("edit-source-modal"),
    editSourceForm: document.getElementById("edit-source-form"),
    editSourceNameInput: document.getElementById("edit-source-name-input"),
    editSourcePathDisplay: document.getElementById("edit-source-path-display"),
    editSourceFeedback: document.getElementById("edit-source-feedback"),
    closeEditSourceBtn: document.getElementById("close-edit-source-btn"),
    cancelEditSourceBtn: document.getElementById("cancel-edit-source-btn"),
    saveEditSourceBtn: document.getElementById("save-edit-source-btn")
  };

  // Helper: Get User Avatar SVG or custom image
  function getUserAvatarSVG(user) {
    if (!user) return "";
    if (user.avatarType === "custom" && user.customAvatarData) {
      return `<img src="${user.customAvatarData}" alt="${user.username}" style="width:100%; height:100%; object-fit:cover;" />`;
    }
    const key = user.avatarKey || "mascot";
    const preset = FLAN_MOCK_DATA.presetAvatars[key] || FLAN_MOCK_DATA.presetAvatars.mascot;
    return preset.svg;
  }

  // Render User Avatar Badge (Top Header)
  function renderUserAvatarBadge(user) {
    if (!user || !el.userAvatar) return;
    el.userAvatar.innerHTML = getUserAvatarSVG(user);
  }

  // =========================================================================
  // UNIFIED MY PROFILE MODAL MANAGER (Accessible HTML5 Dialog)
  // =========================================================================
  let profileTargetUser = null;
  let tempAvatarType = null;
  let tempAvatarKey = null;
  let tempCustomAvatarData = null;
  let lastFocusedElement = null;

  function openProfileModal(user) {
    lastFocusedElement = document.activeElement;
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
    if (typeof el.profileModal.showModal === "function") {
      el.profileModal.showModal();
    } else {
      el.profileModal.classList.remove("hidden");
    }
    el.profileUsernameInput.focus();
  }

  function renderProfileAvatarsGrid() {
    el.profileAvatarsGrid.innerHTML = Object.entries(FLAN_MOCK_DATA.presetAvatars)
      .map(([key, item]) => {
        const isActive = tempAvatarType === "preset" && tempAvatarKey === key && !tempCustomAvatarData;
        return `
          <button type="button" class="preset-avatar-btn ${isActive ? "active" : ""}" data-key="${key}" title="${item.name}" aria-label="Avatar ${item.name}">
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
    if (typeof el.profileModal.close === "function" && el.profileModal.open) {
      el.profileModal.close();
    } else {
      el.profileModal.classList.add("hidden");
    }
    profileTargetUser = null;
    if (lastFocusedElement) {
      lastFocusedElement.focus();
      lastFocusedElement = null;
    }
  }

  // =========================================================================
  // ADD USER MODAL MANAGER (Manage Server - Accessible HTML5 Dialog)
  // =========================================================================
  let newUserSelectedAvatarKey = "mascot";

  function openAddUserModal() {
    lastFocusedElement = document.activeElement;
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
          <button type="button" class="preset-avatar-btn ${isActive ? "active" : ""}" data-key="${key}" title="${item.name}" aria-label="Avatar ${item.name}">
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

    if (typeof el.addUserModal.showModal === "function") {
      el.addUserModal.showModal();
    } else {
      el.addUserModal.classList.remove("hidden");
    }
    el.newUserName.focus();
  }

  function closeAddUserModal() {
    if (typeof el.addUserModal.close === "function" && el.addUserModal.open) {
      el.addUserModal.close();
    } else {
      el.addUserModal.classList.add("hidden");
    }
    if (lastFocusedElement) {
      lastFocusedElement.focus();
      lastFocusedElement = null;
    }
  }

  // Media Detail Edit Modal
  let currentEditingItem = null;
  function openMediaEditModal(item) {
    lastFocusedElement = document.activeElement;
    currentEditingItem = item;
    el.mediaEditFeedback.textContent = "";
    el.mediaEditFeedback.className = "avatar-feedback";

    el.editVideoId.value = item.id;
    el.editVideoTitle.value = item.title || item.proposedTitle || "";
    el.editVideoYear.value = item.releaseYear || item.proposedYear || "";
    el.editVideoType.value = item.type || item.mediaType || "series";
    el.editVideoOverview.value = item.overview || "";

    renderEditFilesTable();

    if (typeof el.mediaEditModal.showModal === "function") {
      el.mediaEditModal.showModal();
    } else {
      el.mediaEditModal.classList.remove("hidden");
    }
    el.editVideoTitle.focus();
  }

  function renderEditFilesTable() {
    if (!currentEditingItem || !currentEditingItem.files) {
      el.editFilesTableWrap.innerHTML = "<p style='font-size:0.85rem; color:#666;'>No playable files detected.</p>";
      return;
    }

    el.editFilesTableWrap.innerHTML = currentEditingItem.files
      .map((f, i) => `
        <div class="edit-file-row" data-file-id="${f.id}" style="display:flex; align-items:center; gap:8px; padding:6px 0; border-bottom:1px solid #ddd;">
          <span style="font-weight:700; opacity:0.6; min-width:24px;">#${i + 1}</span>
          <input type="text" class="tactile-text-input file-title-edit" data-file-id="${f.id}" value="${f.title}" style="flex:1; padding:4px 8px; font-size:0.85rem;" />
          <button type="button" class="tactile-action-btn secondary move-file-up" data-file-id="${f.id}" title="Move Up" style="padding:2px 8px; font-size:0.8rem;" ${i === 0 ? "disabled" : ""}>▲</button>
          <button type="button" class="tactile-action-btn secondary move-file-down" data-file-id="${f.id}" title="Move Down" style="padding:2px 8px; font-size:0.8rem;" ${i === currentEditingItem.files.length - 1 ? "disabled" : ""}>▼</button>
          <label style="display:flex; align-items:center; gap:4px; font-size:0.8rem; cursor:pointer;">
            <input type="checkbox" class="file-hide-check" data-file-id="${f.id}" ${f.isHidden ? "checked" : ""} /> Hide
          </label>
        </div>
      `)
      .join("");

    // Wire Up/Down Reordering
    const upBtns = el.editFilesTableWrap.querySelectorAll(".move-file-up");
    upBtns.forEach((btn) => {
      btn.addEventListener("click", () => {
        const fId = parseInt(btn.getAttribute("data-file-id"));
        const idx = currentEditingItem.files.findIndex((x) => x.id === fId);
        if (idx > 0) {
          const temp = currentEditingItem.files[idx];
          currentEditingItem.files[idx] = currentEditingItem.files[idx - 1];
          currentEditingItem.files[idx - 1] = temp;
          renderEditFilesTable();
        }
      });
    });

    const downBtns = el.editFilesTableWrap.querySelectorAll(".move-file-down");
    downBtns.forEach((btn) => {
      btn.addEventListener("click", () => {
        const fId = parseInt(btn.getAttribute("data-file-id"));
        const idx = currentEditingItem.files.findIndex((x) => x.id === fId);
        if (idx < currentEditingItem.files.length - 1) {
          const temp = currentEditingItem.files[idx];
          currentEditingItem.files[idx] = currentEditingItem.files[idx + 1];
          currentEditingItem.files[idx + 1] = temp;
          renderEditFilesTable();
        }
      });
    });
  }

  function closeMediaEditModal() {
    if (typeof el.mediaEditModal.close === "function" && el.mediaEditModal.open) {
      el.mediaEditModal.close();
    } else {
      el.mediaEditModal.classList.add("hidden");
    }
    currentEditingItem = null;
    if (lastFocusedElement) {
      lastFocusedElement.focus();
      lastFocusedElement = null;
    }
  }

  // Add Storage Source Modal
  function openAddSourceModal() {
    lastFocusedElement = document.activeElement;
    el.addSourceFeedback.textContent = "";
    el.addSourceFeedback.className = "avatar-feedback";
    el.sourceNameInput.value = "";
    el.sourcePathInput.value = "";

    if (typeof el.addSourceModal.showModal === "function") {
      el.addSourceModal.showModal();
    } else {
      el.addSourceModal.classList.remove("hidden");
    }
    el.sourceNameInput.focus();
  }

  function closeAddSourceModal() {
    if (typeof el.addSourceModal.close === "function" && el.addSourceModal.open) {
      el.addSourceModal.close();
    } else {
      el.addSourceModal.classList.add("hidden");
    }
    if (lastFocusedElement) {
      lastFocusedElement.focus();
      lastFocusedElement = null;
    }
  }

  // Direct Web Upload Modal
  function openUploadModal() {
    lastFocusedElement = document.activeElement;
    el.uploadFeedback.textContent = "";
    el.uploadFeedback.className = "avatar-feedback";
    el.uploadTitleInput.value = "";
    el.uploadFileInput.value = "";
    el.uploadProgressWrap.style.display = "none";
    el.uploadProgressFill.style.width = "0%";
    el.uploadPercentText.textContent = "0%";
    el.startUploadBtn.disabled = false;
    el.startUploadBtn.textContent = "Upload Files";

    // Populate target sources
    el.uploadTargetSource.innerHTML = FLAN_MOCK_DATA.storageSources
      .map(
        (s) =>
          `<option value="${s.id}">${s.name} (${s.path}) [${s.mediaType.toUpperCase()}]</option>`
      )
      .join("");

    if (typeof el.uploadMediaModal.showModal === "function") {
      el.uploadMediaModal.showModal();
    } else {
      el.uploadMediaModal.classList.remove("hidden");
    }
    el.uploadTitleInput.focus();
  }

  function closeUploadModal() {
    if (typeof el.uploadMediaModal.close === "function" && el.uploadMediaModal.open) {
      el.uploadMediaModal.close();
    } else {
      el.uploadMediaModal.classList.add("hidden");
    }
    if (lastFocusedElement) {
      lastFocusedElement.focus();
      lastFocusedElement = null;
    }
  }

  // Remove Storage Source Warning Modal Manager
  let pendingRemoveSource = null;

  function openRemoveSourceModal(source) {
    lastFocusedElement = document.activeElement;
    pendingRemoveSource = source;
    el.removeSourceName.textContent = `"${source.name}"`;
    el.removeSourcePath.textContent = source.path;

    if (typeof el.removeSourceModal.showModal === "function") {
      el.removeSourceModal.showModal();
    } else {
      el.removeSourceModal.classList.remove("hidden");
    }
    el.cancelRemoveSourceBtn.focus();
  }

  function closeRemoveSourceModal() {
    if (typeof el.removeSourceModal.close === "function" && el.removeSourceModal.open) {
      el.removeSourceModal.close();
    } else {
      el.removeSourceModal.classList.add("hidden");
    }
    pendingRemoveSource = null;
    if (lastFocusedElement) {
      lastFocusedElement.focus();
      lastFocusedElement = null;
    }
  }

  // Edit Storage Source Modal Manager
  let pendingEditSource = null;

  function openEditSourceModal(source) {
    lastFocusedElement = document.activeElement;
    pendingEditSource = source;
    el.editSourceFeedback.textContent = "";
    el.editSourceFeedback.className = "avatar-feedback";
    el.editSourceNameInput.value = source.name;
    el.editSourcePathDisplay.textContent = source.path;

    if (typeof el.editSourceModal.showModal === "function") {
      el.editSourceModal.showModal();
    } else {
      el.editSourceModal.classList.remove("hidden");
    }
    el.editSourceNameInput.focus();
    el.editSourceNameInput.select();
  }

  function closeEditSourceModal() {
    if (typeof el.editSourceModal.close === "function" && el.editSourceModal.open) {
      el.editSourceModal.close();
    } else {
      el.editSourceModal.classList.add("hidden");
    }
    pendingEditSource = null;
    if (lastFocusedElement) {
      lastFocusedElement.focus();
      lastFocusedElement = null;
    }
  }

  // =========================================================================
  // INGESTION PIPELINE WIZARD (Accessible HTML5 Dialog)
  // =========================================================================
  let pipelineActiveSource = null;
  let pipelineCandidates = [];

  function openIngestionPipeline(source) {
    lastFocusedElement = document.activeElement;
    pipelineActiveSource = source;
    el.pipelineFeedback.textContent = "";
    el.pipelineFeedback.className = "avatar-feedback";
    el.pipelineFilterInput.value = "";

    // Gather pending items for this source (or all sources)
    let rawItems = [];
    if (source) {
      rawItems = FLAN_MOCK_DATA.pendingIngestion.filter((x) => x.sourceId === source.id);
    } else {
      rawItems = [...FLAN_MOCK_DATA.pendingIngestion];
    }

    pipelineCandidates = rawItems.map((item) => ({
      id: item.id,
      sourceId: item.sourceId,
      sourceName: item.sourceName,
      rawFolder: item.rawFolder,
      proposedTitle: item.proposedTitle,
      proposedYear: item.proposedYear,
      proposedAuthor: item.proposedAuthor,
      mediaType: item.mediaType,
      fileCount: item.fileCount,
      coverColor: item.coverColor,
      badge: item.badge,
      overview: item.overview,
      files: (item.files || []).map((f) => ({ ...f })),
      // State in pipeline
      selected: true,
      expanded: false,
      editTitle: item.proposedTitle,
      editMeta: item.proposedYear || item.proposedAuthor || ""
    }));

    if (source) {
      el.pipelineModalTitle.textContent = `Scan Media: ${source.name}`;
      el.pipelineSourceSubtitle.textContent = `Path: ${source.path} • ${pipelineCandidates.length} new candidates discovered`;
    } else {
      el.pipelineModalTitle.textContent = "Scan Media: All Storage Sources";
      el.pipelineSourceSubtitle.textContent = `Indexing all active sources • ${pipelineCandidates.length} total candidates discovered`;
    }

    renderPipelineList();

    if (typeof el.ingestionPipelineModal.showModal === "function") {
      el.ingestionPipelineModal.showModal();
    } else {
      el.ingestionPipelineModal.classList.remove("hidden");
    }
    el.pipelineFilterInput.focus();
  }

  function closeIngestionPipeline() {
    if (typeof el.ingestionPipelineModal.close === "function" && el.ingestionPipelineModal.open) {
      el.ingestionPipelineModal.close();
    } else {
      el.ingestionPipelineModal.classList.add("hidden");
    }
    if (lastFocusedElement) {
      lastFocusedElement.focus();
      lastFocusedElement = null;
    }
  }

  function renderPipelineList() {
    const filterQuery = el.pipelineFilterInput.value.trim().toLowerCase();
    const visibleItems = pipelineCandidates.filter((item) => {
      if (!filterQuery) return true;
      return (
        item.editTitle.toLowerCase().includes(filterQuery) ||
        item.rawFolder.toLowerCase().includes(filterQuery) ||
        item.sourceName.toLowerCase().includes(filterQuery)
      );
    });

    const selectedCount = pipelineCandidates.filter((x) => x.selected).length;
    const totalCount = pipelineCandidates.length;

    el.pipelineSelectAll.checked = totalCount > 0 && selectedCount === totalCount;
    el.pipelineSelectAll.indeterminate = selectedCount > 0 && selectedCount < totalCount;
    el.pipelineSelectedCounter.textContent = `Select All (${selectedCount} of ${totalCount} selected)`;
    el.commitPipelineBtn.textContent = `Ingest Selected Items (${selectedCount})`;
    el.commitPipelineBtn.disabled = selectedCount === 0;

    if (pipelineCandidates.length === 0) {
      el.pipelineItemsList.innerHTML = `
        <div style="padding: 32px 16px; background: #f9f9f9; border: 2px dashed #ccc; border-radius: 8px; text-align: center; color: #555;">
          <div style="font-weight: 800; font-size: 1.1rem; color: #2e7d32; margin-bottom: 6px;">✓ Storage Folder Up to Date</div>
          <div style="font-size: 0.9rem;">All media discovered in this source has already been indexed and ingested into your catalog.</div>
        </div>
      `;
      return;
    }

    if (visibleItems.length === 0) {
      el.pipelineItemsList.innerHTML = `
        <div style="padding: 24px 16px; background: #f9f9f9; border: 2px dashed #ccc; border-radius: 8px; text-align: center; color: #666; font-size: 0.9rem;">
          No candidates match your search filter "<strong>${filterQuery}</strong>".
        </div>
      `;
      return;
    }

    el.pipelineItemsList.innerHTML = visibleItems
      .map(
        (c) => `
        <div class="pipeline-item-card ${c.selected ? "is-selected" : "is-excluded"}" data-id="${c.id}">
          <div class="pipeline-item-header">
            <input type="checkbox" class="pipeline-checkbox item-select-check" data-id="${c.id}" ${c.selected ? "checked" : ""} aria-label="Select ${c.editTitle} for ingestion" />
            <div class="pipeline-item-body">
              <div class="pipeline-item-meta">
                <span class="pipeline-source-badge">${c.sourceName}</span>
                <span class="pipeline-folder-path"><svg style="width:13px;height:13px;vertical-align:-1px;margin-right:4px;" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"></path></svg>${c.rawFolder}</span>
                <span style="font-weight: 700; color: #333;">(${c.fileCount} files)</span>
                <span style="font-size: 0.75rem; font-weight: 800; text-transform: uppercase; background: #6d4ca6; color: #fff; padding: 2px 6px; border-radius: 3px;">${c.badge}</span>
              </div>

              <div class="pipeline-fields-grid">
                <div>
                  <label for="pipeline-title-${c.id}" style="display: block; font-size: 0.75rem; font-weight: 800; color: #555; margin-bottom: 2px;">
                    Catalog Title (Typo & Name Correction)
                  </label>
                  <input type="text" id="pipeline-title-${c.id}" class="tactile-text-input item-title-input" data-id="${c.id}" value="${c.editTitle}" placeholder="Clean title..." style="padding: 6px 10px; font-weight: 700;" />
                </div>
                <div>
                  <label for="pipeline-meta-${c.id}" style="display: block; font-size: 0.75rem; font-weight: 800; color: #555; margin-bottom: 2px;">
                    ${c.mediaType === "book" ? "Author" : "Release Year"}
                  </label>
                  <input type="text" id="pipeline-meta-${c.id}" class="tactile-text-input item-meta-input" data-id="${c.id}" value="${c.editMeta}" placeholder="${c.mediaType === "book" ? "Author..." : "e.g. 2024"}" style="padding: 6px 10px;" />
                </div>
              </div>

              <div class="pipeline-files-details">
                <div style="display: flex; justify-content: space-between; align-items: center;">
                  <button type="button" class="tactile-action-btn secondary item-toggle-files-btn" data-id="${c.id}" style="padding: 3px 8px; font-size: 0.75rem;">
                    ${c.expanded ? "▲ Hide Gathered Files" : `▼ Review ${c.files.length} Gathered Files`}
                  </button>
                  <button type="button" class="tactile-action-btn secondary item-skip-btn" data-id="${c.id}" style="padding: 3px 8px; font-size: 0.75rem; color: #700000;" title="Exclude this item from library">
                    ${c.selected ? "Skip Item" : "Include Item"}
                  </button>
                </div>

                ${
                  c.expanded
                    ? `
                  <div class="pipeline-files-list">
                    ${c.files
                      .map(
                        (f) => `
                      <div class="pipeline-file-row">
                        <div>
                          <strong>${f.title}</strong>
                          <span style="font-family: monospace; color: #666; font-size: 0.75rem; margin-left: 6px;">${f.rawFilename}</span>
                        </div>
                        <span style="font-family: monospace; color: #555;">${f.duration || f.size || "Ready"}</span>
                      </div>
                    `
                      )
                      .join("")}
                  </div>
                `
                    : ""
                }
              </div>
            </div>
          </div>
        </div>
      `
      )
      .join("");

    // Bind item event listeners
    const checkboxes = el.pipelineItemsList.querySelectorAll(".item-select-check");
    checkboxes.forEach((cb) => {
      cb.addEventListener("change", () => {
        const id = parseInt(cb.getAttribute("data-id"));
        const item = pipelineCandidates.find((x) => x.id === id);
        if (item) {
          item.selected = cb.checked;
          const card = el.pipelineItemsList.querySelector(`.pipeline-item-card[data-id="${id}"]`);
          if (card) {
            card.classList.toggle("is-selected", item.selected);
            card.classList.toggle("is-excluded", !item.selected);
          }
          const selCount = pipelineCandidates.filter((x) => x.selected).length;
          const totCount = pipelineCandidates.length;
          el.pipelineSelectAll.checked = totCount > 0 && selCount === totCount;
          el.pipelineSelectAll.indeterminate = selCount > 0 && selCount < totCount;
          el.pipelineSelectedCounter.textContent = `Select All (${selCount} of ${totCount} selected)`;
          el.commitPipelineBtn.textContent = `Ingest Selected Items (${selCount})`;
          el.commitPipelineBtn.disabled = selCount === 0;
        }
      });
    });

    const titleInputs = el.pipelineItemsList.querySelectorAll(".item-title-input");
    titleInputs.forEach((input) => {
      input.addEventListener("input", () => {
        const id = parseInt(input.getAttribute("data-id"));
        const item = pipelineCandidates.find((x) => x.id === id);
        if (item) {
          item.editTitle = input.value;
        }
      });
    });

    const metaInputs = el.pipelineItemsList.querySelectorAll(".item-meta-input");
    metaInputs.forEach((input) => {
      input.addEventListener("input", () => {
        const id = parseInt(input.getAttribute("data-id"));
        const item = pipelineCandidates.find((x) => x.id === id);
        if (item) {
          item.editMeta = input.value;
        }
      });
    });

    const toggleBtns = el.pipelineItemsList.querySelectorAll(".item-toggle-files-btn");
    toggleBtns.forEach((btn) => {
      btn.addEventListener("click", () => {
        const id = parseInt(btn.getAttribute("data-id"));
        const item = pipelineCandidates.find((x) => x.id === id);
        if (item) {
          item.expanded = !item.expanded;
          renderPipelineList();
        }
      });
    });

    const skipBtns = el.pipelineItemsList.querySelectorAll(".item-skip-btn");
    skipBtns.forEach((btn) => {
      btn.addEventListener("click", () => {
        const id = parseInt(btn.getAttribute("data-id"));
        const item = pipelineCandidates.find((x) => x.id === id);
        if (item) {
          item.selected = !item.selected;
          renderPipelineList();
        }
      });
    });
  }

  function commitPipeline() {
    const selected = pipelineCandidates.filter((x) => x.selected);
    if (selected.length === 0) return;

    selected.forEach((c) => {
      if (c.mediaType === "book") {
        FLAN_MOCK_DATA.books.push({
          id: Date.now() + Math.floor(Math.random() * 1000),
          title: c.editTitle || c.proposedTitle,
          author: c.editMeta || c.proposedAuthor || "Unknown",
          overview: c.overview || "",
          coverColor: c.coverColor || "#00695c",
          badge: c.badge || "EPUB",
          metadataLocked: true,
          files: c.files.map((f) => ({
            ...f,
            downloadUrl: f.downloadUrl && f.downloadUrl !== "#" ? f.downloadUrl : `/download/book/${f.id}`
          }))
        });
      } else {
        FLAN_MOCK_DATA.videos.push({
          id: Date.now() + Math.floor(Math.random() * 1000),
          title: c.editTitle || c.proposedTitle,
          releaseYear: parseInt(c.editMeta) || c.proposedYear || null,
          type: c.badge ? c.badge.toLowerCase() : "series",
          overview: c.overview || "",
          coverColor: c.coverColor || "#4a148c",
          badge: c.badge || "Series",
          metadataLocked: true,
          files: c.files.map((f) => ({
            ...f,
            streamUrl: f.streamUrl && f.streamUrl !== "#" ? f.streamUrl : `/stream/video/${f.id}`,
            downloadUrl: f.downloadUrl && f.downloadUrl !== "#" ? f.downloadUrl : `/download/video/${f.id}`
          }))
        });
      }

      // Remove from pending ingestion
      const idx = FLAN_MOCK_DATA.pendingIngestion.findIndex((x) => x.id === c.id);
      if (idx !== -1) {
        FLAN_MOCK_DATA.pendingIngestion.splice(idx, 1);
      }
    });

    el.pipelineFeedback.className = "avatar-feedback success";
    el.pipelineFeedback.textContent = `✓ Successfully ingested ${selected.length} item(s) into your library!`;

    setTimeout(() => {
      closeIngestionPipeline();
      if (state.activeTab === "manage") {
        renderManageView();
      } else if (state.activeTab === "video") {
        renderCatalogView("video");
      } else if (state.activeTab === "books") {
        renderCatalogView("books");
      }
    }, 600);
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
    if (el.topHeader) el.topHeader.classList.remove("hidden");
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

    el.sidebarVideoBtn.removeAttribute("aria-current");
    el.sidebarBooksBtn.removeAttribute("aria-current");
    el.sidebarManageBtn.removeAttribute("aria-current");

    if (state.activeTab === "video") {
      el.sidebarVideoBtn.classList.add("active");
      el.sidebarVideoBtn.setAttribute("aria-current", "page");
    } else if (state.activeTab === "books") {
      el.sidebarBooksBtn.classList.add("active");
      el.sidebarBooksBtn.setAttribute("aria-current", "page");
    } else if (state.activeTab === "manage") {
      el.sidebarManageBtn.classList.add("active");
      el.sidebarManageBtn.setAttribute("aria-current", "page");
    }
  }

  // =========================================================================
  // VIEW 1: 2-STEP SEQUENTIAL LOGIN (Step 1: Profile Picker -> Step 2: PIN)
  // =========================================================================
  function renderLoginScreen(step = "picker", targetUser = null) {
    if (el.topHeader) el.topHeader.classList.add("hidden");
    el.leftSidebar.classList.add("hidden");
    el.mainContent.style.padding = "0";

    if (step === "pin" && targetUser) {
      renderLoginPinState(targetUser);
    } else {
      renderLoginPickerState();
    }
  }

  // STEP 1: Profile Picker State ("Who is watching?")
  function renderLoginPickerState() {
    el.mainContent.innerHTML = `
      <div class="login-pad-viewport">
        <div class="login-pad-card login-picker-card">
          <div class="login-pad-header">
            <h1 class="login-pad-title">Flan Media Server</h1>
            <p class="login-pad-subtitle">Who is watching?</p>
          </div>

          <!-- Household Profile Cards -->
          <div class="login-profiles-grid" role="group" aria-label="Household Profiles">
            ${FLAN_MOCK_DATA.users
              .map(
                (u) => `
                <button
                  type="button"
                  class="login-profile-card"
                  data-username="${u.username}"
                  title="Select ${u.username}"
                >
                  <div class="login-profile-avatar-frame">
                    ${getUserAvatarSVG(u)}
                  </div>
                  <span class="login-profile-username">${u.username}</span>
                  <span class="login-profile-role-pill">${u.role}</span>
                </button>
              `
              )
              .join("")}
          </div>
        </div>
      </div>
    `;

    const profileCards = document.querySelectorAll(".login-profile-card");
    profileCards.forEach((card) => {
      card.addEventListener("click", () => {
        const username = card.dataset.username;
        const found = FLAN_MOCK_DATA.users.find((u) => u.username === username);
        if (!found) return;
        renderLoginScreen("pin", found);
      });
    });

    if (profileCards.length > 0) {
      profileCards[0].focus();
    }
  }

  // STEP 2: Dedicated PIN Entry State
  function renderLoginPinState(user) {
    el.mainContent.innerHTML = `
      <div class="login-pad-viewport">
        <div class="login-pad-card login-pin-card">
          <!-- Active User Badge Header -->
          <div class="login-user-banner">
            <div class="login-user-avatar-badge">
              ${getUserAvatarSVG(user)}
            </div>
            <div class="login-user-info">
              <h2 class="login-user-name">${user.username}</h2>
              <span class="login-profile-role-pill">${user.role}</span>
            </div>
          </div>

          <p class="login-pad-subtitle">Enter your 4-digit PIN</p>

          <!-- Tactile PIN Access Form -->
          <form id="login-form" class="login-pin-form">
            <div class="form-field-group">
              <input
                id="pin-input"
                class="tactile-pin-input"
                type="password"
                maxlength="6"
                inputmode="numeric"
                placeholder="••••"
                autofocus
                autocomplete="current-password"
                aria-label="Enter PIN for ${user.username}"
              />
            </div>

            <button type="submit" id="access-btn" class="tactile-access-btn">
              Access Library →
            </button>

            <button type="button" id="switch-profile-btn" class="tactile-switch-btn">
              ← Switch Profile
            </button>

            <div id="login-feedback" class="login-feedback" role="alert" aria-live="assertive"></div>
          </form>
        </div>
      </div>
    `;

    const form = document.getElementById("login-form");
    const pinInput = document.getElementById("pin-input");
    const switchBtn = document.getElementById("switch-profile-btn");
    const feedback = document.getElementById("login-feedback");

    pinInput.focus();

    switchBtn.addEventListener("click", () => {
      renderLoginScreen("picker");
    });

    form.addEventListener("submit", (e) => {
      e.preventDefault();
      const enteredPin = pinInput.value.trim();

      if (user && user.pin === enteredPin) {
        state.currentUser = user;
        renderUserAvatarBadge(user);
        feedback.textContent = "";
        window.location.hash = "#video";
      } else {
        feedback.className = "login-feedback error";
        feedback.textContent = `Invalid PIN. Try '1234' for mike or '0000' for wesley.`;
        pinInput.classList.add("shake");
        pinInput.value = "";
        setTimeout(() => pinInput.classList.remove("shake"), 400);
      }
    });
  }

  // =========================================================================
  // VIEW 2 & 3: MEDIA CATALOG WITH DYNAMIC HEADER & 7-COL PAGINATED GRID
  // =========================================================================
  function renderCatalogView(type) {
    el.mainContent.style.padding = "24px 32px";

    const allItems =
      type === "video" ? FLAN_MOCK_DATA.videos : FLAN_MOCK_DATA.books;

    // 1. FILTERING
    let filtered = allItems.filter((item) => {
      // Search Query
      if (state.searchQuery.trim().length > 0) {
        const q = state.searchQuery.toLowerCase();
        const matchesTitle = item.title.toLowerCase().includes(q);
        const matchesAuthor = item.author && item.author.toLowerCase().includes(q);
        const matchesOverview = item.overview && item.overview.toLowerCase().includes(q);
        if (!matchesTitle && !matchesAuthor && !matchesOverview) return false;
      }

      // Format Filter
      if (state.selectedFormat !== "all") {
        if (type === "video") {
          if (state.selectedFormat === "movies" && item.type !== "movie") return false;
          if (state.selectedFormat === "series" && item.type !== "series") return false;
        } else {
          if (state.selectedFormat === "epub" && item.badge.toLowerCase() !== "epub") return false;
          if (state.selectedFormat === "pdf" && item.badge.toLowerCase() !== "pdf") return false;
        }
      }

      // Status Filter
      if (state.selectedStatus !== "all") {
        if (type === "video") {
          const inProg = item.files && item.files.some((f) => f.progress > 0 && !f.isFinished);
          if (state.selectedStatus === "in-progress" && !inProg && !item.lastWatched) return false;
          if (state.selectedStatus === "unwatched" && (inProg || item.lastWatched)) return false;
        } else {
          if (state.selectedStatus === "in-progress" && item.status !== "reading") return false;
          if (state.selectedStatus === "unwatched" && item.status !== "unread") return false;
        }
      }

      // Source Filter
      if (state.selectedSource !== "all") {
        if (item.sourceId !== Number(state.selectedSource)) return false;
      }

      return true;
    });

    // 2. SORTING
    let sorted = [...filtered];
    if (state.sortOrder === "az") {
      sorted.sort((a, b) => a.title.localeCompare(b.title));
    } else if (state.sortOrder === "year") {
      sorted.sort((a, b) => (b.releaseYear || 0) - (a.releaseYear || 0));
    } else if (state.sortOrder === "recent") {
      sorted.sort((a, b) => {
        const timeA = a.lastWatched || a.lastRead || 0;
        const timeB = b.lastWatched || b.lastRead || 0;
        if (timeA && timeB) return new Date(timeB) - new Date(timeA);
        if (timeA) return -1;
        if (timeB) return 1;
        return b.id - a.id;
      });
    } else if (state.sortOrder === "random") {
      // Deterministic / seeded shuffle based on item id and fixed hash so pagination stays stable across pages
      sorted.sort((a, b) => {
        const hashA = (a.id * 997 + a.title.length * 37) % 1000;
        const hashB = (b.id * 997 + b.title.length * 37) % 1000;
        return hashA - hashB;
      });
    }

    // 3. CONTINUE WATCHING / JUMP BACK IN SHELF (In-progress items sorted by date)
    let continueShelfHtml = "";
    const inProgressList = allItems.filter((item) => {
      if (type === "video") {
        return (item.files && item.files.some((f) => f.progress > 0 && !f.isFinished)) || !!item.lastWatched;
      } else {
        return item.status === "reading" || !!item.lastRead;
      }
    });

    // Sort in-progress items by last activity date descending
    inProgressList.sort((a, b) => {
      const dateA = new Date(a.lastWatched || a.lastRead || "1970-01-01");
      const dateB = new Date(b.lastWatched || b.lastRead || "1970-01-01");
      return dateB - dateA;
    });

    // If there are in-progress items AND we aren't searching for a specific query that excludes them
    if (inProgressList.length > 0 && state.searchQuery.trim().length === 0 && state.selectedStatus !== "unwatched") {
      const shelfTitle = type === "video" ? "Continue Watching" : "Jump Back In";
      const shelfBadgeText = `${inProgressList.length} IN PROGRESS`;

      const cardsHtml = inProgressList.map((item) => {
        if (type === "video") {
          const activeFile = (item.files && item.files.find((f) => f.progress > 0 && !f.isFinished)) || (item.files && item.files[0]) || {};
          const isSeries = item.type === "series";
          const progressPercent = activeFile.progress || 0;
          const progressLabel = isSeries
            ? `${activeFile.title || 'Next Episode'} - ${activeFile.formattedPos || activeFile.duration || ''}`
            : `${activeFile.formattedPos || ''} / ${activeFile.duration || ''} (${progressPercent}%)`;

          return `
            <div class="continue-card" role="article" aria-label="Resume ${item.title}">
              <div class="continue-card-thumb">
                <svg viewBox="0 0 24 24"><polygon points="5,3 19,12 5,21" fill="none" stroke="currentColor" stroke-width="2"/></svg>
              </div>
              <div class="continue-card-content">
                <div>
                  <div class="continue-card-title" title="${item.title}">${item.title}</div>
                  <div class="continue-card-subtitle" title="${progressLabel}">${progressLabel}</div>
                </div>
                <div>
                  <div class="continue-progress-wrap" aria-label="${progressPercent}% watched">
                    <div class="continue-progress-bar" style="width: ${progressPercent}%;"></div>
                  </div>
                  <button type="button" class="continue-resume-btn resume-video-btn" data-file-id="${activeFile.id}">
                    <svg viewBox="0 0 24 24" width="12" height="12" fill="currentColor" stroke="none"><polygon points="5,3 19,12 5,21"/></svg>
                    <span>Resume</span>
                  </button>
                </div>
              </div>
            </div>
          `;
        } else {
          // Discrete reading status, NO raw percentages for books!
          const readingLabel = item.readingProgress || "Currently Reading";
          return `
            <div class="continue-card" role="article" aria-label="Continue reading ${item.title}">
              <div class="continue-card-thumb">
                <svg viewBox="0 0 24 24"><path d="M4 19.5A2.5 2.5 0 0 1 6.5 17H20"/><path d="M6.5 2H20v20H6.5A2.5 2.5 0 0 1 4 19.5v-15A2.5 2.5 0 0 1 6.5 2z" fill="none" stroke="currentColor" stroke-width="2"/></svg>
              </div>
              <div class="continue-card-content">
                <div>
                  <div class="continue-card-title" title="${item.title}">${item.title}</div>
                  <div class="continue-card-subtitle">${item.author || ''}</div>
                </div>
                <div>
                  <span style="display:inline-block; font-size:0.75rem; font-weight:800; background:#e2d9f3; color:#4a148c; padding:2px 8px; border-radius:3px; border:1px solid #000; margin-top:2px;">
                    ${readingLabel}
                  </span>
                  <div>
                    <button type="button" class="continue-resume-btn resume-book-btn" data-book-id="${item.id}" style="margin-top:6px;">
                      <svg viewBox="0 0 24 24" width="12" height="12" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M2 3h6a4 4 0 0 1 4 4v14a3 3 0 0 0-3-3H2z"/><path d="M22 3h-6a4 4 0 0 0-4 4v14a3 3 0 0 1 3-3h7z"/></svg>
                      <span>Continue</span>
                    </button>
                  </div>
                </div>
              </div>
            </div>
          `;
        }
      }).join("");

      continueShelfHtml = `
        <section class="continue-shelf" aria-label="${shelfTitle}">
          <div class="continue-shelf-header">
            <h2 class="continue-shelf-title">${shelfTitle}</h2>
            <span class="continue-shelf-badge">${shelfBadgeText}</span>
          </div>
          <div class="continue-shelf-grid">
            ${cardsHtml}
          </div>
        </section>
      `;
    }

    // 4. PAGINATION: 21 ITEMS MAX PER PAGE (7x3)
    const ITEMS_PER_PAGE = 21;
    const totalItems = sorted.length;
    const totalPages = Math.max(1, Math.ceil(totalItems / ITEMS_PER_PAGE));
    if (state.currentPage > totalPages) state.currentPage = totalPages;
    if (state.currentPage < 1) state.currentPage = 1;

    const pageStartIndex = (state.currentPage - 1) * ITEMS_PER_PAGE;
    const pageItems = sorted.slice(pageStartIndex, pageStartIndex + ITEMS_PER_PAGE);

    // 5. STATUS STRIP
    const isSearching = state.searchQuery.trim().length > 0;
    const hasActiveFilters = state.selectedFormat !== "all" || state.selectedStatus !== "all" || state.selectedSource !== "all" || state.sortOrder !== "random";
    let statusStripHtml = "";
    if (isSearching || hasActiveFilters) {
      statusStripHtml = `
        <div class="catalog-status-strip filtered">
          <span class="status-indicator-tag">FILTERED</span>
          ${isSearching ? `<span class="status-query-text">"${state.searchQuery}"</span>` : ""}
          <span class="status-match-count">(${totalItems} title${totalItems === 1 ? "" : "s"})</span>
          <button id="clear-filters-btn" class="tactile-clear-btn" title="Reset Filters & Search">
            [x] Reset All
          </button>
        </div>
      `;
    } else {
      const typeLabel = type === "video" ? "ALL VIDEOS & SERIES" : "BOOKS & PUBLICATIONS";
      statusStripHtml = `
        <div class="catalog-status-strip">
          <span class="status-section-name">${typeLabel}</span>
          <span class="status-divider">-</span>
          <span class="status-total-count">${totalItems} TITLES</span>
        </div>
      `;
    }

    // 6. COLLAPSIBLE FILTER DRAWER HTML
    let filterDrawerHtml = "";
    if (state.filterExpanded) {
      const sourcesForType = FLAN_MOCK_DATA.storageSources.filter((s) => s.mediaType === (type === "video" ? "video" : "book"));

      filterDrawerHtml = `
        <div id="filter-drawer" class="filter-drawer" role="region" aria-label="Catalog Filters">
          <!-- Format Row -->
          <div class="filter-drawer-row">
            <span class="filter-drawer-label">Format:</span>
            <div class="filter-pills-group" role="group" aria-label="Format filter">
              <button class="filter-pill ${state.selectedFormat === 'all' ? 'active' : ''}" data-filter-type="format" data-val="all">All</button>
              ${type === 'video' ? `
                <button class="filter-pill ${state.selectedFormat === 'movies' ? 'active' : ''}" data-filter-type="format" data-val="movies">Movies</button>
                <button class="filter-pill ${state.selectedFormat === 'series' ? 'active' : ''}" data-filter-type="format" data-val="series">Series</button>
              ` : `
                <button class="filter-pill ${state.selectedFormat === 'epub' ? 'active' : ''}" data-filter-type="format" data-val="epub">EPUB</button>
                <button class="filter-pill ${state.selectedFormat === 'pdf' ? 'active' : ''}" data-filter-type="format" data-val="pdf">PDF</button>
              `}
            </div>
          </div>

          <!-- Status Row -->
          <div class="filter-drawer-row">
            <span class="filter-drawer-label">Status:</span>
            <div class="filter-pills-group" role="group" aria-label="Status filter">
              <button class="filter-pill ${state.selectedStatus === 'all' ? 'active' : ''}" data-filter-type="status" data-val="all">All</button>
              <button class="filter-pill ${state.selectedStatus === 'in-progress' ? 'active' : ''}" data-filter-type="status" data-val="in-progress">In Progress</button>
              <button class="filter-pill ${state.selectedStatus === 'unwatched' ? 'active' : ''}" data-filter-type="status" data-val="unwatched">${type === 'video' ? 'Unwatched' : 'Unread'}</button>
            </div>
          </div>

          <!-- Source Drive Row -->
          ${sourcesForType.length > 0 ? `
            <div class="filter-drawer-row">
              <span class="filter-drawer-label">Drive:</span>
              <div class="filter-pills-group" role="group" aria-label="Storage source filter">
                <button class="filter-pill ${state.selectedSource === 'all' ? 'active' : ''}" data-filter-type="source" data-val="all">All Drives</button>
                ${sourcesForType.map((src) => `
                  <button class="filter-pill ${String(state.selectedSource) === String(src.id) ? 'active' : ''}" data-filter-type="source" data-val="${src.id}">${src.name}</button>
                `).join('')}
              </div>
            </div>
          ` : ''}

          <!-- Sort Row -->
          <div class="filter-drawer-row">
            <span class="filter-drawer-label">Sort:</span>
            <div class="filter-pills-group" role="group" aria-label="Sort order">
              <button class="filter-pill ${state.sortOrder === 'random' ? 'active' : ''}" data-filter-type="sort" data-val="random">Shuffle / Random</button>
              <button class="filter-pill ${state.sortOrder === 'recent' ? 'active' : ''}" data-filter-type="sort" data-val="recent">Recently Added</button>
              <button class="filter-pill ${state.sortOrder === 'az' ? 'active' : ''}" data-filter-type="sort" data-val="az">Title (A-Z)</button>
              <button class="filter-pill ${state.sortOrder === 'year' ? 'active' : ''}" data-filter-type="sort" data-val="year">Release Year</button>
            </div>
            ${hasActiveFilters ? `
              <button id="drawer-reset-btn" type="button" class="filter-reset-btn">Reset Filters</button>
            ` : ''}
          </div>
        </div>
      `;
    }

    // 7. PAGINATION BAR HTML
    let paginationHtml = "";
    if (totalPages > 1) {
      let pageButtonsHtml = "";
      for (let p = 1; p <= totalPages; p++) {
        pageButtonsHtml += `
          <button class="pagination-btn ${p === state.currentPage ? 'active' : ''}" data-page="${p}" ${p === state.currentPage ? 'aria-current="page"' : ''}>
            ${p}
          </button>
        `;
      }

      const showingEnd = Math.min(pageStartIndex + ITEMS_PER_PAGE, totalItems);
      paginationHtml = `
        <nav class="pagination-bar" aria-label="Catalog Pagination">
          <button class="pagination-btn prev-page-btn" data-page="${state.currentPage - 1}" ${state.currentPage === 1 ? 'disabled aria-disabled="true"' : ''}>
            &larr; Prev
          </button>
          <div class="pagination-pages">
            ${pageButtonsHtml}
          </div>
          <button class="pagination-btn next-page-btn" data-page="${state.currentPage + 1}" ${state.currentPage === totalPages ? 'disabled aria-disabled="true"' : ''}>
            Next &rarr;
          </button>
          <span class="pagination-info">Showing ${pageStartIndex + 1}-${showingEnd} of ${totalItems}</span>
        </nav>
      `;
    }

    // 8. RENDER INTO MAIN CONTENT
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
          <button id="catalog-search-btn" class="search-btn" aria-label="Search ${type === 'video' ? 'videos' : 'books'}" title="Search">
            <svg viewBox="0 0 24 24">
              <circle cx="11" cy="11" r="7" />
              <line x1="16.5" y1="16.5" x2="22" y2="22" />
            </svg>
          </button>
          <button id="filter-toggle-btn" class="filter-toggle-btn ${state.filterExpanded ? 'active' : ''}" aria-expanded="${state.filterExpanded}" aria-controls="filter-drawer" title="Toggle Filters">
            <span>Filter</span>
            <span style="font-size: 0.8rem;">${state.filterExpanded ? '[^]' : '[v]'}</span>
          </button>
        </div>

        <!-- Collapsible Filter Drawer -->
        ${filterDrawerHtml}

        <!-- Compact Tactile Status Strip -->
        ${statusStripHtml}
      </section>

      <!-- Continue Watching / Jump Back In Shelf -->
      ${continueShelfHtml}

      <!-- Catalog Media Grid (7 Columns Max Hard Limit) -->
      <section class="media-grid">
        ${
          pageItems.length === 0
            ? `<div style="grid-column: 1/-1; text-align:center; padding: 40px; font-weight:700; font-size:1.1rem; color: #555;">No media matching your filters.</div>`
            : pageItems
                .map((item) => {
                  return `
            <a class="media-card" href="#detail/${type}/${item.id}" data-id="${item.id}" data-type="${type}" role="article" aria-label="${item.title}${item.badge ? ' (' + item.badge + ')' : ''}">
              <!-- Upper Poster Area -->
              <div class="card-poster">
                <span class="card-badge">${item.badge}</span>
                <div class="card-poster-placeholder">
                  <svg viewBox="0 0 24 24">
                    ${
                      type === "video"
                        ? '<polygon points="5,3 19,12 5,21" fill="none" stroke="currentColor" stroke-width="2"/>'
                        : '<path d="M4 19.5A2.5 2.5 0 0 1 6.5 17H20"/> <path d="M6.5 2H20v20H6.5A2.5 2.5 0 0 1 4 19.5v-15A2.5 2.5 0 0 1 6.5 2z" fill="none" stroke="currentColor" stroke-width="2"/>'
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
            </a>
          `;
                })
                .join("")
        }
      </section>

      <!-- Pagination Bar -->
      ${paginationHtml}
    `;

    // 9. EVENT BINDINGS
    const searchInput = document.getElementById("catalog-search");
    if (searchInput) {
      searchInput.addEventListener("input", (e) => {
        state.searchQuery = e.target.value;
        state.currentPage = 1;
        renderCatalogView(type);
      });
    }

    const filterToggleBtn = document.getElementById("filter-toggle-btn");
    if (filterToggleBtn) {
      filterToggleBtn.addEventListener("click", () => {
        state.filterExpanded = !state.filterExpanded;
        renderCatalogView(type);
      });
    }

    // Filter pills in drawer
    const filterPills = el.mainContent.querySelectorAll(".filter-pill");
    filterPills.forEach((pill) => {
      pill.addEventListener("click", () => {
        const filterType = pill.getAttribute("data-filter-type");
        const val = pill.getAttribute("data-val");

        if (filterType === "format") {
          state.selectedFormat = val;
        } else if (filterType === "status") {
          state.selectedStatus = val;
        } else if (filterType === "source") {
          state.selectedSource = val;
        } else if (filterType === "sort") {
          state.sortOrder = val;
        }
        state.currentPage = 1;
        renderCatalogView(type);
      });
    });

    const resetFiltersBtn = document.getElementById("clear-filters-btn") || document.getElementById("drawer-reset-btn");
    if (resetFiltersBtn) {
      resetFiltersBtn.addEventListener("click", () => {
        state.searchQuery = "";
        state.selectedFormat = "all";
        state.selectedStatus = "all";
        state.selectedSource = "all";
        state.sortOrder = "random";
        state.currentPage = 1;
        renderCatalogView(type);
      });
    }

    const drawerResetBtn = document.getElementById("drawer-reset-btn");
    if (drawerResetBtn && drawerResetBtn !== resetFiltersBtn) {
      drawerResetBtn.addEventListener("click", () => {
        state.selectedFormat = "all";
        state.selectedStatus = "all";
        state.selectedSource = "all";
        state.sortOrder = "random";
        state.currentPage = 1;
        renderCatalogView(type);
      });
    }

    // Resume video buttons
    const resumeVideoBtns = el.mainContent.querySelectorAll(".resume-video-btn");
    resumeVideoBtns.forEach((btn) => {
      btn.addEventListener("click", (e) => {
        e.stopPropagation();
        const fileId = btn.getAttribute("data-file-id");
        window.location.hash = `#watch/${fileId}`;
      });
    });

    // Resume book buttons
    const resumeBookBtns = el.mainContent.querySelectorAll(".resume-book-btn");
    resumeBookBtns.forEach((btn) => {
      btn.addEventListener("click", (e) => {
        e.stopPropagation();
        const bookId = btn.getAttribute("data-book-id");
        window.location.hash = `#detail/books/${bookId}`;
      });
    });

    // Pagination buttons
    const pageBtns = el.mainContent.querySelectorAll(".pagination-btn");
    pageBtns.forEach((btn) => {
      btn.addEventListener("click", () => {
        const targetPage = parseInt(btn.getAttribute("data-page"), 10);
        if (!isNaN(targetPage) && targetPage >= 1 && targetPage <= totalPages && targetPage !== state.currentPage) {
          state.currentPage = targetPage;
          renderCatalogView(type);
          el.mainContent.scrollTo({ top: 0, behavior: "smooth" });
        }
      });
    });

    // Bind Card Click -> Details View
    const cards = el.mainContent.querySelectorAll(".media-card");
    cards.forEach((card) => {
      card.addEventListener("click", (e) => {
        e.preventDefault();
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
            <div style="display: flex; justify-content: space-between; align-items: flex-start; gap: 12px; margin-bottom: 6px;">
              <h1 class="detail-title" style="margin: 0;">${item.title} ${item.releaseYear ? `<span style="font-size: 1.1rem; opacity: 0.7;">(${item.releaseYear})</span>` : ""}</h1>
              ${
                state.currentUser && state.currentUser.role === "admin" && type === "video"
                  ? `<button id="edit-media-btn" class="tactile-action-btn secondary" style="padding: 6px 14px; font-size: 0.85rem; white-space: nowrap; display: inline-flex; align-items: center; gap: 6px;">
                      <svg style="width:14px; height:14px;" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path></svg>
                      Edit Details
                    </button>`
                  : ""
              }
            </div>
            ${item.author ? `<div style="font-weight: 600; color: #555;">Author: ${item.author}</div>` : ""}
            <p class="detail-synopsis">${item.overview}</p>
          </div>
        </div>

        <h2 style="font-size: 1.3rem; font-weight: 800; margin-top: 10px;">
          ${type === "video" ? "PLAYABLE EPISODES & CUTS" : "AVAILABLE VOLUMES & EDITIONS"}
        </h2>

        <div class="playable-items-list">
          ${item.files
            .filter((f) => !f.isHidden)
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
                      <button class="tactile-action-btn play-file-btn" data-file-id="${file.id}" aria-label="Play ${file.title}">
                        ▶ Play
                      </button>
                      <button class="tactile-action-btn vlc download-vlc-btn" data-url="${file.downloadUrl}" data-title="${file.title}" aria-label="Download or stream in VLC: ${file.title}" title="Direct download or stream in VLC">
                        ⬇ VLC / Download
                      </button>
                    `
                      : `
                      ${
                        file.format === "pdf"
                          ? `<button class="tactile-action-btn read-pdf-btn" data-url="${file.readUrl}" data-title="${file.title}" aria-label="Open PDF: ${file.title}">Open PDF</button>`
                          : ""
                      }
                      <button class="tactile-action-btn download download-book-btn" data-url="${file.downloadUrl}" data-title="${file.title}" aria-label="Download ${file.format.toUpperCase()}: ${file.title}">
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

    const editMediaBtn = document.getElementById("edit-media-btn");
    if (editMediaBtn) {
      editMediaBtn.addEventListener("click", () => {
        openMediaEditModal(item);
      });
    }

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

    el.mainContent.innerHTML = `
      <div class="manage-container">
        <h1 style="font-size: 1.8rem; font-weight: 900; letter-spacing: -0.5px;">Manage Server</h1>

        <!-- Media Storage Sources Card -->
        <div class="manage-card">
          <div class="manage-card-header">
            <h2>Media Storage Sources</h2>
            <div class="manage-header-actions">
              <button id="add-source-btn" class="tactile-action-btn secondary" style="padding: 6px 14px; font-size: 0.85rem;">
                + Add Storage Folder
              </button>
              <button id="open-upload-btn" class="tactile-action-btn" style="padding: 6px 14px; font-size: 0.85rem;">
                Upload Files
              </button>
            </div>
          </div>

          <div style="display: flex; flex-direction: column; gap: 10px; margin-top: 12px;">
            ${FLAN_MOCK_DATA.storageSources
              .map(
                (s) => `
              <div class="storage-source-item">
                <div class="storage-source-info">
                  <div class="storage-source-title-row">
                    <span class="storage-source-name">${s.name}</span>
                    <span class="storage-source-badge">${s.mediaType}</span>
                  </div>
                  <div class="storage-source-path-row">
                    ${s.path} • <span style="font-weight: 700; color: #2e7d32;">${s.freeSpace} free</span>
                  </div>
                </div>
                <div class="storage-source-actions">
                  <button class="tactile-action-btn scan-source-btn" data-source-id="${s.id}" title="Scan storage folder" style="padding: 6px 14px; font-size: 0.85rem;">
                    Scan
                  </button>
                  <button class="tactile-action-btn secondary edit-source-btn" data-source-id="${s.id}" title="Edit Storage Folder Name" style="padding: 6px 12px; font-size: 0.85rem;">
                    Edit
                  </button>
                  <button class="tactile-action-btn secondary delete-source-btn" data-source-id="${s.id}" title="Remove Storage Source" style="padding: 6px 10px; font-size: 0.8rem;">
                    Remove
                  </button>
                </div>
              </div>
            `
              )
              .join("")}
          </div>

          <div style="display: flex; gap: 12px; margin-top: 14px; align-items: center;">
            <button id="rescan-btn" class="tactile-action-btn" style="padding: 10px 20px;">
              ⟳ Rescan All Storage Sources
            </button>
            <div id="rescan-feedback" role="alert" aria-live="polite" style="font-weight: 700; color: #2e7d32; font-size: 0.9rem;"></div>
          </div>
        </div>

        <div class="manage-card">
          <div class="manage-card-header">
            <h2>Household Profiles</h2>
            <div class="manage-header-actions">
              <button id="add-user-btn" class="tactile-action-btn" style="padding: 6px 14px; font-size: 0.85rem;">
                + Add New User
              </button>
            </div>
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
              <div class="household-user-item">
                <div class="household-user-info">
                  ${avatarPreview}
                  <div style="min-width: 0;">
                    <div style="font-weight: 700; font-size: 1.05rem;">
                      ${u.username} <span style="font-size: 0.8rem; color: #666;">(${u.role})</span>
                    </div>
                    <span style="font-family: monospace; font-size: 0.85rem; color: #555;">PIN: ••••</span>
                  </div>
                </div>
                <div class="household-user-actions">
                  <button class="tactile-action-btn secondary edit-user-btn" data-username="${u.username}" aria-label="Edit Profile for ${u.username}" style="padding: 6px 12px; font-size: 0.85rem; display: inline-flex; align-items: center; gap: 6px;">
                    <svg style="width:13px; height:13px;" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path></svg>
                    Edit Profile
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

    // Storage source buttons
    const addSourceBtn = document.getElementById("add-source-btn");
    if (addSourceBtn) addSourceBtn.addEventListener("click", openAddSourceModal);

    const openUploadBtn = document.getElementById("open-upload-btn");
    if (openUploadBtn) openUploadBtn.addEventListener("click", openUploadModal);

    const scanSourceBtns = el.mainContent.querySelectorAll(".scan-source-btn");
    scanSourceBtns.forEach((btn) => {
      btn.addEventListener("click", () => {
        const id = parseInt(btn.getAttribute("data-source-id"));
        const source = FLAN_MOCK_DATA.storageSources.find((s) => s.id === id);
        if (source) {
          openIngestionPipeline(source);
        }
      });
    });

    const editSourceBtns = el.mainContent.querySelectorAll(".edit-source-btn");
    editSourceBtns.forEach((btn) => {
      btn.addEventListener("click", () => {
        const id = parseInt(btn.getAttribute("data-source-id"));
        const source = FLAN_MOCK_DATA.storageSources.find((s) => s.id === id);
        if (source) {
          openEditSourceModal(source);
        }
      });
    });

    const deleteSourceBtns = el.mainContent.querySelectorAll(".delete-source-btn");
    deleteSourceBtns.forEach((btn) => {
      btn.addEventListener("click", () => {
        const id = parseInt(btn.getAttribute("data-source-id"));
        const source = FLAN_MOCK_DATA.storageSources.find((s) => s.id === id);
        if (source) {
          openRemoveSourceModal(source);
        }
      });
    });

    // Rescan button
    const rescanBtn = document.getElementById("rescan-btn");
    const feedback = document.getElementById("rescan-feedback");
    if (rescanBtn) {
      rescanBtn.addEventListener("click", () => {
        rescanBtn.textContent = "Scanning...";
        feedback.textContent = "";
        setTimeout(() => {
          rescanBtn.textContent = "⟳ Rescan All Storage Sources";
          if (FLAN_MOCK_DATA.pendingIngestion.length > 0) {
            openIngestionPipeline(null);
          } else {
            feedback.textContent = "✓ Library scan complete: All active storage sources verified and indexed.";
          }
        }, 400);
      });
    }

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
                <strong>Fast Directory Rescanning:</strong> Whenever you add or organize files, go to <strong>Manage Server</strong> and click 
                <code>[ ⟳ Rescan All Media ]</code>. Flan performs a fast directory walk with batch SQLite transactions and in-memory reconciliation.
              </p>
            </section>

            <!-- Section 2 -->
            <section id="sec-playback" class="manual-section-card">
              <div class="manual-section-header">
                <span class="manual-section-num">2</span>
                <h2 class="manual-section-title">Direct Video Playback & VLC Fallback</h2>
              </div>
              <p class="manual-body-text">
                Flan is built for low-power devices, repurposed PCs, and single-board computers (Orange Pi, Raspberry Pi). 
                To ensure maximum responsiveness and zero CPU strain, <strong>Flan never performs server-side video transcoding</strong>.
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

# Run on custom port
PORT=8080 ./flan</code></pre>
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
                  <td class="spec-label">Resource Overhead</td>
                  <td class="spec-val">Lean &amp; Unconstrained (live stats via <code>runtime.ReadMemStats</code>)</td>
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

    // Title Horizontal Scrolling on Hover & Focus (Marquee)
    function handleMarqueeOverflow(card) {
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
    }

    document.addEventListener("mouseover", (e) => {
      handleMarqueeOverflow(e.target.closest(".media-card"));
    });

    document.addEventListener("focusin", (e) => {
      handleMarqueeOverflow(e.target.closest(".media-card"));
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

    document.addEventListener("focusout", (e) => {
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

    // Brand title click (top header) -> Go to #video
    const brandTriggers = new Set([el.brandTitle, ...document.querySelectorAll(".brand-trigger")].filter(Boolean));
    brandTriggers.forEach((btn) => {
      btn.addEventListener("click", () => {
        window.location.hash = "#video";
      });
      btn.addEventListener("keydown", (e) => {
        if (e.key === "Enter" || e.key === " ") {
          e.preventDefault();
          window.location.hash = "#video";
        }
      });
    });

    // Avatar click (top header) -> Opens Unified My Profile Modal
    const avatarTriggers = new Set([el.userAvatar, ...document.querySelectorAll(".user-avatar-trigger")].filter(Boolean));
    avatarTriggers.forEach((btn) => {
      btn.addEventListener("click", () => {
        openProfileModal(state.currentUser);
      });
    });

    // Profile Modal Listeners
    el.closeProfileBtn.addEventListener("click", closeProfileModal);
    el.profileModal.addEventListener("click", (e) => {
      if (e.target === el.profileModal) {
        closeProfileModal();
      }
    });
    el.profileModal.addEventListener("cancel", () => {
      profileTargetUser = null;
      if (lastFocusedElement) {
        lastFocusedElement.focus();
        lastFocusedElement = null;
      }
    });

    el.addUserModal.addEventListener("cancel", () => {
      if (lastFocusedElement) {
        lastFocusedElement.focus();
        lastFocusedElement = null;
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

    // Media Edit Form Submit & Close
    el.closeMediaEditBtn.addEventListener("click", closeMediaEditModal);
    if (el.cancelMediaEditBtn) el.cancelMediaEditBtn.addEventListener("click", closeMediaEditModal);
    el.mediaEditModal.addEventListener("click", (e) => {
      if (e.target === el.mediaEditModal) closeMediaEditModal();
    });

    el.mediaEditForm.addEventListener("submit", (e) => {
      e.preventDefault();
      if (!currentEditingItem) return;

      const newTitle = el.editVideoTitle.value.trim();
      const newYear = parseInt(el.editVideoYear.value) || null;
      const newType = el.editVideoType.value;
      const newOverview = el.editVideoOverview.value.trim();

      if (!newTitle) {
        el.mediaEditFeedback.className = "avatar-feedback error";
        el.mediaEditFeedback.textContent = "Title cannot be empty.";
        return;
      }

      currentEditingItem.title = newTitle;
      currentEditingItem.proposedTitle = newTitle;
      currentEditingItem.releaseYear = newYear;
      currentEditingItem.proposedYear = newYear;
      currentEditingItem.type = newType;
      currentEditingItem.overview = newOverview;
      currentEditingItem.metadataLocked = true;

      // Update episode titles & hidden states from table inputs
      const titleInputs = el.editFilesTableWrap.querySelectorAll(".file-title-edit");
      titleInputs.forEach((input) => {
        const fId = parseInt(input.getAttribute("data-file-id"));
        const file = currentEditingItem.files.find((f) => f.id === fId);
        if (file) {
          file.title = input.value.trim() || file.title;
        }
      });

      const hideChecks = el.editFilesTableWrap.querySelectorAll(".file-hide-check");
      hideChecks.forEach((check) => {
        const fId = parseInt(check.getAttribute("data-file-id"));
        const file = currentEditingItem.files.find((f) => f.id === fId);
        if (file) {
          file.isHidden = check.checked;
        }
      });

      el.mediaEditFeedback.className = "avatar-feedback success";
      el.mediaEditFeedback.textContent = "✓ Changes saved and locked against rescans!";

      setTimeout(() => {
        closeMediaEditModal();
        if (window.location.hash.startsWith("#detail/video/")) {
          renderDetailView("video", currentEditingItem.id);
        } else if (state.activeTab === "manage") {
          renderManageView();
        } else if (state.activeTab === "video") {
          renderCatalogView("video");
        }
      }, 500);
    });

    // Add Source Form Submit & Close
    el.closeAddSourceBtn.addEventListener("click", closeAddSourceModal);
    el.addSourceModal.addEventListener("click", (e) => {
      if (e.target === el.addSourceModal) closeAddSourceModal();
    });

    el.addSourceForm.addEventListener("submit", (e) => {
      e.preventDefault();
      const name = el.sourceNameInput.value.trim();
      const type = el.sourceTypeSelect.value;
      const path = el.sourcePathInput.value.trim();

      if (!name || !path) {
        el.addSourceFeedback.className = "avatar-feedback error";
        el.addSourceFeedback.textContent = "Please fill in all source fields.";
        return;
      }

      const newSource = {
        id: Date.now(),
        name: name,
        mediaType: type,
        path: path,
        freeSpace: "2.0 TB / 4.0 TB",
        isActive: true
      };

      FLAN_MOCK_DATA.storageSources.push(newSource);

      // Auto-discover candidate items for the new source to feed the pipeline
      const mockCandidateId = Date.now();
      if (type === "video") {
        FLAN_MOCK_DATA.pendingIngestion.push({
          id: mockCandidateId,
          sourceId: newSource.id,
          sourceName: newSource.name,
          rawFolder: "neon.genesis.evangelion.1995.1080p.bluray.x265",
          proposedTitle: "Neon Genesis Evangelion",
          proposedYear: 1995,
          mediaType: "video",
          fileCount: 26,
          coverColor: "#4a148c",
          badge: "Series",
          overview: "In 2015, the world stands on the brink of destruction. Humanity's last hope lies in the hands of Nerv, a special United Nations agency.",
          files: [
            { id: mockCandidateId + 1, rawFilename: "NGE.E01.Angel.Attack.mkv", title: "EP 01 - Angel Attack", duration: "24m", progress: 0, positionSeconds: 0, isFinished: false, isHidden: false },
            { id: mockCandidateId + 2, rawFilename: "NGE.E02.Unfamiliar.Ceiling.mkv", title: "EP 02 - Unfamiliar Ceiling", duration: "24m", progress: 0, positionSeconds: 0, isFinished: false, isHidden: false }
          ]
        });
      } else {
        FLAN_MOCK_DATA.pendingIngestion.push({
          id: mockCandidateId,
          sourceId: newSource.id,
          sourceName: newSource.name,
          rawFolder: "Foundation.Isaac.Asimov.1951.epub",
          proposedTitle: "Foundation",
          proposedAuthor: "Isaac Asimov",
          proposedYear: 1951,
          mediaType: "book",
          fileCount: 1,
          coverColor: "#1565c0",
          badge: "EPUB",
          overview: "For twelve thousand years the Galactic Empire has ruled supreme. Now it is dying.",
          files: [
            { id: mockCandidateId + 1, rawFilename: "Foundation - Isaac Asimov.epub", title: "Foundation", format: "epub", size: "1.5 MB" }
          ]
        });
      }

      el.addSourceFeedback.className = "avatar-feedback success";
      el.addSourceFeedback.textContent = `✓ Storage folder validated and .flan-keep initialized! Launching ingestion pipeline...`;

      setTimeout(() => {
        closeAddSourceModal();
        if (state.activeTab === "manage") {
          renderManageView();
        }
        openIngestionPipeline(newSource);
      }, 700);
    });

    // Ingestion Pipeline Modal Event Listeners
    if (el.closePipelineBtn) el.closePipelineBtn.addEventListener("click", closeIngestionPipeline);
    if (el.cancelPipelineBtn) el.cancelPipelineBtn.addEventListener("click", closeIngestionPipeline);
    if (el.ingestionPipelineModal) {
      el.ingestionPipelineModal.addEventListener("click", (e) => {
        if (e.target === el.ingestionPipelineModal) closeIngestionPipeline();
      });
    }

    if (el.pipelineSelectAll) {
      el.pipelineSelectAll.addEventListener("change", () => {
        const checked = el.pipelineSelectAll.checked;
        pipelineCandidates.forEach((item) => {
          item.selected = checked;
        });
        renderPipelineList();
      });
    }

    if (el.pipelineFilterInput) {
      el.pipelineFilterInput.addEventListener("input", () => {
        renderPipelineList();
      });
    }

    if (el.commitPipelineBtn) {
      el.commitPipelineBtn.addEventListener("click", commitPipeline);
    }

    // Direct Web Upload Form Submit & Close
    el.closeUploadBtn.addEventListener("click", closeUploadModal);
    el.uploadMediaModal.addEventListener("click", (e) => {
      if (e.target === el.uploadMediaModal) closeUploadModal();
    });

    el.uploadMediaForm.addEventListener("submit", (e) => {
      e.preventDefault();
      const sourceId = parseInt(el.uploadTargetSource.value);
      const title = el.uploadTitleInput.value.trim();
      const file = el.uploadFileInput.files[0];

      if (!title || !file) {
        el.uploadFeedback.className = "avatar-feedback error";
        el.uploadFeedback.textContent = "Please provide title and select a file.";
        return;
      }

      const source = FLAN_MOCK_DATA.storageSources.find((s) => s.id === sourceId) || FLAN_MOCK_DATA.storageSources[0];
      el.uploadProgressWrap.style.display = "flex";
      el.startUploadBtn.disabled = true;

      let progress = 0;
      const interval = setInterval(() => {
        progress += 25;
        el.uploadProgressFill.style.width = progress + "%";
        el.uploadPercentText.textContent = progress + "%";

        if (progress >= 100) {
          clearInterval(interval);
          el.uploadStatusText.textContent = "Streaming complete! Verifying direct-to-disk write...";

          setTimeout(() => {
            if (source.mediaType === "video") {
              FLAN_MOCK_DATA.videos.push({
                id: Date.now(),
                title: title,
                releaseYear: new Date().getFullYear(),
                type: "movie",
                overview: `Uploaded file: ${file.name} directly into ${source.name}`,
                coverColor: "#00695c",
                badge: "Movie",
                files: [
                  {
                    id: Date.now() + 1,
                    title: "Feature Playback",
                    duration: "Direct Upload",
                    progress: 0,
                    positionSeconds: 0,
                    isFinished: false,
                    isHidden: false,
                    streamUrl: "#",
                    downloadUrl: "#"
                  }
                ]
              });
            } else {
              FLAN_MOCK_DATA.books.push({
                id: Date.now(),
                title: title,
                author: "Direct Upload",
                overview: `Uploaded file: ${file.name} directly into ${source.name}`,
                coverColor: "#d84315",
                badge: file.name.endsWith(".pdf") ? "PDF" : "EPUB",
                files: [
                  {
                    id: Date.now() + 1,
                    title: title,
                    format: file.name.endsWith(".pdf") ? "pdf" : "epub",
                    size: (file.size ? (file.size / (1024 * 1024)).toFixed(1) + " MB" : "12.4 MB"),
                    downloadUrl: "#"
                  }
                ]
              });
            }

            el.uploadFeedback.className = "avatar-feedback success";
            el.uploadFeedback.textContent = `✓ Successfully streamed "${title}" straight to ${source.path}!`;

            setTimeout(() => {
              closeUploadModal();
              if (source.mediaType === "video") {
                window.location.hash = "#video";
              } else {
                window.location.hash = "#books";
              }
            }, 600);
          }, 300);
        }
      }, 150);
    });

    // Remove Storage Source Warning Modal Listeners
    if (el.closeRemoveSourceBtn) el.closeRemoveSourceBtn.addEventListener("click", closeRemoveSourceModal);
    if (el.cancelRemoveSourceBtn) el.cancelRemoveSourceBtn.addEventListener("click", closeRemoveSourceModal);
    if (el.removeSourceModal) {
      el.removeSourceModal.addEventListener("click", (e) => {
        if (e.target === el.removeSourceModal) closeRemoveSourceModal();
      });
    }
    if (el.confirmRemoveSourceBtn) {
      el.confirmRemoveSourceBtn.addEventListener("click", () => {
        if (pendingRemoveSource) {
          const idx = FLAN_MOCK_DATA.storageSources.findIndex((s) => s.id === pendingRemoveSource.id);
          if (idx !== -1) {
            FLAN_MOCK_DATA.storageSources.splice(idx, 1);
          }
          closeRemoveSourceModal();
          renderManageView();
        }
      });
    }

    // Edit Storage Source Modal Listeners
    if (el.closeEditSourceBtn) el.closeEditSourceBtn.addEventListener("click", closeEditSourceModal);
    if (el.cancelEditSourceBtn) el.cancelEditSourceBtn.addEventListener("click", closeEditSourceModal);
    if (el.editSourceModal) {
      el.editSourceModal.addEventListener("click", (e) => {
        if (e.target === el.editSourceModal) closeEditSourceModal();
      });
    }
    if (el.editSourceForm) {
      el.editSourceForm.addEventListener("submit", (e) => {
        e.preventDefault();
        const newName = el.editSourceNameInput.value.trim();
        if (!newName) {
          el.editSourceFeedback.className = "avatar-feedback error";
          el.editSourceFeedback.textContent = "Folder name cannot be blank.";
          return;
        }

        if (pendingEditSource) {
          pendingEditSource.name = newName;
          el.editSourceFeedback.className = "avatar-feedback success";
          el.editSourceFeedback.textContent = `✓ Storage folder name updated!`;

          setTimeout(() => {
            closeEditSourceModal();
            renderManageView();
          }, 400);
        }
      });
    }

    // Help Button (top header) -> Dedicated Software Manual View
    const helpTriggers = new Set([el.helpBtn, ...document.querySelectorAll(".help-trigger-btn")].filter(Boolean));
    helpTriggers.forEach((btn) => {
      btn.addEventListener("click", () => {
        window.location.hash = "#manual";
      });
    });
  }

  // Start app
  initListeners();
  handleRoute();
})();
