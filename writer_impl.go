package markout

import (
	"strings"

	"golang.org/x/exp/slices"
)

// writer_impl implements structured writing to the supported markout backends.
type writer_impl struct {
	on_close func()
	bb       blocks // block-level formatting
	p        printer_impl
}

// Close finalizes writer output.
func (w *writer_impl) CloseEx(postscriptum func(ParagraphWriter)) {
	if w.bb.current_mode() == mtable {
		w.bb.end_table()
	}

	for w.bb.current_mode() == mlist {
		w.bb.list_level_done(w.bb.list_counters())
		w.bb.list_level_out()
	}

	if len(w.bb.sect_counters()) > 0 {
		w.bb.sect_level_out()
	}

	if postscriptum != nil {
		postscriptum(w)
	}

	if w.on_close != nil {
		w.on_close()
		w.on_close = nil
	}

	w.bb.close()
}

func (w *writer_impl) Close() {
	w.CloseEx(nil)
}

func (w *writer_impl) do_print(a any) RawContent {
	w.p.buf.Reset()
	if w.bb.enabled() {
		w.p.Print(a)
	}
	return w.p.buf.Bytes()
}

func (w *writer_impl) do_printf(format string, args ...any) RawContent {
	w.p.buf.Reset()
	if w.bb.enabled() {
		w.p.Printf(format, args...)
	}
	return w.p.buf.Bytes()
}

func (w *writer_impl) handle_section(s RawContent) {
	w.bb.check_mode(mflow)
	cc := w.bb.sect_counters()
	cc[len(cc)-1]++
	w.bb.heading(cc, s)
}

func (w *writer_impl) handle_listitem(s RawContent) {
	w.bb.check_mode(mlist)
	cc := w.bb.list_counters()
	n := len(cc) - 1
	if cc[n] >= 0 {
		cc[n]++
	}
	w.bb.list_item(cc, s)
}

func (w *writer_impl) Para(a any) {
	if w.bb.check_mode(mflow) {
		w.bb.para(w.do_print(a))
	}
}

func (w *writer_impl) Paraf(format string, args ...any) {
	if w.bb.check_mode(mflow) {
		w.bb.para(w.do_printf(format, args...))
	}
}

func (w *writer_impl) BeginSection(a any) {
	if w.bb.check_mode(mflow) {
		w.handle_section(w.do_print(a))
		w.bb.sect_level_in()
	}
}

func (w *writer_impl) BeginSectionf(format string, args ...any) {
	if w.bb.check_mode(mflow) {
		w.handle_section(w.do_printf(format, args...))
		w.bb.sect_level_in()
	}
}

func (w *writer_impl) EndSection() {
	if w.bb.check_mode(mflow) {
		w.bb.sect_level_out()
	}
}

func (w *writer_impl) Section(a any) {
	if w.bb.check_mode(mflow) {
		w.handle_section(w.do_print(a))
	}
}

func (w *writer_impl) Sectionf(format string, args ...any) {
	if w.bb.check_mode(mflow) {
		w.handle_section(w.do_printf(format, args...))
	}
}

func (w *writer_impl) BeginTable(first_column any, other_columns ...any) {
	if w.bb.check_mode(mflow) {
		rr := make([]RawContent, 0, 1+len(other_columns))
		rr = append(rr, slices.Clone(w.do_print(first_column)))
		for _, c := range other_columns {
			rr = append(rr, slices.Clone(w.do_print(c)))
		}
		w.bb.begin_table(rr...)
	}
}

func (w *writer_impl) TableRow(first_cell any, other_cells ...any) {
	if w.bb.check_mode(mtable) && w.bb.enabled() {
		rr := make([]RawContent, 0, 1+len(other_cells))
		rr = append(rr, slices.Clone(w.do_print(first_cell)))
		for _, c := range other_cells {
			rr = append(rr, slices.Clone(w.do_print(c)))
		}
		w.bb.table_row(rr...)
	}
}

func (w *writer_impl) EndTable() {
	if w.bb.check_mode(mtable) {
		w.bb.end_table()
	}
}

func (w *writer_impl) ListTitle(a any) {
	if w.bb.check_mode(mflow) {
		w.bb.list_title(w.do_print(a))
	}
}

func (w *writer_impl) ListTitlef(format string, args ...any) {
	if w.bb.check_mode(mflow) {
		w.bb.list_title(w.do_printf(format, args...))
	}
}

func (w *writer_impl) BeginOList() {
	if w.bb.check_mode(mflow | mlist) {
		w.bb.list_level_in(0)
		w.bb.list_level_start(w.bb.list_counters())
	}
}

func (w *writer_impl) BeginUList() {
	if w.bb.check_mode(mflow | mlist) {
		w.bb.list_level_in(-1)
		w.bb.list_level_start(w.bb.list_counters())
	}
}

func (w *writer_impl) ListItem(a any) {
	if w.bb.check_mode(mlist) {
		w.handle_listitem(w.do_print(a))
	}
}

func (w *writer_impl) ListItemf(format string, args ...any) {
	if w.bb.check_mode(mlist) {
		w.handle_listitem(w.do_printf(format, args...))
	}
}

func (w *writer_impl) EndList() {
	if w.bb.check_mode(mlist) {
		w.bb.list_level_done(w.bb.list_counters())
		w.bb.list_level_out()
	}
}

func (w *writer_impl) Table(columns []any, rows func(TableRowWriter)) {
	if w.bb.check_mode(mflow) && w.bb.enabled() {
		if rows == nil {
			return
		}
		w.BeginTable(columns)
		defer w.EndTable()

		on_row := func(first_cell any, other_cells ...any) {
			w.TableRow(first_cell, other_cells...)
		}

		rows(on_row)
	}
}

func (w *writer_impl) OList(f func(iw ListWriter)) {
	if w.bb.check_mode(mflow|mlist) && w.bb.enabled() {
		if f == nil {
			return
		}
		w.BeginOList()
		defer w.EndList()
		f(w)
	}
}

func (w *writer_impl) UList(f func(iw ListWriter)) {
	if w.bb.check_mode(mflow|mlist) && w.bb.enabled() {
		if f == nil {
			return
		}
		w.BeginUList()
		defer w.EndList()
		f(w)
	}
}

func (w *writer_impl) DisableOutput() {
	w.bb.push_disabled()
}

func (w *writer_impl) EnableOutput() {
	w.bb.pop_disabled()
}

func (w *writer_impl) Codeblock(lang string, lines string) {
	w.bb.codeblock(lang, w.do_print(func(p Printer) {
		for i, ln := range strings.Split(lines, "\n") {
			if ln != "" && ln[len(ln)-1] == '\r' {
				ln = ln[:len(ln)-1]
			}
			if i > 0 {
				p.WriteRawBytes([]byte{'\n'})
			}
			p.CodeblockLine(ln)
		}
	}))
}
