package sixel

import (
	"image/color"
	"strings"
	"testing"
)

// A sixel header is DCS P1 ; P2 ; P3 q. Anything a program can emit has to be
// survivable: a single-parameter header used to index params[1] and take the
// whole terminal down with it.
func TestDecodeHeaderVariants(t *testing.T) {
	tests := []struct {
		name    string
		header  string
		wantErr bool
	}{
		{name: "no parameters", header: ""},
		{name: "P1 only", header: "0"},
		{name: "P1 and P2", header: "0;1"},
		{name: "P1, P2 and P3", header: "0;1;8"},
		{name: "empty parameters", header: ";;"},
		{name: "aspect ratio 2", header: "2"},
		{name: "invalid P1", header: "x", wantErr: true},
		{name: "too many parameters", header: "0;1;8;9", wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := Decode(strings.NewReader(test.header+"q#0;2;0;0;0#0~"), color.Black)
			if test.wantErr {
				if err == nil {
					t.Fatalf("expected an error for header %q", test.header)
				}
				return
			}
			if err != nil {
				t.Fatalf("header %q: %v", test.header, err)
			}
		})
	}
}

func TestDecodeAspectRatioFromP1(t *testing.T) {
	tests := map[string]float64{
		"0": 2, "1": 2, "5": 2, "6": 2,
		"2": 5,
		"3": 3, "4": 3,
		"7": 1, "8": 1, "9": 1,
	}

	for p1, want := range tests {
		t.Run("P1="+p1, func(t *testing.T) {
			d := NewDecoder(strings.NewReader(p1+";1q"), color.Black)
			if err := d.processHeader(); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if d.aspectRatio != want {
				t.Fatalf("P1=%s: aspect ratio = %v, want %v", p1, d.aspectRatio, want)
			}
		})
	}
}

// Repeat counts and raster attributes are program-supplied, so they have to be
// treated as hostile: an unbounded one used to size an image allocation.
func TestDecodeRejectsOversizedImages(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"huge repeat", "q!999999~"},
		{"huge raster width", `q"1;1;99999;10~`},
		{"huge raster height", `q"1;1;10;99999~`},
		{"negative raster", `q"1;1;-5;-5~`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if _, err := Decode(strings.NewReader(test.input), color.Black); err == nil {
				t.Fatal("expected the image to be refused")
			}
		})
	}
}

// Sixel data before any colour register is selected reads back a nil colour,
// which the image writer cannot convert.
func TestDecodeDataBeforeColourSelection(t *testing.T) {
	img, err := Decode(strings.NewReader("q~~~"), color.Black)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if img.Bounds().Empty() {
		t.Fatal("expected a non-empty image")
	}
}
