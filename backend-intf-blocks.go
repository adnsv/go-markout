package markout

type bmode int

const (
	mnone = bmode(0)
	mflow = bmode(1 << iota)
	mtable
	mlist
)

// blocks is an internal interface to be implemented by markout backends
// for structural block level formatting.
type blocks interface {
	current_mode() bmode
	check_mode(bmode) bool

	enabled() bool
	push_disabled() // disable actual output (stacked state)
	pop_disabled()

	close()
	para(RawContent)

	heading(counters []int, s RawContent, aa *Attrs)
	sect_level_in()
	sect_level_out()
	sect_counters() []int

	begin_table(columns ...RawContent)
	table_row(cells ...RawContent)
	end_table()

	list_title(RawContent)
	list_level_in(int) // initial counter 0 for ordered, -1 for unordered
	list_level_out()
	list_counters() []int
	list_item(counters []int, s RawContent)
	list_level_start(counters []int)
	list_level_done(counters []int)

	codeblock(lang string, s RawContent)
}
