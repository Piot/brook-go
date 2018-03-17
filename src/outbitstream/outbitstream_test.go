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

func setup() (*OutBitStream, *outstream.OutStream) {
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
