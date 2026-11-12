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

package securerng

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"strings"
)

// Base64 returns a uniform base64 string with charset defined in RFC 4648 for UTF-8
// from the default [Reader].
// The length (n) denotes number of bytes the base64 string represented, not the
// actual number of characters of the generated base64 string.
func Base64(n uint) string {
	var sb strings.Builder
	i := uint(0)
	buflen := uint(30)
	for i < n {
		var buf []byte
		if i+buflen <= n {
			buf = make([]byte, buflen)
		} else {
			buf = make([]byte, n-i)
		}
		_, err := rand.Read(buf)
		if err != nil {
			panic(err)
		}
		sb.WriteString(base64.StdEncoding.EncodeToString(buf))
		i += buflen
	}
	return sb.String()
}

// Hex returns a uniform hex string in lower case
// from the default [Reader].
// The length (n) denotes number of bytes the hex string represented, not the
// actual number of characters of the generated hex string.
func Hex(n uint) string {
	var sb strings.Builder
	i := uint(0)
	buflen := uint(32)
	for i < n {
		var buf []byte
		if i+buflen <= n {
			buf = make([]byte, buflen)
		} else {
			buf = make([]byte, n-i)
		}
		_, err := rand.Read(buf)
		if err != nil {
			panic(err)
		}
		sb.WriteString(hex.EncodeToString(buf))
		i += buflen
	}
	return sb.String()
}

// String returns a uniform string with specific length (n) and based on
// pre-defined list of characters from the default [Reader].
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
