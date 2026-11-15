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

import (
	"math/big"
)

var (
	BigIntN1  = big.NewInt(-1)
	BigInt0   = big.NewInt(0)
	BigInt1   = big.NewInt(1)
	BigInt2   = big.NewInt(2)
	BigInt3   = big.NewInt(3)
	BigInt4   = big.NewInt(4)
	BigInt7   = big.NewInt(7)
	BigInt16  = big.NewInt(16)
	BigInt32  = big.NewInt(32)
	BigInt64  = big.NewInt(64)
	BigInt255 = big.NewInt(255)
)

// A Uint256 is wrapped *big.Int to intergrate with system that use large number
// and provide fluent apis for type casting.
type Uint256 big.Int

// Cast Uint256 to *big.Int.
func (s *Uint256) BigInt() *big.Int {
	return (*big.Int)(s)
}

// Return underlying value of Bytes as *big.Int.
func (s *Uint256) Value() *big.Int {
	return (*big.Int)(s)
}

// Parse a string into big.Int. Decimal doesn't have prefix. Hexadecimal, Binary,
// and Octal must have prefix 0x, 0b, 0o respectively.
func BigIntFromString(s string) (*big.Int, bool) {
	if s == "" {
		return new(big.Int), true
	}
	var bigint *big.Int
	var ok bool
	if len(s) >= 2 && s[0] == '0' {
		switch s[1] {
		case 'x', 'X':
			bigint, ok = new(big.Int).SetString(s[2:], 16)
		case 'b', 'B':
			bigint, ok = new(big.Int).SetString(s[2:], 2)
		case 'o', 'O':
			bigint, ok = new(big.Int).SetString(s[2:], 8)
		default:
			bigint, ok = new(big.Int).SetString(s, 10)
		}
	} else {
		bigint, ok = new(big.Int).SetString(s, 10)
	}
	return bigint, ok
}
