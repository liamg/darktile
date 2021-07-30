// +build darwin

package clipboard

import (
	"git.wow.st/gmp/clip"
	"golang.org/x/xerrors"
)

func set(text string) error {
	ok := clip.Set(text)
	if !ok {
		return xerrors.New("nothing to set string")
	}
	return nil
}

func get() (string, error) {
	return clip.Get(), nil
}
