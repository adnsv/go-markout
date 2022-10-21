package markout

import (
	"bytes"
)

type imode int

const (
	iflow = imode(1 << iota)
	ilink
)

// inlines is an internal interface to be implemented by markout backends
// for inline content formatting.
type inlines interface {
	close() error

	current_mode() imode
	check_mode(imode)
	check_not_mode(imode)

	put_str(*bytes.Buffer, string)
	put_raw(*bytes.Buffer, RawContent)
	code_str(*bytes.Buffer, string)
	code_raw(*bytes.Buffer, RawContent)
	begin_link(*bytes.Buffer, RawContent)
	end_link(*bytes.Buffer)
	begin_styled(b *bytes.Buffer, sty Style)
	end_styled(*bytes.Buffer)
	simple_link(b *bytes.Buffer, caption RawContent, url RawContent)
}
