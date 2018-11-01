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

	"github.com/piot/brook-go/src/inbitstream"
	"github.com/piot/brook-go/src/outstream"
)

// OutBitStreamImpl : Read bit stream
type OutBitStreamImpl struct {
	octetWriter   *outstream.OutStream
	remainingBits uint
	ac            uint32
	bitPosition   uint
}

// New : Creates an input bit stream
func New(octetWriter *outstream.OutStream) *OutBitStreamImpl {
	stream := OutBitStreamImpl{octetWriter: octetWriter, ac: 0, remainingBits: 32, bitPosition: 0}
	return &stream
}

func NewWithOption(octetWriter *outstream.OutStream, useDebug bool) OutBitStream {
	impl := New(octetWriter)
	if useDebug {
		return NewDebugStream(impl)
	}
	return impl
}

func maskFromCount(count uint) uint32 {
	return (1 << uint(count)) - 1
}

func (s *OutBitStreamImpl) writeOctets() {
	s.octetWriter.WriteUint32(s.ac)
	s.ac = 0
	s.remainingBits = 32
}

// Tell :
func (s *OutBitStreamImpl) Tell() uint {
	return s.bitPosition
}

// Close :
func (s *OutBitStreamImpl) Close() {
	if s.remainingBits != 32 {
		ac := s.ac
		bitsWritten := 32 - s.remainingBits
		octetCount := ((bitsWritten - 1) / 8) + 1
		outOctets := make([]byte, octetCount)
		for i := uint(0); i < octetCount; i++ {
			out := byte((ac & 0xff000000) >> 24)
			ac <<= 8
			outOctets[i] = out
		}
		s.octetWriter.WriteOctets(outOctets)
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

func (s *OutBitStreamImpl) writeRest(v uint32, count uint, bitsToKeepFromLeft uint) {
	ov := v

	ov >>= uint(count - bitsToKeepFromLeft)
	ov &= maskFromCount(bitsToKeepFromLeft)
	ov <<= s.remainingBits - bitsToKeepFromLeft
	s.remainingBits -= bitsToKeepFromLeft
	s.bitPosition += bitsToKeepFromLeft
	s.ac |= ov
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
