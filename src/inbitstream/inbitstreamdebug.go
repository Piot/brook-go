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

import "fmt"

type InBitStreamDebug struct {
	stream InBitStream
}

func debugTypeValueToString(expectedType int) string {
	switch expectedType {
	case 1:
		return "uint16"
	case 2:
		return "int16"
	case 3:
		return "uint32"
	case 4:
		return "uint64"
	case 5:
		return "uint8"
	case 6:
		return "signed"
	case 7:
		return "unsigned"
	}

	return "unknown"
}

func NewDebugStream(stream InBitStream) *InBitStreamDebug {
	return &InBitStreamDebug{stream: stream}
}

func (i *InBitStreamDebug) Tell() uint {
	tryInfo, tryInfoWorked := i.stream.(InBitStreamInfo)
	if !tryInfoWorked {
		return 0
	}
	return tryInfo.Tell()
}

func (i *InBitStreamDebug) IsEOF() bool {
	return i.stream.IsEOF()
}

func (i *InBitStreamDebug) checkType(expectedType int, expectedBitCount uint) error {
	t, tErr := i.internalRead(4)
	if tErr != nil {
		return tErr
	}
	if int(t) != expectedType {
		return fmt.Errorf("Expected %v but received %v (%v vs %v)", debugTypeValueToString(expectedType), debugTypeValueToString(int(t)), expectedType, t)
	}

	bitCount, bitCountErr := i.internalRead(7)
	if bitCountErr != nil {
		return bitCountErr
	}
	if uint(bitCount) != expectedBitCount {
		return fmt.Errorf("Expected %v count but received %v bitcount (type:%v %v)", expectedBitCount, bitCount, debugTypeValueToString(expectedType), expectedType)
	}

	// fmt.Printf("Verified %v %v\n", expectedType, expectedBitCount)
	return nil
}

func (i *InBitStreamDebug) internalRead(count uint) (uint32, error) {
	return i.stream.ReadBits(count)
}

func (i *InBitStreamDebug) ReadRawBits(count uint) (uint32, error) {
	return i.internalRead(count)
}

func (i *InBitStreamDebug) ReadBits(count uint) (uint32, error) {
	checkErr := i.checkType(7, count)
	if checkErr != nil {
		return 0, checkErr
	}
	return i.stream.ReadBits(count)
}

// ReadSignedBits : Read signed bits from stream
func (i *InBitStreamDebug) ReadSignedBits(count uint) (int32, error) {
	checkErr := i.checkType(6, count)
	if checkErr != nil {
		return 0, checkErr
	}
	return i.stream.ReadSignedBits(count)
}

// ReadUint64 : Read unsigned 64-bit from stream
func (i *InBitStreamDebug) ReadUint64() (uint64, error) {
	checkErr := i.checkType(4, 64)
	if checkErr != nil {
		return 0, checkErr
	}
	return i.stream.ReadUint64()
}

// ReadUint32 : Read unsigned 32-bit from stream
func (i *InBitStreamDebug) ReadUint32() (uint32, error) {
	checkErr := i.checkType(3, 32)
	if checkErr != nil {
		return 0, checkErr
	}
	return i.stream.ReadUint32()
}

// ReadUint16 : Read unsigned 16-bit from stream
func (i *InBitStreamDebug) ReadUint16() (uint16, error) {
	checkErr := i.checkType(1, 16)
	if checkErr != nil {
		return 0, checkErr
	}
	return i.stream.ReadUint16()
}

// ReadInt16 : Read unsigned 16-bit from stream
func (i *InBitStreamDebug) ReadInt16() (int16, error) {
	checkErr := i.checkType(2, 16)
	if checkErr != nil {
		return 0, checkErr
	}
	return i.stream.ReadInt16()
}

// ReadUint8 : Read unsigned 8-bit from stream
func (i *InBitStreamDebug) ReadUint8() (uint8, error) {
	checkErr := i.checkType(5, 8)
	if checkErr != nil {
		return 0, checkErr
	}
	return i.stream.ReadUint8()
}
