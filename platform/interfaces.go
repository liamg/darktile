package platform

import (
	"io"
)

// Process represents a child process by pid or HPROCESS in a platform-independent way
type Process interface {
	io.Closer

	Wait() error
	// TODO: make useful stuff here
}

// Pty represents a pseudo-terminal either by pty/tty file pair or by HCON
type Pty interface {
	io.ReadWriteCloser

	Resize(x int, y int) error
	CreateGuestProcess(imagePath string) (Process, error)
	GetPlatformDependentSettings() PlatformDependentSettings
}
