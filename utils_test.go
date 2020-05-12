package sieve

import (
	"testing"
)

func TestUtils(t *testing.T) {
	tables := []struct {
		Data []byte
		Pos  int
	}{
		{[]byte("  "), 2},
		{[]byte(" "), 1},
		{[]byte(""), 0},
	}

	for _, tb := range tables {
		if pos := skipWhiteSpace(tb.Data, 0); pos != tb.Pos {
			t.Errorf("skipWhiteSpace for %#v got: {%v}, want: {%v}", string(tb.Data), pos, tb.Pos)
		}
	}
}
