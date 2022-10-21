package markout

import (
	"bytes"
	"testing"

	"golang.org/x/exp/slices"
)

func Test_md_scramble(t *testing.T) {
	tests := []struct {
		arg  string
		want RawContent
	}{
		{"", RawContent("")},
		{"abc", RawContent(`abc`)},
		{"\\n", RawContent(`\\n`)},
		{"{}", RawContent(`\{\}`)},
		{"a+b", RawContent(`a\+b`)},
		{"<a>", RawContent(`\<a\>`)},
		{"## head", RawContent(`\#\# head`)},
	}
	for _, tt := range tests {
		t.Run(tt.arg, func(t *testing.T) {
			b := bytes.Buffer{}
			md_scramble(&b, tt.arg)
			if got := b.Bytes(); !slices.Equal(got, tt.want) {
				t.Errorf("scramble() = %v, want %v", got, tt.want)
			}
		})
	}
}
