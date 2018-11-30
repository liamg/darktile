package config

import (
	"fmt"
	"strings"

	"github.com/go-gl/glfw/v3.2/glfw"
)

type KeyCombination struct {
	mods glfw.ModifierKey
	char rune
}

type KeyMod string

const (
	ctrl  KeyMod = "ctrl"
	alt   KeyMod = "alt"
	shift KeyMod = "shift"
	super KeyMod = "super"
)

var modMap = map[KeyMod]glfw.ModifierKey{
	ctrl:  glfw.ModControl,
	alt:   glfw.ModAlt,
	shift: glfw.ModShift,
	super: glfw.ModSuper,
}

// keyStr e.g. "ctrl + alt + a"
func parseKeyCombination(keyStr string) (*KeyCombination, error) {

	var mods glfw.ModifierKey
	var key rune

	keys := strings.Split(keyStr, "+")
	for _, k := range keys {
		k = strings.ToLower(strings.TrimSpace(k))
		mod, ok := modMap[KeyMod(k)]
		if ok {
			mods = mods + mod
			continue
		}

		if key > 0 {
			return nil, fmt.Errorf("Multiple non-modifier keys specified in keyboard shortcut")
		}

		key = rune(k[0])
	}

	if key == 0 {
		return nil, fmt.Errorf("No non-modifier key specified in keyboard shortcut")
	}

	if mods == 0 {
		return nil, fmt.Errorf("No modifier key specified in keyboard shortcut")
	}

	return &KeyCombination{
		mods: mods,
		char: key,
	}, nil
}

func (combi KeyCombination) Match(pressedMods glfw.ModifierKey, pressedChar rune) bool {
	return pressedChar == combi.char && pressedMods == combi.mods
}

func (keyMapConfig KeyMappingConfig) GenerateActionMap() (map[UserAction]*KeyCombination, error) {
	m := map[UserAction]*KeyCombination{}
	for actionStr, keyStr := range keyMapConfig {
		combi, err := parseKeyCombination(keyStr)
		if err != nil {
			return nil, err
		}
		m[UserAction(actionStr)] = combi
	}

	return m, nil
}
