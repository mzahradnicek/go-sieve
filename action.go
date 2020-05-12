package sieve

import (
	"bytes"
	"errors"
	"strings"
)

var action = []struct {
	Name string
	Src  []byte
}{ // dont change order here !!!!
	{"fileinto_copy", []byte(`fileinto :copy`)}, // string		// fileinto
	{"fileinto", []byte(`fileinto`)},            // string		// fileinto
	{"redirect_copy", []byte(`redirect :copy`)}, // string
	{"redirect", []byte(`redirect`)},            // string
	{"reject", []byte(`reject`)},                // string
	{"discard", []byte(`discard`)},
	{"setflag", []byte(`setflag`)},       // string		// imap4flags
	{"addflag", []byte(`addflag`)},       // string		// imap4flags
	{"removeflag", []byte(`removeflag`)}, // string	// imap4flags
	// {"vacation", []byte(``)},
	// {"set", []byte(``)},
	// {"notify", []byte(``)},
	// {"keep", []byte(``)},
	// {"stop", []byte(`stop`)},
}

type Action struct {
	Type   string     `json:"type"`
	Values StringType `json:"values"`
}

func (a *Action) Scan(data []byte) (int, error) {
	// get test command
	pos := 0
	l := len(data)

	if l == 0 {
		return 0, errors.New("Unexpected zero length test")
	}

	for _, act := range action {
		if bytes.HasPrefix(data, act.Src) {
			a.Type = act.Name
			pos = pos + len(act.Src)
			break
		}
	}

	if a.Type == "" {
		return 0, errors.New("Unknow action")
	}

	if a.Type == "discard" { // has no params
		return pos, nil
	}

	pos = skipWhiteSpace(data, pos)

	if i, err := a.Values.Scan(data[pos:]); err != nil {
		return 0, err
	} else {
		pos = pos + i
	}

	return pos, nil
}

func (a Action) String() string {
	var res []string

	for _, act := range action {
		if a.Type == act.Name {
			res = append(res, string(act.Src))
			break
		}
	}

	if len(a.Values) > 0 {
		res = append(res, a.Values.String())
	}

	return strings.Join(res, " ")
}

func (a Action) GetRequires() StringType {
	switch a.Type {
	case "vacation", "reject":
		return StringType{stringBase(a.Type)}

	case "fileinto":
		return StringType{"fileinto"}

	case "fileinto_copy":
		return StringType{"fileinto", "copy"}

	case "redirect_copy":
		return StringType{"copy"}

	case "setflag", "addflag", "removeflag":
		return StringType{"imap4flags"}
	}

	return StringType{}
}
