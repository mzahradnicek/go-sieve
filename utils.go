package sieve

import "errors"

func skipWhiteSpace(data []byte, pos int) int {
	l := len(data)

	if l == 0 {
		return pos
	}

	for pos < l && (data[pos] == byte(' ') || data[pos] == byte('\t') || data[pos] == byte('\r') || data[pos] == byte('\n')) {
		pos = pos + 1
	}

	return pos
}

func isWhiteSpace(b byte) bool {
	return b == byte(' ') || b == byte('\n') || b == byte('\t') || b == byte('\r')
}

func errorMsg(msg string, data []byte) error {
	l := len(data)
	if l > 20 {
		l = 19
	}

	return errors.New(msg + `"` + string(data[:l]) + `"`)
}
