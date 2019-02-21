//go:generate goversioninfo -icon=aminal.ico

/*
Looks at directory "Versions" next to this executable. Finds the latest version
and runs the executable with the same name as this executable in that directory.
Eg.:
  Aminal.exe (=launcher.exe)
  Versions/
    1.0.0/
      Aminal.exe
    1.0.1
      Aminal.exe
-> Launches Versions/1.0.1/Aminal.exe.
*/

package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"github.com/liamg/aminal/windows/winutil"
)

type Version struct {
	number [3]int
	name string
}
type Versions []Version

func main() {
	executable, err := winutil.GetExecutablePath()
	check(err)
	executableDir, executableName := filepath.Split(executable)
	versionsDir := filepath.Join(executableDir, "Versions")
	latestVersion, err := getLatestVersion(versionsDir)
	check(err)
	target := filepath.Join(versionsDir, latestVersion, executableName)
	cmd := exec.Command(target, os.Args[1:]...)
	check(cmd.Start())
}

func getLatestVersion(versionsDir string) (string, error) {
	potentialVersions, err := ioutil.ReadDir(versionsDir)
	if err != nil {
		return "", err
	}
	var versions Versions
	for _, file := range potentialVersions {
		if !file.IsDir() {
			continue
		}
		version, err := parseVersionString(file.Name())
		if err != nil {
			continue
		}
		versions = append(versions, version)
	}
	if len(versions) == 0 {
		errMsg := fmt.Sprintf("No valid version in %s.", versionsDir)
		return "", errors.New(errMsg)
	}
	sort.Sort(versions)
	return versions[len(versions)-1].String(), nil
}

func parseVersionString(version string) (Version, error) {
	var result Version
	result.name = version
	err := error(nil)
	version = strings.TrimSuffix(version, "-SNAPSHOT")
	parts := strings.Split(version, ".")
	if len(parts) != len(result.number) {
		err = errors.New("Wrong number of parts.")
	} else {
		for i, partStr := range parts {
			result.number[i], err = strconv.Atoi(partStr)
			if err != nil {
				break
			}
		}
	}
	return result, err
}

func (arr Versions) Len() int {
	return len(arr)
}

func (arr Versions) Less(i, j int) bool {
	for k, left := range arr[i].number {
		right := arr[j].number[k]
		if left > right {
			return false
		} else if left < right {
			return true
		}
	}
	fmt.Printf("%s < %s\n", arr[i], arr[j])
	return true
}

func (arr Versions) Swap(i, j int) {
	tmp := arr[j]
	arr[j] = arr[i]
	arr[i] = tmp
}

func (version Version) String() string {
	return version.name
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}