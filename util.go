package ws

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"unsafe"
)




func bsplit3(bts []byte, sep byte) (b1, b2, b3 []byte) {
	a := bytes.IndexByte(bts, sep)
	b := bytes.IndexByte(bts[a+1:], sep)
	if a == -1 || b == -1 {
		return bts, nil, nil
	}
	b += a + 1
	return bts[:a], bts[a+1 : b], bts[b+1:]
}

func readLine(br *bufio.Reader) ([]byte, error) {
	var line []byte
	for {
		bts, err := br.ReadSlice('\n')
		if errors.Is(err, bufio.ErrBufferFull) {
			// Copy bytes because next read will discard them.
			line = append(line, bts...)
			continue
		}

		// Avoid copy of single read.
		if line == nil {
			line = bts
		} else {
			line = append(line, bts...)
		}

		if err != nil {
			return line, err
		}

		// Size of line is at least 1.
		// In other case bufio.ReadSlice() returns error.
		n := len(line)

		// Cut '\n' or '\r\n'.
		if n > 1 && line[n-2] == '\r' {
			line = line[:n-2]
		} else {
			line = line[:n-1]
		}

		return line, nil
	}
}

func Bytes2String(b []byte) string {
	return unsafe.String(&b[0], len(b))
}


func StringToBytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

func BytesToString(b []byte) string {
	return unsafe.String(&b[0], len(b))
}

func  Equal(a, b []byte) bool {
    return BytesToString(a) == BytesToString(b)
}


// asciiToInt converts bytes to int.
func asciiToInt(bts []byte) (ret int, err error) {
	// ASCII numbers all start with the high-order bits 0011.
	// If you see that, and the next bits are 0-9 (0000 - 1001) you can grab those
	// bits and interpret them directly as an integer.
	var n int
	if n = len(bts); n < 1 {
		return 0, fmt.Errorf("converting empty bytes to int")
	}
	for i := 0; i < n; i++ {
		if bts[i]&0xf0 != 0x30 {
			return 0, fmt.Errorf("%s is not a numeric character", string(bts[i]))
		}
		ret += int(bts[i]&0xf) * pow(10, n-i-1)
	}
	return ret, nil
}


// pow for integers implementation.
// See Donald Knuth, The Art of Computer Programming, Volume 2, Section 4.6.3
func pow(a, b int) int {
	p := 1
	for b > 0 {
		if b&1 != 0 {
			p *= a
		}
		b >>= 1
		a *= a
	}
	return p
}