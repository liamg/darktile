package hinters

import (
	"sort"
	"sync"

	"github.com/liamg/darktile/internal/app/darktile/termutil"
)

type HinterRegistration struct {
	Priority Priority
	Hinter   Hinter
}

type Priority uint8

const (
	PriorityNone     = 0
	PriorityVeryLow  = 48
	PriorityLow      = 96
	PriorityMedium   = 128
	PriorityHigh     = 192
	PriorityVeryHigh = 224
	PriorityCritical = 255
)

type Hinter interface {
	// Match should return the index in the text of the matched occurrence
	Match(text string, cursorIndex int) (matched bool, offset int, length int)
	// Activate fires when mouseover happens afer a match - takes raw coords
	Activate(api HintAPI, match string, start termutil.Position, end termutil.Position) error
	Deactivate(api HintAPI) error
	Click(api HintAPI) error
}

var hintLock sync.RWMutex
var hinters []HinterRegistration

func register(h Hinter, p Priority) {
	hintLock.Lock()
	defer hintLock.Unlock()
	hinters = append(hinters, HinterRegistration{
		Priority: p,
		Hinter:   h,
	})
}

func All() []Hinter {
	hintLock.RLock()
	defer hintLock.RUnlock()

	var output []Hinter

	sort.Slice(hinters, func(i, j int) bool {
		return hinters[i].Priority > hinters[j].Priority
	})

	for _, hinter := range hinters {
		output = append(output, hinter.Hinter)
	}

	return output
}
