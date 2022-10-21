package markout

import (
	"bytes"
	"encoding"
	"fmt"
	"reflect"
	"strconv"

	"golang.org/x/exp/slices"
)

// Printer is an interface used in callbacks and custom marshalers for writing
// inline content.
type Printer interface {
	// Low-level string content
	WriteString(string)
	WriteRawBytes([]byte)

	// Code spans
	CodeString(string)
	CodeRawBytes([]byte)

	// Inline links
	BeginLink(url string)
	EndLink()

	// Quoted spans
	BeginStyled(Style)
	EndStyled()

	// High-level api
	Print(any)
	Printf(format string, args ...any)
	SimpleLink(a any, url string)
	Styled(Style, any)
}

type url_filter = func(url string) []byte

// printer_impl implements Printer interface for writing into
// bytes.Buffer
type printer_impl struct {
	buf        *bytes.Buffer
	ii         inlines
	url_filter url_filter
}

func (p *printer_impl) WriteString(s string) {
	p.ii.check_mode(iflow)
	p.ii.put_str(p.buf, s)
}
func (p *printer_impl) WriteRawBytes(s []byte) {
	p.ii.check_mode(iflow)
	p.ii.put_raw(p.buf, s)
}
func (p *printer_impl) CodeString(s string) {
	p.ii.check_mode(iflow)
	p.ii.code_str(p.buf, s)
}
func (p *printer_impl) CodeRawBytes(s []byte) {
	p.ii.check_mode(iflow)
	p.ii.code_raw(p.buf, s)
}
func (p *printer_impl) BeginLink(url string) {
	p.ii.check_not_mode(ilink)
	if p.url_filter != nil {
		p.ii.begin_link(p.buf, RawContent(p.url_filter(url)))
	} else {
		p.ii.begin_link(p.buf, RawContent(url))
	}
}
func (p *printer_impl) EndLink() {
	p.ii.check_mode(ilink)
	p.ii.end_link(p.buf)
}
func (p *printer_impl) BeginStyled(sty Style) {
	p.ii.check_mode(iflow)
	p.ii.begin_styled(p.buf, sty)
}
func (p *printer_impl) EndStyled() {
	p.ii.check_mode(iflow)
	p.ii.end_styled(p.buf)
}
func (p *printer_impl) Print(a any) {
	p.ii.check_mode(iflow)
	to_buffer(p.buf, p.ii, p.url_filter, a)
}
func (p *printer_impl) Printf(format string, args ...any) {
	p.ii.check_mode(iflow)
	scratch := bytes.Buffer{}
	p.ii.put_str(&scratch, format)
	fmt_raw := scratch.String()
	args_raw := fmt_args(&scratch, p.ii, p.url_filter, args...)
	fmt.Fprintf(p.buf, fmt_raw, args_raw...)
}
func (p *printer_impl) SimpleLink(a any, url string) {
	p.ii.check_not_mode(ilink)
	scratch := bytes.Buffer{}
	to_buffer(&scratch, p.ii, p.url_filter, a)
	if p.url_filter != nil {
		p.ii.simple_link(p.buf, scratch.Bytes(), RawContent(p.url_filter(url)))
	} else {
		p.ii.simple_link(p.buf, scratch.Bytes(), RawContent(url))
	}
}
func (p *printer_impl) Styled(sty Style, a any) {
	p.BeginStyled(sty)
	p.Print(a)
	p.EndStyled()
}

var (
	inlineMarshalerType = reflect.TypeOf((*InlineMarshaler)(nil)).Elem()
	textMarshalerType   = reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem()
)

// handle_print_marshaler handles writing of markout.Marshaler; in addition it also
// supports writing of types that implement encoding.TextMarshaler.
func handle_print_marshaler(p Printer, val reflect.Value, typ reflect.Type) (bool, error) {
	can_interface := val.CanInterface()
	if can_interface && typ.Implements(inlineMarshalerType) {
		return true, val.Interface().(InlineMarshaler).MarshalMarkoutInline(p)
	}
	can_addr := val.CanAddr()
	if can_addr {
		pv := val.Addr()
		if pv.CanInterface() && pv.Type().Implements(inlineMarshalerType) {
			return true, pv.Interface().(InlineMarshaler).MarshalMarkoutInline(p)
		}
	}
	if can_interface && typ.Implements(textMarshalerType) {
		t, err := val.Interface().(encoding.TextMarshaler).MarshalText()
		if err == nil {
			p.WriteString(string(t))
		}
		return true, err
	}
	if can_addr {
		pv := val.Addr()
		if pv.CanInterface() && pv.Type().Implements(textMarshalerType) {
			t, err := pv.Interface().(encoding.TextMarshaler).MarshalText()
			if err == nil {
				p.WriteString(string(t))
			}
			return true, err
		}
	}
	return false, nil
}

// handle_print_stringlike handles writing of strings, character arrays, and character slices.
func handle_print_stringlike(p Printer, val reflect.Value, typ reflect.Type) (bool, error) {
	k := val.Kind()
	switch k {
	case reflect.String:
		p.WriteString(val.String())
		return true, nil
	case reflect.Array:
		if typ.Elem().Kind() != reflect.Uint8 {
			break
		}
		// [...]byte
		var bytes []byte
		if val.CanAddr() {
			bytes = val.Slice(0, val.Len()).Bytes()
		} else {
			bytes = make([]byte, val.Len())
			reflect.Copy(reflect.ValueOf(bytes), val)
		}
		p.WriteString(string(bytes))
		return true, nil
	case reflect.Slice:
		k := typ.Elem().Kind()
		if k != reflect.Uint8 {
			break
		}
		// []byte
		p.WriteString(string(val.Bytes()))
		return true, nil
	}
	return false, nil
}

func handle_print_simple(p Printer, val reflect.Value, typ reflect.Type) (bool, error) {
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p.WriteRawBytes(RawContent(strconv.FormatInt(val.Int(), 10)))
		return true, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p.WriteRawBytes(RawContent(strconv.FormatUint(val.Uint(), 10)))
		return true, nil
	case reflect.Float32, reflect.Float64:
		p.WriteRawBytes(RawContent(strconv.FormatFloat(val.Float(), 'g', -1, val.Type().Bits())))
		return true, nil
	case reflect.Bool:
		p.WriteRawBytes(RawContent(strconv.FormatBool(val.Bool())))
		return true, nil
	}
	return false, nil
}

func print_any(p Printer, a any) error {

	// handle RawStr first
	switch v := a.(type) {
	case RawContent:
		p.WriteRawBytes(v)
		return nil
	case link_wrapper:
		p.SimpleLink(v.caption, v.url)
		return nil
	case style_wrapper:
		p.Styled(v.sty, v.content)
		return nil
	case Callback:
		v(p)
		return nil
	}

	val := reflect.ValueOf(a)

	if !val.IsValid() {
		return nil
	}

	for val.Kind() == reflect.Interface || val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil
		}
		val = val.Elem()
	}

	typ := val.Type()

	// handle InlineMarshaler and encoding.TextMarshaler
	handled, err := handle_print_marshaler(p, val, typ)
	if !handled {
		switch v := a.(type) {
		case error:
			p.WriteString(v.Error())
			handled = true
		case fmt.Stringer:
			p.WriteString(v.String())
			handled = true
		case string:
			p.WriteString(v)
			handled = true
		}
	}
	if !handled {
		handled, err = handle_print_stringlike(p, val, typ)
	}
	if !handled {
		handled, err = handle_print_simple(p, val, typ)
	}
	if handled {
		return err
	} else {
		return &UnsupportedTypeError{typ}
	}
}

// UnsupportedTypeError is returned when Marshal encounters a type
// that cannot be converted into markout.
type UnsupportedTypeError struct {
	Type reflect.Type
}

func (e *UnsupportedTypeError) Error() string {
	return "markout: unsupported type: " + e.Type.String()
}

const marshalErr = "#ERR"

// to_buffer converts an argument to RawStr and writes it into b.
func to_buffer(b *bytes.Buffer, ii inlines, uf url_filter, a any) {
	p := printer_impl{b, ii, uf}
	if err := print_any(&p, a); err != nil {
		p.WriteString(marshalErr)
	}
}

func fmt_args(scratch *bytes.Buffer, ii inlines, uf url_filter, args ...any) []any {
	r := slices.Clone(args)
	p := printer_impl{scratch, ii, uf}

	for i := range r {
		scratch.Reset()

		switch v := r[i].(type) {
		case RawContent:
			p.WriteRawBytes(v)
			r[i] = scratch.String()
			continue
		case link_wrapper:
			p.SimpleLink(v.caption, v.url)
			r[i] = scratch.String()
			continue
		case style_wrapper:
			p.Styled(v.sty, v.content)
			r[i] = scratch.String()
			continue
		case Callback:
			v(&p)
			r[i] = scratch.String()
			continue
		}

		val := reflect.ValueOf(r[i])
		if !val.IsValid() {
			continue
		}
		for val.Kind() == reflect.Interface || val.Kind() == reflect.Ptr {
			if val.IsNil() {
				break
			}
			val = val.Elem()
		}
		typ := val.Type()

		handled, err := handle_print_marshaler(&p, val, typ)
		if !handled {
			switch v := r[i].(type) {
			case error:
				p.WriteString(v.Error())
				handled = true
			case fmt.Stringer:
				p.WriteString(v.String())
				handled = true
			case string:
				p.WriteString(v)
				handled = true
			}
		}
		if !handled {
			handled, err = handle_print_stringlike(&p, val, typ)
		}
		if handled {
			if err != nil {
				p.WriteString("#ERR")
			}
			r[i] = RawContent(scratch.String())
		}
	}
	return r
}
