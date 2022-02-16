package engine

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

import (
	"bytes"
	"fmt"
)

// GoString converts a Resource object to Go source code.
// This is useful if serialising to a Go file for compilation into a binary.
func (r *StaticResource) GoString() string {
	var buffer bytes.Buffer

	buffer.WriteString("&gui.StaticResource{\n")
	buffer.WriteString("\tStaticName: \"" + r.StaticName + "\",\n")
	buffer.WriteString("\tStaticContent: []byte{\n\t\t")
	for i, v := range r.StaticContent {
		if i > 0 {
			buffer.WriteString(", ")
		}

		buffer.WriteString(fmt.Sprint(v))
	}
	buffer.WriteString("}}")

	return buffer.String()
}
