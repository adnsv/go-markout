package markout

import (
	"bytes"
	"fmt"
	"os"
)

func ExampleNewMultiWriter() {
	buf := bytes.Buffer{}
	html_w := NewHTML(&buf, HTMLOptions{})
	cout_w := NewTXT(os.Stdout, TXTOptions{})
	null_w := NewNULL()

	multi_w := NewMultiWriter(cout_w, html_w, null_w)

	w := Writer(multi_w)
	w.ListTitle("List:")
	w.List(Unordered, func(li ListWriter) {
		li.ListItem("item 1")
		li.ListItem("item 2")
	})
	w.Table([]any{"h1", "h2"}, func(row TableRowWriter) {
		row("c1", "c2")
		row("c3", "cell4")
	})

	cout_w.Close()
	html_w.Close()

	fmt.Printf("************\n")
	os.Stdout.Write(buf.Bytes())
	// Output:
	// List:
	// * item 1
	// * item 2
	//
	// h1 h2
	// -- -----
	// c1 c2
	// c3 cell4
	//
	// ************
	// <html>
	// <body>
	// <p>List:</p>
	// <ul>
	//   <li>item 1</li>
	//   <li>item 2</li>
	// </ul>
	//
	// <table>
	// <thead><tr><th>h1</th><th>h2</th></tr></thead>
	// <tbody>
	// <tr><td>c1</td><td>c2</td></tr>
	// <tr><td>c3</td><td>cell4</td></tr>
	// </tbody>
	// </table>
	//
	// </body>
	// </html>
}
