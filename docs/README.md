# Documentation Index

Technical documentation for Flan Media Server.

---

## Reading Paths

To navigate the documentation efficiently based on your role or objective:

### 1. Backend Development

If you are contributing Go code or adding server capabilities:

1. [Directory Structure & Architecture](directories.md): Package layout, Controller-Service-Repository-Model layers, and Java-to-Go concept mapping.
2. [Master System Architecture](design.md): System constraints, HTTP routes, authorization matrix, and data flow.
3. [Database Schema & Queries](database.md): 7-table relational schema, connection pooling, and migrations.
4. [Storage Architecture](storage.md): Mount validation, filesystem layout, and media conventions.
5. [Local Metadata Engine](scraper.md): Directory traversal and database reconciliation algorithms.
6. [Testing Strategy](testing.md): Guidelines for writing tests with in-memory SQLite and mock filesystems.

### 2. Operations & Deployment

If you are deploying Flan on a home server, VPS, or single-board computer:

1. [Compilation & Deployment Guide](compilation.md): Build flags, architecture targets (AMD64, ARM64, ARMv7), and systemd service setup.
2. [Storage Architecture](storage.md): Drive layout, mount recovery, and USB disk safety.
3. [Authentication Security & Rate Limiting](rate-limiting.md): PIN lockout safeguards and bcrypt CPU throttling.
4. [Security Threat Model](threat-model.md): Network exposure, authentication mechanisms, and privilege boundaries.

### 3. Web Client & UI Design

If you are editing templates, CSS styles, or frontend interactions:

1. [Design System Foundations](client/design-system.md): Layout structure, color palette, typography, and CSS variables.
2. [Page Templates & Wireframes](client/pages.md): Specifications for all 11 application views.
3. [Component Specifications](client/components.md): Modals, cards, form inputs, and notification banners.
4. [Accessibility Guide](client/accessibility.md): WCAG 2.1 AA compliance, keyboard navigation, and focus management.
5. [Mobile Responsiveness](client/responsiveness.md): Screen breakpoints and mobile bottom navigation behavior.

---

## Document Directory

### Architecture & Backend

* **[Master System Architecture](design.md):** Overall system design, hardware target expectations, HTTP route definitions, and authorization rules.
* **[Directory Structure & Architecture](directories.md):** Detailed guide to packages in `internal/`, constructor dependency injection, and separation of concerns.
* **[Database Schema & Queries](database.md):** Complete DDL for the 7-table SQLite database, WAL mode configuration, and migration strategies.
* **[Storage Architecture & Drive Resiliency](storage.md):** Partition decoupling between flash storage and bulk media, `.flan-keep` markers, and drive-disconnect defenses.
* **[Local-First Metadata Engine](scraper.md):** Specification for the offline crawler that discovers titles and posters without external network access.

### Operations & Security

* **[Compilation & Homelab Deployment](compilation.md):** Cross-compilation instructions for x86 and ARM SBCs, performance considerations, and production systemd service files.
* **[Authentication Security & Rate Limiting](rate-limiting.md):** Specifications for the dual-key PIN lockout, bcrypt CPU protection gate, and direct streaming.
* **[Security Threat Model](threat-model.md):** Analysis of local network threats, token validation, CSRF defenses, and authentication boundaries.
* **[Testing Strategy](testing.md):** Testing methodology using Go's `testing/fstest` and in-memory SQLite instances.

### Web Client & User Interface

* **[Design System](client/design-system.md):** Core design tokens, high-contrast palette, typography, and CSS rules.
* **[Component Specifications](client/components.md):** Structural markup and behavioral specifications for UI components.
* **[Page Templates](client/pages.md):** Wireframes and layout specifications for all 11 HTML views.
* **[Accessibility Guide](client/accessibility.md):** Standards for screen readers, keyboard-only operation, and contrast compliance.
* **[Mobile Responsiveness](client/responsiveness.md):** Viewport adaptation rules and bottom navigation transformations.

### System Diagrams

* **[Data Flow & Architecture](diagrams/data-flow.md):** Text and Mermaid sequence diagrams for requests, streaming, and background scanning.
* **[User Journeys & Flows](diagrams/user-flows.md):** Step-by-step state diagrams for user login, initial setup, and media playback.
