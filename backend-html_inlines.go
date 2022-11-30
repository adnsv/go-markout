package markout

import (
	"bytes"
	"fmt"
)

type html_inlines struct {
	base_inlines
	in_link bool
}

func (ii *html_inlines) current_mode() imode {
	return pick(ii.in_link, iflow, iflow|ilink)
}

func (ii *html_inlines) check_mode(wanted imode) {
	m := ii.current_mode()
	if wanted&m == 0 {
		panic(err_inline_mode)
	}
}

func (ii *html_inlines) check_not_mode(not_wanted imode) {
	m := ii.current_mode()
	if not_wanted&m != 0 {
		panic(err_inline_mode)
	}
}

func (ii *html_inlines) close() error {
	if ii.in_link {
		return errUnpairedBeginLink
	} else if len(ii.style_stack) > 0 {
		return errUnpairedBeginStyled
	} else {
		return nil
	}
}

func (ii *html_inlines) put_raw(b *bytes.Buffer, s RawContent) {
	b.Write(s)
}
func (ii *html_inlines) put_str(b *bytes.Buffer, s string) {
	html_scramble(b, s)
}
func (ii *html_inlines) code_raw(b *bytes.Buffer, s RawContent) {
	b.WriteString("<code>")
	b.Write(s)
	b.WriteString("</code>")
}
func (ii *html_inlines) code_str(b *bytes.Buffer, s string) {
	b.WriteString("<code>")
	html_scramble(b, s)
	b.WriteString("</code>")
}
func (ii *html_inlines) codeblock_line(b *bytes.Buffer, s string) {
	html_scramble(b, s)
}
func (ii *html_inlines) begin_styled(b *bytes.Buffer, sty Style) {
	ii.start_styled(sty)
	switch sty {
	case SingleQuotedStyle:
		b.WriteString(ii.quote_specs[0])
	case DoubleQuotedStyle:
		b.WriteString(ii.quote_specs[2])
	case StrongStyle:
		b.WriteString("<strong>")
	case EmphasizedStyle:
		b.WriteString("<em>")
	}
}
func (ii *html_inlines) end_styled(b *bytes.Buffer) {
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
func (ii *html_inlines) begin_link(b *bytes.Buffer, url RawContent) {
	fmt.Fprintf(b, "<a href=\"%s\">", url)
	ii.in_link = true
}
func (ii *html_inlines) end_link(b *bytes.Buffer) {
	ii.in_link = false
	b.WriteString("</a>")
}

func html_scramble(b *bytes.Buffer, s string) {
	i, o, n := 0, 0, len(s)
	if n <= 0 {
		return
	}
	for i < n {
		c := s[i]
		i++
		if c == '<' {
			b.WriteString(s[o : i-1])
			b.WriteString("&lt;")
			o = i
		} else if c == '>' {
			b.WriteString(s[o : i-1])
			b.WriteString("&gt;")
			o = i
		}
	}
	b.WriteString(s[o:n])
}
