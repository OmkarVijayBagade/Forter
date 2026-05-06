#!/bin/bash
# Setup script to make distribution live
# Run this after tagging v1.0.0

set -e

echo "🚀 Setting up forter distribution..."
echo ""

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Step 1: Create Homebrew Tap
echo -e "${BLUE}Step 1: Creating Homebrew Tap...${NC}"
TAP_DIR="$HOME/homebrew-forter"
mkdir -p "$TAP_DIR"
cd "$TAP_DIR"

if [ ! -d .git ]; then
    git init
fi

# Create formula with placeholders
cat > forter.rb << 'FORMULA'
class Forter < Formula
  desc "Fast Organized Terminal Explorer - TUI file organizer"
  homepage "https://github.com/OmkarVijayBagade/forter"
  version "1.0.0"

  if OS.mac? && Hardware::CPU.intel?
    url "https://github.com/OmkarVijayBagade/forter/releases/download/v1.0.0/forter_v1.0.0_darwin_amd64.tar.gz"
    sha256 "PLACEHOLDER"
  elsif OS.mac? && Hardware::CPU.arm?
    url "https://github.com/OmkarVijayBagade/forter/releases/download/v1.0.0/forter_v1.0.0_darwin_arm64.tar.gz"
    sha256 "PLACEHOLDER"
  elsif OS.linux? && Hardware::CPU.intel?
    url "https://github.com/OmkarVijayBagade/forter/releases/download/v1.0.0/forter_v1.0.0_linux_amd64.tar.gz"
    sha256 "PLACEHOLDER"
  elsif OS.linux? && Hardware::CPU.arm?
    url "https://github.com/OmkarVijayBagade/forter/releases/download/v1.0.0/forter_v1.0.0_linux_arm64.tar.gz"
    sha256 "PLACEHOLDER"
  end

  license "MIT"

  def install
    bin.install "forter"
  end

  test do
    system "#{bin}/forter", "--version"
  end
end
FORMULA

# Create README
cat > README.md << 'README'
# Homebrew Tap for Forter

\`\`\`bash
brew tap OmkarVijayBagade/forter
brew install forter
\`\`\`

https://github.com/OmkarVijayBagade/forter
README

git add .
git commit -m "Initial formula for forter v1.0.0" || true

echo -e "${GREEN}✓ Homebrew tap created at: $TAP_DIR${NC}"
echo ""

# Step 2: Create Scoop Bucket
echo -e "${BLUE}Step 2: Creating Scoop Bucket...${NC}"
SCOOP_DIR="$HOME/scoop-forter"
mkdir -p "$SCOOP_DIR"
cd "$SCOOP_DIR"

if [ ! -d .git ]; then
    git init
fi

cat > forter.json << 'SCOOP'
{
  "version": "1.0.0",
  "description": "Fast Organized Terminal Explorer - TUI file organizer",
  "homepage": "https://github.com/OmkarVijayBagade/forter",
  "license": "MIT",
  "architecture": {
    "64bit": {
      "url": "https://github.com/OmkarVijayBagade/forter/releases/download/v1.0.0/forter_v1.0.0_windows_amd64.zip",
      "hash": "PLACEHOLDER"
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
        "url": "https://github.com/OmkarVijayBagade/forter/releases/download/v$version/forter_v$version_windows_amd64.zip"
      }
    }
  }
}
SCOOP

cat > README.md << 'README'
# Scoop Bucket for Forter

\`\`\`powershell
scoop bucket add forter https://github.com/OmkarVijayBagade/scoop-forter
scoop install forter
\`\`\`

https://github.com/OmkarVijayBagade/forter
README

git add .
git commit -m "Initial manifest for forter v1.0.0" || true

echo -e "${GREEN}✓ Scoop bucket created at: $SCOOP_DIR${NC}"
echo ""

# Step 3: Instructions
echo -e "${YELLOW}📋 Next Steps:${NC}"
echo ""
echo "1. Push the main forter repo (merge develop to main):"
echo "   cd /Users/omkarvijaybagade/Desktop/Forter"
echo "   git checkout main"
echo "   git merge develop"
echo "   git push origin main"
echo ""
echo "2. Create GitHub repos for package managers:"
echo "   - https://github.com/new → Name: homebrew-forter"
echo "   - https://github.com/new → Name: scoop-forter"
echo ""
echo "3. Push tap and bucket:"
echo "   cd $HOME/homebrew-forter"
echo "   git remote add origin https://github.com/OmkarVijayBagade/homebrew-forter.git"
echo "   git push -u origin main"
echo ""
echo "   cd $HOME/scoop-forter"
echo "   git remote add origin https://github.com/OmkarVijayBagade/scoop-forter.git"
echo "   git push -u origin main"
echo ""
echo "4. Tag v1.0.0 release (triggers GitHub Actions):"
echo "   cd /Users/omkarvijaybagade/Desktop/Forter"
echo "   git tag v1.0.0"
echo "   git push origin v1.0.0"
echo ""
echo -e "${GREEN}Done! After release builds, update SHA256s in formula/manifest.${NC}"
