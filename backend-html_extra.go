package markout

import (
	"bytes"
	"fmt"
)

// things related exclusively to HTML file structure

func (bb *html_blocks) begin_html() {
	bb.putblock(raw_bytes("<html>"))
	bb.want_nextln()
}

func (bb *html_blocks) end_html() {
	bb.want_nextln()
	bb.putblock(raw_bytes("</html>"))
	bb.want_nextln()
}

func (bb *html_blocks) head(title, style raw_bytes) {
	if len(title) == 0 && len(style) == 0 {
		return
	}
	bb.putblock(raw_bytes("<head>"))
	bb.want_nextln()
	if len(title) > 0 {
		bb.putblock_ex(1, "<title>", title, "</title>")
		bb.want_nextln()
	}
	if len(style) > 0 {
		bb.putblock_ex(1, "<style>", style, "</style>")
		bb.want_nextln()
	}
	bb.putblock(raw_bytes("</head>"))
	bb.want_nextln()
}

func (bb *html_blocks) begin_body() {
	bb.putblock(raw_bytes("<body>"))
	bb.want_nextln()
}

func (bb *html_blocks) end_body() {
	bb.putblock(raw_bytes("</body>"))
	bb.want_nextln()
}

func (ii *html_inlines) simple_link(b *bytes.Buffer, caption raw_bytes, url raw_bytes) {
	fmt.Fprintf(b, "<a href=\"%s\">", url)
	if len(caption) == 0 {
		b.Write(url)
	} else {
		b.Write(caption)
	}
	b.WriteString("</a>")
}
