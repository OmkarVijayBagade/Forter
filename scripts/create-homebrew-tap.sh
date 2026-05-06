#!/bin/bash
# Script to create and push homebrew tap

set -e

TAP_DIR="$HOME/homebrew-forter"
mkdir -p "$TAP_DIR"
cd "$TAP_DIR"

# Initialize if not already a git repo
if [ ! -d .git ]; then
    git init
fi

# Create the formula
cat > forter.rb << 'EOF'
class Forter < Formula
  desc "Fast Organized Terminal Explorer - TUI file organizer"
  homepage "https://github.com/OmkarVijayBagade/forter"
  version "1.0.0"

  # macOS AMD64 (Intel)
  if OS.mac? && Hardware::CPU.intel?
    url "https://github.com/OmkarVijayBagade/forter/releases/download/v1.0.0/forter_v1.0.0_darwin_amd64.tar.gz"
    sha256 "PLACEHOLDER_SHA256_DARWIN_AMD64"
  end

  # macOS ARM64 (Apple Silicon)
  if OS.mac? && Hardware::CPU.arm?
    url "https://github.com/OmkarVijayBagade/forter/releases/download/v1.0.0/forter_v1.0.0_darwin_arm64.tar.gz"
    sha256 "PLACEHOLDER_SHA256_DARWIN_ARM64"
  end

  # Linux AMD64
  if OS.linux? && Hardware::CPU.intel?
    url "https://github.com/OmkarVijayBagade/forter/releases/download/v1.0.0/forter_v1.0.0_linux_amd64.tar.gz"
    sha256 "PLACEHOLDER_SHA256_LINUX_AMD64"
  end

  # Linux ARM64
  if OS.linux? && Hardware::CPU.arm?
    url "https://github.com/OmkarVijayBagade/forter/releases/download/v1.0.0/forter_v1.0.0_linux_arm64.tar.gz"
    sha256 "PLACEHOLDER_SHA256_LINUX_ARM64"
  end

  license "MIT"

  def install
    bin.install "forter"
  end

  test do
    system "#{bin}/forter", "--version"
  end
end
EOF

# Create README
cat > README.md << 'EOF'
# Homebrew Tap for Forter

## Installation

```bash
brew tap OmkarVijayBagade/forter
brew install forter
```

## Updating

```bash
brew update
brew upgrade forter
```

## More Info

See the main repo: https://github.com/OmkarVijayBagade/forter
EOF

git add .
git commit -m "Initial formula for forter v1.0.0"

echo "Homebrew tap created at: $TAP_DIR"
echo "Next steps:"
echo "1. Create GitHub repo: https://github.com/new"
echo "2. Name: homebrew-forter"
echo "3. Run: git remote add origin https://github.com/OmkarVijayBagade/homebrew-forter.git"
echo "4. Run: git push -u origin main"
