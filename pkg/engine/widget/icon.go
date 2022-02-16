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
	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/canvas"
	"github.com/bhojpur/gui/pkg/engine/internal/widget"
	"github.com/bhojpur/gui/pkg/engine/theme"
)

type iconRenderer struct {
	widget.BaseRenderer
	raster *canvas.Image

	image *Icon
}

func (i *iconRenderer) MinSize() gui.Size {
	size := theme.IconInlineSize()
	return gui.NewSize(size, size)
}

func (i *iconRenderer) Layout(size gui.Size) {
	if len(i.Objects()) == 0 {
		return
	}

	i.Objects()[0].Resize(size)
}

func (i *iconRenderer) Refresh() {
	if i.image.Resource == i.image.cachedRes {
		return
	}

	i.image.propertyLock.RLock()
	i.raster.Resource = i.image.Resource
	i.image.cachedRes = i.image.Resource
	i.image.propertyLock.RUnlock()

	canvas.Refresh(i.image.super())
}

// Icon widget is a basic image component that load's its resource to match the theme.
type Icon struct {
	BaseWidget

	Resource  gui.Resource // The resource for this icon
	cachedRes gui.Resource
}

// SetResource updates the resource rendered in this icon widget
func (i *Icon) SetResource(res gui.Resource) {
	i.Resource = res
	i.Refresh()
}

// MinSize returns the size that this widget should not shrink below
func (i *Icon) MinSize() gui.Size {
	i.ExtendBaseWidget(i)
	return i.BaseWidget.MinSize()
}

// CreateRenderer is a private method to Bhojpur GUI which links this widget to its renderer
func (i *Icon) CreateRenderer() gui.WidgetRenderer {
	i.ExtendBaseWidget(i)
	i.propertyLock.RLock()
	defer i.propertyLock.RUnlock()

	img := canvas.NewImageFromResource(i.Resource)
	img.FillMode = canvas.ImageFillContain
	r := &iconRenderer{image: i, raster: img}
	r.SetObjects([]gui.CanvasObject{img})
	i.cachedRes = i.Resource
	return r
}

// NewIcon returns a new icon widget that displays a themed icon resource
func NewIcon(res gui.Resource) *Icon {
	icon := &Icon{}
	icon.ExtendBaseWidget(icon)
	icon.SetResource(res) // force the image conversion

	return icon
}
