// +build linux

package platform

import (
	"os/exec"
)

func LaunchTarget(target string) error {
	return exec.Command("xdg-open", target).Run()
}
