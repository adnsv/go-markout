package markout

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type md_blocks struct {
	base_blocks
}

func (bb *md_blocks) para(s RawContent) {
	bb.putblock(s)
	bb.want_emptyln()
}

func (bb *md_blocks) heading(counters []int, s RawContent, aa *Attrs) {
	if bb.enabled() {
		level := len(counters)
		b := bytes.Buffer{}
		wrepeat(&b, level, []byte("########"))
		b.WriteByte(' ')
		b.Write(s)
		if aa != nil {
			segments := []string{}
			if aa.Identifier != "" {
				segments = append(segments, "#"+aa.Identifier)
			}
			for _, c := range aa.Classes {
				segments = append(segments, "."+c)
			}
			if len(aa.KeyVals) > 0 {
				kvs := []string{}
				for k, v := range aa.KeyVals {
					kvs = append(kvs, fmt.Sprintf("%s=%s", k, v))
				}
				sort.Strings(kvs)
				segments = append(segments, kvs...)
			}
			if len(segments) > 0 {
				b.WriteString(" {")
				b.WriteString(strings.Join(segments, " "))
				b.WriteByte('}')
			}
		}
		bb.putblock(b.Bytes())
	}
	bb.want_emptyln()
}

func (bb *md_blocks) list_title(s RawContent) {
	bb.putblock(s)
	bb.want_emptyln()
}

func (bb *md_blocks) list_level_start(counters []int) {
}

func (bb *md_blocks) list_level_done(counters []int) {
	if len(counters) == 1 {
		bb.want_emptyln()
	} else {
		bb.want_nextln()
	}
}

func (bb *md_blocks) list_item(counters []int, s RawContent) {
	if bb.enabled() {
		level := len(counters)
		counter := counters[level-1]
		if counter < 0 {
			// unordered
			bb.putblock_ex(level-1, "- ", s, "")
		} else {
			// ordered
			bb.putblock_ex(level-1, strconv.FormatInt(int64(counter), 10)+". ", s, "")
		}
	}
	bb.want_nextln()
}

func (bb *md_blocks) end_table() {
	if len(bb.table) > 1 {
		ww := bb.table.measure_cells(bb.table[0], nil)

		eol := []byte{'\n'}
		rdecor := table_decor{nil, []byte("-|-"), nil}
		cdecor := table_decor{nil, []byte(" | "), nil}
		rule := []byte("--------")

		bb.do_nextline()
		bb.table.print_row(bb.out, bb.table[0], &cdecor, ww)
		bb.out.Write(eol)
		bb.table.print_rule(bb.out, rule, &rdecor, ww)
		for _, row := range bb.table[1:] {
			bb.out.Write(eol)
			bb.table.print_row(bb.out, row, &cdecor, nil)
		}
		bb.table = bb.table[:0]
	}
	bb.want_emptyln()
}

func (bb *md_blocks) codeblock(lang string, s RawContent) {
	bb.want_emptyln()
	bb.do_nextline()
	bb.out.Write([]byte("```"))
	if lang != "" {
		bb.out.Write([]byte(lang))
	}
	bb.out.Write([]byte{'\n'})
	bb.out.Write(s)
	bb.out.Write([]byte("\n```"))
	bb.want_emptyln()
}
