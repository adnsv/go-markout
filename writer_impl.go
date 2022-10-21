package markout

import "golang.org/x/exp/slices"

// fw_impl implements structured writing to the supported markout backends.
type fw_impl struct {
	on_close func() error
	blocks   blocks // block-level formatting
	p        printer_impl
}

// Close finalizes writer output.
func (w *fw_impl) Close() error {
	if w.on_close != nil {
		err := w.on_close()
		w.on_close = nil
		return err
	} else {
		return nil
	}
}

func (w *fw_impl) do_print(a any) RawContent {
	w.p.buf.Reset()
	if w.blocks.enabled() {
		w.p.Print(a)
	}
	return w.p.buf.Bytes()
}

func (w *fw_impl) do_printf(format string, args ...any) RawContent {
	w.p.buf.Reset()
	if w.blocks.enabled() {
		w.p.Printf(format, args...)
	}
	return w.p.buf.Bytes()
}

func (w *fw_impl) do_section(s RawContent) {
	w.blocks.check_mode(mflow)
	cc := w.blocks.sect_counters()
	cc[len(cc)-1]++
	w.blocks.heading(cc, s)
}

func (w *fw_impl) do_listitem(s RawContent) {
	w.blocks.check_mode(mlist)
	cc := w.blocks.list_counters()
	n := len(cc) - 1
	if cc[n] >= 0 {
		cc[n]++
	}
	w.blocks.list_item(cc, s)
}

func (w *fw_impl) Para(a any) {
	w.blocks.check_mode(mflow)
	w.blocks.para(w.do_print(a))
}

func (w *fw_impl) Paraf(format string, args ...any) {
	w.blocks.check_mode(mflow)
	w.blocks.para(w.do_printf(format, args...))
}

func (w *fw_impl) BeginSection(a any) {
	w.do_section(w.do_print(a))
	w.blocks.sect_level_in()
}

func (w *fw_impl) BeginSectionf(format string, args ...any) {
	w.do_section(w.do_printf(format, args...))
	w.blocks.sect_level_in()
}

func (w *fw_impl) EndSection() {
	w.blocks.check_mode(mflow)
	w.blocks.sect_level_out()
}

func (w *fw_impl) Section(a any) {
	w.do_section(w.do_print(a))
}

func (w *fw_impl) Sectionf(format string, args ...any) {
	w.do_section(w.do_printf(format, args...))
}

func (w *fw_impl) BeginTable(first_column any, other_columns ...any) {
	w.blocks.check_mode(mflow)
	rr := make([]RawContent, 0, 1+len(other_columns))
	rr = append(rr, slices.Clone(w.do_print(first_column)))
	for _, c := range other_columns {
		rr = append(rr, slices.Clone(w.do_print(c)))
	}
	w.blocks.begin_table(rr...)
}

func (w *fw_impl) TableRow(first_cell any, other_cells ...any) {
	w.blocks.check_mode(mtable)
	if w.blocks.enabled() {
		rr := make([]RawContent, 0, 1+len(other_cells))
		rr = append(rr, slices.Clone(w.do_print(first_cell)))
		for _, c := range other_cells {
			rr = append(rr, slices.Clone(w.do_print(c)))
		}
		w.blocks.table_row(rr...)
	}
}

func (w *fw_impl) EndTable() {
	w.blocks.check_mode(mtable)
	w.blocks.end_table()
}

func (w *fw_impl) ListTitle(a any) {
	w.blocks.check_mode(mflow)
	w.blocks.list_title(w.do_print(a))
}

func (w *fw_impl) ListTitlef(format string, args ...any) {
	w.blocks.check_mode(mflow)
	w.blocks.list_title(w.do_printf(format, args...))
}

func (w *fw_impl) BeginOList() {
	w.blocks.check_mode(mflow | mlist)
	w.blocks.list_level_in(0)
	w.blocks.list_level_start(w.blocks.list_counters())
}

func (w *fw_impl) BeginUList() {
	w.blocks.check_mode(mflow | mlist)
	w.blocks.list_level_in(-1)
	w.blocks.list_level_start(w.blocks.list_counters())
}

func (w *fw_impl) ListItem(a any) {
	w.do_listitem(w.do_print(a))
}

func (w *fw_impl) ListItemf(format string, args ...any) {
	w.do_listitem(w.do_printf(format, args...))
}

func (w *fw_impl) EndList() {
	w.blocks.check_mode(mlist)
	w.blocks.list_level_done(w.blocks.list_counters())
	w.blocks.list_level_out()
}

func (w *fw_impl) Table(columns []any, rows func(TableWriter)) {
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

func (w *fw_impl) OList(f func(iw ListWriter)) {
	if f == nil {
		return
	}
	w.BeginOList()
	defer w.EndList()
	f(w)
}

func (w *fw_impl) UList(f func(iw ListWriter)) {
	if f == nil {
		return
	}
	w.BeginUList()
	defer w.EndList()
	f(w)
}

func (w *fw_impl) DisableOutput() {
	w.blocks.push_disabled()
}

func (w *fw_impl) EnableOutput() {
	w.blocks.pop_disabled()
}
