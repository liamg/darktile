SHELL := /bin/bash
BINARY := aminal
FONTPATH := ./gui/packed-fonts
GEN_SRC_DIR := ./generated-src
VERSION_MAJOR := 0
VERSION_MINOR := 9
VERSION_PATCH := 0
VERSION := ${VERSION_MAJOR}.${VERSION_MINOR}.${VERSION_PATCH}

.PHONY: build
build: 
	./build.sh `git describe --tags`

.PHONY: test
test:
	go test -v ./...
	go vet -v

.PHONY: install
install: build
	go install -ldflags "-X github.com/liamg/aminal/version.Version=`git describe --tags`"

.PHONY: install-tools
install-tools:
	which dep || curl -L https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
	which packr || go get -u github.com/gobuffalo/packr/packr

.PHONY:	build-linux
build-linux:
	mkdir -p bin/linux
	GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o bin/linux/${BINARY}-linux-amd64 -ldflags "-X github.com/liamg/aminal/version.Version=${CIRCLE_TAG}"

.PHONY:	build-darwin
build-darwin:
	mkdir -p bin/darwin
	xgo -x -v -ldflags "-X github.com/liamg/aminal/version.Version=${CIRCLE_TAG}" --targets=darwin/amd64 -out bin/darwin/${BINARY} .

.PHONY:	package-debian
package-debian: build-linux
	./scripts/package-debian.sh "${CIRCLE_TAG}" bin/linux/${BINARY}-linux-amd64

.PHONY:	build-linux-travis
build-linux-travis:
	mkdir -p bin/linux
	GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o bin/linux/${BINARY}-linux-amd64 -ldflags "-X github.com/liamg/aminal/version.Version=${TRAVIS_TAG}"

.PHONY: windows-cross-compile-travis
windows-cross-compile-travis:
	mkdir -p bin/windows
	x86_64-w64-mingw32-windres -o aminal.syso aminal.rc
	GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CXX=x86_64-w64-mingw32-g++ CC=x86_64-w64-mingw32-gcc go build -o bin/windows/${BINARY}-windows-amd64.exe -ldflags "-X github.com/liamg/aminal/version.Version=${TRAVIS_TAG}"

.PHONY:	build-windows
build-windows:
	windres -o aminal.syso aminal.rc
	go build -o ${BINARY}-windows-amd64.exe

.PHONY: launcher-windows
launcher-windows: build-windows
	if exist "${GEN_SRC_DIR}\launcher" rmdir /S /Q "${GEN_SRC_DIR}\launcher"
	xcopy "windows\launcher\*.*" "${GEN_SRC_DIR}\launcher" /K /H /Y /Q /I
	powershell -Command "(gc ${GEN_SRC_DIR}\launcher\versioninfo.json) -creplace 'VERSION_MAJOR', '${VERSION_MAJOR}' | Out-File -Encoding default ${GEN_SRC_DIR}\launcher\versioninfo.json"
	powershell -Command "(gc ${GEN_SRC_DIR}\launcher\versioninfo.json) -creplace 'VERSION_MINOR', '${VERSION_MINOR}' | Out-File -Encoding default ${GEN_SRC_DIR}\launcher\versioninfo.json"
	powershell -Command "(gc ${GEN_SRC_DIR}\launcher\versioninfo.json) -creplace 'VERSION_PATCH', '${VERSION_PATCH}' | Out-File -Encoding default ${GEN_SRC_DIR}\launcher\versioninfo.json"
	powershell -Command "(gc ${GEN_SRC_DIR}\launcher\versioninfo.json) -creplace 'VERSION', '${VERSION}' | Out-File -Encoding default ${GEN_SRC_DIR}\launcher\versioninfo.json"
	powershell -Command "(gc ${GEN_SRC_DIR}\launcher\versioninfo.json) -creplace 'YEAR', (Get-Date -UFormat '%Y') | Out-File -Encoding default ${GEN_SRC_DIR}\launcher\versioninfo.json"
	copy aminal.ico "${GEN_SRC_DIR}\launcher" /Y
	go generate "${GEN_SRC_DIR}\launcher"
	if exist "bin\windows\launcher" rmdir /S /Q "bin\windows\launcher"
	mkdir "bin\windows\launcher\Versions\${VERSION}"
	go build -o "bin\windows\launcher\${BINARY}.exe" -ldflags "-H windowsgui" "${GEN_SRC_DIR}\launcher"
	copy ${BINARY}-windows-amd64.exe "bin\windows\launcher\Versions\${VERSION}\${BINARY}.exe" /Y

.PHONY: installer-windows
installer-windows:
	if exist "${GEN_SRC_DIR}\installer" rmdir /S /Q "${GEN_SRC_DIR}\installer"
	xcopy "windows\installer\*.*" "${GEN_SRC_DIR}\installer" /K /H /Y /Q /I
	powershell -Command "(gc ${GEN_SRC_DIR}\installer\installer.go) -creplace 'VERSION', '${VERSION}' | Out-File -Encoding default ${GEN_SRC_DIR}\installer\installer.go"
	go-bindata -prefix "bin\windows\launcher" -o "${GEN_SRC_DIR}/installer/data/data.go" "./bin/windows/launcher/..."
	powershell -Command "(gc ${GEN_SRC_DIR}\installer\data\data.go) -creplace 'package main', 'package data' | Out-File -Encoding default ${GEN_SRC_DIR}\installer\data\data.go"
	go build -o bin/windows/AminalSetup.exe -ldflags "-H windowsgui" "${GEN_SRC_DIR}/installer/installer.go"
	rem If an .exe name contains "installer", "setup" etc., then at least Windows 10 automatically
	rem opens a UAC prompt upon opening it. To avoid this, we add a compatibility manifest to the .exe.
	mt -manifest windows\installer\AminalSetup.exe.manifest -outputresource:bin\windows\AminalSetup.exe;1

.PHONY:	build-darwin-native-travis
build-darwin-native-travis:
	mkdir -p bin/darwin
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build -o bin/darwin/${BINARY}-darwin-amd64 -ldflags "-X github.com/liamg/aminal/version.Version=${TRAVIS_TAG}"

