package widget

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
	"fmt"
	"testing"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/data/binding"
	"github.com/bhojpur/gui/pkg/engine/test"
	"github.com/stretchr/testify/assert"
)

func TestNewProgressBarWithData(t *testing.T) {
	val := binding.NewFloat()
	val.Set(0.4)

	label := NewProgressBarWithData(val)
	waitForBinding()
	assert.Equal(t, 0.4, label.Value)
}

func TestProgressBar_Binding(t *testing.T) {
	bar := NewProgressBar()
	assert.Equal(t, 0.0, bar.Value)

	val := binding.NewFloat()
	val.Set(0.1)
	bar.Bind(val)
	waitForBinding()
	assert.Equal(t, 0.1, bar.Value)

	val.Set(0.4)
	waitForBinding()
	assert.Equal(t, 0.4, bar.Value)

	bar.Unbind()
	waitForBinding()
	assert.Equal(t, 0.4, bar.Value)
}

func TestProgressBar_SetValue(t *testing.T) {
	bar := NewProgressBar()

	assert.Equal(t, 0.0, bar.Min)
	assert.Equal(t, 1.0, bar.Max)
	assert.Equal(t, 0.0, bar.Value)

	bar.SetValue(.5)
	assert.Equal(t, .5, bar.Value)
}

func TestProgressBar_TextFormatter(t *testing.T) {
	bar := NewProgressBar()
	formatted := false

	bar.SetValue(0.2)
	assert.Equal(t, false, formatted)

	formatter := func() string {
		formatted = true
		return fmt.Sprintf("%.2f out of %.2f", bar.Value, bar.Max)
	}
	bar.TextFormatter = formatter

	bar.SetValue(0.4)

	assert.Equal(t, true, formatted)
}

func TestProgressRenderer_Layout(t *testing.T) {
	bar := NewProgressBar()
	bar.Resize(gui.NewSize(100, 10))

	render := test.WidgetRenderer(bar).(*progressRenderer)
	assert.Equal(t, float32(0), render.bar.Size().Width)

	bar.SetValue(.5)
	assert.Equal(t, float32(50), render.bar.Size().Width)

	bar.SetValue(1)
	assert.Equal(t, bar.Size().Width, render.bar.Size().Width)
}

func TestProgressRenderer_Layout_Overflow(t *testing.T) {
	bar := NewProgressBar()
	bar.Resize(gui.NewSize(100, 10))

	render := test.WidgetRenderer(bar).(*progressRenderer)
	bar.SetValue(1)
	assert.Equal(t, bar.Size().Width, render.bar.Size().Width)

	bar.SetValue(1.2)
	assert.Equal(t, bar.Size().Width, render.bar.Size().Width)
}

func TestProgressRenderer_ApplyTheme(t *testing.T) {
	bar := NewProgressBar()
	render := test.WidgetRenderer(bar).(*progressRenderer)

	textSize := render.label.TextSize
	customTextSize := textSize
	test.WithTestTheme(t, func() {
		render.applyTheme()
		customTextSize = render.label.TextSize
	})

	assert.NotEqual(t, textSize, customTextSize)
}
