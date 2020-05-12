package sieve

import (
	"bytes"
	"errors"
	"strings"
)

var scopes = map[string]bool{
	"allof": true,
	"anyof": true,
	"true":  true,
}

type Filter struct {
	Name    string
	Enabled bool
	Scope   string

	Rules   []Rule
	Actions []Action
}

func (f *Filter) Scan(data []byte) (int, error) {
	pos := 0
	l := len(data)

	if l == 0 {
		return 0, errors.New("Unexpected zero length test")
	}

	// get name
	if data[pos] != '#' {
		return 0, errors.New("Rule has no name")
	}

	pos = pos + 8
	endPos := bytes.IndexByte(data[pos:], ']')
	f.Name = string(data[pos : endPos+pos])

	pos = skipWhiteSpace(data, pos+endPos+1) // add end index + ']'

	// if
	endPos = bytes.IndexByte(data[pos:], ' ')
	if endPos == -1 || string(data[pos:pos+endPos]) != "if" {
		return 0, errors.New("Required reserved word \"if\" not find")
	}

	pos = skipWhiteSpace(data, pos+endPos)

	// scope
	endPos = bytes.IndexByte(data[pos:], ' ')

	if endPos == -1 {
		return 0, errors.New("Unexpected end of input near scope")
	}

	f.Scope = string(data[pos : pos+endPos])
	if _, ok := scopes[f.Scope]; !ok {
		return 0, errors.New("Scope " + f.Scope + " is not supported")
	}

	pos = skipWhiteSpace(data, pos+endPos)
	if f.Scope == "true" {
		goto processActions
	}

	// read rules
	if data[pos] != '(' {
		return 0, errors.New("After scope expected \"(\"")
	}

	pos++

	for pos < l {
		pos = skipWhiteSpace(data, pos)

		// process rule
		rule := &Rule{}
		if plen, err := rule.Scan(data[pos:]); err != nil {
			return 0, err
		} else {
			f.Rules = append(f.Rules, *rule)
			pos = skipWhiteSpace(data, pos+plen)
		}

		if data[pos] != ',' {
			break
		}
		pos++
	}

	if data[pos] != ')' {
		return 0, errors.New("Unexpected end of scope, expect \")\"")
	}

	pos++

processActions:

	pos = skipWhiteSpace(data, pos)

	// read actions
	if data[pos] != '{' {
		return 0, errors.New("After scope expected actions block begin \"{\"")
	}

	pos++

	for pos < l {
		pos = skipWhiteSpace(data, pos)

		if data[pos] == '}' {
			break
		}

		// process action
		action := &Action{}
		if plen, err := action.Scan(data[pos:]); err != nil {
			return 0, err
		} else {
			f.Actions = append(f.Actions, *action)
			pos = skipWhiteSpace(data, pos+plen)
		}

		if data[pos] != ';' {
			return 0, errors.New("Unexpected end of action - no semicolon")
		}

		pos++
	}

	if data[pos] != '}' {
		return 0, errors.New("Unexpected end of block, expect \"}\"")
	}

	pos++

	return pos, nil
}

func (f Filter) String() string {

	if len(f.Actions) == 0 {
		return "NO ACTIONS"
	}

	res := []string{"# rule:[" + f.Name + "]"}

	rules := "if " + f.Scope
	var rls []string

	if f.Scope == "true" {
		goto processActions
	}

	if len(f.Rules) == 0 {
		return "NO RULES"
	}

	rules = rules + " ("

	for _, r := range f.Rules {
		rls = append(rls, r.String())
	}

	rules = rules + strings.Join(rls, ",") + ")"

processActions:
	res = append(res, rules, "{")

	for _, a := range f.Actions {
		res = append(res, "\t"+a.String()+";")
	}

	res = append(res, "}")

	return strings.Join(res, "\r\n")
}

func (f Filter) GetRequires() StringType {
	var res StringType

	for _, r := range f.Rules {
		if rres := r.GetRequires(); len(rres) > 0 {
			res = append(res, rres...)
		}
	}

	for _, a := range f.Actions {
		if ares := a.GetRequires(); len(ares) > 0 {
			res = append(res, ares...)
		}
	}

	return res
}
