package sieve

import "testing"

func TestNumberType(t *testing.T) {
	tables := []struct {
		Data   []byte
		Parsed NumberType
		Pos    int
		Err    bool
	}{

		{[]byte("123"), "123", 3, false},
		{[]byte("123K some text"), "123K", 4, false},
		{[]byte("123M"), "123M", 4, false},
		{[]byte("123G"), "123G", 4, false},
		{[]byte("123G"), "123G", 4, false},

		// wrong
		{[]byte("C123G"), "123G", 0, true},
		{[]byte("1M23G"), "123G", 0, true},
		{[]byte(" 1M23G"), "123G", 0, true},
	}

	for _, tb := range tables {
		var num NumberType

		if pos, err := num.Scan(tb.Data); pos != tb.Pos || (num != tb.Parsed && err == nil) || (err == nil) == tb.Err {
			t.Errorf("NumberType \"%+v\" for \"%v\" got:  {%v, %v}, want: {%v, %v}", num, string(tb.Data), pos, err, tb.Pos, tb.Err)
		}
	}
}
