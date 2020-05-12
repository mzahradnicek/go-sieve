package sieve

import (
	"reflect"
	"testing"
)

func TestActionScan(t *testing.T) {
	tables := []struct {
		Data   []byte
		Parsed Action
		Pos    int
		Err    bool
	}{
		{[]byte("fileinto \"INBOX.Drafts\""), Action{Type: "fileinto", Values: StringType{"INBOX.Drafts"}}, 23, false},
		{[]byte("fileinto :copy \"INBOX.Drafts\""), Action{Type: "fileinto_copy", Values: StringType{"INBOX.Drafts"}}, 29, false},

		{[]byte("redirect \"INBOX.Drafts\""), Action{Type: "redirect", Values: StringType{"INBOX.Drafts"}}, 23, false},
		{[]byte("redirect :copy \"INBOX.Drafts\""), Action{Type: "redirect_copy", Values: StringType{"INBOX.Drafts"}}, 29, false},

		{[]byte("reject \"Message of discard\""), Action{Type: "reject", Values: StringType{"Message of discard"}}, 27, false},
		{[]byte("discard"), Action{Type: "discard"}, 7, false},

		{[]byte("setflag [\"\\\\Seen\",\"\\\\Answered\",\"MojFlag\"]"), Action{Type: "setflag", Values: StringType{"\\Seen", "\\Answered", "MojFlag"}}, 41, false},
		{[]byte("addflag [\"\\\\Seen\",\"\\\\Answered\",\"MojFlag\"]"), Action{Type: "addflag", Values: StringType{"\\Seen", "\\Answered", "MojFlag"}}, 41, false},
		{[]byte("removeflag [\"\\\\Seen\",\"\\\\Answered\",\"MojFlag\"]"), Action{Type: "removeflag", Values: StringType{"\\Seen", "\\Answered", "MojFlag"}}, 44, false},
	}

	for _, tb := range tables {
		var str Action

		if pos, err := str.Scan(tb.Data); pos != tb.Pos || (!reflect.DeepEqual(str, tb.Parsed) && err == nil) || (err == nil) == tb.Err {
			t.Errorf("Action %#v GOT: %#v, pos: %v, err: %v, WANT: %#v, pos: %v, err: %v", string(tb.Data), str, pos, err, tb.Parsed, tb.Pos, tb.Err)
		}
	}
}

func TestActionString(t *testing.T) {
	tables := []struct {
		Expect string
		Src    Action
	}{
		{"fileinto \"INBOX.Drafts\"", Action{Type: "fileinto", Values: StringType{"INBOX.Drafts"}}},
		{"fileinto :copy \"INBOX.Drafts\"", Action{Type: "fileinto_copy", Values: StringType{"INBOX.Drafts"}}},

		{"redirect \"INBOX.Drafts\"", Action{Type: "redirect", Values: StringType{"INBOX.Drafts"}}},
		{"redirect :copy \"INBOX.Drafts\"", Action{Type: "redirect_copy", Values: StringType{"INBOX.Drafts"}}},

		{"reject \"Message of discard\"", Action{Type: "reject", Values: StringType{"Message of discard"}}},
		{"discard", Action{Type: "discard"}},

		{"setflag [\"\\\\Seen\",\"\\\\Answered\",\"MojFlag\"]", Action{Type: "setflag", Values: StringType{"\\Seen", "\\Answered", "MojFlag"}}},
		{"addflag [\"\\\\Seen\",\"\\\\Answered\",\"MojFlag\"]", Action{Type: "addflag", Values: StringType{"\\Seen", "\\Answered", "MojFlag"}}},
		{"removeflag [\"\\\\Seen\",\"\\\\Answered\",\"MojFlag\"]", Action{Type: "removeflag", Values: StringType{"\\Seen", "\\Answered", "MojFlag"}}},
	}

	for _, tb := range tables {
		if res := tb.Src.String(); res != tb.Expect {
			t.Errorf("Test for %#v GOT: %v, WANT: %#v", tb.Src, res, tb.Expect)
		}
	}
}
