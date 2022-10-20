package wcwidth

import "testing"

func TestStringCells(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want int
	}{
		{"empty", "", 0},
		{"zero-width", "\u200B", 0},
		{"single latin", "a", 1},
		{"single wide", "常", 2},
		{"multiple latin", "abcd", 4},
		{"multiple wide", "常用漢字", 8},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringCells(tt.arg); got != tt.want {
				t.Errorf("StringWidth() = %v, want %v", got, tt.want)
			}
		})
	}
}
