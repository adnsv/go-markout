package markout

type null_impl struct {
}

func NewNULL() Writer {
	return &null_impl{}
}

func (w *null_impl) Close()                                                  {}
func (w *null_impl) CloseEx(func(ParagraphWriter))                           {}
func (w *null_impl) DisableOutput()                                          {}
func (w *null_impl) EnableOutput()                                           {}
func (w *null_impl) Para(a any)                                              {}
func (w *null_impl) Paraf(format string, args ...any)                        {}
func (w *null_impl) BeginSection(a any)                                      {}
func (w *null_impl) BeginSectionf(format string, args ...any)                {}
func (w *null_impl) EndSection()                                             {}
func (w *null_impl) Section(a any)                                           {}
func (w *null_impl) Sectionf(format string, args ...any)                     {}
func (w *null_impl) BeginTable(first_column any, other_columns ...any)       {}
func (w *null_impl) TableRow(first_cell any, other_cells ...any)             {}
func (w *null_impl) EndTable()                                               {}
func (w *null_impl) Table(columns []any, rows func(callback TableRowWriter)) {}
func (w *null_impl) ListTitle(a any)                                         {}
func (w *null_impl) ListTitlef(format string, args ...any)                   {}
func (w *null_impl) BeginOList()                                             {}
func (w *null_impl) BeginUList()                                             {}
func (w *null_impl) EndList()                                                {}
func (w *null_impl) ListItem(a any)                                          {}
func (w *null_impl) ListItemf(format string, args ...any)                    {}
func (w *null_impl) OList(func(ListWriter))                                  {}
func (w *null_impl) UList(func(ListWriter))                                  {}
