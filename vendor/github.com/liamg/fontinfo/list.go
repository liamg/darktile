package fontinfo

// List all otf/ttf fonts installed on the system
func List() ([]Font, error) {
	return Match()
}
