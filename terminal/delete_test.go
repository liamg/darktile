package terminal

import "testing"

func TestDelete(t *testing.T) {
	terminal := &Terminal{
		lines: []Line{
			{
				Cells: []Cell{
					{
						r: 'a',
					},
					{
						r: 'b',
					},
					{
						r: 'c',
					},
					{
						r: 'd',
					},
					{
						r: 'e',
					},
				},
			},
			{
				Cells: []Cell{
					{
						r: 'f',
					},
					{
						r: 'g',
					},
					{
						r: 'h',
					},
					{
						r: 'i',
					},
					{
						r: 'j',
					},
				},
			},
			{
				Cells: []Cell{
					{
						r: 'k',
					},
					{
						r: 'l',
					},
					{
						r: 'm',
					},
					{
						r: 'n',
					},
					{
						r: 'o',
					},
				},
			},
		},
	}

	terminal.position = Position{
		Col:  3,
		Line: 1,
	}

	if err := terminal.delete(2); err != nil {
		t.Errorf("Delete failed: %s", err)
	}

	if len(terminal.lines) != 3 {
		t.Errorf("No. of lines has changed by deleting characters")
	}

	if "fgh" != terminal.lines[1].String() {
		t.Errorf("Unexpected string after deletion: %s", terminal.lines[1].String())
	}
	if "abcde" != terminal.lines[0].String() {
		t.Errorf("Unexpected string after deletion: %s", terminal.lines[0].String())
	}

	if "klmno" != terminal.lines[2].String() {
		t.Errorf("Unexpected string after deletion: %s", terminal.lines[2].String())
	}
}
