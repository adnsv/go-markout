package markout

// MultiWriter implements a writer that funnels its output to other writers.
type MultiWriter struct {
	targets map[Writer]struct{}
}

// NewMultiWriter constructs a MultiWriter that writes into targets.
func NewMultiWriter(targets ...Writer) *MultiWriter {
	r := &MultiWriter{targets: map[Writer]struct{}{}}
	for _, w := range targets {
		r.targets[w] = struct{}{}
	}
	return r
}

func (w *MultiWriter) DisableOutput() {
	for t := range w.targets {
		t.DisableOutput()
	}
}

func (w *MultiWriter) EnableOutput() {
	for t := range w.targets {
		t.EnableOutput()
	}
}

func (w *MultiWriter) Close() error {
	w.targets = map[Writer]struct{}{}
	return nil
}

func (w *MultiWriter) Para(a any) {
	for t := range w.targets {
		t.Para(a)
	}
}

func (w *MultiWriter) Paraf(format string, args ...any) {
	for t := range w.targets {
		t.Paraf(format, args...)
	}
}

func (w *MultiWriter) BeginSection(a any) {
	for t := range w.targets {
		t.BeginSection(a)
	}
}

func (w *MultiWriter) BeginSectionf(format string, args ...any) {
	for t := range w.targets {
		t.BeginSectionf(format, args...)
	}
}

func (w *MultiWriter) EndSection() {
	for t := range w.targets {
		t.EndSection()
	}
}

func (w *MultiWriter) Section(a any) {
	for t := range w.targets {
		t.Section(a)
	}
}

func (w *MultiWriter) Sectionf(format string, args ...any) {
	for t := range w.targets {
		t.Sectionf(format, args...)
	}
}

func (w *MultiWriter) BeginTable(first_column any, other_columns ...any) {
	for t := range w.targets {
		t.BeginTable(first_column, other_columns...)
	}
}

func (w *MultiWriter) TableRow(first_cell any, other_cells ...any) {
	for t := range w.targets {
		t.TableRow(first_cell, other_cells...)
	}
}

func (w *MultiWriter) EndTable() {
	for t := range w.targets {
		t.EndTable()
	}
}

func (w *MultiWriter) BeginUList() {
	for t := range w.targets {
		t.BeginUList()
	}
}

func (w *MultiWriter) BeginOList() {
	for t := range w.targets {
		t.BeginOList()
	}
}

func (w *MultiWriter) ListTitle(a any) {
	for t := range w.targets {
		t.ListTitle(a)
	}
}

func (w *MultiWriter) ListTitlef(format string, args ...any) {
	for t := range w.targets {
		t.ListTitlef(format, args...)
	}
}

func (w *MultiWriter) ListItem(a any) {
	for t := range w.targets {
		t.ListItem(a)
	}
}

func (w *MultiWriter) ListItemf(format string, args ...any) {
	for t := range w.targets {
		t.ListItemf(format, args...)
	}
}

func (w *MultiWriter) EndList() {
	for t := range w.targets {
		t.EndList()
	}
}

func (w *MultiWriter) Table(columns []any, rows func(cb TableWriter)) {
	if len(columns) == 0 || rows == nil {
		return
	}
	for t := range w.targets {
		t.BeginTable(columns[0], columns[1:]...)
		defer t.EndTable()
	}

	on_row := func(first_cell any, other_cells ...any) {
		w.TableRow(first_cell, other_cells...)
	}

	rows(on_row)
}

func (w *MultiWriter) OList(items func(ListWriter)) {
	if items == nil {
		return
	}
	for t := range w.targets {
		t.BeginOList()
		defer t.EndList()
	}
	items(w)
}

func (w *MultiWriter) UList(items func(ListWriter)) {
	if items == nil {
		return
	}
	for t := range w.targets {
		t.BeginUList()
		defer t.EndList()
	}
	items(w)
}
