class Portly < Formula
  desc "CLI for managing ports — list, inspect, and kill by port number"
  homepage "https://github.com/mohokh67/portly"
  version "0.1.0"
  license "MIT"

  on_macos do
    on_intel do
      url "https://github.com/mohokh67/portly/releases/download/v#{version}/portly_#{version}_darwin_amd64.tar.gz"
      sha256 "REPLACE_WITH_SHA256_DARWIN_AMD64"
    end
    on_arm do
      url "https://github.com/mohokh67/portly/releases/download/v#{version}/portly_#{version}_darwin_arm64.tar.gz"
      sha256 "REPLACE_WITH_SHA256_DARWIN_ARM64"
    end
  end

  on_linux do
    on_intel do
      url "https://github.com/mohokh67/portly/releases/download/v#{version}/portly_#{version}_linux_amd64.tar.gz"
      sha256 "REPLACE_WITH_SHA256_LINUX_AMD64"
    end
    on_arm do
      url "https://github.com/mohokh67/portly/releases/download/v#{version}/portly_#{version}_linux_arm64.tar.gz"
      sha256 "REPLACE_WITH_SHA256_LINUX_ARM64"
    end
  end

  def install
    bin.install "portly"
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/portly --version")
  end
end
