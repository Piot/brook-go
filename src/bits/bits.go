package bits

import (
	"fmt"
	"regexp"
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

func FromString(bitString string) []byte {
	bitString = stripEverythingExceptZeroAndOne(bitString)
	fmt.Println(bitString)
	expression, _ := regexp.Compile("[0|1]{8}")
	matches := expression.FindAllString(bitString, -1)
	b := make([]byte, len(matches))
	for index, octetString := range matches {
		octetIntValue, _ := strconv.ParseUint(octetString, 2, 8)
		b[index] = byte(octetIntValue)
	}

	return b
}
