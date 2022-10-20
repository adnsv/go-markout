package markout

// Link creates a wrapper for a named links.
func Link(caption any, url string) link_wrapper {
	return link_wrapper{caption: caption, url: url}
}

// URL creates a wrapped for unnamed link.
func URL(url string) link_wrapper {
	return link_wrapper{url: url}
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

// Callback is a funcional inline content builder that can be used for more
// complex inline formatting.
type Callback = func(Printer)

// InlineMarshaler is the interface implemented by types that support marshal custom
// marshaling into markout targets (as inline fragments).
type InlineMarshaler interface {
	MarshalMarkoutInline(Printer) error
}

// Style provides decorations for inline output.
type Style int

const (
	SingleQuotedStyle = Style(iota)
	DoubleQuotedStyle
	EmphasizedStyle
	StrongStyle
)

// raw_bytes is the content that is ready for output, no additional
// scrambling required.
type raw_bytes []byte

type link_wrapper struct {
	caption any
	url     string
}

type style_wrapper struct {
	sty     Style
	content any
}
