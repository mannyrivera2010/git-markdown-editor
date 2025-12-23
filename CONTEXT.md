Here is the **Master Context File** for Version 15.

This document serves as the architectural blueprint for the project. It describes the system design, file structure, and the responsibility of every component without repeating the source code.

---

# Project: Git-Backed Todo App (Version 15)

## 1. System Overview

* **Architecture:** Monolithic Go web server using Server-Side Rendering (SSR).
* **Database:** Local filesystem. The state is stored in Markdown files (default: `todo.md`).
* **Transaction Log:** A local Git repository tracks every change (Add, Toggle, Edit, Delete) as a commit.
* **Frontend Interactivity:** HTML5 combined with **HTMX** for partial page swaps (no full reloads) and **Tailwind CSS** for styling.
* **Authentication:** Session-based auth using Argon2 password hashing stored in a local `shadow_file`.

## 2. Directory Structure

```text
/mytodo
├── go.mod                # Go module definition (dependencies: goldmark, argon2)
├── main.go               # Application entry point, routing, and HTTP handlers
├── auth.go               # Authentication logic (Shadow file management)
├── git.go                # Wrapper for executing system Git commands
├── store.go              # File I/O, Markdown parsing, and File Tree logic
├── renderer.go           # HTML generation (Markdown -> HTML, Diff views, History)
├── index.html            # The master HTML template
└── static/
    ├── style.css         # CSS overrides for raw Markdown rendering
    └── script.js         # Client-side DOM manipulation (button injection)

```

---

## 3. Component Responsibilities

### Backend Components

**`main.go`**

* **Server Setup:** Initializes the HTTP server on port `:8080`.
* **Dependency Injection:** Wires the Store, VCS, Renderer, and Auth services together.
* **Middleware:** Implements a `protect()` middleware that checks for the `session_token` cookie before allowing access to protected routes.
* **Routing:** Maps HTTP endpoints (e.g., `/add`, `/toggle`, `/tree`, `/diff`) to specific handler functions.

**`store.go`**

* **Persistence Layer:** Handles reading and writing to the local `.md` files.
* **Parsing Logic:** Contains specific logic to identify and manipulate ````list` and ````table` blocks within the Markdown text.
* **File Management:** Implements `GetFileTree` (recursive/flat), `CreateFile`, and `DeleteFile` to manage the workspace.
* **Concurrency:** Uses a `sync.Mutex` to ensure atomic file writes.

**`git.go`**

* **VCS Interface:** Abstract interface for version control operations.
* **Git Wrapper:** Executes raw shell commands (`git add`, `git commit`, `git push`, `git show`, `git log`) to sync state.
* **Diff Engine:** Fetches the raw text diff of the `HEAD` commit for the UI.

**`auth.go`**

* **Identity Provider:** Manages the `shadow_file` (flat-file user database).
* **Security:** Handles Argon2 password hashing and salt generation.
* **Verification:** Iterates through the shadow file to validate credentials during login.

**`renderer.go`**

* **View Layer:** Converts raw text/markdown into HTML fragments.
* **Markdown Engine:** Uses the `goldmark` library to render standard Markdown.
* **Specialized Renderers:**
* `RenderTable`: Manually parses pipe-separated tables to inject editable HTML inputs and delete buttons.
* `RenderDiff`: Parses Git diff output to color-code additions (green) and deletions (red).
* `RenderFileTree`: Generates the HTML list for the file explorer sidebar.



### Frontend Components

**`index.html`**

* **Layout:** Defines the main 4-column grid layout (File Tree | Main Workspace | Git History).
* **HTMX Attributes:** Most interactions (forms, buttons) use `hx-post`, `hx-target`, and `hx-swap` to update specific parts of the DOM without refreshing the page.

**`static/script.js`**

* **DOM Injection:** Scans the rendered Markdown HTML for list items (`<li>`).
* **Interactive Buttons:** Dynamically injects "Done/Undo" and "Delete" buttons next to every task item, attaching HTMX attributes to them on the fly.

**`static/style.css`**

* **Markdown Styling:** Provides specific styles for the `markdown-body` class (e.g., hiding raw checkboxes, strikethrough for completed items) that Tailwind utilities cannot easily target.

---

## 4. Key Workflows

1. **Login:** User submits credentials -> `auth.go` validates against `shadow_file` -> Cookie set -> Redirect to `/`.
2. **Add Task:** User types task -> `main.go` calls `store.Add()` -> Appends to `todo.md` -> `git.go` commits -> `renderer.go` returns updated HTML fragment -> HTMX swaps the list.
3. **File Management:** User toggles "Recursive" checkbox -> `store.GetFileTree()` scans directory -> Returns updated file tree HTML.
4. **Syncing:** User clicks "Push" -> `git.go` executes `git push`.