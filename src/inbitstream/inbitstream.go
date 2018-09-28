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

type InBitStream interface {
	ReadBits(count uint) (uint32, error)
	ReadRawBits(count uint) (uint32, error)
	// ReadSignedBits : Read signed bits from stream
	ReadSignedBits(count uint) (int32, error)
	// ReadUint64 : Read unsigned 64-bit from stream
	ReadUint64() (uint64, error)
	// ReadUint32 : Read unsigned 32-bit from stream
	ReadUint32() (uint32, error)
	// ReadUint16 : Read unsigned 16-bit from stream
	ReadUint16() (uint16, error)
	// ReadInt16 : Read unsigned 16-bit from stream
	ReadInt16() (int16, error)
	// ReadUint8 : Read unsigned 8-bit from stream
	ReadUint8() (uint8, error)

	RemainingBitCount() uint
}
