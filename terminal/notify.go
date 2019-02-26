package terminal

func newNotifier() *notifier {
	return &notifier{
		C: make(chan struct{}),
	}
}

type notifier struct {
	C chan struct{}
}

// Notify is used to signal an event in a non-blocking way.
func (n *notifier) Notify() {
	select {
	case n.C <- struct{}{}:
	default:
	}
}
