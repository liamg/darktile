#!/usr/bin/env bash

set -e

echo "Determining platform..."
platform=$(uname | tr '[:upper:]' '[:lower:]')
echo "Finding latest release..."
asset=$(curl --silent https://api.github.com/repos/liamg/darktile/releases/latest | jq -r ".assets[] | select(.name | contains(\"${platform}\")) | .url")
echo "Downloading latest release for your platform..."
curl -s -L -H "Accept: application/octet-stream" "${asset}" --output /tmp/darktile
echo "Installing darktile..."
chmod +x /tmp/darktile
installdir="${HOME}/bin/"
if [ "$EUID" -eq 0 ]; then
  installdir="/usr/local/bin/"
fi
mkdir -p $installdir
mv /tmp/darktile "${installdir}/darktile"
which darktile &> /dev/null || (echo "Please add ${installdir} to your PATH to complete installation!" && exit 1)
echo "Installation complete!"