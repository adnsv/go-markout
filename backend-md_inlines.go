package markout

import (
	"bytes"
	"strings"

	"golang.org/x/exp/slices"
)

const want_backslash_escaped = "\\`*_{}[]()#+-.!<>"

type md_inlines struct {
	base_inlines
	pending_link *RawContent
}

func md_scramble(b *bytes.Buffer, s string) {
	i, o, n := 0, 0, len(s)
	if n <= 0 {
		return
	}
	for i < n {
		c := s[i]
		i++
		if c == '\n' {
			b.WriteString(s[o : i-1])
			b.WriteString("\\n")
			o = i
		} else if strings.IndexByte(want_backslash_escaped, c) >= 0 {
			b.WriteString(s[o : i-1])
			b.WriteByte('\\')
			b.WriteByte(c)
			o = i
		}
	}
	b.WriteString(s[o:n])
}

func md_scramble_code(b *bytes.Buffer, s string) {
	// note: will fail if both `</code>` and a backtick are in s
	if strings.IndexByte(s, '`') >= 0 {
		b.WriteString("<code>")
		b.WriteString(s)
		b.WriteString("</code>")
	} else {
		b.WriteByte('`')
		b.WriteString(s)
		b.WriteByte('`')
	}
}

func (ii *md_inlines) current_mode() imode {
	return pick(ii.pending_link != nil, iflow, iflow|ilink)
}

func (ii *md_inlines) check_mode(wanted imode) {
	m := ii.current_mode()
	if wanted&m == 0 {
		panic(err_inline_mode)
	}
}

func (ii *md_inlines) check_not_mode(not_wanted imode) {
	m := ii.current_mode()
	if not_wanted&m != 0 {
		panic(err_inline_mode)
	}
}

func (ii *md_inlines) close() error {
	if ii.pending_link != nil {
		return errUnpairedBeginLink
	} else if len(ii.style_stack) > 0 {
		return errUnpairedBeginStyled
	} else {
		return nil
	}
}

func (ii *md_inlines) put_raw(b *bytes.Buffer, s RawContent) {
	b.Write(s)
}
func (ii *md_inlines) put_str(b *bytes.Buffer, s string) {
	md_scramble(b, s)
}
func (ii *md_inlines) code_raw(b *bytes.Buffer, s RawContent) {
	b.WriteByte('`')
	b.Write(s)
	b.WriteByte('`')
}
func (ii *md_inlines) code_str(b *bytes.Buffer, s string) {
	md_scramble_code(b, s)
}
func (ii *md_inlines) codeblock_line(b *bytes.Buffer, s string) {
	b.Write([]byte(s)) // todo: deal with "```"
}
func (ii *md_inlines) begin_styled(b *bytes.Buffer, sty Style) {
	ii.start_styled(sty)
	switch sty {
	case SingleQuotedStyle:
		b.WriteString(ii.quote_specs[0])
	case DoubleQuotedStyle:
		b.WriteString(ii.quote_specs[2])
	case StrongStyle:
		b.WriteString("<strong>") // "**" or "__" is not reliable enough
	case EmphasizedStyle:
		b.WriteString("<em>") // "*" or "_" is not reliable enough
	}
}
func (ii *md_inlines) end_styled(b *bytes.Buffer) {
	sty := ii.finish_styled()
	switch sty {
	case SingleQuotedStyle:
		b.WriteString(ii.quote_specs[1])
	case DoubleQuotedStyle:
		b.WriteString(ii.quote_specs[3])
	case StrongStyle:
		b.WriteString("</strong>")
	case EmphasizedStyle:
		b.WriteString("</em>")
	}
}
func (ii *md_inlines) begin_link(b *bytes.Buffer, url RawContent) {
	b.WriteByte('[')
	ii.pending_link = &url
}
func (ii *md_inlines) end_link(b *bytes.Buffer) {
	b.WriteString("](")
	b.Write(*ii.pending_link)
	b.WriteByte(')')
	ii.pending_link = nil
}
func (ii *md_inlines) simple_link(b *bytes.Buffer, caption RawContent, url RawContent) {
	if len(caption) == 0 || slices.Equal(caption, url) {
		b.WriteByte('[')
		b.Write(url)
		b.WriteByte(']')
	} else {
		b.WriteByte('[')
		b.Write(caption)
		b.WriteString("](")
		b.Write(url)
		b.WriteByte(')')
	}
}
