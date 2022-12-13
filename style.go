package markout

// Link creates a wrapper for a named links.
func Link(caption any, url string) link_wrapper {
	return link_wrapper{caption: caption, url: url}
}

// URL creates a wrapper for an unnamed link.
func URL(url string) link_wrapper {
	return link_wrapper{url: url}
}

func Code(s string) codespan {
	return codespan(s)
}

// SingleQuoted creates a wrapper for single quoted inline spans.
func SingleQuoted(a any) style_wrapper {
	return style_wrapper{sty: SingleQuotedStyle, content: a}
}

// DoubleQuoted creates a wrapper for double quoted inline spans.
func DoubleQuoted(a any) style_wrapper {
	return style_wrapper{sty: DoubleQuotedStyle, content: a}
}

// DoubleQuoted creates a wrapper for emphasized (italic) inline spans.
func Emphasized(a any) style_wrapper {
	return style_wrapper{sty: EmphasizedStyle, content: a}
}

// DoubleQuoted creates a wrapper for strong-formatted (bold) inline spans.
func Strong(a any) style_wrapper {
	return style_wrapper{sty: StrongStyle, content: a}
}

// Callback is a funcional inline content builder that can be used for complex
// inline formatting.
type Callback = func(Printer)

// InlineMarshaler is the interface implemented by types that support marshal
// custom marshaling into markout targets (as inline fragments).
type InlineMarshaler interface {
	MarshalMarkoutInline(Printer) error
}

// Style provides decorations for spans in inline output. Styles can be nested.
type Style int

// Accepted Style values:
const (
	SingleQuotedStyle = Style(iota) // a fragment surrounded by single quotation marks
	DoubleQuotedStyle               // a fragment surrounded by double quotation marks
	EmphasizedStyle                 // emphasized fragment, typically rendered in italic type
	StrongStyle                     // strong fragment, typically rendered in bold type
)

// RawContent is the the sequence of bytes that is written out to a target
// 'as-is'. No additional scrambling or escaping is performed.
type RawContent []byte

// Codespan formats inline span
type codespan string

type link_wrapper struct {
	caption any
	url     string
}

type style_wrapper struct {
	sty     Style
	content any
}
