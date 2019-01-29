// +build darwin

package platform

import (
	"os/exec"
)

func LaunchTarget(target string) error {
	return exec.Command("open", target).Run()
}