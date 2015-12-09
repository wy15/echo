package netstring

import (
	"bytes"
	"errors"
	"strconv"
)

var (
	ErrNoColon             = errors.New("Can't found colon")
	ErrNsLenBytesConv      = errors.New("nsLenBytes conv error")
	ErrNoComma             = errors.New("Can't found comma")
	ErrNsLenNotEqaulOrgLen = errors.New("nsLen does not equal to orgLen")
)

/*
utf-8
<31 32 3a 68 65 6c 6c 6f 20 77 6f 72 6c 64 21 2c>
12:hello world!,
<30 3a 2c>
0:,
*/
func Marshall(org []byte) []byte {
	orgLen := len(org)
	if orgLen == 0 {
		return nil
	}
	lenString := strconv.Itoa(orgLen)
	return append(append(append([]byte(lenString), ':'), org...), ',')
}

func Unmarshall(ns []byte) ([]byte, error) {
	if ns == nil {
		return nil, errors.New("ns is nil")
	}
	colonPos := bytes.Index(ns, []byte(":"))
	if colonPos == -1 {
		return nil, ErrNoColon
	}
	nsLenBytes := ns[:colonPos]
	nsLen, err := strconv.Atoi(string(nsLenBytes))
	if err != nil {
		return nil, ErrNsLenBytesConv
	}
	commaPos := bytes.LastIndex(ns, []byte(","))
	if commaPos == -1 {
		return nil, ErrNoComma
	}
	if commaPos-1 == colonPos {
		return nil, errors.New("org is nil")
	}

	if nsLen == commaPos-colonPos-1 {
		return ns[colonPos+1 : commaPos], nil
	} else {
		return nil, ErrNsLenNotEqaulOrgLen
	}
}
