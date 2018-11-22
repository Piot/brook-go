/*

MIT License

Copyright (c) 2017 Peter Bjorklund

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

*/

package bits

import (
	"strconv"
	"strings"
)

func stripEverythingExceptZeroAndOne(str string) string {
	return strings.Map(func(r rune) rune {
		if r != '0' && r != '1' {
			return -1
		}
		return r
	}, str)
}

func FromString(bitString string) ([]byte, uint) {
	bitString = stripEverythingExceptZeroAndOne(bitString)
	bitCount := len(bitString)
	remaining := bitCount % 8
	if remaining != 0 {
		bitString += strings.Repeat("0", 8-remaining)
	}

	octetCount := len(bitString) / 8
	b := make([]byte, octetCount)
	for i := 0; i < octetCount; i++ {
		octetString := bitString[i*8 : i*8+8]
		octetIntValue, _ := strconv.ParseUint(octetString, 2, 8)
		b[i] = byte(octetIntValue)
	}

	return b, uint(bitCount)
}

func ToString(octets []byte) string {
	s := ""
	for i := 0; i < len(octets); i++ {
		o := octets[i]
		for j := 0; j < 8; j++ {
			testMask := uint(1 << uint(7-j))
			if uint(o)&testMask != 0 {
				s += "1"
			} else {
				s += "0"
			}

		}
	}

	return s
}
