package loginshell

import (
    "testing"
    "os"
    "fmt"
)

func TestShell(t *testing.T) {
    shell, err := Shell()
    if err != nil {
        t.Error(err)
    }

    currentShell := os.Getenv("SHELL")
    if shell != currentShell {
        t.Error(fmt.Sprintf("Output: %s, Current login shell: %s", shell, currentShell))
    }
}
