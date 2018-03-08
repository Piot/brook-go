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

// Package outstream ...
package outstream

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
)

// OutStream : Write to octet stream
type OutStream struct {
	buffer bytes.Buffer
}

// New : Creates an output stream
func New() *OutStream {
	stream := OutStream{}
	return &stream
}

// Feed : Adds octets to stream
func (s *OutStream) Feed(octets []byte) error {
	lengthWritten, err := s.buffer.Write(octets)
	if err != nil {
		return err
	}

	if lengthWritten != len(octets) {
		err := errors.New("couldn't write all octets")
		return err
	}
	if false {
		hexPayload := hex.Dump(s.buffer.Bytes())
		fmt.Printf("Buffer is now:%s\n", hexPayload)
	}
	return nil
}

// WriteUint32 : Writes an unsigned 32-bit integer to stream
func (s *OutStream) WriteUint32(v uint32) error {
	temp := make([]byte, 4)
	binary.BigEndian.PutUint32(temp, v)
	s.Feed(temp)
	return nil
}

// WriteUint64 : Writes an unsigned 64-bit integer to stream
func (s *OutStream) WriteUint64(v uint64) error {
	temp := make([]byte, 8)
	binary.BigEndian.PutUint64(temp, v)
	s.Feed(temp)
	return nil
}

// WriteUint16 : Writes an unsigned 16-bit integer to stream
func (s *OutStream) WriteUint16(v uint16) error {
	temp := make([]byte, 2)
	binary.BigEndian.PutUint16(temp, v)
	s.Feed(temp)
	return nil
}

// WriteUint8 : Writes an octet to stream
func (s *OutStream) WriteUint8(v uint8) error {
	temp := make([]byte, 1)
	temp[0] = v
	s.Feed(temp)
	return nil
}

// WriteOctets : Writes octets to stream
func (s *OutStream) WriteOctets(octets []byte) error {
	return s.Feed(octets)
}

// Octets : Gets the written octets
func (s *OutStream) Octets() []byte {
	return s.buffer.Bytes()
}

// String : Outputs a debug string of the stream
func (s *OutStream) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("[outstream ")
	buffer.WriteString(fmt.Sprintf("size: %d", s.buffer.Len()))
	buffer.WriteString("]")
	return buffer.String()
}
