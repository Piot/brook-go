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

package inbitstream

import (
	"testing"

	"github.com/piot/brook-go/src/instream"
)

func setup() *InBitStream {
	octets := []byte{0xca, 0xfe, 0xde, 0xad, 0xc0, 0xde, 0xff, 0x00}
	reader := instream.New(octets)

	bitstream := New(&reader)
	return bitstream
}

func TestReadTenBits(t *testing.T) {
	bitstream := setup()
	found, firstErr := bitstream.ReadBits(10)
	if firstErr != nil {
		t.Error(firstErr)
	}
	const expected = 1 + 2 + 8 + 32 + 256 + 512
	if found != expected {
		t.Errorf("Expected %d found %d", expected, found)
	}
}

func TestReadMoreThanThirtyBits(t *testing.T) {
	bitstream := setup()

	_, skipErr := bitstream.ReadBits(12)
	if skipErr != nil {
		t.Error(skipErr)
	}

	found, firstErr := bitstream.ReadBits(32)
	if firstErr != nil {
		t.Error(firstErr)
	}

	const expected = 0xedeadc0d
	if found != expected {
		t.Errorf("Expected %X found %X", expected, found)
	}
}

func TestReadTooMuch(t *testing.T) {
	bitstream := setup()

	_, skipErr := bitstream.ReadBits(33)
	if skipErr == nil {
		t.Errorf("Expected error")
	}
}

func TestReadTooFar(t *testing.T) {
	bitstream := setup()

	for i := 0; i < 10; i++ {
		_, skipErr := bitstream.ReadBits(5)
		if skipErr != nil {
			t.Error(skipErr)
		}
	}

	_, skip2Err := bitstream.ReadBits(8)
	if skip2Err != nil {
		t.Error(skip2Err)
	}

	_, readTooFarErr := bitstream.ReadBits(7)
	if readTooFarErr == nil {
		t.Errorf("Expected error")
	}
}
