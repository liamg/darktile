// +build windows

package main

import (
	"bufio"
	"errors"
	"golang.org/x/sys/windows/registry"
	"os"
	"os/user"
	"strings"
	"path/filepath"
	"github.com/liamg/aminal/windows/winutil"
	"github.com/liamg/aminal/generated-src/installer/data"
	"text/template"
	"io/ioutil"
	"os/exec"
	"syscall"
	"flag"
)

const Version = "VERSION"
const ProductId = `{35B0CF1E-FBB0-486F-A1DA-BE3A41DDC780}`

func main() {
	doInstallPtr := flag.Bool("install", false, "Install Aminal")
	doUpdatePtr := flag.Bool("update", false, "Update Aminal")
	flag.Parse()

	var installDir string
	isUserInstall := strings.HasPrefix(os.Args[0], os.Getenv("LOCALAPPDATA"))
	if *doInstallPtr {
		installDir = getInstallDirWhenManagedByOmaha()
		extractAssets(installDir)
		createRegistryKeysForUninstaller(installDir, isUserInstall)
		updateVersionInRegistry(isUserInstall)
		createStartMenuShortcut(installDir, isUserInstall)
		launchAminal(installDir)
	} else if *doUpdatePtr {
		installDir = getInstallDirWhenManagedByOmaha()
		extractAssets(installDir)
		updateVersionInRegistry(isUserInstall)
		removeOldVersions(installDir)
	} else {
		// Offline installer.
		// We don't know whether we're being executed with Admin privileges.
		// It's also not easy to determine. So perform user install:
		isUserInstall = true
		installDir = getDefaultInstallDir()
		extractAssets(installDir)
		createRegistryKeysForUninstaller(installDir, isUserInstall)
		createStartMenuShortcut(installDir, isUserInstall)
		launchAminal(installDir)
	}
}

func getInstallDirWhenManagedByOmaha() string {
	executablePath, err := winutil.GetExecutablePath()
	check(err)
	result := executablePath
	prevResult := ""
	for filepath.Base(result) != "Aminal" {
		prevResult = result
		result = filepath.Dir(result)
		if result == prevResult {
			break
		}
	}
	if result == prevResult {
		msg := "Could not find parent directory 'Aminal' above " + executablePath
		check(errors.New(msg))
	}
	return result
}

func getDefaultInstallDir() string {
	localAppData := os.Getenv("LOCALAPPDATA")
	if localAppData == "" {
		panic("Environment variable LOCALAPPDATA is not set.")
	}
	return filepath.Join(localAppData, "Aminal")
}

func extractAssets(installDir string) {
	for _, relPath := range data.AssetNames() {
		bytes, err := data.Asset(relPath)
		check(err)
		absPath := filepath.Join(installDir, relPath)
		check(os.MkdirAll(filepath.Dir(absPath), 0755))
		f, err := os.OpenFile(absPath, os.O_CREATE, 0755)
		check(err)
		defer f.Close()
		w := bufio.NewWriter(f)
		_, err = w.Write(bytes)
		check(err)
		w.Flush()
	}
}

func createRegistryKeysForUninstaller(installDir string, isUserInstall bool) {
	regRoot := getRegistryRoot(isUserInstall)
	uninstKey := `Software\Microsoft\Windows\CurrentVersion\Uninstall\Aminal`
	writeRegStr(regRoot, uninstKey, "", installDir)
	writeRegStr(regRoot, uninstKey, "DisplayName", "Aminal")
	writeRegStr(regRoot, uninstKey, "Publisher", "Liam Galvin")
	uninstaller := filepath.Join(installDir, "uninstall.exe")
	uninstString := `"` + uninstaller + `"`
	if isUserInstall {
		uninstString += " /CurrentUser"
	} else {
		uninstString += " /AllUsers"
	}
	writeRegStr(regRoot, uninstKey, "UninstallString", uninstString)
}

func updateVersionInRegistry(isUserInstall bool) {
	regRoot := getRegistryRoot(isUserInstall)
	updateKey := `Software\Aminal\Update\Clients\` + ProductId
	writeRegStr(regRoot, updateKey, "pv", Version + ".0")
	writeRegStr(regRoot, updateKey, "name", "Aminal")
}

func getRegistryRoot(isUserInstall bool) registry.Key {
	if isUserInstall {
		return registry.CURRENT_USER
	}
	return registry.LOCAL_MACHINE
}

func writeRegStr(regRoot registry.Key, keyPath string, valueName string, value string) {
	const mode = registry.WRITE|registry.WOW64_32KEY
	key, _, err := registry.CreateKey(regRoot, keyPath, mode)
	check(err)
	defer key.Close()
	check(key.SetStringValue(valueName, value))
}

func createStartMenuShortcut(installDir string, isUserInstall bool) {
	startMenuDir := getStartMenuDir(isUserInstall)
	linkPath := filepath.Join(startMenuDir, "Programs", "Aminal.lnk")
	targetPath := filepath.Join(installDir, "Aminal.exe")
	createShortcut(linkPath, targetPath)
}

func getStartMenuDir(isUserInstall bool) string {
	if isUserInstall {
		usr, err := user.Current()
		check(err)
		return usr.HomeDir + `\AppData\Roaming\Microsoft\Windows\Start Menu`
	} else {
		return os.Getenv("ProgramData") + `\Microsoft\Windows\Start Menu`
	}
}

func createShortcut(linkPath, targetPath string) {
	type Shortcut struct {
		LinkPath string
		TargetPath string
	}
	tmpl := template.New("createLnk.vbs")
	tmpl, err := tmpl.Parse(`Set oWS = WScript.CreateObject("WScript.Shell")
sLinkFile = "{{.LinkPath}}"
Set oLink = oWS.CreateShortcut(sLinkFile)
oLink.TargetPath = "{{.TargetPath}}"
oLink.Save
WScript.Quit 0`)
	check(err)
	tmpDir, err := ioutil.TempDir("", "Aminal")
	check(err)
	createLnk := filepath.Join(tmpDir, "createLnk.vbs")
	defer os.RemoveAll(tmpDir)
	f, err := os.Create(createLnk)
	check(err)
	defer f.Close()
	w := bufio.NewWriter(f)
	shortcut := Shortcut{linkPath, targetPath}
	check(tmpl.Execute(w, shortcut))
	w.Flush()
	f.Close()
	cmd := exec.Command("cscript", f.Name())
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	check(cmd.Run())
}

func launchAminal(installDir string) {
	cmd := exec.Command(getAminalExePath(installDir))
	check(cmd.Start())
}

func getAminalExePath(installDir string) string {
	return filepath.Join(installDir, "Aminal.exe")
}

func removeOldVersions(installDir string) {
	versionsDir := filepath.Join(installDir, "Versions")
	versions, err := ioutil.ReadDir(versionsDir)
	check(err)
	for _, version := range versions {
		if version.Name() == Version {
			continue
		}
		versionPath := filepath.Join(versionsDir, version.Name())
		// Try deleting the main executable first. We do this to prevent a
		// version that is still running from being deleted.
		mainExecutable := filepath.Join(versionPath, "Aminal.exe")
		err = os.Remove(mainExecutable)
		if err == nil {
			// Remove the rest:
			check(os.RemoveAll(versionPath))
		}
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}