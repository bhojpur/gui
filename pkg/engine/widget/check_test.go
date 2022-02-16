package widget_test

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

	"github.com/stretchr/testify/assert"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/data/binding"
	"github.com/bhojpur/gui/pkg/engine/layout"
	"github.com/bhojpur/gui/pkg/engine/test"
	"github.com/bhojpur/gui/pkg/engine/widget"
)

func TestCheck_Binding(t *testing.T) {
	c := widget.NewCheck("", nil)
	c.SetChecked(true)
	assert.Equal(t, true, c.Checked)

	val := binding.NewBool()
	c.Bind(val)
	waitForBinding()
	assert.Equal(t, false, c.Checked)

	err := val.Set(true)
	assert.Nil(t, err)
	waitForBinding()
	assert.Equal(t, true, c.Checked)

	c.SetChecked(false)
	v, err := val.Get()
	assert.Nil(t, err)
	assert.Equal(t, false, v)

	c.Unbind()
	waitForBinding()
	assert.Equal(t, false, c.Checked)
}

func TestCheck_Layout(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	for name, tt := range map[string]struct {
		text     string
		checked  bool
		disabled bool
	}{
		"checked": {
			text:    "Test",
			checked: true,
		},
		"unchecked": {
			text: "Test",
		},
		"checked_disabled": {
			text:     "Test",
			checked:  true,
			disabled: true,
		},
		"unchecked_disabled": {
			text:     "Test",
			disabled: true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			check := &widget.Check{
				Text:    tt.text,
				Checked: tt.checked,
			}
			if tt.disabled {
				check.Disable()
			}

			window := test.NewWindow(gui.NewContainerWithLayout(layout.NewCenterLayout(), check))
			window.Resize(check.MinSize().Max(gui.NewSize(150, 200)))

			test.AssertRendersToMarkup(t, "check/layout_"+name+".xml", window.Canvas())

			window.Close()
		})
	}
}

func TestNewCheckWithData(t *testing.T) {
	val := binding.NewBool()
	err := val.Set(true)
	assert.Nil(t, err)

	c := widget.NewCheckWithData("", val)
	waitForBinding()
	assert.Equal(t, true, c.Checked)

	c.SetChecked(false)
	v, err := val.Get()
	assert.Nil(t, err)
	assert.Equal(t, false, v)
}
