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

// Close :
func (s *OutBitStream) Close() {
	if s.remainingBits != 32 {
		s.writeOctets()
	}
}

func (s *OutBitStream) writeRest(v uint32, count uint, bitsToKeepFromLeft uint) {
	ov := v

	ov >>= uint(count - bitsToKeepFromLeft)
	ov &= maskFromCount(bitsToKeepFromLeft)
	ov <<= s.remainingBits - bitsToKeepFromLeft
	s.remainingBits -= bitsToKeepFromLeft
	s.ac |= ov
}

// WriteBits : Write bits to stream
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

// WriteSignedBits : Write signed bits to stream
func (s *OutBitStream) WriteSignedBits(v int32, count uint) error {
	sign := 0

	if v < 0 {
		sign = 1
		v = -v
	}

	s.WriteBits(uint32(sign), 1)
	s.WriteBits(uint32(v), count-1)
	return nil
}

// WriteInt32 : Write bits to stream
func (s *OutBitStream) WriteInt32(v int32) error {
	return s.WriteSignedBits(int32(v), 32)
}

// WriteUint32 : Write bits to stream
func (s *OutBitStream) WriteUint32(v uint32) error {
	return s.WriteBits(v, 32)
}

// WriteUint64 : Write bits to stream
func (s *OutBitStream) WriteUint64(v uint64) error {
	upper := uint32(v >> 32)
	s.WriteBits(upper, 32)
	lower := uint32(v & 0xffffffff)
	return s.WriteBits(lower, 32)
}

// WriteUint16 : Write bits to stream
func (s *OutBitStream) WriteUint16(v uint16) error {
	return s.WriteBits(uint32(v), 16)
}

// WriteInt16 : Write bits to stream
func (s *OutBitStream) WriteInt16(v int16) error {
	return s.WriteSignedBits(int32(v), 16)
}

// WriteUint8 : Write bits from stream
func (s *OutBitStream) WriteUint8(v uint8) error {
	return s.WriteBits(uint32(v), 8)
}
