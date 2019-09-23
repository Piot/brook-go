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
	"github.com/piot/brook-go/src/inbitstream"
)

type OutBitStream interface {

	// Get the bit position in the stream
	Tell() uint

	// Rewinds to the specified position and marks that as the end
	Rewind(position uint) error

	// Close flushes the latest changes and closes the stream
	Close()

	// Close flushes the latest changes and closes the stream
	Octets() []byte

	// Copy octets
	CopyOctets(target []byte) uint

	// WriteBitsFromStream copy bits from another stream
	WriteBitsFromStream(in inbitstream.InBitStream, bitCount uint) error

	// WriteBits : Write bits to stream
	WriteBits(v uint32, count uint) error

	// WriteRawBits only for internal use
	WriteRawBits(v uint32, count uint) error

	// WriteSignedBits : Write signed bits to stream
	WriteSignedBits(v int32, count uint) error

	// WriteUint64 : Write bits to stream
	WriteUint64(v uint64) error

	// WriteInt32 : Write bits to stream
	WriteInt32(v int32) error

	// WriteUint32 : Write bits to stream
	WriteUint32(v uint32) error

	// WriteUint16 : Write bits to stream
	WriteUint16(v uint16) error

	// WriteInt16 : Write bits to stream
	WriteInt16(v int16) error

	// WriteUint8 : Write bits from stream
	WriteUint8(v uint8) error
}
