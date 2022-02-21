//go:build fribidi && !js
// +build fribidi,!js

package text

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

//#cgo CPPFLAGS: -I/usr/include/fribidi
//#cgo LDFLAGS: -L/usr/lib -lfribidi
/*
#include <fribidi.h>
*/
import "C"
import (
	"unsafe"
)

// Bidi maps the string from its logical order to the visual order to correctly display mixed LTR/RTL text. It returns a mapping of rune positions.
func Bidi(text string) (string, []int) {
	str := []rune(text)
	pbaseDir := C.FriBidiParType(C.FRIBIDI_PAR_ON) // neutral direction
	visualStr := make([]rune, len(str))
	positionsL2V := make([]C.FriBidiStrIndex, len(str))
	positionsV2L := make([]C.FriBidiStrIndex, len(str))
	embeddingLevels := make([]C.FriBidiLevel, len(str))
	C.fribidi_log2vis(
		// input
		(*C.FriBidiChar)(unsafe.Pointer(&str[0])),
		C.FriBidiStrIndex(len(str)),
		&pbaseDir,

		// output
		(*C.FriBidiChar)(unsafe.Pointer(&visualStr[0])),
		&positionsL2V[0],
		&positionsV2L[0],
		&embeddingLevels[0],
	)
	text = string(visualStr)

	mapV2L := make([]int, len(positionsV2L))
	for i, pos := range positionsV2L {
		mapV2L[i] = int(pos)
	}
	return text, mapV2L
}
