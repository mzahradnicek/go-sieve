package sieve

import (
	"reflect"
	"testing"
)

func TestScriptScan(t *testing.T) {
	tables := []struct {
		Data   []byte
		Parsed Script
		Err    bool
	}{
		{[]byte("require [\"klobaa asa\"]\r\n# rule:[Klobasy do spajze]\nif allof (exists \"From\", exists \"Subject\" )\n{\ndiscard;\n}\r\n# rule:[Salamova pochutka]\r\nif allof (exists \"From\",not header :contains [\"Subject\",\"To\"] \"salamka\")\r\n{\r\n\tdiscard;\r\n}\r\n  "), Script{
			Filter{Name: "Klobasy do spajze", Scope: "allof", Rules: []Rule{
				Rule{Type: "from", Operator: "exists", TargetString: StringType{"From"}},
				Rule{Type: "subject", Operator: "exists", TargetString: StringType{"Subject"}},
			}, Actions: []Action{
				Action{Type: "discard"},
			}},
			Filter{Name: "Salamova pochutka", Scope: "allof", Rules: []Rule{
				Rule{Type: "from", Operator: "exists", TargetString: StringType{"From"}},
				Rule{Type: "...", Operator: "notcontains", QueryString: StringType{"salamka"}, TargetString: StringType{"Subject", "To"}},
			}, Actions: []Action{
				Action{Type: "discard"},
			}},
		}, false},
	}

	for _, tb := range tables {
		var str Script

		if err := str.Scan(tb.Data); (!reflect.DeepEqual(str, tb.Parsed) && err == nil) || (err == nil) == tb.Err {
			t.Errorf("Filter %#v GOT: %#v, err: %v, WANT: %#v, err: %v", string(tb.Data), str, err, tb.Parsed, tb.Err)
		}
	}
}

func TestScriptString(t *testing.T) {
	tables := []struct {
		Expect string
		Src    Script
	}{
		{"require [\"fileinto\",\"copy\"];\r\n# rule:[Klobasy do spajze]\r\nif allof (exists \"From\",exists \"Subject\")\r\n{\r\n\tfileinto :copy \"INBOX.asdf\";\r\n\tfileinto \"INBOX.asdf\";\r\n}\r\n# rule:[Salamova pochutka]\r\nif allof (exists \"From\",not header :contains [\"Subject\",\"To\"] \"salamka\")\r\n{\r\n\tdiscard;\r\n}", Script{
			Filter{Name: "Klobasy do spajze", Scope: "allof", Rules: []Rule{
				Rule{Type: "from", Operator: "exists", TargetString: StringType{"From"}},
				Rule{Type: "subject", Operator: "exists", TargetString: StringType{"Subject"}},
			}, Actions: []Action{
				Action{Type: "fileinto_copy", Values: StringType{"INBOX.asdf"}},
				Action{Type: "fileinto", Values: StringType{"INBOX.asdf"}},
			}},
			Filter{Name: "Salamova pochutka", Scope: "allof", Rules: []Rule{
				Rule{Type: "from", Operator: "exists", TargetString: StringType{"From"}},
				Rule{Type: "...", Operator: "notcontains", QueryString: StringType{"salamka"}, TargetString: StringType{"Subject", "To"}},
			}, Actions: []Action{
				Action{Type: "discard"},
			}},
		}},
	}

	for _, tb := range tables {
		if res := tb.Src.String(); res != tb.Expect {
			t.Errorf("Script for %#v GOT: %v, WANT: %v", tb.Src, res, tb.Expect)
		}
	}
}
