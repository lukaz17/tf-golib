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

package pseudorng

import "strings"

// Base64 returns a pseudo-random base64 string with charset defined in RFC 4648 for UTF-8
// from the default [Source].
// The length (n) denotes number of bytes the base64 string represented, not the
// actual number of characters of the generated base64 string.
func Base64(n uint) string {
	charset := []rune{
		'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		'+', '/',
	}
	ret := String(n+ceilUint(n, 3), charset)
	switch n % 3 {
	case 1:
		return ret + "=="
	case 2:
		return ret + "="
	default:
		return ret
	}
}

// Hex returns a pseudo-random hex string in lower case
// from the default [Source].
// The length (n) denotes number of bytes the hex string represented, not the
// actual number of characters of the generated hex string.
func Hex(n uint) string {
	charset := []rune{
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		'a', 'b', 'c', 'd', 'e', 'f',
	}
	return String(n*2, charset)
}

// String returns a pseudo-random string with specific length (n) and based on
// pre-defined list of characters from the default [Source].
// The length of charset should be less than 2³¹ for reliability.
// The probability for appearances of each character cannot be customized.
func String(n uint, charset []rune) string {
	var sb strings.Builder
	charcount := len(charset)
	for i := uint(0); i < n; i++ {
		sb.WriteRune(charset[Intn(charcount)])
	}
	return sb.String()
}

// ceilUint returns ceiling result of x / y as same math.Ceil.
// See https://stackoverflow.com/a/2745086 for more information.
func ceilUint(x, y uint) uint {
	if x == 0 {
		return 0
	}
	return 1 + ((x - 1) / y)
}
