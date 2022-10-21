package markout

import (
	"io"

	"github.com/adnsv/go-markout/wcwidth"
)

type table_grid [][]RawContent

func (t table_grid) measure_cell(s RawContent) int {
	return wcwidth.StringCells(string(s))
}

func (t table_grid) measure_cells(cc []RawContent, col_widths *[]int) (widths []int) {
	r := make([]int, len(cc))
	for i := range r {
		w := t.measure_cell(cc[i])
		r[i] = w
		if col_widths != nil {
			if i < len(*col_widths) {
				if w > (*col_widths)[i] {
					(*col_widths)[i] = w
				}
			} else {
				*col_widths = append(*col_widths, w)
			}
		}
	}
	return r
}

type table_decor struct {
	l, c, r []byte
}

func (t table_grid) print_row(w io.Writer, cc []RawContent, decor *table_decor, col_widths []int) {
	w.Write(decor.l)
	for i := range cc {
		if i > 0 {
			w.Write(decor.c)
		}
		col := 0
		if i < len(col_widths) {
			col = col_widths[i]
		}
		adv := t.measure_cell(cc[i])
		if col < adv {
			col = adv
		}
		w.Write(cc[i])
		if col > adv && i+1 < len(cc) {
			wrepeat(w, col-adv, nil)
		}
	}
	w.Write(decor.r)
}

func (t table_grid) print_rule(w io.Writer, rule []byte, decor *table_decor, col_widths []int) {
	w.Write(decor.l)
	for i := range col_widths {
		if i > 0 {
			w.Write(decor.c)
		}
		wrepeat(w, col_widths[i], rule)
	}
	w.Write(decor.r)
}
