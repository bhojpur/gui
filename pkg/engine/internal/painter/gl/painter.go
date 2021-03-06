package gl

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

// It provides a full Bhojpur GUI render implementation using system OpenGL libraries.

import (
	"image"
	"math"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/internal/driver"
)

// Painter defines the functionality of our OpenGL based renderer
type Painter interface {
	// Init tell a new painter to initialise, usually called after a context is available
	Init()
	// Capture requests that the specified canvas be drawn to an in-memory image
	Capture(gui.Canvas) image.Image
	// Clear tells our painter to prepare a fresh paint
	Clear()
	// Free is used to indicate that a certain canvas object is no longer needed
	Free(gui.CanvasObject)
	// Paint a single gui.CanvasObject but not its children.
	Paint(gui.CanvasObject, gui.Position, gui.Size)
	// SetFrameBufferScale tells us when we have more than 1 framebuffer pixel for each output pixel
	SetFrameBufferScale(float32)
	// SetOutputSize is used to change the resolution of our output viewport
	SetOutputSize(int, int)
	// StartClipping tells us that the following paint actions should be clipped to the specified area.
	StartClipping(gui.Position, gui.Size)
	// StopClipping stops clipping paint actions.
	StopClipping()
}

// Declare conformity to Painter interface
var _ Painter = (*glPainter)(nil)

type glPainter struct {
	canvas      gui.Canvas
	context     driver.WithContext
	program     Program
	lineProgram Program
	texScale    float32
	pixScale    float32 // pre-calculate scale*texScale for each draw
}

func (p *glPainter) SetFrameBufferScale(scale float32) {
	p.texScale = scale
	p.pixScale = p.canvas.Scale() * p.texScale
}

func (p *glPainter) Clear() {
	p.glClearBuffer()
}

func (p *glPainter) StartClipping(pos gui.Position, size gui.Size) {
	x := p.textureScale(pos.X)
	y := p.textureScale(p.canvas.Size().Height - pos.Y - size.Height)
	w := p.textureScale(size.Width)
	h := p.textureScale(size.Height)
	p.glScissorOpen(int32(x), int32(y), int32(w), int32(h))
}

func (p *glPainter) StopClipping() {
	p.glScissorClose()
}

func (p *glPainter) Paint(obj gui.CanvasObject, pos gui.Position, frame gui.Size) {
	if obj.Visible() {
		p.drawObject(obj, pos, frame)
	}
}

func (p *glPainter) Free(obj gui.CanvasObject) {
	p.freeTexture(obj)
}

func (p *glPainter) textureScale(v float32) float32 {
	if p.pixScale == 1.0 {
		return float32(math.Round(float64(v)))
	}

	return float32(math.Round(float64(v * p.pixScale)))
}

// NewPainter creates a new GL based renderer for the provided canvas.
// If it is a master painter it will also initialise OpenGL
func NewPainter(c gui.Canvas, ctx driver.WithContext) Painter {
	p := &glPainter{canvas: c, context: ctx}
	p.SetFrameBufferScale(1.0)

	glInit()

	return p
}
