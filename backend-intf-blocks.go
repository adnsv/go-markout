package markout

type bmode int

const (
	mflow = bmode(1 << iota)
	mtable
	mlist
)

// blocks is an internal interface to be implemented by markout backends
// for structural block level formatting.
type blocks interface {
	close() error

	check_mode(bmode)

	enabled() bool
	push_disabled() // disable actual output (stacked state)
	pop_disabled()

	para(raw_bytes)

	heading(counters []int, s raw_bytes)
	sect_level_in()
	sect_level_out()
	sect_counters() []int

	begin_table(columns ...raw_bytes)
	table_row(cells ...raw_bytes)
	end_table()

	list_title(raw_bytes)
	list_level_in(int) // initial counter 0 for ordered, -1 for unordered
	list_level_out()
	list_counters() []int
	list_item(counters []int, s raw_bytes)
	list_level_start(counters []int)
	list_level_done(counters []int)
	// todo: code blocks, etc.
}
