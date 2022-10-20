package markout

import (
	"bytes"

	"golang.org/x/exp/slices"
)

type txt_inlines struct {
	base_inlines
	pending_link *raw_bytes
}

func txt_scramble(s string) raw_bytes {
	return raw_bytes(s)
}

func (ii *txt_inlines) current_mode() imode {
	return pick(ii.pending_link != nil, iflow, iflow|ilink)
}

func (ii *txt_inlines) check_mode(wanted imode) {
	m := ii.current_mode()
	if wanted&m == 0 {
		panic(err_inline_mode)
	}
}

func (ii *txt_inlines) check_not_mode(not_wanted imode) {
	m := ii.current_mode()
	if not_wanted&m != 0 {
		panic(err_inline_mode)
	}
}

func (ii *txt_inlines) close() error {
	if ii.pending_link != nil {
		return errUnpairedBeginLink
	} else if len(ii.style_stack) > 0 {
		return errUnpairedBeginStyled
	} else {
		return nil
	}
}

func (ii *txt_inlines) put_raw(b *bytes.Buffer, s raw_bytes) {
	b.Write(s)
}
func (ii *txt_inlines) put_str(b *bytes.Buffer, s string) {
	b.Write(txt_scramble(s))
}
func (ii *txt_inlines) code_raw(b *bytes.Buffer, s raw_bytes) {
	b.WriteByte('`')
	b.Write(s)
	b.WriteByte('`')
}
func (ii *txt_inlines) code_str(b *bytes.Buffer, s string) {
	b.WriteByte('`')
	b.Write(txt_scramble(s))
	b.WriteByte('`')
}
func (ii *txt_inlines) begin_styled(b *bytes.Buffer, sty Style) {
	ii.start_styled(sty)
	switch sty {
	case SingleQuotedStyle:
		b.WriteString(ii.quote_specs[0])
	case DoubleQuotedStyle:
		b.WriteString(ii.quote_specs[2])
	case StrongStyle:
		b.WriteString("**")
	case EmphasizedStyle:
		b.WriteString("*")
	}
}
func (ii *txt_inlines) end_styled(b *bytes.Buffer) {
	sty := ii.finish_styled()
	switch sty {
	case SingleQuotedStyle:
		b.WriteString(ii.quote_specs[1])
	case DoubleQuotedStyle:
		b.WriteString(ii.quote_specs[3])
	case StrongStyle:
		b.WriteString("**")
	case EmphasizedStyle:
		b.WriteString("*")
	}
}
func (ii *txt_inlines) begin_link(b *bytes.Buffer, url raw_bytes) {
	b.WriteByte('[')
	ii.pending_link = &url
}
func (ii *txt_inlines) end_link(b *bytes.Buffer) {
	b.WriteString("](")
	b.Write(*ii.pending_link)
	b.WriteByte(')')
	ii.pending_link = nil
}

func (ii *txt_inlines) simple_link(b *bytes.Buffer, caption raw_bytes, url raw_bytes) {
	if len(caption) == 0 || slices.Equal(caption, url) {
		b.Write(url)
	} else {
		b.WriteByte('[')
		b.Write(caption)
		b.WriteString("](")
		b.Write(url)
		b.WriteByte(')')
	}
}
