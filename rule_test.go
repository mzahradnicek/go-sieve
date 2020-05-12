package sieve

import (
	"reflect"
	"testing"
)

func TestRuleScan(t *testing.T) {
	tables := []struct {
		Data   []byte
		Parsed Rule
		Pos    int
		Err    bool
	}{
		// header
		{[]byte("header :contains \"From\" \"salama\""), Rule{Type: "from", Operator: "contains", TargetString: StringType{"From"}, QueryString: StringType{"salama"}}, 32, false},
		{[]byte("header :contains \"To\" \"salama\""), Rule{Type: "to", Operator: "contains", TargetString: StringType{"To"}, QueryString: StringType{"salama"}}, 30, false},
		{[]byte("header :contains \"Subject\" \"salama\""), Rule{Type: "subject", Operator: "contains", TargetString: StringType{"Subject"}, QueryString: StringType{"salama"}}, 35, false},
		{[]byte("header :matches [\"Subject\",\"From\"] \"salama\""), Rule{Type: "...", Operator: "matches", TargetString: StringType{"Subject", "From"}, QueryString: StringType{"salama"}}, 43, false},

		// body
		{[]byte("body :text :contains \"salama\""), Rule{Type: "body", Operator: "contains", QueryString: StringType{"salama"}}, 29, false},
		{[]byte("body :text :is \"klobasa\""), Rule{Type: "body", Operator: "is", QueryString: StringType{"klobasa"}}, 24, false},

		// exists
		{[]byte("exists \"From\""), Rule{Type: "from", Operator: "exists", TargetString: StringType{"From"}}, 13, false},
		{[]byte("exists [\"From\",\"To\"]"), Rule{Type: "...", Operator: "exists", TargetString: StringType{"From", "To"}}, 20, false},

		// not
		{[]byte("not header :matches [\"Subject\",\"From\"] \"salama\""), Rule{Type: "...", Operator: "notmatches", TargetString: StringType{"Subject", "From"}, QueryString: StringType{"salama"}}, 47, false},
		{[]byte("not body :text :is \"klobasa\""), Rule{Type: "body", Operator: "notis", QueryString: StringType{"klobasa"}}, 28, false},
	}

	for _, tb := range tables {
		var str Rule

		if pos, err := str.Scan(tb.Data); pos != tb.Pos || (!reflect.DeepEqual(str, tb.Parsed) && err == nil) || (err == nil) == tb.Err {
			t.Errorf("Rule %#v GOT: %#v, pos: %v, err: %v, WANT: %#v, pos: %v, err: %v", string(tb.Data), str, pos, err, tb.Parsed, tb.Pos, tb.Err)
		}
	}
}

func TestRuleString(t *testing.T) {
	tables := []struct {
		Expect string
		Src    Rule
	}{
		// header
		{"header :contains \"From\" \"salama\"", Rule{Type: "from", Operator: "contains", QueryString: StringType{"salama"}}},
		{"header :contains \"To\" \"salama\"", Rule{Type: "to", Operator: "contains", QueryString: StringType{"salama"}}},
		{"header :contains \"Subject\" \"salama\"", Rule{Type: "subject", Operator: "contains", QueryString: StringType{"salama"}}},
		{"header :matches [\"Subject\",\"From\"] [\"salama\",\"klobasa\"]", Rule{Type: "...", Operator: "matches", TargetString: StringType{"Subject", "From"}, QueryString: StringType{"salama", "klobasa"}}},

		// body
		{"body :text :contains \"salama\"", Rule{Type: "body", Operator: "contains", QueryString: StringType{"salama"}}},
		{"body :text :is \"klobasa\"", Rule{Type: "body", Operator: "is", QueryString: StringType{"klobasa"}}},

		// exists
		{"exists \"From\"", Rule{Type: "from", Operator: "exists"}},
		{"exists [\"From\",\"To\"]", Rule{Type: "...", Operator: "exists", TargetString: StringType{"From", "To"}}},

		// not
		{"not header :matches [\"Subject\",\"From\"] \"salama\"", Rule{Type: "...", Operator: "notmatches", TargetString: StringType{"Subject", "From"}, QueryString: StringType{"salama"}}},
		{"not body :text :is \"klobasa\"", Rule{Type: "body", Operator: "notis", QueryString: StringType{"klobasa"}}},
	}

	for _, tb := range tables {
		if res := tb.Src.String(); res != tb.Expect {
			t.Errorf("Rule for %#v GOT: %v, WANT: %#v", tb.Src, res, tb.Expect)
		}
	}
}
