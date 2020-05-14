package sieve

import (
	"bytes"
	"errors"
	"strings"
)

type Rule struct {
	Type         string     `json:"type"`
	Operator     string     `json:"operator"`
	TargetString StringType `json:"target_string,omitempty"`
	QueryString  StringType `json:"query_string,omitempty"`
	NumValue     NumberType `json:"numeric_value,omitempty"`
}

var typesTempl = map[string]bool{
	"from":    true,
	"to":      true,
	"subject": true,
}

var operators = []struct {
	Name string
	Src  []byte
}{
	{"contains", []byte(`:contains`)},
	{"is", []byte(`:is`)},
	{"matches", []byte(`:matches`)},
	{"regex", []byte(`:regex`)},
	{"count-gt", []byte(`:count "gt"`)},
	{"count-ge", []byte(`:count "ge"`)},
	{"count-lt", []byte(`:count "lt"`)},
	{"count-le", []byte(`:count "le"`)},
	{"count-eq", []byte(`:count "eq"`)},
	{"count-ne", []byte(`:count "ne"`)},
	{"value-gt", []byte(`:value "gt"`)},
	{"value-ge", []byte(`:value "ge"`)},
	{"value-lt", []byte(`:value "lt"`)},
	{"value-le", []byte(`:value "le"`)},
	{"value-eq", []byte(`:value "eq"`)},
	{"value-ne", []byte(`:value "ne"`)},
	{"over", []byte(`:over`)},
	{"under", []byte(`:under`)},
}

func (t *Rule) Scan(data []byte) (int, error) {
	// get test command
	pos := 0
	l := len(data)

	if l == 0 {
		return 0, errors.New("Unexpected zero length test")
	}

	testName := ""

getRuleName:
	for ; pos < l; pos++ {
		if isWhiteSpace(data[pos]) {
			// move to next argument
			pos = skipWhiteSpace(data, pos)

			// check if there is not
			if testName == "not" {
				t.Operator = "not"
				testName = ""
				goto getRuleName
			}

			break
		}
		testName = testName + string(data[pos])
	}

	switch testName {
	case "header":
		// read operator
		if i, err := t.getOperator(data[pos:]); err != nil {
			return 0, err
		} else {
			pos = pos + i
		}

		// target string
		pos = skipWhiteSpace(data, pos)

		if i, err := t.processTarget(data[pos:]); err != nil {
			return 0, err
		} else {
			pos = pos + i
		}

		// query string
		pos = skipWhiteSpace(data, pos)
		if i, err := t.QueryString.Scan(data[pos:]); err != nil {
			return 0, err
		} else {
			pos = pos + i
		}

	case "body":
		t.Type = "body"

		// skip :text
		pos = pos + 6

		// read operator
		if i, err := t.getOperator(data[pos:]); err != nil {
			return 0, err
		} else {
			pos = pos + i
		}

		// query string
		pos = skipWhiteSpace(data, pos)
		if i, err := t.QueryString.Scan(data[pos:]); err != nil {
			return 0, err
		} else {
			pos = pos + i
		}

	case "exists":
		t.Operator = t.Operator + "exists"

		// target string
		pos = skipWhiteSpace(data, pos)

		if i, err := t.processTarget(data[pos:]); err != nil {
			return 0, err
		} else {
			pos = pos + i
		}
	}

	return pos, nil
}

func (t *Rule) getOperator(data []byte) (int, error) {
	for _, o := range operators {
		if bytes.HasPrefix(data, o.Src) {
			t.Operator = t.Operator + o.Name
			return len(o.Src), nil
		}
	}

	return 0, errors.New("Unknown operator")
}

func (t Rule) getOperatorSrc(name string) string {
	for _, o := range operators {
		if o.Name == name {
			return string(o.Src)
		}
	}

	return ""
}

func (t *Rule) processTarget(data []byte) (int, error) {
	pos := 0
	if i, err := t.TargetString.Scan(data); err != nil {
		return 0, err
	} else {
		pos = pos + i
	}

	t.Type = "..."
	if len(t.TargetString) == 1 {
		tstr := strings.ToLower(string(t.TargetString[0]))

		if typesTempl[tstr] {
			t.Type = tstr
		}
	}

	return pos, nil
}

func (t Rule) String() string {
	var res []string

	if strings.HasPrefix(t.Operator, "not") {
		res = append(res, "not")
		t.Operator = t.Operator[3:]
	}

	// preprocess type
	if typesTempl[t.Type] {
		t.TargetString = StringType{stringBase(strings.Title(t.Type))}
	}

	switch {
	case t.Type == "body":
		res = append(res, "body :text", t.getOperatorSrc(t.Operator), t.QueryString.String())
	case t.Operator == "exists":
		res = append(res, "exists", t.TargetString.String())
	default: // header
		res = append(res, "header", t.getOperatorSrc(t.Operator), t.TargetString.String(), t.QueryString.String())
	}

	return strings.Join(res, " ")
}

func (t Rule) GetRequires() StringType {
	var res StringType

	// Type, Operator
	if strings.HasPrefix(t.Operator, "count") || strings.HasPrefix(t.Operator, "value") {
		res = append(res, "relational")
	}

	if t.Type == "body" {
		res = append(res, "body")
	}

	return res
}
