
### Windows dev env setup and build instructions:

1. Setup choco package manager https://chocolatey.org/docs/installation
2. Use `choco` to install golang and mingw
```choco install golang mingw```


### Setting aminal GoLang build env and directories structures for the project:

```
cd %YOUR_PROJECT_WORKING_DIR%
mkdir go\src\github.com\liamg
cd go\src\github.com\liamg
git clone https://github.com/liamg/aminal.git

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

### Building an installer for automatic updates:

In addition to the above commands:

```
go get github.com/josephspurrier/goversioninfo/cmd/goversioninfo
go get golang.org/x/sys/windows
go get github.com/jteeuwen/go-bindata/...
```

Install NSIS and place it on the `PATH`:
```
choco install nsis
set PATH=%PATH%;%ProgramFiles(x86)%\NSIS\Bin
```

Ensure `signtool.exe` is on the `PATH`. For instance, on Windows 10:

```
set PATH=%PATH%;%ProgramFiles(x86)%\Windows Kits\10\bin\10.0.17763.0\x64
```

Copy your code signing certificate to `go\src\github.com\liamg\windows\codesigning_certificate.pfx`.

Set the `WINDOWS_CODESIGNING_CERT_PW` to the password of your code signing certificate:

```
set WINDOWS_CODESIGNING_CERT_PW=PASSWORD
```

Compile Aminal and build the installer:

```
mingw32-make installer-windows
```

This produces several files in `bin/windows`. Their purpose is explained below.

### How Aminal's automatic update mechanism works (on Windows)

Aminal uses a technology called [Google Omaha](https://github.com/google/omaha) for automatic updates on Windows. It's the same technology which Google use to auto-update Chrome. For a quick introduction, see [this Google Omaha tutorial](https://fman.io/blog/google-omaha-tutorial/).

Aminal has an online installer. This is a 1 MB executable, which was created using Google Omaha. Suppose it's called `OnlineInstaller.exe`. When a user runs it, the following things happen:

 * `OnlineInstaller.exe` installs Aminal's version of Google Omaha on the user's system. In particular, this creates two `AminalUpdateTask...` tasks in the Windows Task Scheduler. They're set up to run once per day in the background and check for updates.
 * `OnlineInstaller.exe` contacts Aminal's update server and asks "what is the latest version?". The server responds with the URL to an _offline_ installer and some command line parameters. Say this offline installer is called `install-aminal-0.9.0.exe`.
 * `OnlineInstaller.exe` downloads `install-aminal-0.9.0.exe` and invokes it with the given command line arguments, typically `-install`.
 * `install-aminal-0.9.0.exe` performs the following steps:
   * It installs Aminal 0.9.0 into `%LOCALAPPDATA%\Aminal`.
   * It sets some registry keys in `HKCU\Software\Aminal`. This lets the `AminalUpdateTask...` tasks know which version of Aminal is installed.
   * It creates a Start menu shortcut for starting Aminal.

When the update tasks run, you will see `AminalUpdate.exe` in the Windows Task Manager. They use the registry to send the current Aminal version to the update server and ask "is there a new version?". If yes, the server again responds with the URL to an `.exe` and some command line parameters. In Aminal's current setup, this too is `install-aminal-0.9.0.exe` (if the user had, say, Aminal 0.8.9 installed). But this time the command line flag is `-update`. The `.exe` again installs the current version of Aminal and updates the registry.

The offline installer `install-aminal-0.9.0.exe` is actually what's produced by the `mingw32-make installer-windows` command in the previous section. It is placed at `bin/windows/AminalSetup.exe` and supports the command line flags `-install`, `-update` or none. In the last case (i.e. when invoked without arguments), it acts as a normal offline installer to be invoked by the user, and does not set Omaha's registry keys. The source code for this installer lies in `windows/installer/installer.go`.

Due to the asynchronous nature of the update tasks, it can happen that Aminal is running while a new version is being downloaded / installed. To prevent this from breaking Aminal's running instance, Aminal's install dir `%LOCALAPPDATA%\Aminal` contains the following hierarchy:

 * `Aminal.exe`
 * `Versions/`
   * `0.9.0/`
     * `Aminal.exe`

The top-level `Aminal.exe` is a launcher that always invokes the latest version. When an update is downloaded / installed, it is placed in a new subfolder of `Versions/`. For instance:

 * `Aminal.exe`
 * `Versions/`
   * `0.9.0/`
   * `0.9.1/`

The next time the top-level `Aminal.exe` is invoked, it runs `Versions/0.9.1/Aminal.exe`.

The code for this top-level launcher is in `windows/launcher/launcher.go`. Its binary (and the current version subdirectory) is produced in `bin/windows` when you do `mingw32-make installer-windows`.

#### Forcing updates

By default, Aminal's (/Omaha's) update tasks only run once per day. To force an immediate update, perform the following steps:

 * Delete the registry key `HKCU\Software\Aminal\Updated\LastChecked`.
 * Run the task `AminalUpdateTask...UA` in the Windows Task Scheduler. Press `F5`. You'll see its result change to `0x41301`. This means it's currently running. You'll also see `AminalUpdate.exe` in the Task _Scheduler_. Keep refreshing with `F5` until both disappear.

#### Uninstalling Aminal

The installer above adds an uninstaller to the user's _Add or remove programs_ panel in Windows. When the user goes to the Control Panel and uninstalls Aminal this way, the install directory and start menu entries are removed. Further, the registry key `HKCU\Software\Aminal\Clients\{35B0CF1E-FBB0-486F-A1DA-BE3A41DDC780}` is removed. What's not removed immediately is Omaha. (So you'll still see the update tasks in the Task Scheduler.) But! The next time the update tasks run, they realize by looking at the registry that Aminal is no longer installed and uninstall themselves (and Omaha).

To work around some potential permission issues, the uninstaller is not implemented in Go (like the installer and launcher above). But via NSIS. The source code is in `windows/Uninstaller.nsi`. It's an "installer" whose sole purpose is to generate `bin/windows/Aminal/uninstall.exe`.

#### Releasing a new version via automatic updates

To release a new version, update the `VERSION` fields in `Makefile`. Then, invoke `mingw32-make installer-windows`. Log into the Omaha update server, add a new Version for Aminal that mirrors the one you set in `Makefile`. You'll need to add a trailing `.0` to the version number on the server, because Omaha uses four-tuples for versions (`0.9.1.0` instead of `0.9.1`). As the "File" for the version, upload `bin/windows/AminalSetup.exe`. Add two "Action"s: One for Event _install_, one for event _update_. In both cases, instruct the server to _Run_ `AminalSetup.exe`. For the install Event, supply Arguments `-install`. For the update Event, supply Arguments `-update`.

Once you have done this, the background update tasks on your users' systems will (over the next 24 hours) download and install the new version of Aminal from the server.