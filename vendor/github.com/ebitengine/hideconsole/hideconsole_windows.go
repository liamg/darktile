// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2024 The Ebitengine Authors

package hideconsole

import (
	"fmt"

	"golang.org/x/sys/windows"
)

var (
	kernel32 = windows.NewLazySystemDLL("kernel32.dll")
	user32   = windows.NewLazySystemDLL("user32.dll")

	procFreeConsoleWindow = kernel32.NewProc("FreeConsole")
	procGetConsoleWindow  = kernel32.NewProc("GetConsoleWindow")
)

func freeConsole() error {
	r, _, e := procFreeConsoleWindow.Call()
	if int32(r) == 0 {
		if e != nil && e != windows.ERROR_SUCCESS {
			return fmt.Errorf("FreeConsole failed: %w", e)
		}
		return fmt.Errorf("FreeConsole returned 0")
	}
	return nil
}

func getConsoleWindow() windows.HWND {
	r, _, _ := procGetConsoleWindow.Call()
	return windows.HWND(r)
}

// hideConsoleWindow will hide the console window that is showing when
// compiling on Windows without specifying the '-ldflags "-Hwindowsgui"' flag.
func hideConsoleWindow() {
	// In Xbox, GetWindowThreadProcessId might not exist.
	if user32.NewProc("GetWindowThreadProcessId").Find() != nil {
		return
	}

	pid := windows.GetCurrentProcessId()

	// Get the process ID of the console's creator.
	var cpid uint32
	if _, err := windows.GetWindowThreadProcessId(getConsoleWindow(), &cpid); err != nil {
		// Even if closing the console fails, this is not harmful.
		// Ignore error.
		return
	}

	if pid == cpid {
		// The current process created its own console. Hide this.
		// Ignore error.
		_ = freeConsole()
	}
}

func init() {
	hideConsoleWindow()
}
