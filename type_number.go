package sieve

import (
	"errors"
)

// Number type
type NumberType string

func (n NumberType) String() string {
	return string(n)
}

func (n *NumberType) Scan(data []byte) (int, error) {
	pos := 0
	l := len(data)

	if l == 0 {
		return 0, errors.New("Unexpected zero length number")
	}

	for pos < l {
		if (data[pos] < '0' || data[pos] > '9') && ((data[pos] != 'K' && data[pos] != 'M' && data[pos] != 'G') || (pos+1 < l && !isWhiteSpace(data[pos+1]))) {
			return 0, errorMsg("Number has wrong character ", data)
		}

		*n = *n + NumberType(data[pos])

		pos++
		if pos == l || isWhiteSpace(data[pos]) {
			break
		}
	}

	return pos, nil
}
