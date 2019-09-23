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
)

// InBitStreamImpl : Read bit stream
type InBitStreamImpl struct {
	octets                []byte
	remainingBits         uint
	data                  uint32
	remainingBitsInStream uint
	position              uint
	tell                  uint
	octetReadPosition     int
}

// New : Creates an input bit stream
func New(octets []byte, bitCount uint) *InBitStreamImpl {
	stream := InBitStreamImpl{octets: octets, data: 0, remainingBits: 0, remainingBitsInStream: bitCount, position: 0}
	return &stream
}

func NewWithOption(octets []byte, bitCount uint, useDebugStream bool) InBitStream {
	s := New(octets, bitCount)
	if useDebugStream {
		return NewDebugStream(s)
	}
	return s
}

func (s *InBitStreamImpl) Octets() []byte {
	return s.octets
}

func (s *InBitStreamImpl) Seek(position uint) error {
	if position != 0 {
		return fmt.Errorf("Only seek to zero now")
	}
	s.remainingBits = 0
	s.octetReadPosition = 0
	s.position = 0
	s.data = 0
	s.fill()
	return nil
}

func (s *InBitStreamImpl) IsEOF() bool {
	return s.remainingBitsInStream == 0
}

func maskFromCount(count uint) uint32 {
	return (1 << uint(count)) - 1
}

func (s *InBitStreamImpl) readOnce(bitsToRead uint) (uint32, error) {
	if bitsToRead == 0 {
		return 0, nil
	}

	if bitsToRead > s.remainingBitsInStream {
		return 0, &EOFError{Count: bitsToRead, Tell: s.tell}
	}

	s.remainingBitsInStream -= bitsToRead
	mask := maskFromCount(bitsToRead)
	shiftPos := uint(s.remainingBits - bitsToRead)
	bits := (s.data >> shiftPos) & mask
	s.tell += bitsToRead
	s.remainingBits -= bitsToRead
	return bits, nil
}

func (s *InBitStreamImpl) Tell() uint {
	return s.tell
}

func (s *InBitStreamImpl) fill() error {
	maxOctetsToRead := int(4)
	newData := uint32(0)
	remainingOctetCount := len(s.octets) - s.octetReadPosition
	if remainingOctetCount <= 0 {
		return &EOFError{}
	}
	octetCountToRead := maxOctetsToRead
	if octetCountToRead > remainingOctetCount {
		octetCountToRead = remainingOctetCount
	}
	for i := 0; i < octetCountToRead; i++ {
		octet := s.octets[s.octetReadPosition+i]
		rotateCount := uint((3 - i) * 8)
		rotatedOctet := uint32(octet) << rotateCount
		newData |= rotatedOctet
	}
	s.octetReadPosition += octetCountToRead
	s.data = newData
	s.remainingBits = 32
	return nil
}

func (s *InBitStreamImpl) ReadRawBits(count uint) (uint32, error) {
	return s.ReadBits(count)
}

func (s *InBitStreamImpl) Skip(count uint) error {
	dwordCount := count / 32
	restBitCount := count % 32
	for i := uint(0); i < dwordCount; i++ {
		_, dwordErr := s.ReadRawBits(32)
		if dwordErr != nil {
			return dwordErr
		}
	}
	_, err := s.ReadRawBits(restBitCount)
	return err
}

// ReadBits : Read bits from stream
func (s *InBitStreamImpl) ReadBits(count uint) (uint32, error) {
	if count > 32 {
		return 0, fmt.Errorf("Max 32 bits to read")
	}

	if count > s.remainingBitsInStream {
		return 0, &EOFError{Count: count, Tell: s.tell}
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
func (s *InBitStreamImpl) ReadSignedBits(count uint) (int32, error) {
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
func (s *InBitStreamImpl) ReadUint64() (uint64, error) {
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
func (s *InBitStreamImpl) ReadUint32() (uint32, error) {
	v, err := s.ReadBits(32)
	return uint32(v), err
}

// ReadUint16 : Read unsigned 16-bit from stream
func (s *InBitStreamImpl) ReadUint16() (uint16, error) {
	v, err := s.ReadBits(16)
	return uint16(v), err
}

// ReadInt16 : Read unsigned 16-bit from stream
func (s *InBitStreamImpl) ReadInt16() (int16, error) {
	v, err := s.ReadSignedBits(16)
	return int16(v), err
}

// ReadUint8 : Read unsigned 8-bit from stream
func (s *InBitStreamImpl) ReadUint8() (uint8, error) {
	v, err := s.ReadBits(8)
	return uint8(v), err
}

func (s *InBitStreamImpl) String() string {
	return fmt.Sprintf("[inbitstream buf:%v]", s.octets)
}
