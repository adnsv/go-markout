package markout

import (
	"bytes"
	"io"
)

// SectionWriter is an interface that supports structured sectioning of a
// document.
type SectionWriter = interface {
	// BeginSection writes the section heading block and increments section level
	// counter. Each BeginSection() call must be followed by matching EndSection()
	BeginSection(a any)
	BeginSectionf(format string, args ...any)
	BeginAttrSection(aa Attrs, a any)
	BeginAttrSectionf(aa Attrs, format string, args ...any)

	// EndSection decrements section level counter.
	EndSection()

	// Section writes section heading without incrementing section level counter.
	Section(a any)
	Sectionf(format string, args ...any)
	AttrSection(aa Attrs, a any)
	AttrSectionf(aa Attrs, format string, args ...any)
}

// TableRowWriter is a callback for writing table rows.
type TableRowWriter = func(first_cell any, other_cells ...any)

// ParagraphWriter interface supports writing plain paragraphs.
type ParagraphWriter = interface {
	// Para writes paragraph block.
	Para(a any)
	Paraf(format string, args ...any)
}

// TableWriter interface supports writing of tabular data.
type TableWriter = interface {
	// BeginTable starts table mode.
	//   - only TableRow() calls are supported in table mode.
	//   - use EndTable() to exit from the table mode.
	BeginTable(first_column any, other_columns ...any)
	TableRow(first_cell any, other_cells ...any)
	EndTable()

	// Callback-based table writing method that begins a new table, writes rows
	// into it with callback, then calls EndTable.
	Table(columns []any, rows func(callback TableRowWriter))
}

// ListWriter is an interface for writing items and child lists into list
// blocks.
type ListWriter interface {
	// ListTitle writes paragraph block that acts as a title preceeding a list.
	ListTitle(a any)
	ListTitlef(format string, args ...any)

	// BeginOList and BeginUList begins a list block (ordered or unourdered). If
	// writer is already in list mode, this begins a child list. Each BeginList
	// must be matched with EndList.
	BeginOList()
	BeginUList()
	EndList()

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

type CodeblockWriter interface {
	Codeblock(lang string, lines string)
}

// Writer is a high level interface for writing markout documents in all
// supported formats.
type Writer interface {
	SectionWriter
	ParagraphWriter
	ListWriter
	TableWriter
	CodeblockWriter

	Close()
	CloseEx(ps func(ParagraphWriter))

	DisableOutput()
	EnableOutput()
}

// Attr can be used to add id, classes, and attributes to section headings.
type Attrs struct {
	Identifier string
	Classes    []string
	KeyVals    map[string]string
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
	return &writer_impl{
		bb: bb,
		p:  printer_impl{ii: ii, buf: &bytes.Buffer{}, url_filter: opts.URLFilter},
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

	r := &writer_impl{
		bb: bb,
		p:  printer_impl{ii: ii, buf: &bytes.Buffer{}, url_filter: opts.URLFilter},
		on_close: func() {
			bb.end_body()
			bb.end_html()
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
	return &writer_impl{
		bb: bb,
		p:  printer_impl{ii: ii, buf: &bytes.Buffer{}, url_filter: opts.URLFilter},
	}
}
