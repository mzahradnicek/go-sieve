package sieve

import (
	"reflect"
	"testing"
)

func TestStringBase(t *testing.T) {
	tables := []struct {
		Data   []byte
		Parsed stringBase
		Pos    int
		Err    bool
	}{

		{[]byte("\"123\""), "123", 5, false},
		{[]byte("\"123K some text\""), "123K some text", 16, false},
		{[]byte("\"123K so\\\" text\""), "123K so\" text", 16, false},
		{[]byte("\"123K so\\\\ text\""), "123K so\\ text", 16, false},
		{[]byte("\"123K so\\x text\""), "123K so\\x text", 16, false},

		// wrong
		{[]byte("1M23G\""), "123G", 0, true}, // dont start with quote
		{[]byte("\"C123G"), "123G", 0, true}, // dont end with quote
	}

	for _, tb := range tables {
		var str stringBase

		if pos, err := str.Scan(tb.Data); pos != tb.Pos || (str != tb.Parsed && err == nil) || (err == nil) == tb.Err {
			t.Errorf("stringBase \"%v\" for %v got:  {%v, %v}, want: {%v, %v}", string(str), string(tb.Data), pos, err, tb.Pos, tb.Err)
		}
	}
}

func TestStringType(t *testing.T) {
	tables := []struct {
		Data   []byte
		Parsed StringType
		Pos    int
		Err    bool
	}{
		{[]byte("\"123\""), []stringBase{"123"}, 5, false},
		{[]byte("[\"123\"]"), []stringBase{"123"}, 7, false},
		{[]byte("[ \"123\"]"), []stringBase{"123"}, 8, false},
		{[]byte("[ \"123\",\"klobasa\"]"), []stringBase{"123", "klobasa"}, 18, false},
		{[]byte("[ \"123\"  ,\"klobasa\"]"), []stringBase{"123", "klobasa"}, 20, false},
		{[]byte("[ \"123\"  ,  \"klobasa\"]"), []stringBase{"123", "klobasa"}, 22, false},
		{[]byte("[ \"123\"  ,  \"klobasa\" ]"), []stringBase{"123", "klobasa"}, 23, false},

		// wrong
		{[]byte("x[ \"123\""), []stringBase{}, 0, true},
		{[]byte("[ \"123\""), []stringBase{}, 0, true},
		{[]byte("[ \"123\","), []stringBase{}, 0, true},
		{[]byte("[ \"123\", "), []stringBase{}, 0, true},
	}

	for _, tb := range tables {
		var str StringType

		if pos, err := str.Scan(tb.Data); pos != tb.Pos || (!reflect.DeepEqual(str, tb.Parsed) && err == nil) || (err == nil) == tb.Err {
			t.Errorf("stringBase %#v for %v got:  {%v, %v}, want: {%v, %v}", str, string(tb.Data), pos, err, tb.Pos, tb.Err)
		}
	}
}
