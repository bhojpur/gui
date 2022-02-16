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
	"testing"

	"github.com/stretchr/testify/assert"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/internal/cache"
)

type extendedSlider struct {
	Slider
}

func newExtendedSlider() *extendedSlider {
	slider := &extendedSlider{}
	slider.ExtendBaseWidget(slider)
	slider.Min = 0
	slider.Max = 10
	return slider
}

func TestSlider_Extended_Value(t *testing.T) {
	slider := newExtendedSlider()
	slider.Resize(slider.MinSize().Add(gui.NewSize(20, 0)))
	objs := cache.Renderer(slider).Objects()
	assert.Equal(t, 3, len(objs))
	thumb := objs[2]
	thumbPos := thumb.Position()

	slider.Value = 2
	slider.Refresh()
	assert.Greater(t, thumb.Position().X, thumbPos.X)
	assert.Equal(t, thumbPos.Y, thumb.Position().Y)
}

func TestSlider_Extended_Drag(t *testing.T) {
	slider := newExtendedSlider()
	objs := cache.Renderer(slider).Objects()
	assert.Equal(t, 3, len(objs))
	thumb := objs[2]
	thumbPos := thumb.Position()

	drag := &gui.DragEvent{Dragged: gui.NewDelta(10, 2)}
	slider.Dragged(drag)
	assert.Greater(t, thumbPos.X, thumb.Position().X)
	assert.Equal(t, thumbPos.Y, thumb.Position().Y)
}
