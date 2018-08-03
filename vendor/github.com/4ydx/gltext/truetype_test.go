package gltext

import (
	"golang.org/x/image/math/fixed"
	"os"
	"testing"
)

func TestRuneRangesSort(t *testing.T) {
	rr := make(RuneRanges, 0)
	r := RuneRange{Low: 400, High: 500}
	rr = append(rr, r)
	r = RuneRange{Low: 32, High: 127}
	rr = append(rr, r)

	if !rr.Validate() {
		t.Error("Not validating.")
	}

	previousMax := rune(0)
	for _, r := range rr {
		if r.Low < previousMax {
			t.Error("Unsorted")
		}
		if r.Low == previousMax {
			t.Error("Overlap")
		}
		previousMax = r.High
	}
}

func TestRuneRangesOverlap(t *testing.T) {
	rr := make(RuneRanges, 0)
	r := RuneRange{Low: 40, High: 50}
	rr = append(rr, r)
	r = RuneRange{Low: 30, High: 40}
	rr = append(rr, r)

	if rr.Validate() {
		t.Error("Expecting invalidity due to overlap.")
	}
}

func TestRuneRangesLowHigh(t *testing.T) {
	rr := make(RuneRanges, 0)
	r := RuneRange{Low: 40, High: 39}
	rr = append(rr, r)

	if rr.Validate() {
		t.Error("Expecting invalidity.")
	}
}

func TestGetGlyphIndex(t *testing.T) {
	runeRanges := make(RuneRanges, 0)

	r := RuneRange{Low: 30, High: 40}
	runeRanges = append(runeRanges, r)
	r = RuneRange{Low: 100, High: 400}
	runeRanges = append(runeRanges, r)

	if !runeRanges.Validate() {
		t.Error("Not validating properly.")
	}

	index := runeRanges.GetGlyphIndex(30)
	if index != 0 {
		t.Error("Bad index", index)
	}
	index = runeRanges.GetGlyphIndex(40)
	if index != 10 {
		t.Error("Bad index", index)
	}

	index = runeRanges.GetGlyphIndex(100)
	if index != 11 {
		t.Error("Bad index", index)
	}
	index = runeRanges.GetGlyphIndex(390)
	if index != 301 {
		t.Error("Bad index", index)
	}
}

func TestGetGlyphIndexEdge(t *testing.T) {
	IsDebug = true

	runeRanges := RuneRanges{{Low: 32, High: 40}}
	if !runeRanges.Validate() {
		t.Error("Not validating properly.")
	}

	char := ' '
	index := runeRanges.GetGlyphIndex(char)
	if index != 0 {
		t.Error("Bad index", index)
	}
	char = '('
	index = runeRanges.GetGlyphIndex(char)
	if index != 8 {
		t.Error("Bad index", index)
	}

	fd, err := os.Open("font/font_1_honokamin.ttf")
	if err != nil {
		panic(err)
	}
	defer fd.Close()

	scale := fixed.Int26_6(24)
	runesPerRow := fixed.Int26_6(3)
	config, err := NewTruetypeFontConfig(fd, scale, runeRanges, runesPerRow)
	if err != nil {
		panic(err)
	}
	config.Name = "font_1_honokamin"

	// save png for manual inspection
	err = config.Save("fontconfigs")
	if err != nil {
		panic(err)
	}
}
