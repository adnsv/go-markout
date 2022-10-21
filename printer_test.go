package markout

import (
	"bytes"
	"testing"
)

type MyStruct struct {
	v int
}

// MarshalMarkoutInline implements InlineMarshaler support for MyStruct.
func (m MyStruct) MarshalMarkoutInline(w Printer) error {
	w.BeginStyled(SingleQuotedStyle)
	w.WriteString("q")
	w.CodeString("c[]<&>")
	w.EndStyled()
	w.BeginLink("url")
	w.WriteString("a")
	w.EndLink()
	w.CodeString("c`d") // handle backtick within the code
	return nil
}

func Test_to_buffer(t *testing.T) {
	buf := bytes.Buffer{}
	html_w := html_inlines{}
	md_w := md_inlines{}
	text_w := txt_inlines{}
	html_w.setup_quotation_marks(TypographicalQuotes)
	md_w.setup_quotation_marks(TypographicalQuotes)
	text_w.setup_quotation_marks(TypographicalQuotes)

	testms := &MyStruct{42}

	tests := []struct {
		name      string
		a         any
		want_html string
		want_md   string
		want_text string
	}{
		{"empty string", "", "", "", ""},
		{"simple string", "abc", "abc", "abc", "abc"},
		{"string with specials", "<a>", "&lt;a&gt;", "\\<a\\>", "<a>"},
		{"string with backslash", "a\\c", "a\\c", "a\\\\c", "a\\c"},
		{"int", 42, "42", "42", "42"},
		{"bool", true, "true", "true", "true"},
		{"floating point", 3.14159, "3.14159", "3.14159", "3.14159"},
		{"raw", RawContent("<>@#$%`@"), "<>@#$%`@", "<>@#$%`@", "<>@#$%`@"},
		{"simple-link-full", Link("c", "url"), "<a href=\"url\">c</a>", "[c](url)", "[c](url)"},
		{"simple-link-urls", Link("url", "url"), "<a href=\"url\">url</a>", "[url]", "url"},
		{"simple-link-nocaption", Link("", "url"), "<a href=\"url\">url</a>", "[url]", "url"},
		{"custom marshaler", &testms,
			"‘q<code>c[]&lt;&&gt;</code>’<a href=\"url\">a</a><code>c`d</code>",
			"‘q`c[]<&>`’[a](url)<code>c`d</code>",
			"‘q`c[]<&>`’[a](url)`c`d`"},
	}
	for _, tt := range tests {
		buf.Reset()
		to_buffer(&buf, &html_w, nil, tt.a)
		got_html := buf.String()
		buf.Reset()
		to_buffer(&buf, &md_w, nil, tt.a)
		got_md := buf.String()
		buf.Reset()
		to_buffer(&buf, &text_w, nil, tt.a)
		got_text := buf.String()

		t.Run("html: "+tt.name, func(t *testing.T) {
			if got_html != tt.want_html {
				t.Errorf("HTML fmt_any() returned %q, want %q", got_html, tt.want_html)
			}
		})
		t.Run("markdown: "+tt.name, func(t *testing.T) {
			if got_md != tt.want_md {
				t.Errorf("MD fmt_any() returned %q, want %q", got_md, tt.want_md)
			}
		})
		t.Run("text: "+tt.name, func(t *testing.T) {
			if got_text != tt.want_text {
				t.Errorf("TXT fmt_any() returned %q, want %q", got_text, tt.want_text)
			}
		})
	}
}
