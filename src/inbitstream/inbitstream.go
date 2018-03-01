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

// Package inbitstream ...
package inbitstream

import (
	"fmt"

	"github.com/piot/brook-go/src/instream"
)

// InBitStream : Read bit stream
type InBitStream struct {
	octetReader   *instream.InStream
	remainingBits uint
	data          uint32
}

// New : Creates an input bit stream
func New(octetReader *instream.InStream) *InBitStream {
	stream := InBitStream{octetReader: octetReader, data: 0, remainingBits: 0}
	return &stream
}

func maskFromCount(count uint) uint32 {
	return (1 << uint(count)) - 1
}

func (s *InBitStream) readOnce(bitsToRead uint) (uint32, error) {
	if bitsToRead == 0 {
		return 0, nil
	}

	if bitsToRead > s.remainingBits {
		return 0, fmt.Errorf("Read passed end of stream")
	}
	mask := maskFromCount(bitsToRead)
	shiftPos := uint(s.remainingBits - bitsToRead)
	bits := (s.data >> shiftPos) & mask
	s.remainingBits -= bitsToRead
	return bits, nil
}

func (s *InBitStream) fill() error {
	octetsToRead := uint(4)

	newData := uint32(0)
	for i := uint(0); i < octetsToRead; i++ {
		newData <<= 8
		octet, readOctetErr := s.octetReader.ReadUint8()
		if readOctetErr != nil {
			return readOctetErr
		}
		newData |= uint32(octet)
	}

	s.data = newData
	s.remainingBits = octetsToRead * 8
	return nil
}

// ReadBits : Read bits from stream
func (s *InBitStream) ReadBits(count uint) (uint32, error) {
	if count > 32 {
		return 0, fmt.Errorf("Max 32 bits to read")
	}

	if count > s.remainingBits {
		secondCount := uint(count - s.remainingBits)
		v, firstErr := s.readOnce(s.remainingBits)
		if firstErr != nil {
			return 0, firstErr
		}
		fillErr := s.fill()
		if fillErr != nil {
			return 0, fillErr
		}
		v <<= secondCount
		v2, _ := s.readOnce(secondCount)
		v |= v2
		return v, nil
	}
	return s.readOnce(count)
}
