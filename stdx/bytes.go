// Copyright (c) 2025 T-Force I/O
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package stdx

import "encoding/hex"

// A Bytes is wrapped byte slice to provide fluent api for type casting.
type Bytes []byte

// Cast Bytes as byte slice.
func (s Bytes) ByteArr() []byte {
	return s
}

// Cast Bytes to a Hex.
func (s Bytes) Hex() *Hex {
	return NewHex(s, false)
}

// Cast Bytes to a hex string without "0x" prefix.
func (s Bytes) HexStr() string {
	return hex.EncodeToString(s)
}

// Return underlying value of Bytes as byte slice.
func (s Bytes) Value() []byte {
	return s
}

// A Hex is another presetnation of hex string to provide fluent api for type casting.
type Hex struct {
	value     Bytes
	hasPrefix bool
}

// Create new instance of Hex using Bytes.
// Must indicate whether that Hex has prefix of "0x".
func NewHex(value Bytes, hasPrefix bool) *Hex {
	return &Hex{value, hasPrefix}
}

// Cast Hex as Bytes.
func (s *Hex) Bytes() Bytes {
	return s.value
}

// Cast Hex as byte slice.
func (s *Hex) ByteArr() []byte {
	return s.value
}

// Return underlying value of Hex as string.
func (s *Hex) Value() string {
	if s.hasPrefix {
		return "0x" + s.value.HexStr()
	}
	return s.value.HexStr()
}
