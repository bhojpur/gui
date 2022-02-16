package dialog

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
	"testing"

	"github.com/bhojpur/gui/pkg/engine/test"
	"github.com/stretchr/testify/assert"
)

func TestEntryDialog_Confirm(t *testing.T) {
	value := ""
	ed := NewEntryDialog("Test", "message", func(v string) {
		value = v
	}, test.NewWindow(nil))
	ed.Show()
	test.Type(ed.entry, "123")
	test.Tap(ed.confirm)

	assert.Equal(t, value, "123", "Control form should be confirmed with no validation")
}

func TestEntryDialog_Dismiss(t *testing.T) {
	value := "123"
	ed := NewEntryDialog("Test", "message", func(v string) {
		value = v
	}, test.NewWindow(nil))
	ed.Show()
	test.Type(ed.entry, "XYZ")
	test.Tap(ed.cancel)

	assert.Equal(t, value, "123", "Control form should not change value on dismiss")
}
