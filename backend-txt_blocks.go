package markout

import (
	"bytes"
	"strconv"

	"github.com/adnsv/go-markout/wcwidth"
)

type txt_blocks struct {
	base_blocks
	numbered_sections   bool
	underlined_sections bool
	listitem_prefix     string
}

func (bb *txt_blocks) para(s RawContent) {
	bb.putblock(s)
	bb.want_emptyln()
}

func (bb *txt_blocks) heading(counters []int, s RawContent, _ *Attrs) {
	if bb.enabled() {
		level := len(counters)

		if bb.numbered_sections || (bb.underlined_sections && level <= 2) {
			b := bytes.Buffer{}
			if bb.numbered_sections {
				for _, i := range counters {
					b.WriteString(strconv.Itoa(i))
					b.WriteByte('.')
				}
				b.WriteByte(' ')
				b.Write(s)
			}
			if bb.underlined_sections && level <= 2 {
				width := wcwidth.StringCells(b.String())
				b.WriteByte('\n')
				if level == 1 {
					wrepeat(&b, width, []byte("========"))
				} else {
					wrepeat(&b, width, []byte("--------"))
				}
			}
			s = b.Bytes()
		}

		bb.putblock(s)
	}

	bb.want_emptyln()
}

func (bb *txt_blocks) list_title(s RawContent) {
	bb.putblock(s)
	bb.want_nextln()
}

func (bb *txt_blocks) list_level_start(counters []int) {
}

func (bb *txt_blocks) list_level_done(counters []int) {
	if len(counters) == 1 {
		bb.want_emptyln()
	} else {
		bb.want_nextln()
	}
}

func (bb *txt_blocks) list_item(counters []int, s RawContent) {
	if bb.enabled() {
		level := len(counters)
		counter := counters[level-1]
		if counter < 0 {
			// unordered
			bb.putblock_ex(level-1, bb.listitem_prefix, s, "")
		} else {
			// ordered
			bb.putblock_ex(level-1, strconv.FormatInt(int64(counter), 10)+". ", s, "")
		}
	}
	bb.want_nextln()
}

func (bb *txt_blocks) end_table() {
	if len(bb.table) > 1 {
		cols := []int{}
		for r := range bb.table {
			bb.table.measure_cells(bb.table[r], &cols)
		}
		eol := []byte{'\n'}
		decor := table_decor{nil, []byte{' '}, nil}
		rule := []byte("--------")
		bb.do_nextline()
		bb.table.print_row(bb.out, bb.table[0], &decor, cols)
		bb.out.Write(eol)
		bb.table.print_rule(bb.out, rule, &decor, cols)
		for _, row := range bb.table[1:] {
			bb.out.Write(eol)
			bb.table.print_row(bb.out, row, &decor, cols)
		}
	}
	bb.table = bb.table[:0]
	bb.want_emptyln()
}

func (bb *txt_blocks) codeblock(lang string, s RawContent) {
	bb.want_emptyln()
	bb.do_nextline()
	bb.out.Write(s)
	bb.want_emptyln()
}
