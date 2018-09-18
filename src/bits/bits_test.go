package bits

import (
	"testing"
)

func TestBitString(t *testing.T) {
	input := "1111    0100 01110101 01111000 10101000"
	octets := FromString(input)
	if len(octets) != 4 {
		t.Errorf("Unexpected length")
	}
	if octets[0] != 0xf4 {
		t.Errorf("Wrong encoding %v", octets)
	}
}
