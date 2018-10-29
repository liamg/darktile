package hints

type Hint struct {
	Word        string
	StartX      uint16
	StartY      uint16
	Line        string
	Description string
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
