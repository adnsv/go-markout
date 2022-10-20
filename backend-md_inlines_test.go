package markout

import (
	"bytes"
	"testing"

	"golang.org/x/exp/slices"
)

func Test_md_scramble(t *testing.T) {
	tests := []struct {
		arg  string
		want raw_bytes
	}{
		{"", raw_bytes("")},
		{"abc", raw_bytes(`abc`)},
		{"\\n", raw_bytes(`\\n`)},
		{"{}", raw_bytes(`\{\}`)},
		{"a+b", raw_bytes(`a\+b`)},
		{"<a>", raw_bytes(`\<a\>`)},
		{"## head", raw_bytes(`\#\# head`)},
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
