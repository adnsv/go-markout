package markout

import (
	"bytes"
	"fmt"
)

// things related exclusively to HTML file structure

func (bb *html_blocks) begin_html() {
	bb.putblock(RawContent("<html>"))
	bb.want_nextln()
}

func (bb *html_blocks) end_html() {
	bb.want_nextln()
	bb.putblock(RawContent("</html>"))
	bb.want_nextln()
}

func (bb *html_blocks) head(title, style RawContent) {
	if len(title) == 0 && len(style) == 0 {
		return
	}
	bb.putblock(RawContent("<head>"))
	bb.want_nextln()
	if len(title) > 0 {
		bb.putblock_ex(1, "<title>", title, "</title>")
		bb.want_nextln()
	}
	if len(style) > 0 {
		bb.putblock_ex(1, "<style>", style, "</style>")
		bb.want_nextln()
	}
	bb.putblock(RawContent("</head>"))
	bb.want_nextln()
}

func (bb *html_blocks) begin_body() {
	bb.putblock(RawContent("<body>"))
	bb.want_nextln()
}

func (bb *html_blocks) end_body() {
	bb.putblock(RawContent("</body>"))
	bb.want_nextln()
}

func (ii *html_inlines) simple_link(b *bytes.Buffer, caption RawContent, url RawContent) {
	fmt.Fprintf(b, "<a href=\"%s\">", url)
	if len(caption) == 0 {
		b.Write(url)
	} else {
		b.Write(caption)
	}
	b.WriteString("</a>")
}
