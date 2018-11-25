package hints

type Hint struct {
	Word             string
	StartX           uint16
	StartY           uint16
	Line             string
	Description      string
	BackgroundColour [3]float32
	ForegroundColour [3]float32
}

type hinter func(word string, context string, wordX uint16, wordY uint16) *Hint

var hinters = []hinter{}

func Get(word string, context string, wordX uint16, wordY uint16) *Hint {
	for _, exp := range hinters {
		if h := exp(word, context, wordX, wordY); h != nil {
			return h
		}
	}
	return nil
}

func NewHint(word string, context string, wordX uint16, wordY uint16) *Hint {
	return &Hint{
		Line:             context,
		Word:             word,
		StartX:           wordX,
		StartY:           wordY,
		BackgroundColour: [3]float32{0, 0, 0},
		ForegroundColour: [3]float32{0.2, 1, 0.2},
	}
}
