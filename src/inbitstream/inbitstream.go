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
	"io"

	"github.com/piot/brook-go/src/instream"
)

// InBitStream : Read bit stream
type InBitStream struct {
	octetReader           *instream.InStream
	remainingBits         uint
	data                  uint32
	remainingBitsInStream uint
	position              uint
}

// New : Creates an input bit stream
func New(octetReader *instream.InStream, bitCount uint) *InBitStream {
	stream := InBitStream{octetReader: octetReader, data: 0, remainingBits: 0, remainingBitsInStream: bitCount, position: 0}
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

	s.remainingBitsInStream -= bitsToRead
	mask := maskFromCount(bitsToRead)
	shiftPos := uint(s.remainingBits - bitsToRead)
	bits := (s.data >> shiftPos) & mask
	s.remainingBits -= bitsToRead
	return bits, nil
}

func (s *InBitStream) fill() error {
	maxOctetsToRead := uint(4)
	newData := uint32(0)
	octetsRead := 0
	for i := uint(0); i < maxOctetsToRead; i++ {
		octet, readOctetErr := s.octetReader.ReadUint8()
		if readOctetErr != nil {
			if readOctetErr == io.EOF {
				break
			}
			return readOctetErr
		}
		octetsRead++
		octetValue := uint32(octet)
		newData <<= 8
		newData |= octetValue
	}

	s.data = newData
	s.remainingBits = uint(octetsRead * 8)
	return nil
}

// ReadBits : Read bits from stream
func (s *InBitStream) ReadBits(count uint) (uint32, error) {
	if count > 32 {
		return 0, fmt.Errorf("Max 32 bits to read")
	}

	if count > s.remainingBitsInStream {
		return 0, fmt.Errorf("End of stream")
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
		v2, secondCountErr := s.readOnce(secondCount)
		if secondCountErr != nil {
			return 0, secondCountErr
		}
		v |= v2
		return v, nil
	}
	return s.readOnce(count)
}

// ReadSignedBits : Read signed bits from stream
func (s *InBitStream) ReadSignedBits(count uint) (int32, error) {
	sign, signErr := s.ReadBits(1)
	if signErr != nil {
		return 0, signErr
	}

	v, vErr := s.ReadBits(count - 1)
	if vErr != nil {
		return 0, vErr
	}

	signed := int32(v)

	if sign != 0 {
		signed = -signed
	}

	return signed, nil
}

// ReadUint64 : Read unsigned 64-bit from stream
func (s *InBitStream) ReadUint64() (uint64, error) {
	upper, err := s.ReadBits(32)
	if err != nil {
		return 0, err
	}
	r := uint64(upper) << 32
	lower, lowerErr := s.ReadBits(32)
	r |= uint64(lower)
	return r, lowerErr
}

// ReadUint32 : Read unsigned 32-bit from stream
func (s *InBitStream) ReadUint32() (uint32, error) {
	v, err := s.ReadBits(32)
	return uint32(v), err
}

// ReadUint16 : Read unsigned 16-bit from stream
func (s *InBitStream) ReadUint16() (uint16, error) {
	v, err := s.ReadBits(16)
	return uint16(v), err
}

// ReadInt16 : Read unsigned 16-bit from stream
func (s *InBitStream) ReadInt16() (int16, error) {
	v, err := s.ReadSignedBits(16)
	return int16(v), err
}

// ReadUint8 : Read unsigned 8-bit from stream
func (s *InBitStream) ReadUint8() (uint8, error) {
	v, err := s.ReadBits(8)
	return uint8(v), err
}

func (s *InBitStream) String() string {
	return fmt.Sprintf("[inbitstream buf:%v]", s.octetReader)
}
