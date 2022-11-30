package markout

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type html_blocks struct {
	base_blocks
	list_title_class string
}

func (bb *html_blocks) para(s RawContent) {
	bb.putblock_ex(0, "<p>", s, "</p>")
	bb.want_emptyln()
}

func (bb *html_blocks) heading(counters []int, s RawContent, aa *Attrs) {
	level := len(counters)
	tagname := "h" + strconv.Itoa(level)
	t := fmt.Sprintf("<%s", tagname)
	if aa != nil {
		if aa.Identifier != "" {
			t += fmt.Sprintf(" id=\"%s\"", aa.Identifier)
		}
		if len(aa.Classes) > 0 {
			t += fmt.Sprintf(" class=\"%s\"", strings.Join(aa.Classes, " "))
		}
		if len(aa.KeyVals) > 0 {
			kk := make([]string, 0, len(aa.KeyVals))
			for k := range aa.KeyVals {
				kk = append(kk, k)
			}
			sort.Strings(kk)
			for _, k := range kk {
				t += fmt.Sprintf(" %s=\"%s\"", k, aa.KeyVals[k])
			}
		}
	}
	t += ">"
	bb.putblock_ex(0, t, s, "</"+tagname+">")
	bb.want_emptyln()
}

func (bb *html_blocks) list_title(s RawContent) {
	if bb.list_title_class != "" {
		bb.putblock_ex(0, "<p classs=\""+bb.list_title_class+"\">", s, "</p>")
	} else {
		bb.putblock_ex(0, "<p>", s, "</p>")
	}
	bb.want_nextln()
}

func (bb *html_blocks) list_level_start(counters []int) {
	if bb.enabled() {
		n := len(counters) - 1
		bb.putblock_ex(n, pick(counters[n] >= 0, "<ul>", "<ol>"), []byte{}, "")
	}
	bb.want_nextln()
}

func (bb *html_blocks) list_level_done(counters []int) {
	if bb.enabled() {
		n := len(counters) - 1
		bb.putblock_ex(n, pick(counters[n] >= 0, "</ul>", "</ol>"), []byte{}, "")
	}
	if len(counters) == 1 {
		bb.want_emptyln()
	} else {
		bb.want_nextln()
	}
}

func (bb *html_blocks) list_item(counters []int, s RawContent) {
	if bb.enabled() {
		n := len(counters) - 1
		bb.putblock_ex(n+1, "<li>", s, "</li>")
	}
	bb.want_nextln()
}

func (bb *html_blocks) end_table() {
	if len(bb.table) > 1 {
		bb.do_nextline()
		bb.out.Write([]byte("<table>\n<thead><tr>"))
		for _, c := range bb.table[0] {
			bb.out.Write([]byte("<th>"))
			bb.out.Write(c)
			bb.out.Write([]byte("</th>"))
		}
		bb.out.Write([]byte("</tr></thead>\n<tbody>"))
		for _, row := range bb.table[1:] {
			bb.out.Write([]byte("\n<tr>"))
			for _, c := range row {
				bb.out.Write([]byte("<td>"))
				bb.out.Write(c)
				bb.out.Write([]byte("</td>"))
			}
			bb.out.Write([]byte("</tr>"))
		}
		bb.out.Write([]byte("\n</tbody>\n</table>"))
		bb.table = bb.table[:0]
	}
	bb.want_emptyln()
}

func (bb *html_blocks) codeblock(lang string, s RawContent) {
	bb.want_emptyln()
	bb.do_nextline()
	bb.out.Write([]byte("<pre"))
	if lang != "" {
		bb.out.Write([]byte(" lang=\""))
		bb.out.Write([]byte(lang))
		bb.out.Write([]byte{'"'})
	}
	bb.out.Write([]byte(">\n"))
	bb.out.Write(s)
	bb.out.Write([]byte("\n</pre>"))
	bb.want_emptyln()
}
