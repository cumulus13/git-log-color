class GitRemoteColor < Formula
  desc "A wrapper around `git log`"
  homepage "https://github.com/cumulus13/git-log-color"
  url "https://github.com/cumulus13/git-log-color/releases/download/v1.0.6/git-log-color-darwin-amd64"
  version "1.0.2"
  sha256 "PUT_REAL_SHA256_HERE"

  def install
    bin.install "git-log-color-darwin-amd64" => "git-log-color"
  end
end