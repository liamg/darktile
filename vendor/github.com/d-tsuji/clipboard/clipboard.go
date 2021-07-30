package clipboard

// Get returns the current text data of the clipboard.
func Get() (string, error) {
	return get()
}

// Set sets the current text data of the clipboard.
func Set(text string) error {
	return set(text)
}
