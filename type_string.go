package sieve

import (
	"errors"
	"strings"
)

type stringBase string

// Escape string
func (s stringBase) String() string {
	res := ""
	for i, l := 0, len(s); i < l; i++ {
		if s[i] == '"' || s[i] == '\\' {
			res += `\`
		}

		res += string(s[i])
	}

	return "\"" + res + "\""
}

func (s *stringBase) Scan(data []byte) (int, error) {
	pos := 1
	l := len(data)

	if l == 0 {
		return 0, errors.New("Unexpected zero string")
	}

	if data[0] != '"' {
		return 0, errorMsg("String dont start with quote: ", data)
	}

	for ; pos+1 < l; pos++ {
		if data[pos] == byte('"') {
			break
		}

		// escaping
		if data[pos] == '\\' && (data[pos+1] == '\\' || data[pos+1] == '"') {
			pos++
		}

		*s = *s + stringBase(data[pos])
	}

	// check end of string, if has double quotes
	if data[pos] != byte('"') {
		return 0, errorMsg("String parsing unexpected end near: ", data)
	}

	pos++ // add one because of last quotes

	return pos, nil
}

type StringType []stringBase

func (s StringType) String() string {
	l := len(s)
	if l > 1 {
		var p []string
		for _, v := range s {
			p = append(p, v.String())
		}
		return "[" + strings.Join(p, ",") + "]"
	} else if l == 1 {
		return s[0].String()
	}

	return ""
}

func (s *StringType) Scan(data []byte) (int, error) {
	l := len(data)
	if l == 0 {
		return 0, errors.New("String parsing has zero size")
	}

	pos := 0
	isArray := false

	if data[pos] == '[' {
		pos++
		isArray = true
		pos = skipWhiteSpace(data, pos)
	}

	// read stringBase
readBase:
	var sb stringBase
	if dpos, err := sb.Scan(data[pos:]); err != nil {
		return 0, err
	} else {
		pos = pos + dpos
		*s = append(*s, sb)
	}

	if isArray {
		pos = skipWhiteSpace(data, pos)
		if pos >= l {
			return 0, errors.New("Unexpected end of string array")
		}

		if data[pos] == ',' {
			pos++
			pos = skipWhiteSpace(data, pos)
			if pos >= l {
				return 0, errors.New("Unexpected end of string array")
			}

			goto readBase
		}

		pos = skipWhiteSpace(data, pos)
		if data[pos] != ']' {
			return 0, errors.New("String array must end with \"]\"")
		}
		pos++
	}

	return pos, nil
}
