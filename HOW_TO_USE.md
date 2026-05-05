# How to Use `forter`

A step-by-step guide to using the TUI file organizer effectively.

## Table of Contents

1. [Installation](#installation)
2. [Basic Usage](#basic-usage)
3. [TUI Navigation](#tui-navigation)
4. [File Selection](#file-selection)
5. [Organization Modes](#organization-modes)
6. [Category Management](#category-management)
7. [Undo Operations](#undo-operations)
8. [Configuration](#configuration)
9. [Tips & Tricks](#tips--tricks)
10. [Troubleshooting](#troubleshooting)

---

## Installation

### Option 1: Install via Homebrew (Recommended)

```bash
# Add the tap (replace with your actual tap when published)
brew tap OmkarVijayBagade/forter
brew install forter
```

### Option 2: Build from Source

```bash
# Clone the repository
git clone https://github.com/OmkarVijayBagade/forter.git
cd forter

# Build the binary
make build

# Install to your system
make install
```

### Verify Installation

```bash
forter --version
```

---

## Basic Usage

### 1. Organize Current Directory

```bash
forter
```

This opens the TUI with all files in your current directory.

### 2. Organize a Specific Directory

```bash
forter ~/Downloads
forter ~/Desktop
forter /path/to/your/folder
```

### 3. Recursive Organization

Include all subdirectories:

```bash
forter -r ~/Documents
# or
forter --recursive ~/Documents
```

### 4. Dry Run Mode (Preview Changes)

See what would happen without moving any files:

```bash
forter --dry-run ~/Downloads
```

---

## TUI Navigation

When you launch `forter`, you'll see an interactive interface with three main areas:

```
┌─────────────────────────┬──────────────────┐
│ 📁 forter               │                  │
│ ~/Downloads | Files: 45 │                  │
├─────────────────────────┼──────────────────┤
│                         │                  │
│  ☑ file1.pdf            │ Categories       │
│  ☐ file2.jpg            │                  │
│  ☐ file3.mp4            │   All 45         │
│  ☑ file4.zip            │ ▸ Documents 12   │
│  ☐ file5.mp3            │   Images 8       │
│                         │   Videos 5       │
│                         │   Audio 3        │
│                         │   Archives 7     │
│                         │   Code 10        │
│                         │                  │
│                         │ Preview          │
│                         │ Name: file1.pdf  │
│                         │ Size: 2.5 MB     │
│                         │ Cat: Documents   │
│                         │                  │
│                         │ Target:          │
│                         │ Documents/PDF/   │
├─────────────────────────┴──────────────────┤
│ ↑↓ navigate │ tab switch │ space toggle ││
└───────────────────────────────────────────┘
```

### Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `↑` or `k` | Move cursor up |
| `↓` or `j` | Move cursor down |
| `Tab` or `←`/`→` | Switch between panels |
| `h`/`l` | Switch panels (vim-style) |
| `Space` or `Enter` | Toggle file selection / Confirm action |
| `a` | Select all files |
| `c` | Clear all selections |
| `o` | Start organization (when files selected) |
| `d` | Toggle dry run mode |
| `u` | Undo last operation |
| `q` or `Esc` or `Ctrl+C` | Quit |

---

## File Selection

### Selecting Individual Files

1. Navigate to a file using `↑`/`↓`
2. Press `Space` or `Enter` to select it
3. A `☑` will appear next to selected files

### Selecting All Files

Press `a` to select all visible files. Press `a` again or `c` to deselect all.

### Selecting by Category

1. Press `Tab` to switch to the **Categories** panel
2. Navigate to the desired category with `↑`/`↓`
3. Press `Space` or `Enter` to filter files by that category
4. Press `Tab` to return to files panel
5. Press `a` to select all filtered files

### Bulk Operations

- **Select multiple files**: Use `Space` on each file, or `a` for all
- **Preview destination**: The right panel shows where each file will be moved
- **Check total size**: The header shows total selected files and size

---

## Organization Modes

### Normal Mode (Move Files)

1. Select files you want to organize
2. Press `o`
3. Confirm with `Enter`
4. Files will be moved to `Category/Extension/` subdirectories

Example result:
```
Downloads/
├── Documents/
│   ├── PDF/
│   │   ├── report.pdf
│   │   └── invoice.pdf
│   └── DOCX/
│       └── letter.docx
├── Images/
│   ├── JPG/
│   │   ├── photo1.jpg
│   │   └── photo2.jpg
│   └── PNG/
│       └── screenshot.png
└── Videos/
    └── MP4/
        └── vacation.mp4
```

### Dry Run Mode (Preview Only)

Test your organization plan without making changes:

1. Press `d` to toggle dry run mode (header will show `[DRY RUN]`)
2. Select files and press `o`
3. You'll see a preview of what would happen
4. No files are actually moved

This is useful for:
- Testing category assignments
- Verifying destination paths
- Checking for potential duplicates

---

## Category Management

### Default Categories

Files are automatically categorized based on their extensions:

| Category | Extensions |
|----------|-----------|
| **Documents** | pdf, doc, docx, txt, rtf, odt, xls, xlsx, ppt, pptx, csv, md |
| **Images** | jpg, jpeg, png, gif, bmp, svg, webp, ico, tiff, raw, psd |
| **Videos** | mp4, avi, mkv, mov, wmv, flv, webm, m4v, mpg, mpeg, 3gp |
| **Audio** | mp3, wav, flac, aac, ogg, m4a, wma, opus |
| **Archives** | zip, rar, 7z, tar, gz, bz2, xz, tgz, tbz, iso |
| **Code** | go, py, js, ts, jsx, tsx, html, css, scss, sass, java, c, cpp, h, hpp, rs, rb, php, swift, kt, json, xml, yaml, yml, sql, sh, bash, zsh |
| **Others** | Everything else |

### Customizing Categories

Create `~/.forter.yaml`:

```yaml
# Add custom extensions to existing categories
categories:
  - name: "Documents"
    extensions: ["pdf", "doc", "docx", "txt", "custom"]
  
  - name: "Code"
    extensions: ["go", "py", "js", "myext"]

# Or create custom mappings for specific extensions
custom_mappings:
  "log": "Archives"
  "bak": "Archives"
  "tmp": "Others"
```

Reload the TUI after editing the config file.

---

## Undo Operations

Made a mistake? Undo the last organization:

1. Press `u`
2. All files from the last operation will be moved back to their original locations
3. Works for both normal and dry run modes

**Note**: Undo only works for the most recent operation. Restarting the tool clears the undo history.

---

## Configuration

### Config File Location

`~/.forter.yaml`

### Full Configuration Example

```yaml
# Default settings
default_category: "Others"
skip_hidden: true      # Skip files starting with .
skip_dirs: true        # Don't show directories in file list

# Category definitions
categories:
  - name: "Documents"
    extensions: 
      - "pdf"
      - "doc"
      - "docx"
      - "txt"
      - "md"
      - "rtf"
      - "odt"
    description: "Document files"

  - name: "Images"
    extensions:
      - "jpg"
      - "jpeg"
      - "png"
      - "gif"
      - "webp"
    description: "Image files"

  - name: "Videos"
    extensions:
      - "mp4"
      - "avi"
      - "mkv"
      - "mov"
    description: "Video files"

  - name: "Audio"
    extensions:
      - "mp3"
      - "wav"
      - "flac"
      - "aac"
    description: "Audio files"

  - name: "Archives"
    extensions:
      - "zip"
      - "rar"
      - "7z"
      - "tar"
      - "gz"
    description: "Archive files"

  - name: "Code"
    extensions:
      - "go"
      - "py"
      - "js"
      - "ts"
      - "html"
      - "css"
    description: "Source code files"

# Override specific extensions
custom_mappings:
  "log": "Archives"
  "bak": "Archives"
  "old": "Archives"
```

---

## Tips & Tricks

### 1. Organize Downloads Folder Regularly

```bash
# Create an alias
alias clean-downloads='forter ~/Downloads'

# Or a function with date
forter-downloads() {
    forter ~/Downloads
    echo "Downloads organized on $(date)"
}
```

### 2. Handle Duplicates Gracefully

The tool automatically renames duplicates:
- `file.txt` → `file.txt`
- `file.txt` (duplicate) → `file (1).txt`
- Another duplicate → `file (2).txt`

### 3. Use Dry Run for Unknown Directories

Always use `--dry-run` first when organizing unfamiliar directories:

```bash
forter --dry-run ~/SomeRandomFolder
```

### 4. Filter Before Selecting

1. Filter by category first (Tab → select category)
2. Then use `a` to select all filtered files
3. This prevents accidentally selecting unwanted file types

### 5. Check Total Size

The header shows the total size of selected files. Useful for:
- Checking if you're moving large amounts of data
- Verifying selections before organizing

### 6. Keyboard Shortcuts Reference

Print this and keep it handy:

```
┌─────────────────────────────────────┐
│ NAVIGATION                          │
│   ↑/k     Move up                   │
│   ↓/j     Move down                 │
│   Tab     Switch panel              │
│   h/l     Switch panel (vim)        │
│                                     │
│ SELECTION                           │
│   Space   Toggle selection          │
│   Enter   Toggle/Confirm            │
│   a       Select all                │
│   c       Clear selection           │
│                                     │
│ ACTIONS                             │
│   o       Forter selected         │
│   d       Toggle dry run            │
│   u       Undo last                 │
│   q/Esc   Quit                      │
└─────────────────────────────────────┘
```

---

## Troubleshooting

### Issue: "No files found"

**Solution**: Check if:
- The directory path is correct
- Hidden files are being skipped (check `skip_hidden` in config)
- You have permission to read the directory

### Issue: "Permission denied" when organizing

**Solution**:
- Check file/directory permissions
- Run with appropriate user permissions
- Some system files may be locked - these will be skipped automatically

### Issue: TUI looks broken / misaligned

**Solution**:
- Resize your terminal window (minimum 80x24 recommended)
- Ensure your terminal supports Unicode
- Try a different terminal emulator

### Issue: Cannot undo

**Solution**:
- Undo only works for the most recent operation
- If you've restarted the tool, undo history is cleared
- Check that destination files still exist

### Issue: Wrong category assignment

**Solution**:
- Add custom mapping to `~/.forter.yaml`
- Example: `custom_mappings: {"ext": "Category"}`

### Issue: Build fails with "module not found"

**Solution**:
```bash
cd forter
go mod download
go mod tidy
make build
```

---

## Quick Reference Card

```bash
# Basic commands
forter                    # Current directory
forter ~/path             # Specific directory
forter -r ~/path          # Recursive
forter --dry-run          # Preview mode
forter --version          # Show version

# TUI keys
↑/k, ↓/j     Navigate
Tab/h/l      Switch panel
Space/Enter  Toggle/Confirm
a            Select all
c            Clear
o            Forter
d            Dry run toggle
u            Undo
q/Esc        Quit
```

---

## Getting Help

- **GitHub Issues**: https://github.com/OmkarVijayBagade/forter/issues
- **Discussions**: https://github.com/OmkarVijayBagade/forter/discussions
- **Documentation**: See README.md for development docs

---

Happy organizing! 📁✨
