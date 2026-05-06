# Distribution Guide: Publishing Forter to Package Managers

Complete guide to distribute `forter` via Homebrew, Scoop, APT, and other package managers.

---

## Table of Contents

1. [Homebrew (macOS/Linux)](#homebrew-macoslinux)
2. [GitHub Releases](#github-releases)
3. [Scoop (Windows)](#scoop-windows)
4. [APT (Debian/Ubuntu)](#apt-debianubuntu)
5. [YUM/DNF (RHEL/CentOS/Fedora)](#yumdnf-rhelcentosfedora)
6. [Arch User Repository (AUR)](#arch-user-repository-aur)
7. [Go Install](#go-install)
8. [Snapcraft](#snapcraft)
9. [Flathub](#flathub)

---

## Homebrew (macOS/Linux)

Homebrew is the most popular package manager for macOS and works on Linux too.

### Step 1: Create a Homebrew Tap Repository

```bash
# Create a new GitHub repo for your tap
github.com/OmkarVijayBagade/homebrew-forter

# Clone it locally
git clone https://github.com/OmkarVijayBagade/homebrew-forter.git
cd homebrew-forter
```

### Step 2: Create the Formula

Create file: `forter.rb`

```ruby
class Forter < Formula
  desc "Fast Organized Terminal Explorer - TUI file organizer"
  homepage "https://github.com/OmkarVijayBagade/forter"
  version "1.0.0"
  
  # macOS AMD64 (Intel)
  if OS.mac? && Hardware::CPU.intel?
    url "https://github.com/OmkarVijayBagade/forter/releases/download/v1.0.0/forter_1.0.0_darwin_amd64.tar.gz"
    sha256 "PLACEHOLDER_SHA256"
  end
  
  # macOS ARM64 (Apple Silicon)
  if OS.mac? && Hardware::CPU.arm?
    url "https://github.com/OmkarVijayBagade/forter/releases/download/v1.0.0/forter_1.0.0_darwin_arm64.tar.gz"
    sha256 "PLACEHOLDER_SHA256"
  end
  
  # Linux AMD64
  if OS.linux? && Hardware::CPU.intel?
    url "https://github.com/OmkarVijayBagade/forter/releases/download/v1.0.0/forter_1.0.0_linux_amd64.tar.gz"
    sha256 "PLACEHOLDER_SHA256"
  end

  license "MIT"

  def install
    bin.install "forter"
    # Install shell completions if available
    # bash_completion.install "completions/forter.bash" => "forter"
    # zsh_completion.install "completions/forter.zsh" => "_forter"
    # fish_completion.install "completions/forter.fish"
  end

  test do
    system "#{bin}/forter", "--version"
  end
end
```

### Step 3: Commit and Push Formula

```bash
git add forter.rb
git commit -m "Add forter formula v1.0.0"
git push origin main
```

### Step 4: Users Can Now Install

```bash
# Add your tap
brew tap OmkarVijayBagade/forter

# Install
brew install forter

# Upgrade
brew upgrade forter
```

### Step 5: Automate with GitHub Actions

Create `.github/workflows/update-formula.yml` in the tap repo:

```yaml
name: Update Formula

on:
  repository_dispatch:
    types: [release]

jobs:
  update:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Update formula
        run: |
          VERSION=${{ github.event.client_payload.version }}
          # Update URLs and SHA256s
          # Use script to calculate SHA256s from release artifacts
```

---

## GitHub Releases

The foundation for all package managers. Build and upload binaries.

### Step 1: Create Release Script

Create `scripts/release.sh`:

```bash
#!/bin/bash
set -e

VERSION=$1
if [ -z "$VERSION" ]; then
  echo "Usage: ./release.sh v1.0.0"
  exit 1
fi

# Clean and build
rm -rf dist/
mkdir -p dist

# Build for all platforms
PLATFORMS=(
  "darwin/amd64"
  "darwin/arm64"
  "linux/amd64"
  "linux/arm64"
  "windows/amd64"
)

for platform in "${PLATFORMS[@]}"; do
  GOOS=${platform%/*}
  GOARCH=${platform#*/}
  
  output="dist/forter_${VERSION}_${GOOS}_${GOARCH}"
  if [ "$GOOS" = "windows" ]; then
    output="${output}.exe"
  fi
  
  echo "Building for $GOOS/$GOARCH..."
  GOOS=$GOOS GOARCH=$GOARCH go build \
    -ldflags "-X main.Version=$VERSION -X main.Commit=$(git rev-parse --short HEAD) -X main.BuildDate=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
    -o $output \
    cmd/forter/main.go
done

# Create archives
cd dist
for binary in forter_*; do
  if [[ "$binary" == *.exe ]]; then
    zip "${binary%.exe}.zip" "$binary"
  else
    tar -czf "${binary}.tar.gz" "$binary"
  fi
done
cd ..

echo "Release artifacts ready in dist/"
```

### Step 2: Make Script Executable and Run

```bash
chmod +x scripts/release.sh
./scripts/release.sh v1.0.0
```

### Step 3: Create GitHub Release

Option A: Using `gh` CLI
```bash
gh release create v1.0.0 \
  --title "Release v1.0.0" \
  --notes "Initial stable release of forter" \
  dist/*.tar.gz dist/*.zip
```

Option B: Manual on GitHub
1. Go to https://github.com/OmkarVijayBagade/forter/releases
2. Click "Draft a new release"
3. Choose tag: `v1.0.0`
4. Title: `Release v1.0.0`
5. Upload all files from `dist/`
6. Publish release

---

## Scoop (Windows)

Scoop is a command-line installer for Windows.

### Step 1: Create Scoop Bucket Repo

Create new repo: `github.com/OmkarVijayBagade/scoop-forter`

### Step 2: Create Manifest File

Create `forter.json`:

```json
{
  "version": "1.0.0",
  "description": "Fast Organized Terminal Explorer - TUI file organizer",
  "homepage": "https://github.com/OmkarVijayBagade/forter",
  "license": "MIT",
  "architecture": {
    "64bit": {
      "url": "https://github.com/OmkarVijayBagade/forter/releases/download/v1.0.0/forter_1.0.0_windows_amd64.zip",
      "hash": "PLACEHOLDER_SHA256"
    }
  },
  "bin": "forter.exe",
  "checkver": {
    "github": {
      "owner": "OmkarVijayBagade",
      "repo": "forter"
    }
  },
  "autoupdate": {
    "architecture": {
      "64bit": {
        "url": "https://github.com/OmkarVijayBagade/forter/releases/download/v$version/forter_$version_windows_amd64.zip"
      }
    }
  }
}
```

### Step 3: Users Install via Scoop

```powershell
# Add bucket
scoop bucket add forter https://github.com/OmkarVijayBagade/scoop-forter

# Install
scoop install forter

# Update
scoop update forter
```

---

## APT (Debian/Ubuntu)

Create a DEB package for Debian-based distributions.

### Step 1: Create Package Structure

```bash
mkdir -p forter-deb/DEBIAN
mkdir -p forter-deb/usr/local/bin
mkdir -p forter-deb/usr/share/doc/forter

# Copy binary
cp dist/forter_1.0.0_linux_amd64 forter-deb/usr/local/bin/forter
chmod +x forter-deb/usr/local/bin/forter

# Copy docs
cp README.md LICENSE forter-deb/usr/share/doc/forter/
```

### Step 2: Create Control File

Create `forter-deb/DEBIAN/control`:

```
Package: forter
Version: 1.0.0
Section: utils
Priority: optional
Architecture: amd64
Depends: 
Maintainer: Omkar Vijay Bagade <your.email@example.com>
Description: Fast Organized Terminal Explorer
 TUI file organizer built with Go and Bubble Tea.
 Organizes files into categories with interactive interface.
```

### Step 3: Build DEB Package

```bash
dpkg-deb --build forter-deb
mv forter-deb.deb forter_1.0.0_amd64.deb
```

### Step 4: Create APT Repository (Optional)

Use GitHub Pages + Aptly or reprepro to host APT repo.

Quick method using Packagecloud or Gemfury:
1. Sign up at https://packagecloud.io or https://gemfury.com
2. Upload the `.deb` file
3. Users add your repo and install

---

## YUM/DNF (RHEL/CentOS/Fedora)

Create RPM package for Red Hat-based distributions.

### Step 1: Create RPM Spec File

Create `forter.spec`:

```spec
Name:           forter
Version:        1.0.0
Release:        1%{?dist}
Summary:        Fast Organized Terminal Explorer
License:        MIT
URL:            https://github.com/OmkarVijayBagade/forter
Source0:        https://github.com/OmkarVijayBagade/forter/releases/download/v%{version}/forter_%{version}_linux_amd64.tar.gz

%description
TUI file organizer built with Go and Bubble Tea.

%prep
%setup -q -c

%install
mkdir -p %{buildroot}/usr/local/bin
install -m 755 forter %{buildroot}/usr/local/bin/forter

%files
/usr/local/bin/forter

%changelog
* Mon May 6 2026 Omkar Vijay Bagade <your.email@example.com> - 1.0.0-1
- Initial release
```

### Step 2: Build RPM

```bash
# Using rpmbuild
rpmbuild -ba forter.spec

# Or use Docker
docker run -v "$PWD:/workspace" -w /workspace \
  fedora:latest \
  bash -c "dnf install -y rpm-build && rpmbuild -ba forter.spec"
```

---

## Arch User Repository (AUR)

PKGBUILD for Arch Linux users.

### Step 1: Create PKGBUILD

```bash
# Maintainer: Omkar Vijay Bagade <your.email@example.com>
pkgname=forter
pkgver=1.0.0
pkgrel=1
pkgdesc="Fast Organized Terminal Explorer - TUI file organizer"
arch=('x86_64' 'aarch64')
url="https://github.com/OmkarVijayBagade/forter"
license=('MIT')
depends=('glibc')
makedepends=('go')
source=("$pkgname-$pkgver.tar.gz::https://github.com/OmkarVijayBagade/forter/archive/v$pkgver.tar.gz")
sha256sums=('SKIP')

build() {
  cd "$pkgname-$pkgver"
  go build -ldflags "-X main.Version=$pkgver" -o forter cmd/forter/main.go
}

package() {
  cd "$pkgname-$pkgver"
  install -Dm755 forter "$pkgdir/usr/bin/forter"
  install -Dm644 README.md "$pkgdir/usr/share/doc/forter/README.md"
  install -Dm644 LICENSE "$pkgdir/usr/share/licenses/forter/LICENSE"
}
```

### Step 2: Submit to AUR

```bash
# Create AUR account at https://aur.archlinux.org/
# Upload PKGBUILD
# Users install with: yay -S forter
```

---

## Go Install

Simplest method for Go users.

Users can install directly:

```bash
go install github.com/OmkarVijayBagade/forter/cmd/forter@latest
```

Ensure your repo has proper Go module tags.

---

## Snapcraft

Universal Linux package format.

### Step 1: Create snap/snapcraft.yaml

```yaml
name: forter
base: core22
version: '1.0.0'
summary: Fast Organized Terminal Explorer
description: |
  TUI file organizer built with Go and Bubble Tea.
  Organizes files into categories with interactive interface.

grade: stable
confinement: strict

apps:
  forter:
    command: bin/forter
    plugs:
      - home
      - removable-media

parts:
  forter:
    plugin: go
    source: .
    source-type: git
    build-packages:
      - golang-go
    override-build: |
      go build -ldflags "-X main.Version=$SNAPCRAFT_PROJECT_VERSION" \
        -o $SNAPCRAFT_PART_INSTALL/bin/forter \
        cmd/forter/main.go
```

### Step 2: Build and Publish

```bash
# Install snapcraft
sudo snap install snapcraft --classic

# Build
snapcraft

# Upload to Snap Store
snapcraft upload --release=stable forter_1.0.0_amd64.snap
```

---

## Flathub

For GUI app stores (optional, mostly for GUI apps).

TUI apps can use Flathub but it's less common. Prefer Snap or native packages.

---

## Automated Release with GitHub Actions

Create `.github/workflows/release.yml` in your forter repo:

```yaml
name: Release

on:
  push:
    tags:
      - 'v*'

permissions:
  contents: write

jobs:
  build:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        include:
          - goos: darwin
            goarch: amd64
          - goos: darwin
            goarch: arm64
          - goos: linux
            goarch: amd64
          - goos: linux
            goarch: arm64
          - goos: windows
            goarch: amd64
            extension: .exe

    steps:
      - uses: actions/checkout@v4
      
      - uses: actions/setup-go@v5
        with:
          go-version: '1.22'
      
      - name: Build
        env:
          GOOS: ${{ matrix.goos }}
          GOARCH: ${{ matrix.goarch }}
        run: |
          output="forter_${{ github.ref_name }}_${{ matrix.goos }}_${{ matrix.goarch }}${{ matrix.extension }}"
          go build -ldflags "-X main.Version=${{ github.ref_name }}" -o $output cmd/forter/main.go
          if [ "${{ matrix.goos }}" = "windows" ]; then
            zip "${output%.exe}.zip" "$output"
          else
            tar -czf "${output}.tar.gz" "$output"
          fi
      
      - name: Upload artifacts
        uses: actions/upload-artifact@v4
        with:
          name: binaries
          path: forter_*

  release:
    needs: build
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Download artifacts
        uses: actions/download-artifact@v4
        with:
          name: binaries
          path: dist/
      
      - name: Create Release
        uses: softprops/action-gh-release@v1
        with:
          files: dist/*
          generate_release_notes: true
```

---

## Quick Start Checklist

1. **Tag release**: `git tag v1.0.0 && git push origin v1.0.0`
2. **GitHub Actions** builds and releases automatically
3. **Update Homebrew formula** with new SHA256s
4. **Update Scoop manifest** with new URL and hash
5. **Build DEB/RPM** for Linux distributions
6. **Announce** on Reddit, Twitter, Hacker News, etc.

---

## Useful Tools

- **goreleaser**: https://goreleaser.com/ - Automates releases to multiple platforms
- **fury**: https://fury.co/ - Host APT/YUM repos
- **snapcraft**: https://snapcraft.io/ - Build and publish snaps

---

## Next Steps

1. Set up **goreleaser** for automated multi-platform releases
2. Create **Homebrew tap** (most important for macOS users)
3. Create **Scoop bucket** for Windows users
4. Submit to **AUR** for Arch users
5. Build **DEB/RPM** for enterprise Linux users

Happy distributing! 🚀
