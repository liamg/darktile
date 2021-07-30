package fontinfo

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Font represents a font file on disk
type Font struct {
	Family string
	Style  string
	Path   string
}

var validExtensions = []string{
	".ttf",
	".otf",
}

var fontDirs = []string{
	"~/.fonts",
	"~/.local/share/fonts",
	"/usr/local/share/fonts",
	"/usr/share/fonts",
	filepath.Join(os.Getenv("XDG_DATA_HOME"), "fonts"),
	filepath.Join(os.Getenv("XDG_DATA_DIRS"), "fonts"),
}

// Match finds all fonts installed on the system which match the provided matchers
func Match(matchers ...matcher) ([]Font, error) {

	var fonts []Font
	meta := make(map[string]*fontMetadata)

	for _, dir := range fontDirs {

		if info, err := os.Stat(dir); os.IsNotExist(err) {
			continue
		} else if !info.IsDir() {
			continue
		}

		if err := filepath.WalkDir(dir, func(path string, info fs.DirEntry, err error) error {
			if _, ok := meta[path]; ok {
				return nil
			}
			ext := filepath.Ext(path)
			for _, valid := range validExtensions {
				if strings.EqualFold(ext, valid) {
					f, err := os.Open(path)
					if err != nil {
						return err
					}
					defer f.Close()
					metadata, err := readMetadata(f)
					if err != nil {
						return err
					}
					for _, match := range matchers {
						if !match(metadata) {
							return nil
						}
					}
					meta[path] = metadata
					return nil
				}
			}
			return nil
		}); err != nil {
			return nil, err
		}
	}

	for path, metadata := range meta {
		fonts = append(fonts, Font{
			Family: metadata.FontFamily,
			Style:  metadata.FontStyle,
			Path:   path,
		})
	}

	return fonts, nil
}
