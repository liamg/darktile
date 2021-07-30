package font

import (
	"fmt"
	"image"
	"math"
	"os"

	"github.com/liamg/darktile/internal/app/darktile/packed"
	"github.com/liamg/fontinfo"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

type Style uint8

const (
	Regular Style = iota
	Bold
	Italic
	BoldItalic
)

type StyleName string

const (
	StyleRegular    StyleName = "Regular"
	StyleBold       StyleName = "Bold"
	StyleItalic     StyleName = "Italic"
	StyleBoldItalic StyleName = "Bold Italic"
)

type Manager struct {
	family         string
	regularFace    font.Face
	boldFace       font.Face
	italicFace     font.Face
	boldItalicFace font.Face
	size           float64
	dpi            float64
	charSize       image.Point
	fontDotDepth   int
}

func NewManager() *Manager {
	return &Manager{
		size: 16,
		dpi:  72,
	}
}

func (m *Manager) CharSize() image.Point {
	return m.charSize
}

func (m *Manager) IncreaseSize() {
	m.SetSize(m.size + 1)
}

func (m *Manager) DecreaseSize() {
	if m.size < 2 {
		return
	}
	m.SetSize(m.size - 1)
}

func (m *Manager) DotDepth() int {
	return m.fontDotDepth
}

func (m *Manager) DPI() float64 {
	return m.dpi
}

func (m *Manager) SetDPI(dpi float64) error {
	if dpi <= 0 {
		return fmt.Errorf("DPI must be >0")
	}
	m.dpi = dpi
	return nil
}

func (m *Manager) SetSize(size float64) error {
	m.size = size
	if m.regularFace != nil {
		// effectively reload fonts at new size
		m.SetFontByFamilyName(m.family)
	}
	return nil
}

func (m *Manager) loadFontFace(path string) (font.Face, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	fnt, err := opentype.ParseReaderAt(f)
	if err != nil {
		return nil, err
	}
	return m.createFace(fnt)
}

func (m *Manager) createFace(f *opentype.Font) (font.Face, error) {
	return opentype.NewFace(f, &opentype.FaceOptions{
		Size:    m.size,
		DPI:     m.dpi,
		Hinting: font.HintingFull,
	})
}

func (m *Manager) SetFontByFamilyName(name string) error {

	m.family = name

	if name == "" {
		return m.loadDefaultFonts()
	}

	fonts, err := fontinfo.Match(fontinfo.MatchFamily(name))
	if err != nil {
		return err
	}

	if len(fonts) == 0 {
		return fmt.Errorf("could not find font with family '%s'", name)
	}

	for _, fontMeta := range fonts {
		switch StyleName(fontMeta.Style) {
		case StyleRegular:
			m.regularFace, err = m.loadFontFace(fontMeta.Path)
			if err != nil {
				return err
			}
		case StyleBold:
			m.boldFace, err = m.loadFontFace(fontMeta.Path)
			if err != nil {
				return err
			}
		case StyleItalic:
			m.italicFace, err = m.loadFontFace(fontMeta.Path)
			if err != nil {
				return err
			}
		case StyleBoldItalic:
			m.boldItalicFace, err = m.loadFontFace(fontMeta.Path)
			if err != nil {
				return err
			}
		}
	}

	if m.regularFace == nil {
		return fmt.Errorf("could not find regular style for font family '%s'", name)
	}

	return m.calcMetrics()
}

func (m *Manager) calcMetrics() error {

	face := m.regularFace

	var prevAdvance int
	for ch := rune(32); ch <= 126; ch++ {
		adv26, ok := face.GlyphAdvance(ch)
		if ok && adv26 > 0 {
			advance := int(adv26)
			if prevAdvance > 0 && prevAdvance != advance {
				return fmt.Errorf("the specified font is not monospaced: %d 0x%X=%d", prevAdvance, ch, advance)
			}
			prevAdvance = advance
		}
	}

	if prevAdvance == 0 {
		return fmt.Errorf("failed to calculate advance width for font face")
	}

	metrics := face.Metrics()

	m.charSize.X = int(math.Round(float64(prevAdvance) / m.dpi))
	m.charSize.Y = int(math.Round(float64(metrics.Height) / m.dpi))
	m.fontDotDepth = int(math.Round(float64(metrics.Ascent) / m.dpi))

	return nil
}

func (m *Manager) loadDefaultFonts() error {

	regular, err := opentype.Parse(packed.MesloLGSNFRegularTTF)
	if err != nil {
		return err
	}
	m.regularFace, err = m.createFace(regular)
	if err != nil {
		return err
	}

	bold, err := opentype.Parse(packed.MesloLGSNFBoldTTF)
	if err != nil {
		return err
	}
	m.boldFace, err = m.createFace(bold)
	if err != nil {
		return err
	}

	italic, err := opentype.Parse(packed.MesloLGSNFItalicTTF)
	if err != nil {
		return err
	}
	m.italicFace, err = m.createFace(italic)
	if err != nil {
		return err
	}

	boldItalic, err := opentype.Parse(packed.MesloLGSNFBoldItalicTTF)
	if err != nil {
		return err
	}
	m.boldItalicFace, err = m.createFace(boldItalic)
	if err != nil {
		return err
	}

	return m.calcMetrics()
}

func (m *Manager) RegularFontFace() font.Face {
	return m.regularFace
}

func (m *Manager) BoldFontFace() font.Face {
	if m.boldFace == nil {
		return m.RegularFontFace()
	}
	return m.boldFace
}

func (m *Manager) ItalicFontFace() font.Face {
	if m.italicFace == nil {
		return m.RegularFontFace()
	}
	return m.italicFace
}

func (m *Manager) BoldItalicFontFace() font.Face {
	if m.boldItalicFace == nil {
		if m.boldFace == nil {
			return m.ItalicFontFace()
		}
		return m.BoldFontFace()
	}
	return m.boldItalicFace
}
