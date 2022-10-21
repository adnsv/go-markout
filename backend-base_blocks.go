package markout

import (
	"errors"
	"io"
)

type base_blocks struct {
	out             io.Writer
	disable_counter int
	eols            int   // number of eols pending
	sect_levels     []int // section level counters
	list_levels     []int // list level counters (-1 for unordered levels)
	table           table_grid
}

func (bb *base_blocks) current_mode() bmode {
	if len(bb.table) > 0 {
		return mtable
	} else if len(bb.list_levels) > 0 {
		return mlist
	} else if len(bb.sect_levels) > 0 {
		return mflow
	} else {
		return mnone // aborted
	}
}

func (bb *base_blocks) check_mode(wanted bmode) bool {
	m := bb.current_mode()
	if m == mnone {
		return false
	}
	if wanted&m != 0 {
		return true
	} else {
		panic("markout: command is not allowed within the current block mode")
	}
}

func (bb *base_blocks) close() {
	bb.table = bb.table[:0]
	bb.sect_levels = bb.sect_levels[:0]
	bb.list_levels = bb.list_levels[:0]
	bb.do_nextline()
	bb.eols = 0
}

func (bb *base_blocks) want_nextln() {
	if bb.eols == 0 || bb.enabled() {
		bb.eols = 1
	}
}

func (bb *base_blocks) want_emptyln() {
	bb.eols = 2
}

func wrepeat(w io.Writer, nbytes int, sequence []byte) {
	if nbytes == 0 {
		return
	}
	if len(sequence) == 0 {
		const space_sequence = "                "
		sequence = []byte(space_sequence)
	}
	n := len(sequence)
	for nbytes >= n {
		w.Write(sequence)
		nbytes -= n
	}
	w.Write(sequence[:nbytes])
}

func (bb *base_blocks) enabled() bool {
	return bb.disable_counter == 0
}

func (bb *base_blocks) push_disabled() {
	bb.disable_counter++
}

func (bb *base_blocks) pop_disabled() {
	if bb.disable_counter > 0 {
		bb.disable_counter--
	}
}

func (bb *base_blocks) begin_table(columns ...RawContent) {
	hh := []RawContent{}
	hh = append(hh, columns...)
	bb.table = append(bb.table, hh)
}

func (bb *base_blocks) table_row(cells ...RawContent) {
	if bb.enabled() {
		cc := []RawContent{}
		cc = append(cc, cells...)
		bb.table = append(bb.table, cc)
	}
}

func (bb *base_blocks) do_nextline() {
	n := bb.eols
	if n > 2 {
		n = 2
	}
	bb.out.Write([]byte("\n\n")[:bb.eols])
}

func (bb *base_blocks) putblock(s RawContent) {
	if bb.enabled() {
		bb.do_nextline()
		bb.out.Write(s)
	}
}

func (bb *base_blocks) putblock_ex(indent_level int, prefix string, s RawContent, postfix string) {
	if bb.enabled() {
		bb.do_nextline()
		wrepeat(bb.out, 2*indent_level, nil)
		bb.out.Write([]byte(prefix))
		bb.out.Write(s)
		bb.out.Write([]byte(postfix))
	}
}

func (bb *base_blocks) sect_level_in() {
	bb.sect_levels = append(bb.sect_levels, 0)
}

func (bb *base_blocks) sect_level_out() {
	n := len(bb.sect_levels)
	bb.sect_levels = bb.sect_levels[:n-1]
}

func (bb *base_blocks) sect_counters() []int {
	return bb.sect_levels
}

func (bb *base_blocks) list_level_in(initial int) {
	bb.list_levels = append(bb.list_levels, initial)
}

func (bb *base_blocks) list_level_out() {
	n := len(bb.list_levels)
	bb.list_levels = bb.list_levels[:n-1]
}

func (bb *base_blocks) list_counters() []int {
	return bb.list_levels
}

var errSectLevel = errors.New("markout: invalid section level (unpaired section level calls)")
var errListLevel = errors.New("markout: invalid list level (unpaired list level calls)")

func pick[T any](use_second bool, first, second T) T {
	if use_second {
		return second
	} else {
		return first
	}
}
