// +build windows

package platform

import (
	"github.com/MaxRis/w32"
)

func LaunchTarget(target string) error {
	return w32.ShellExecute(0, "", target, "", "", w32.SW_SHOW)
}