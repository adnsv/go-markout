package markout

import (
	"bytes"
	"io"
)

// TableWriter is a callback for writing table rows.
type TableWriter = func(first_cell any, other_cells ...any)

// Writer is a high level interface for writing markout documents in all
// supported formats.
type Writer interface {
	ListWriter
	io.Closer

	DisableOutput()
	EnableOutput()

	// Para writes paragraph block.
	Para(a any)
	Paraf(format string, args ...any)

	// BeginSection writes the section heading block and increments section level
	// counter. Each BeginSection() call must be followed by matching EndSection()
	BeginSection(a any)
	BeginSectionf(format string, args ...any)

	// EndSection decrements section level counter.
	EndSection()

	// Section writes section heading without incrementing section level counter.
	Section(a any)
	Sectionf(format string, args ...any)

	// BeginTable starts table mode.
	//   - only TableRow() calls are supported while in tablt mode.
	//   - use EndTable() to exit from the table mode.
	BeginTable(first_column any, other_columns ...any)
	TableRow(first_cell any, other_cells ...any)
	EndTable()

	// Callback-based table writing method that begins a new table, writes rows
	// into it with callback, then calls EndTable.
	Table(columns []any, rows func(callback TableWriter))

	// ListTitle writes paragraph block that acts as a title preceeding a list.
	ListTitle(a any)
	ListTitlef(format string, args ...any)

	// BeginOList and BeginUList begins a list block (ordered or unourdered). If
	// writer is already in list mode, this begins a child list. Each BeginList
	// must be matched with EndList.
	BeginOList()
	BeginUList()
	EndList()
}

// ListWriter is an interface for writing items and child lists into list
// blocks.
type ListWriter interface {
	// ListItem writes an item into the list block. Must be called within the
	// BeginList()/EndList() fragment.
	ListItem(a any)

	// ListItemf is a version of ListItem() with built-in formatting.
	ListItemf(format string, args ...any)

	// Callback-based list writing methods that automatically wrap items
	// BeginList/EndList blocks.
	OList(func(ListWriter))
	UList(func(ListWriter))
}

// Quotation marks
const (
	ASCIIQuotes         = `'|'|"|"`
	TypographicalQuotes = `‘|’|“|”`
)

type TXTOptions struct {
	PutBOM             bool
	QuotationMarks     string // pipe-separated single and double quotes (defaults to '|'|"|")
	ListItemPrefix     string // content inserted before each list item (defaults to `* `)
	UnderlinedSections bool
	NumberedSections   bool
	URLFilter          url_filter
}

// NewTxt creates a new markout writer targeting plain text output.
func NewTXT(out io.Writer, opts TXTOptions) Writer {
	ii := &txt_inlines{}
	ii.setup_quotation_marks(opts.QuotationMarks)
	bb := &txt_blocks{}
	bb.out = out
	bb.listitem_prefix = opts.ListItemPrefix
	bb.underlined_sections = opts.UnderlinedSections
	bb.numbered_sections = opts.NumberedSections
	if bb.listitem_prefix == "" {
		bb.listitem_prefix = "* "
	}
	if opts.PutBOM {
		bb.out.Write(RawContent("uFEFF"))
	}

	bb.sect_level_in()
	return &fw_impl{
		blocks:   bb,
		p:        printer_impl{ii: ii, buf: &bytes.Buffer{}, url_filter: opts.URLFilter},
		on_close: bb.close,
	}
}

type HTMLOptions struct {
	PutBOM         bool
	QuotationMarks string
	Title          string
	Style          string
	ListTitleClass string
	URLFilter      url_filter
}

// NewHtml creates a new markout writer targeting html output.
func NewHTML(out io.Writer, opts HTMLOptions) Writer {
	ii := &html_inlines{}
	ii.setup_quotation_marks(opts.QuotationMarks)
	bb := &html_blocks{}
	bb.out = out
	bb.list_title_class = opts.ListTitleClass

	r := &fw_impl{
		blocks: bb,
		p:      printer_impl{ii: ii, buf: &bytes.Buffer{}, url_filter: opts.URLFilter},
		on_close: func() error {
			err := bb.close()
			if err != nil {
				return err
			}
			bb.end_body()
			bb.end_html()
			return nil
		},
	}

	if opts.PutBOM {
		bb.out.Write(RawContent("uFEFF"))
	}
	bb.begin_html()
	bb.head(r.do_print(opts.Title), RawContent(opts.Style))
	bb.begin_body()
	bb.sect_level_in()

	return r
}

type MDOptions struct {
	PutBOM         bool
	QuotationMarks string // pipe-separated single and double quotes (defaults to '|'|"|")
	URLFilter      url_filter
}

// NewMD creates a new markout writer targeting markdown output.
func NewMD(out io.Writer, opts MDOptions) Writer {
	ii := &md_inlines{}
	ii.setup_quotation_marks(opts.QuotationMarks)
	bb := &md_blocks{}
	bb.out = out
	if opts.PutBOM {
		bb.out.Write(RawContent("uFEFF"))
	}
	bb.sect_level_in()
	return &fw_impl{
		blocks:   bb,
		p:        printer_impl{ii: ii, buf: &bytes.Buffer{}, url_filter: opts.URLFilter},
		on_close: bb.close,
	}
}
