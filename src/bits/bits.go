package bits

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

func stripWhitespace(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, str)
}

func FromString(bitString string) []byte {
	bitString = stripWhitespace(bitString)
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
