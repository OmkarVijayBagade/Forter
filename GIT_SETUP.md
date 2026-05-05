# Git Setup Instructions

Complete workflow to create the GitHub repository and set up clean branching.

## Step 1: Create GitHub Repository

Go to https://github.com/new and create a new repository:

- **Repository name**: `forter`
- **Description**: "A high-performance TUI file organizer built with Go and Bubble Tea"
- **Visibility**: Public (or Private)
- **Initialize**: ❌ DO NOT initialize with README (we have our own)

## Step 2: Initialize Local Repository

```bash
# Navigate to your project
cd /Users/omkarvijaybagade/Desktop/Forter

# Initialize git
git init

# Add all files
git add .

# Create initial commit
git commit -m "Initial commit: TUI file organizer with Bubble Tea

Features:
- Interactive file browser with vim-like navigation
- Automatic file categorization (Documents, Images, Videos, etc.)
- Bulk selection and preview
- Dry-run mode for safe testing
- Duplicate handling with auto-renaming
- Undo support
- Cross-platform (macOS, Linux)"

# Connect to GitHub (replace with your actual URL)
git remote add origin https://github.com/OmkarVijayBagade/forter.git

# Push to main
git branch -M main
git push -u origin main
```

## Step 3: Create Development Branch

```bash
# Create and switch to develop branch
git checkout -b develop

# Push develop branch to GitHub
git push -u origin develop
```

## Step 4: Feature Branch Workflow

For adding HOW_TO_USE.md (clean workflow):

```bash
# From develop, create feature branch
git checkout develop
git checkout -b feature/add-usage-docs

# Add the HOW_TO_USE.md file (already created)
git add HOW_TO_USE.md
git add GIT_SETUP.md

# Commit
git commit -m "docs: add comprehensive usage documentation

- Add HOW_TO_USE.md with step-by-step guide
- Include keyboard shortcuts reference
- Add troubleshooting section
- Include configuration examples"

# Push feature branch
git push -u origin feature/add-usage-docs
```

## Step 5: Create Pull Request and Merge

### Option A: Using GitHub CLI (gh)

```bash
# Install gh if not present: brew install gh

# Authenticate (if first time)
gh auth login

# Create PR
gh pr create \
  --title "docs: Add comprehensive usage documentation" \
  --body "## Changes\n\n- Added HOW_TO_USE.md with detailed usage instructions\n- Added keyboard shortcuts reference\n- Added troubleshooting section\n- Added configuration examples\n\n## Checklist\n\n- [x] Documentation added\n- [x] Examples provided\n- [x] Keyboard shortcuts documented" \
  --base develop \
  --head feature/add-usage-docs

# Merge PR to develop
gh pr merge --merge
```

### Option B: Manual on GitHub

1. Go to https://github.com/OmkarVijayBagade/forter
2. Click "Pull requests" → "New pull request"
3. Set:
   - **base**: `develop`
   - **compare**: `feature/add-usage-docs`
4. Add title: "docs: Add comprehensive usage documentation"
5. Add description from above
6. Click "Create pull request"
7. After tests pass, click "Merge pull request"
8. Confirm merge

## Step 6: Merge to Main (Release)

```bash
# Switch to develop and pull latest
git checkout develop
git pull origin develop

# Create PR from develop to main (using gh)
gh pr create \
  --title "Release: Initial stable version with documentation" \
  --body "## Release v1.0.0\n\n### Features\n- TUI file organizer\n- Category-based organization\n- Dry-run mode\n- Undo support\n\n### Documentation\n- README.md\n- HOW_TO_USE.md" \
  --base main \
  --head develop

# Or merge directly (if you're working solo)
git checkout main
git merge develop --no-ff -m "Release v1.0.0: Initial stable version

Includes:
- Complete TUI implementation
- File scanner with concurrent workers
- Organization engine with duplicate handling
- Comprehensive documentation"

git push origin main
```

## Step 7: Clean Up

```bash
# Delete feature branch locally
git branch -d feature/add-usage-docs

# Delete feature branch on remote
git push origin --delete feature/add-usage-docs

# Verify clean state
git branch -a
```

## Branch Strategy (Git Flow)

```
main (production)
  ↑
develop (integration)
  ↑
feature/* (features)
  ↑
hotfix/* (urgent fixes)
```

### Quick Commands Cheat Sheet

```bash
# Start new feature
git checkout develop
git checkout -b feature/my-feature
git push -u origin feature/my-feature
# ... make changes ...
git add .
git commit -m "feat: add new feature"
git push
# Create PR on GitHub, merge, then:
git checkout develop
git pull origin develop
git branch -d feature/my-feature
git push origin --delete feature/my-feature

# Release to main
git checkout main
git merge develop --no-ff
git tag -a v1.0.0 -m "Release version 1.0.0"
git push origin main --tags
```

## Optional: Protect Main Branch

On GitHub repository settings:

1. Go to **Settings** → **Branches**
2. Add rule for `main`:
   - ☑ Require pull request reviews
   - ☑ Require status checks (if you have CI)
   - ☑ Include administrators
3. Add rule for `develop`:
   - ☑ Require pull request reviews

This ensures all changes go through proper review.

## Quick Start Summary

```bash
# One-time setup
git init
git add .
git commit -m "Initial commit"
git remote add origin https://github.com/OmkarVijayBagade/forter.git
git push -u origin main

# Branch setup
git checkout -b develop
git push -u origin develop

# Feature workflow
git checkout -b feature/docs
git add HOW_TO_USE.md
git commit -m "Add usage docs"
git push -u origin feature/docs
# ... create PR and merge ...

# Release
git checkout main
git merge develop
git push
```

Done! Your repository is now set up with a clean branching strategy.
