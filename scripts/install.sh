#!/usr/bin/env sh
set -eu

repo="JinRudy/reprogate"
bin_dir="${BIN_DIR:-$HOME/.local/bin}"
version="${REPROGATE_VERSION:-latest}"

fail() {
  echo "reprogate install: $*" >&2
  exit 1
}

command -v curl >/dev/null 2>&1 || fail "curl is required"
command -v uname >/dev/null 2>&1 || fail "uname is required"

os="$(uname -s | tr '[:upper:]' '[:lower:]')"
arch="$(uname -m)"

case "$os" in
  linux|darwin) ;;
  *) fail "unsupported OS: $os. Download a binary from https://github.com/$repo/releases" ;;
esac

case "$arch" in
  x86_64|amd64) arch="amd64" ;;
  arm64|aarch64) arch="arm64" ;;
  *) fail "unsupported architecture: $arch" ;;
esac

asset="reprogate-$os-$arch"
base_url="https://github.com/$repo/releases"

if [ "$version" = "latest" ]; then
  url="$base_url/latest/download/$asset"
else
  url="$base_url/download/$version/$asset"
fi

mkdir -p "$bin_dir"
tmp_file="$(mktemp)"
trap 'rm -f "$tmp_file"' EXIT HUP INT TERM

curl -fL --retry 3 --retry-delay 1 --connect-timeout 10 --max-time 120 "$url" -o "$tmp_file"
install -m 0755 "$tmp_file" "$bin_dir/reprogate"

echo "reprogate installed to $bin_dir/reprogate"

case ":$PATH:" in
  *":$bin_dir:"*) ;;
  *) echo "Add $bin_dir to PATH if reprogate is not found." ;;
esac
