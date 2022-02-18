package settings

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

	"github.com/bhojpur/gui/pkg/engine/canvas"
	"github.com/bhojpur/gui/pkg/engine/test"
	"github.com/bhojpur/gui/pkg/engine/theme"
	"github.com/bhojpur/gui/pkg/engine/widget"
	"github.com/stretchr/testify/assert"
)

func TestChooseScale(t *testing.T) {
	s := &Settings{}
	s.guiSettings.Scale = 1.0
	buttons := s.makeScaleButtons()

	test.Tap(buttons[0].(*widget.Button))
	assert.Equal(t, float32(0.5), s.guiSettings.Scale)
	assert.Equal(t, widget.HighImportance, buttons[0].(*widget.Button).Importance)
	assert.Equal(t, widget.MediumImportance, buttons[2].(*widget.Button).Importance)
}

func TestMakeScaleButtons(t *testing.T) {
	s := &Settings{}
	s.guiSettings.Scale = 1.0
	buttons := s.makeScaleButtons()

	assert.Equal(t, 5, len(buttons))
	assert.Equal(t, widget.MediumImportance, buttons[0].(*widget.Button).Importance)
	assert.Equal(t, widget.HighImportance, buttons[2].(*widget.Button).Importance)
}

func TestMakeScalePreviews(t *testing.T) {
	s := &Settings{}
	s.guiSettings.Scale = 1.0
	previews := s.makeScalePreviews(1.0)

	assert.Equal(t, 5, len(previews))
	assert.Equal(t, theme.TextSize(), previews[2].(*canvas.Text).TextSize)

	s.appliedScale(1.5)
	assert.Equal(t, theme.TextSize()/1.5, previews[2].(*canvas.Text).TextSize)
}
