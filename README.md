# Git Markdown Editor

A simple markdown editor with Git integration.

## Features

- Edit markdown files
- Git integration (commit, push, pull, diff)
- Basic markdown extensions (lists, tables)

## Installation

### Dependencies

- Go (1.18 or later)
- Git

### Git LFS (Large File Storage)

This editor supports uploading images, which are handled using Git LFS. Git LFS is a separate dependency and needs to be installed on your system for image uploads to work correctly.

To install Git LFS, please refer to the official Git LFS documentation: [https://git-lfs.com/](https://git-lfs.com/)

### Running the application

1. Clone the repository:
   `git clone https://github.com/your-username/git-markdown-editor.git`
2. Navigate to the `markdown-editor` directory:
   `cd git-markdown-editor/markdown-editor`
3. Run the application:
   `go run main.go`

The application will be available at `http://localhost:8080`.

## Usage

- **Toolbox:** Use the toolbox to insert lists, tables, or upload pictures.
- **Files:** Browse and manage your markdown files.
- **Git:** Commit, push, pull, and view diffs of your changes.
- **Editing:** Click on a file to view and edit its content. Click the "Source" button to edit the raw markdown.

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License

[MIT](LICENSE)
