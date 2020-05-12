package sieve

import (
	"bytes"
	"errors"
	"strings"
)

type Script []Filter

func (fs *Script) Scan(data []byte) error {
	pos := 0
	l := len(data)

	if l == 0 {
		return errors.New("Unexpected zero length of script")
	}

	// ignore require
	if string(data[0:7]) == "require" {
		pos = skipWhiteSpace(data, bytes.IndexByte(data, '\n'))
	}

	// read filters one by one
	for pos < l {
		filter := &Filter{}
		if rl, err := filter.Scan(data[pos:]); err != nil {
			return err
		} else {
			pos = skipWhiteSpace(data, pos+rl)
		}
		*fs = append(*fs, *filter)
	}

	return nil
}

func (fs Script) String() string {
	var res []string
	if req := fs.generateRequire(); req != "" {
		res = append(res, req)
	}
	for _, f := range fs {
		res = append(res, f.String())
	}

	return strings.Join(res, "\r\n")
}

func (fs Script) generateRequire() string {
	var res StringType

	for _, f := range fs {
		res = append(res, f.GetRequires()...)
	}

	if len(res) == 0 {
		return ""
	}

	// deduplicate
	var dd map[stringBase]bool
	var ddRes StringType

	dd = make(map[stringBase]bool)

	for _, r := range res {
		if _, ok := dd[r]; ok {
			continue
		}

		ddRes = append(ddRes, r)
		dd[r] = true
	}

	return "require " + ddRes.String() + ";"
}
