package markout

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
)

func ExampleNewTXT() {
	buf := bytes.Buffer{}
	w := NewTXT(&buf, TXTOptions{
		PutBOM:             false,
		NumberedSections:   true,
		UnderlinedSections: true})
	w.Para("Para")
	w.ListTitle("list:")
	w.BeginList(Unordered)
	w.ListItem("list item")
	w.ListItem(3.14)
	w.ListItem(true)
	w.ListItem(42)
	w.ListItem("subitems:")
	w.BeginList(Ordered)
	w.ListItem("subitem1")
	w.ListItem(Emphasized("subitem2"))
	w.ListItem(SingleQuoted("subitem3"))
	w.ListItem(Emphasized(SingleQuoted("subitem4")))
	w.ListItem(URL("subitem5"))
	w.ListItem(Link("a", "b"))
	w.EndList()
	w.ListItem("last")
	w.EndList()
	w.BeginTable("th1", "th2")
	w.TableRow("tcell", "another cell")
	w.EndTable()
	w.BeginSection("SECTION")
	w.BeginSection("SUBSECTION")
	w.Section("First SubSubSection")
	w.Section("Second SubSubSection")
	w.Section("Third SubSubSection")
	w.EndSection()
	w.Section("ANOTHER SUBSECTION")
	w.EndSection()
	w.Codeblock("go", "codeblock\ncontent")
	w.Close()
	out := buf.String()
	fmt.Println(out)
	// Output:
	// Para
	//
	// list:
	// * list item
	// * 3.14
	// * true
	// * 42
	// * subitems:
	//   1. subitem1
	//   2. *subitem2*
	//   3. 'subitem3'
	//   4. *'subitem4'*
	//   5. subitem5
	//   6. [a](b)
	// * last
	//
	// th1   th2
	// ----- ------------
	// tcell another cell
	//
	// 1. SECTION
	// ==========
	//
	// 1.1. SUBSECTION
	// ---------------
	//
	// 1.1.1. First SubSubSection
	//
	// 1.1.2. Second SubSubSection
	//
	// 1.1.3. Third SubSubSection
	//
	// 1.2. ANOTHER SUBSECTION
	// -----------------------
	//
	// codeblock
	// content

}

func ExampleNewMD() {
	buf := bytes.Buffer{}
	w := NewMD(&buf, MDOptions{PutBOM: false})
	w.BeginSection("Section")
	w.Para("Para")
	w.BeginList(Unordered)
	w.ListItem("list item")
	w.ListItem(3.14)
	w.ListItem(true)
	w.ListItem(42)
	w.ListItem("subitems:")
	w.BeginList(Ordered)
	w.ListItem("subitem1")
	w.ListItem(Emphasized("subitem2"))
	w.EndList()
	w.ListItem("last")
	w.EndList()

	w.Para("Broad:")
	w.BeginList(Broad)
	w.ListItem("broad1")
	w.ListItem("broad2")
	w.ListItem("broad3")
	w.BeginList(Tight)
	w.ListItem("tight1")
	w.ListItem("tight2")
	w.ListItem("tight3")
	w.EndList()
	w.ListItem("broad4")
	w.EndList()

	w.AttrSection(Attrs{Identifier: "ident", Classes: []string{"cls"}}, "Subsection")

	w.Para(func(p Printer) {
		p.Print("Inline formatting: ")
		p.BeginStyled(DoubleQuotedStyle)
		p.Styled(EmphasizedStyle, "Hello")
		p.Print(", ")
		p.Styled(StrongStyle, "World!")
		p.EndStyled()
	})

	w.BeginTable("th", "thead")
	w.TableRow("tcell", "tcell")
	w.EndTable()
	w.Codeblock("go", "codeblock\ncontent")
	w.EndSection()
	/*w.List(Ordered|Broad, func(ListWriter) {
		w.ListItem(func(w ParagraphWriter) {
			w.Para("Advanced list items")
			w.Para("Can be composed")
			w.Para("From multiple paragraph blocks")
		})
		w.ListItem("Second item")
	})
	*/
	w.Close()
	out := buf.String()
	fmt.Println(out)
	// Output:
	// # Section
	//
	// Para
	//
	// - list item
	// - 3.14
	// - true
	// - 42
	// - subitems:
	//   1. subitem1
	//   2. <em>subitem2</em>
	// - last
	//
	// Broad:
	//
	// - broad1
	//
	// - broad2
	//
	// - broad3
	//
	//   - tight1
	//   - tight2
	//   - tight3
	//
	// - broad4
	//
	// ## Subsection {#ident .cls}
	//
	// Inline formatting: "<em>Hello</em>, <strong>World\!</strong>"
	//
	// |th | thead
	// |---|------
	// |tcell | tcell
	//
	// ```go
	// codeblock
	// content
	// ```
}

func ExampleNewHTML() {
	buf := bytes.Buffer{}
	w := NewHTML(&buf, HTMLOptions{PutBOM: false})
	w.Para("Para")
	w.BeginList(Unordered)
	w.ListItem("list item")
	w.ListItem(3.14)
	w.ListItem(true)
	w.ListItem(42)
	w.ListItem("subitems:")
	w.BeginList(Ordered)
	w.ListItem("subitem1")
	w.ListItem(Emphasized("subitem2"))
	w.ListItem(SingleQuoted("subitem3"))
	w.ListItem(Emphasized(SingleQuoted("subitem4")))
	w.EndList()
	w.ListItem("last")
	w.EndList()
	w.BeginTable("thead", "thead")
	w.TableRow("tcell", "tcell")
	w.EndTable()
	w.BeginSection("Section")
	w.BeginAttrSection(Attrs{Identifier: "ident", Classes: []string{"cls"}}, "SubSection")
	w.Section("SubSubSection")
	w.EndSection()
	w.EndSection()
	w.Codeblock("go", "codeblock\ncontent")
	w.List(Ordered|Broad, func(ListWriter) {
		w.ListItem(func(w ParagraphWriter) {
			w.Para("Advanced list items")
			w.Para("Can be composed")
			w.Para("From multiple paragraph blocks")
		})
		w.ListItem("Second item")
	})
	w.Close()
	out := buf.String()
	fmt.Println(out)
	// Output:
	// <html>
	// <body>
	// <p>Para</p>
	//
	// <ul>
	//   <li>list item</li>
	//   <li>3.14</li>
	//   <li>true</li>
	//   <li>42</li>
	//   <li>subitems:</li>
	//   <ol>
	//     <li>subitem1</li>
	//     <li><em>subitem2</em></li>
	//     <li>'subitem3'</li>
	//     <li><em>'subitem4'</em></li>
	//   </ol>
	//   <li>last</li>
	// </ul>
	//
	// <table>
	// <thead><tr><th>thead</th><th>thead</th></tr></thead>
	// <tbody>
	// <tr><td>tcell</td><td>tcell</td></tr>
	// </tbody>
	// </table>
	//
	// <h1>Section</h1>
	//
	// <h2 id="ident" class="cls">SubSection</h2>
	//
	// <h3>SubSubSection</h3>
	//
	// <pre lang="go">
	// codeblock
	// content
	// </pre>
	//
	// <ol>
	//   <li><p>Advanced list items</p>
	//     <p>Can be composed</p>
	//     <p>From multiple paragraph blocks</p></li>
	//   <li><p>Second item</p></li>
	// </ol>
	//
	// </body>
	// </html>
}

func ExampleURL() {
	url_filter := func(url string) []byte {
		return []byte(filepath.Base(url))
	}

	w := NewTXT(os.Stdout, TXTOptions{URLFilter: url_filter})
	w.Paraf("Link: %s", Link("test", "../../path.txt"))
	w.Paraf("URL: %s", URL("../../path.txt"))
	w.Close()
	// Output:
	// Link: [test](path.txt)
	//
	// URL: path.txt
}
