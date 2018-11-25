package config

import (
	"fmt"
	"strings"

	"github.com/go-gl/glfw/v3.2/glfw"
)

type KeyCombination struct {
	mods glfw.ModifierKey
	key  glfw.Key
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

var keyMap = map[string]glfw.Key{
	"a": glfw.KeyA,
	"b": glfw.KeyB,
	"c": glfw.KeyC,
	"d": glfw.KeyD,
	"e": glfw.KeyE,
	"f": glfw.KeyF,
	"g": glfw.KeyG,
	"h": glfw.KeyH,
	"i": glfw.KeyI,
	"j": glfw.KeyJ,
	"k": glfw.KeyK,
	"l": glfw.KeyL,
	"m": glfw.KeyM,
	"n": glfw.KeyN,
	"o": glfw.KeyO,
	"p": glfw.KeyP,
	"q": glfw.KeyQ,
	"r": glfw.KeyR,
	"s": glfw.KeyS,
	"t": glfw.KeyT,
	"u": glfw.KeyU,
	"v": glfw.KeyV,
	"w": glfw.KeyW,
	"x": glfw.KeyX,
	"y": glfw.KeyY,
	"z": glfw.KeyZ,
	";": glfw.KeySemicolon,
}

// keyStr e.g. "ctrl + alt + a"
func parseKeyCombination(keyStr string) (*KeyCombination, error) {

	var mods glfw.ModifierKey
	var key *glfw.Key

	keys := strings.Split(keyStr, "+")
	for _, k := range keys {
		k = strings.ToLower(strings.TrimSpace(k))
		mod, ok := modMap[KeyMod(k)]
		if ok {
			mods = mods + mod
			continue
		}
		mappedKey, ok := keyMap[k]
		if ok {
			if key != nil {
				return nil, fmt.Errorf("Multiple non-modifier keys specified in keyboard shortcut")
			}
			key = &mappedKey
			continue
		}

		return nil, fmt.Errorf("Unknown key '%s' in configured keyboard shortcut", k)
	}

	if key == nil {
		return nil, fmt.Errorf("No non-modifier key specified in keyboard shortcut")
	}

	if mods == 0 {
		return nil, fmt.Errorf("No modifier key specified in keyboard shortcut")
	}

	return &KeyCombination{
		mods: mods,
		key:  *key,
	}, nil
}

func (combi KeyCombination) Match(pressedMods glfw.ModifierKey, pressedKey glfw.Key) bool {
	return pressedKey == combi.key && pressedMods == combi.mods
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
