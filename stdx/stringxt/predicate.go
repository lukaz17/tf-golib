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

package stringxt

import (
	"regexp"
	"strings"
)

// A Predicate holds all the requirements to test a string for matching.
type Predicate struct {
	Prefix string
	Suffix string
	Regexp string
}

// Test whether string str satisfied all requirements in the Predicate.
// In the case Regexp is invalid, the result is always false.
func (p *Predicate) Match(str string) (bool, error) {
	result := true
	var err error
	if p.Prefix != "" {
		if !strings.HasPrefix(str, p.Prefix) {
			result = false
		}
	}
	if p.Suffix != "" {
		if !strings.HasSuffix(str, p.Suffix) {
			result = false
		}
	}
	if p.Regexp != "" {
		r, err := regexp.MatchString(p.Regexp, str)
		if !r || err != nil {
			result = false
		}
	}
	return result, err
}
