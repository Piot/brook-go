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

package outbitstream

import (
	"github.com/piot/brook-go/src/inbitstream"
)

type OutBitStreamDebug struct {
	stream OutBitStream
}

func NewDebugStream(stream OutBitStream) *OutBitStreamDebug {
	return &OutBitStreamDebug{stream: stream}
}

func (o *OutBitStreamDebug) Tell() uint {
	return o.stream.Tell()
}

func (o *OutBitStreamDebug) Rewind(position uint) error {
	return o.stream.Rewind(position)
}

func (o *OutBitStreamDebug) Close() {
	o.stream.Close()
}

func (o *OutBitStreamDebug) WriteBitsFromStream(in inbitstream.InBitStream, bitCount uint) error {
	return o.stream.WriteBitsFromStream(in, bitCount)
}

func (o *OutBitStreamDebug) WriteBits(v uint32, count uint) error {
	o.writeType(7, count)
	return o.stream.WriteBits(v, count)
}

func (o *OutBitStreamDebug) WriteRawBits(v uint32, count uint) error {
	return o.stream.WriteRawBits(v, count)
}

func (o *OutBitStreamDebug) WriteSignedBits(v int32, count uint) error {
	o.writeType(6, count)
	return o.stream.WriteSignedBits(v, count)
}

func (o *OutBitStreamDebug) WriteInt32(v int32) error {
	o.writeType(9, 32)
	return o.stream.WriteInt32(v)
}

func (o *OutBitStreamDebug) WriteUint32(v uint32) error {
	o.writeType(3, 32)
	return o.stream.WriteUint32(v)
}

func (o *OutBitStreamDebug) WriteUint64(v uint64) error {
	o.writeType(4, 64)
	return o.stream.WriteUint64(v)
}

func (o *OutBitStreamDebug) WriteUint16(v uint16) error {
	o.writeType(1, 16)
	return o.stream.WriteUint16(v)
}

func (o *OutBitStreamDebug) WriteInt16(v int16) error {
	o.writeType(2, 16)
	return o.stream.WriteInt16(v)
}

func (o *OutBitStreamDebug) WriteUint8(v uint8) error {
	o.writeType(5, 8)
	return o.stream.WriteUint8(v)
}

func (o *OutBitStreamDebug) writeType(t int, bitCount uint) {
	o.stream.WriteBits(uint32(t), 4)
	o.stream.WriteBits(uint32(bitCount), 7)
}

func (o *OutBitStreamDebug) Octets() []byte {
	return o.stream.Octets()
}

func (o *OutBitStreamDebug) CopyOctets(target []byte) uint {
	return o.stream.CopyOctets(target)
}
