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

// Package instream ...
package instream

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
)

// InStream : Read octet stream
type InStream struct {
	buffer   bytes.Buffer
	position int
}

// New : Creates an input stream
func New(octets []byte) InStream {
	buf := bytes.NewBuffer(octets)
	stream := InStream{buffer: *buf, position: 0}

	return stream
}

func (stream *InStream) Tell() int {
	return stream.position
}

// IsEOF : Checks if the input stream is empty
func (stream *InStream) IsEOF() bool {
	return stream.buffer.Len() == 0
}

// Read : Reads octets from the stream
func (stream *InStream) Read(octetCount int) ([]byte, error) {
	tempBuffer := make([]byte, octetCount)
	lengthWritten, err := stream.buffer.Read(tempBuffer)
	if err != nil {
		return nil, err
	}

	if lengthWritten != octetCount {
		err := errors.New("Couldn't read all octets")
		return nil, err
	}
	stream.position += octetCount
	if false {
		hexPayload := hex.Dump(stream.buffer.Bytes())
		fmt.Printf("Buffer is now:%s", hexPayload)
	}
	return tempBuffer, nil
}

// ReadOctets : Reads octets from the stream
func (stream *InStream) ReadOctets(octetCount int) ([]byte, error) {
	return stream.Read(octetCount)
}

// ReadUint64 reads an unsigned 64-bit integer from the stream
func (stream *InStream) ReadUint64() (uint64, error) {
	temp, err := stream.Read(8)
	if err != nil {
		return 0, err
	}
	v := binary.BigEndian.Uint64(temp)
	return v, nil
}

// ReadUint32 : Reads an unsigned 32-bit integer from the stream
func (stream *InStream) ReadUint32() (uint32, error) {
	temp, err := stream.Read(4)
	if err != nil {
		return 0, err
	}
	v := binary.BigEndian.Uint32(temp)
	return v, nil
}

// ReadUint16 : Reads an unsigned 16-bit integer from the stream
func (stream *InStream) ReadUint16() (uint16, error) {
	temp, err := stream.Read(2)
	if err != nil {
		return 0, err
	}
	v := binary.BigEndian.Uint16(temp)
	return v, nil
}

// ReadUint8 : Reads an octet from the stream
func (stream *InStream) ReadUint8() (uint8, error) {
	v, err := stream.Read(1)
	if err != nil {
		return 0, err
	}
	return v[0], nil
}

func (stream *InStream) String() string {
	return fmt.Sprintf("[instream buffer size:%d]", stream.buffer.Len())
}
