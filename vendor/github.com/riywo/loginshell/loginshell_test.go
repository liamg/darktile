package loginshell

import (
    "testing"
    "os"
    "fmt"
    "runtime"
)

func TestShell(t *testing.T) {
    shell, err := Shell()
    if err != nil {
        t.Error(err)
    }

    if runtime.GOOS == "windows" {
        if shell == "" {
            t.Error("Output is empty!")
        }
    } else {
        currentShell := os.Getenv("SHELL")
        if shell != currentShell {
            t.Error(fmt.Sprintf("Output: %s, Current login shell: %s", shell, currentShell))
        }
    }
}
