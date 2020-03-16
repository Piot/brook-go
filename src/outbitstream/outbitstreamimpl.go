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
	"encoding/binary"
	"fmt"

	"github.com/piot/brook-go/src/inbitstream"
)

// OutBitStreamImpl : Read bit stream
type OutBitStreamImpl struct {
	bitsInAccumulator uint
	ac                uint32
	bitPosition       uint
	octetPosition     uint
	octetArray        []byte
}

// New : Creates an input bit stream
func New(octetCount int) *OutBitStreamImpl {
	if (octetCount % 4) != 0 {
		octetCount += 4 - octetCount%4
	}
	stream := OutBitStreamImpl{octetArray: make([]byte, octetCount)}
	return &stream
}

func NewWithOption(octetCount int, useDebug bool) OutBitStream {
	impl := New(octetCount)
	if useDebug {
		return NewDebugStream(impl)
	}
	return impl
}

func maskFromCount(count uint) uint32 {
	return (1 << uint(count)) - 1
}

func (s *OutBitStreamImpl) writeAccumulatorToArray() error {
	if s.octetPosition+4 >= uint(len(s.octetArray)) {
		return fmt.Errorf("write accumulator: octet positions outside octet array (%v out of %v)",
			s.octetPosition+4, len(s.octetArray))
	}

	unusedBitCount := 32 - s.bitsInAccumulator
	dwordToWrite := s.ac << unusedBitCount
	binary.BigEndian.PutUint32(s.octetArray[s.octetPosition:s.octetPosition+4], dwordToWrite)

	return nil
}

// Tell :
func (s *OutBitStreamImpl) Tell() uint {
	return s.bitPosition
}

// Rewind :
func (s *OutBitStreamImpl) Rewind(position uint) error {
	if s.octetPosition+4 < uint(len(s.octetArray)) {
		flushErr := s.writeAccumulatorToArray()
		if flushErr != nil {
			return fmt.Errorf("rewind: %w", flushErr)
		}
	}
	s.bitPosition = position
	s.octetPosition = position / 32
	if s.octetPosition+4 > uint(len(s.octetArray)) {
		return fmt.Errorf("seeked too far %v vs %v", position, len(s.octetArray)*8)
	}
	a := binary.BigEndian.Uint32(s.octetArray[s.octetPosition : s.octetPosition+4])
	bitCountToUse := position % 32
	bitCountToFlush := 32 - bitCountToUse
	a >>= bitCountToFlush
	s.bitsInAccumulator = bitCountToUse
	s.ac = a
	return nil
}

// Close :
func (s *OutBitStreamImpl) Close() {
	if s.bitsInAccumulator != 0 {
		s.writeAccumulatorToArray()
	}
}

func (s *OutBitStreamImpl) WriteBitsFromStream(in inbitstream.InBitStream, bitCount uint) error {
	lastBitCount := uint(bitCount % 32)
	for i := uint(0); i < bitCount/32; i++ {
		data, readErr := in.ReadRawBits(32)
		if readErr != nil {
			return readErr
		}
		writeErr := s.WriteRawBits(data, 32)
		if writeErr != nil {
			return writeErr
		}
	}
	data, lastReadErr := in.ReadRawBits(lastBitCount)
	if lastReadErr != nil {
		return lastReadErr
	}
	lastWriteErr := s.WriteRawBits(data, lastBitCount)
	if lastWriteErr != nil {
		return lastWriteErr
	}
	return nil
}

func (s *OutBitStreamImpl) addBitsToAccumulator(v uint32, count uint) {
	ov := v
	ov &= maskFromCount(count)
	s.ac <<= count
	s.ac |= ov
	s.bitsInAccumulator += count
	s.bitPosition += count
	if s.bitsInAccumulator > 32 {
		panic("wrong logic in bitstream")
	}
}

// WriteRawBits : Write bits to stream
func (s *OutBitStreamImpl) WriteRawBits(v uint32, count uint) error {
	return s.WriteBits(v, count)
}

// WriteBits : Write bits to stream
func (s *OutBitStreamImpl) WriteBits(v uint32, count uint) error {
	if count > 32 {
		return fmt.Errorf("Max 32 bits to write")
	}

	bitCountLeftInAc := 32 - s.bitsInAccumulator
	if count > bitCountLeftInAc {
		shiftValue := count - bitCountLeftInAc
		ov := v >> shiftValue
		s.addBitsToAccumulator(ov, bitCountLeftInAc)
		flushErr := s.writeAccumulatorToArray()
		if flushErr != nil {
			return fmt.Errorf("WriteBits: %w", flushErr)
		}
		s.octetPosition += 4
		s.bitsInAccumulator = 0
		s.ac = 0
		s.addBitsToAccumulator(v, count-bitCountLeftInAc)
	} else {
		s.addBitsToAccumulator(v, count)
	}

	return nil
}

// WriteSignedBits : Write signed bits to stream
func (s *OutBitStreamImpl) WriteSignedBits(v int32, count uint) error {
	sign := uint32(0)
	var uv uint32
	if v < 0 {
		sign = 1
		uv = uint32(-v)
	} else {
		uv = uint32(v)
	}

	signWriteErr := s.WriteBits(uint32(sign), 1)
	if signWriteErr != nil {
		return signWriteErr
	}
	valueWriteErr := s.WriteBits(uv, count-1)
	if valueWriteErr != nil {
		return valueWriteErr
	}
	return nil
}

// WriteInt32 : Write bits to stream
func (s *OutBitStreamImpl) WriteInt32(v int32) error {
	return s.WriteSignedBits(int32(v), 32)
}

// WriteUint32 : Write bits to stream
func (s *OutBitStreamImpl) WriteUint32(v uint32) error {
	return s.WriteBits(v, 32)
}

// WriteUint64 : Write bits to stream
func (s *OutBitStreamImpl) WriteUint64(v uint64) error {
	upper := uint32(v >> 32)
	s.WriteBits(upper, 32)
	lower := uint32(v & 0xffffffff)
	return s.WriteBits(lower, 32)
}

// WriteUint16 : Write bits to stream
func (s *OutBitStreamImpl) WriteUint16(v uint16) error {
	return s.WriteBits(uint32(v), 16)
}

// WriteInt16 : Write bits to stream
func (s *OutBitStreamImpl) WriteInt16(v int16) error {
	return s.WriteSignedBits(int32(v), 16)
}

// WriteUint8 : Write bits from stream
func (s *OutBitStreamImpl) WriteUint8(v uint8) error {
	return s.WriteBits(uint32(v), 8)
}

func (s *OutBitStreamImpl) Octets() []byte {
	s.writeAccumulatorToArray()
	octetCountWrittenTo := (s.bitPosition + 7) / 8
	return s.octetArray[0:octetCountWrittenTo]
}

func (s *OutBitStreamImpl) CopyOctets(target []byte) uint {
	s.writeAccumulatorToArray()
	octetCountWrittenTo := (s.bitPosition + 7) / 8
	copy(target, s.octetArray[0:octetCountWrittenTo])
	return octetCountWrittenTo
}
