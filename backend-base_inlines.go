package markout

import (
	"errors"
	"strings"
)

type base_inlines struct {
	quote_specs [4]string // quotation marks [<single_open>, <single_close>, <double_open>, <double_close>]
	style_stack []Style
}

func (ii *base_inlines) setup_quotation_marks(quote_fmt string) {
	// quote_fmt is a pipe-separated singles followed by doubles
	qq := strings.Split(quote_fmt, "|")
	if len(qq) != 4 {
		// fallback to defaults
		qq = []string{`'`, `'`, `"`, `"`}
	}
	copy(ii.quote_specs[:], qq)
}

func (ii *base_inlines) start_styled(sty Style) {
	ii.style_stack = append(ii.style_stack, sty)
}

func (ii *base_inlines) finish_styled() Style {
	n := len(ii.style_stack)
	if n == 0 {
		panic(err_unpaired_endstyle)
	}
	r := ii.style_stack[n-1]
	ii.style_stack = ii.style_stack[:n-1]
	return r
}

const err_unpaired_endstyle = "markout: unpaired EndStyled call"
const err_inline_mode = "markout: command is not allowed within the current inline mode"

var errUnpairedBeginStyled = errors.New("markout: unpaired BeginStyled call")
var errUnpairedBeginLink = errors.New("markout: unpaired BeginLink call")
