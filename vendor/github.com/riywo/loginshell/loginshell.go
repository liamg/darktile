package loginshell

import (
    "runtime"
    "fmt"
    "os"
    "os/exec"
    "os/user"
    "errors"
    "regexp"
    "strings"
)

func Shell() (string, error) {
    switch runtime.GOOS {
        case "linux":
          return LinuxShell()
        case "darwin":
          return DarwinShell()
    }
    return "", errors.New("Undefined GOOS: " + runtime.GOOS)
}

func LinuxShell() (string, error) {
    user, err := user.Current()
    if err != nil { return "", err }
    out, err := exec.Command("getent", "passwd", user.Uid).Output()
    if err != nil { return "", err }

    ent := strings.Split(strings.TrimSuffix(string(out), "\n"), ":")
    return ent[6], nil
}

func DarwinShell() (string, error) {
    dir := "Local/Default/Users/" + os.Getenv("USER")
    out, err := exec.Command("dscl", "localhost", "-read", dir, "UserShell").Output()
    if err != nil { return "", err }

    re := regexp.MustCompile("UserShell: (/[^ ]+)\n")
    matched := re.FindStringSubmatch(string(out))
    shell := matched[1]
    if shell == "" { return "", errors.New(fmt.Sprintf("Invalid output: %s", string(out))) }
    return shell, nil
}
