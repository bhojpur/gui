package chart

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import "strings"

// SplitCSV splits a corpus by the `,`, dropping leading or trailing whitespace unless quoted.
func SplitCSV(text string) (output []string) {
	if len(text) == 0 {
		return
	}

	var state int
	var word []rune
	var opened rune
	for _, r := range text {
		switch state {
		case 0: // word
			if isQuote(r) {
				opened = r
				state = 1
			} else if isCSVDelim(r) {
				output = append(output, strings.TrimSpace(string(word)))
				word = nil
			} else {
				word = append(word, r)
			}
		case 1: // we're in a quoted section
			if matchesQuote(opened, r) {
				state = 0
				continue
			}
			word = append(word, r)
		}
	}

	if len(word) > 0 {
		output = append(output, strings.TrimSpace(string(word)))
	}
	return
}

func isCSVDelim(r rune) bool {
	return r == rune(',')
}

func isQuote(r rune) bool {
	return r == '"' || r == '\'' || r == '“' || r == '”' || r == '`'
}

func matchesQuote(a, b rune) bool {
	if a == '“' && b == '”' {
		return true
	}
	if a == '”' && b == '“' {
		return true
	}
	return a == b
}
