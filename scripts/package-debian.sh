#!/bin/bash
set -e

VERSION=${1//v}
BINARYFILE=$2

if [[ "$VERSION" == "" ]]; then
    echo "No version specified"
    exit 1
fi

if [[ "$BINARYFILE" == "" ]]; then
    echo "No path to binary specified"
    exit 1
fi

rm -rf package
mkdir -p package/DEBIAN
mkdir -p package/usr/local/bin/
cp $BINARYFILE package/usr/local/bin/aminal
chmod +x package/usr/local/bin/aminal

cat > package/DEBIAN/control <<- EOM
Package: aminal
Version: $VERSION
Maintainer: Liam Galvin
Architecture: amd64
Description: A Modern Terminal Emulator
EOM

dpkg-deb --build package
rm -rf package

mkdir -p bin/debian
mv package.deb bin/debian/${BINARYFILE}.deb
