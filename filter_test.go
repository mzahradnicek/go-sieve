package sieve

import (
	"reflect"
	"testing"
)

func TestFilterScan(t *testing.T) {
	tables := []struct {
		Data   []byte
		Parsed Filter
		Pos    int
		Err    bool
	}{
		{[]byte("# rule:[Klobasy do spajze]\nif allof (exists \"From\", exists \"Subject\" )\n{\ndiscard;\n}"), Filter{Name: "Klobasy do spajze", Scope: "allof", Rules: []Rule{
			Rule{Type: "from", Operator: "exists", TargetString: StringType{"From"}},
			Rule{Type: "subject", Operator: "exists", TargetString: StringType{"Subject"}},
		}, Actions: []Action{
			Action{Type: "discard"},
		}}, 83, false},

		// errors
		// {[]byte("# rule:[Klobasy do spajze]\n"), Filter{Name: "Klobasy do spajze"}, 0, true},
	}

	for _, tb := range tables {
		var str Filter

		if pos, err := str.Scan(tb.Data); pos != tb.Pos || (!reflect.DeepEqual(str, tb.Parsed) && err == nil) || (err == nil) == tb.Err {
			t.Errorf("Filter %#v GOT: %#v, pos: %v, err: %v, WANT: %#v, pos: %v, err: %v", string(tb.Data), str, pos, err, tb.Parsed, tb.Pos, tb.Err)
		}
	}
}

func TestFilterString(t *testing.T) {
	tables := []struct {
		Expect string
		Src    Filter
	}{

		// first rule
		{"# rule:[Klobasy do spajze]\r\nif allof (exists \"From\",exists \"Subject\")\r\n{\r\n\tdiscard;\r\n}", Filter{Name: "Klobasy do spajze", Scope: "allof", Rules: []Rule{
			Rule{Type: "from", Operator: "exists", TargetString: StringType{"From"}},
			Rule{Type: "subject", Operator: "exists", TargetString: StringType{"Subject"}},
		}, Actions: []Action{
			Action{Type: "discard"},
		}}},

		// second rule
		{"# rule:[Klobasy do spajze]\r\nif allof (exists \"From\",not header :contains [\"Subject\",\"To\"] \"salamka\")\r\n{\r\n\tdiscard;\r\n}", Filter{Name: "Klobasy do spajze", Scope: "allof", Rules: []Rule{
			Rule{Type: "from", Operator: "exists", TargetString: StringType{"From"}},
			Rule{Type: "...", Operator: "notcontains", QueryString: StringType{"salamka"}, TargetString: StringType{"Subject", "To"}},
		}, Actions: []Action{
			Action{Type: "discard"},
		}}},

		// third rule
		{"# rule:[Klobasy do spajze]\r\nif allof (exists \"From\")\r\n{\r\n\tfileinto :copy \"INBOX.MyFolder\";\r\n\tdiscard;\r\n}", Filter{Name: "Klobasy do spajze", Scope: "allof", Rules: []Rule{
			Rule{Type: "from", Operator: "exists", TargetString: StringType{"From"}},
		}, Actions: []Action{
			Action{Type: "fileinto_copy", Values: StringType{"INBOX.MyFolder"}},
			Action{Type: "discard"},
		}}},
	}

	for _, tb := range tables {
		if res := tb.Src.String(); res != tb.Expect {
			t.Errorf("Filter for %#v GOT: %v, WANT: %v", tb.Src, res, tb.Expect)
		}
	}
}
