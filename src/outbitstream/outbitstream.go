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

// Package outbitstream ...
package outbitstream

import (
	"fmt"

	"github.com/piot/brook-go/src/outstream"
)

// OutBitStream : Read bit stream
type OutBitStream struct {
	octetWriter   *outstream.OutStream
	remainingBits uint
	ac            uint32
}

// New : Creates an input bit stream
func New(octetWriter *outstream.OutStream) *OutBitStream {
	stream := OutBitStream{octetWriter: octetWriter, ac: 0, remainingBits: 32}
	return &stream
}

func maskFromCount(count uint) uint32 {
	return (1 << uint(count)) - 1
}

func (s *OutBitStream) writeOctets() {
	s.octetWriter.WriteUint32(s.ac)
	s.ac = 0
	s.remainingBits = 32
}

func (s *OutBitStream) writeRest(v uint32, count uint, bitsToKeepFromLeft uint) {
	ov := v

	ov >>= uint(count - bitsToKeepFromLeft)
	ov &= maskFromCount(bitsToKeepFromLeft)
	ov <<= s.remainingBits - bitsToKeepFromLeft
	s.remainingBits -= bitsToKeepFromLeft
	s.ac |= ov
}

// WriteBits : Write bits from stream
func (s *OutBitStream) WriteBits(v uint32, count uint) error {
	if count > 32 {
		return fmt.Errorf("Max 32 bits to write")
	}

	if count > s.remainingBits {
		firstWriteCount := s.remainingBits
		s.writeRest(v, count, firstWriteCount)
		s.writeOctets()
		s.writeRest(v, count-firstWriteCount, count-firstWriteCount)
	} else {
		s.writeRest(v, count, count)
	}

	return nil
}
