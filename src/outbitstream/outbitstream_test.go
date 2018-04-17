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

package outbitstream

import (
	"encoding/binary"
	"testing"

	"github.com/piot/brook-go/src/outstream"
)

func setup() (OutBitStream, *outstream.OutStream) {
	writer := outstream.New()

	bitstream := New(writer)
	return bitstream, writer
}

func TestWriteMoreThanThirtyBits(t *testing.T) {
	bitstream, octetWriter := setup()
	firstErr := bitstream.WriteBits(0xcafed, 20)
	if firstErr != nil {
		t.Error(firstErr)
	}

	secondErr := bitstream.WriteBits(0xbeef, 16)
	if secondErr != nil {
		t.Error(secondErr)
	}

	octets := octetWriter.Octets()
	if len(octets) != 4 {
		t.Errorf("Wrong length:%d", len(octets))
	}

	readFromBuffer := binary.BigEndian.Uint32(octets)
	const expected uint32 = 0xcafedbee
	if readFromBuffer != expected {
		t.Errorf("Expected %d got %d", expected, readFromBuffer)
	}

}

// swapUint32 converts a uint16 to network byte order and back.
func swapUint32(n uint32) uint32 {
	return (n&0x000000FF)<<24 | (n&0x0000FF00)<<8 |
		(n&0x00FF0000)>>8 | (n&0xFF000000)>>24
}

func printBits(n uint32) string {
	s := ""
	for i := 0; i < 32; i++ {
		bit := (n & 0x80000000) != 0
		n <<= 1
		if i%4 == 0 {
			s += " "
		}
		if bit {
			s += "1"
		} else {
			s += "0"
		}

	}
	return s
}

func TestWriteMoreThanThirtyBitsDebug(t *testing.T) {
	_bitstream, octetWriter := setup()
	bitstream := NewDebugStream(_bitstream)
	firstErr := bitstream.WriteBits(0xcafed, 20)
	if firstErr != nil {
		t.Error(firstErr)
	}

	secondErr := bitstream.WriteBits(0xbeef, 16)
	if secondErr != nil {
		t.Error(secondErr)
	}

	octets := octetWriter.Octets()
	if len(octets) != 4 {
		t.Errorf("Wrong length:%d", len(octets))
	}

	readFromBuffer := binary.BigEndian.Uint32(octets)
	const expected uint32 = 0x6532BFB5
	if readFromBuffer != expected {
		t.Errorf("Expected %d got %v %08X", expected, printBits((readFromBuffer)), readFromBuffer)
	}

}

func checkOctetLength(t *testing.T, octetWriter *outstream.OutStream, expectedLength int) {
	octets := octetWriter.Octets()
	if len(octets) != expectedLength {
		t.Errorf("Wrong length:%d expected %d", len(octets), expectedLength)
	}

}

func TestOctetLength(t *testing.T) {
	bitstream, octetWriter := setup()
	checkOctetLength(t, octetWriter, 0)
	firstErr := bitstream.WriteBits(0xcafe, 32)
	if firstErr != nil {
		t.Error(firstErr)
	}

	secondErr := bitstream.WriteBits(0x3, 2)
	if secondErr != nil {
		t.Error(secondErr)
	}

	bitstream.Close()

	checkOctetLength(t, octetWriter, 5)

	if octetWriter.Octets()[4] != 0xc0 {
		t.Errorf("Didn't work")
	}

}
