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

package outstream

import (
	"bytes"
	"testing"
)

func Test16and32bit(t *testing.T) {
	octets := []byte{0xca, 0xfe, 0xde, 0xad, 0xc0, 0xde}
	stream := New()

	err1 := stream.WriteUint16(51966)
	if err1 != nil {
		t.Error(err1)
	}

	err2 := stream.WriteUint32(3735929054)
	if err2 != nil {
		t.Error(err1)
	}

	createdOctets := stream.Octets()
	if !bytes.Equal(octets, createdOctets) {
		t.Errorf("Not equal")
	}
}
