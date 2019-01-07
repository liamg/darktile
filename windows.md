
### Windows dev env setup and build instructions:

1. Setup choco package manager https://chocolatey.org/docs/installation
2. Use `choco` to install golang and mingw
```choco install golang mingw```


### Setting aminal GoLang build env and directories structures for the project:

```
cd %YOUR_PROJECT_WORKING_DIR%
mkdir go\src\github.com\liamg
cd go\src\github.com\liamg
git clone git@github.com:jumptrading/aminal-mirror.git
move aminal-mirror aminal

set GOPATH=%YOUR_PROJECT_WORKING_DIR%\go
set GOBIN=%GOPATH%/bin
set PATH=%GOBIN%;%PATH%

cd aminal
go get
windres -o aminal.syso aminal.rc
go build
go install
```

Look for the aminal.exe built binary under your %GOBIN% path

