package gui

import (
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	KeyPressDelayNS  = 500_000_000
	KeyPressRepeatNS = 30_000_000
	KeyPressResetNS  = 60_000_000
)

type keyState struct {
	mu   sync.Mutex
	keys map[ebiten.Key]press
}

func newKeyState() *keyState {
	return &keyState{
		keys: make(map[ebiten.Key]press),
	}
}

type press struct {
	at        int64
	repeating bool
}

func (k *keyState) AnythingPressed() bool {
	k.mu.Lock()
	defer k.mu.Unlock()
	return len(k.keys) > 0
}

func (k *keyState) RepeatPressed(key ebiten.Key) bool {
	now := time.Now().UnixNano()
	k.mu.Lock()
	defer k.mu.Unlock()

	if ebiten.IsKeyPressed(key) {

		event, ok := k.keys[key]
		if !ok {
			k.keys[key] = press{at: now}
			return true
		}

		since := now - event.at
		if !event.repeating && since > int64(KeyPressDelayNS) {
			k.keys[key] = press{at: now, repeating: true}
			return true
		} else if event.repeating && since > int64(KeyPressRepeatNS) {
			k.keys[key] = press{at: now, repeating: true}
			return true
		}

		return false
	}

	delete(k.keys, key)
	return false
}
