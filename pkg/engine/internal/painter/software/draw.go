package software

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
	"image"
	"math"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/canvas"
	"github.com/bhojpur/gui/pkg/engine/internal"
	"github.com/bhojpur/gui/pkg/engine/internal/painter"
	"github.com/bhojpur/gui/pkg/engine/theme"

	"github.com/goki/freetype/truetype"
	"golang.org/x/image/draw"
)

type gradient interface {
	Generate(int, int) image.Image
	Size() gui.Size
}

func drawCircle(c gui.Canvas, circle *canvas.Circle, pos gui.Position, base *image.NRGBA, clip image.Rectangle) {
	pad := painter.VectorPad(circle)
	scaledWidth := internal.ScaleInt(c, circle.Size().Width+pad*2)
	scaledHeight := internal.ScaleInt(c, circle.Size().Height+pad*2)
	scaledX, scaledY := internal.ScaleInt(c, pos.X-pad), internal.ScaleInt(c, pos.Y-pad)
	bounds := clip.Intersect(image.Rect(scaledX, scaledY, scaledX+scaledWidth, scaledY+scaledHeight))

	raw := painter.DrawCircle(circle, pad, func(in float32) float32 {
		return float32(math.Round(float64(in) * float64(c.Scale())))
	})

	// the clip intersect above cannot be negative, so we may need to compensate
	offX, offY := 0, 0
	if scaledX < 0 {
		offX = -scaledX
	}
	if scaledY < 0 {
		offY = -scaledY
	}
	draw.Draw(base, bounds, raw, image.Point{offX, offY}, draw.Over)
}

func drawGradient(c gui.Canvas, g gradient, pos gui.Position, base *image.NRGBA, clip image.Rectangle) {
	bounds := g.Size()
	width := internal.ScaleInt(c, bounds.Width)
	height := internal.ScaleInt(c, bounds.Height)
	tex := g.Generate(width, height)
	drawTex(internal.ScaleInt(c, pos.X), internal.ScaleInt(c, pos.Y), width, height, base, tex, clip)
}

func drawImage(c gui.Canvas, img *canvas.Image, pos gui.Position, base *image.NRGBA, clip image.Rectangle) {
	bounds := img.Size()
	if bounds.IsZero() {
		return
	}
	width := internal.ScaleInt(c, bounds.Width)
	height := internal.ScaleInt(c, bounds.Height)
	scaledX, scaledY := internal.ScaleInt(c, pos.X), internal.ScaleInt(c, pos.Y)

	origImg := painter.PaintImage(img, c, width, height)

	if img.FillMode == canvas.ImageFillContain {
		imgAspect := painter.GetAspect(img)
		objAspect := float32(width) / float32(height)

		if objAspect > imgAspect {
			newWidth := int(float32(height) * imgAspect)
			scaledX += (width - newWidth) / 2
			width = newWidth
		} else if objAspect < imgAspect {
			newHeight := int(float32(width) / imgAspect)
			scaledY += (height - newHeight) / 2
			height = newHeight
		}
	}

	drawPixels(scaledX, scaledY, width, height, img.ScaleMode, base, origImg, clip)
}

func drawPixels(x, y, width, height int, mode canvas.ImageScale, base *image.NRGBA, origImg image.Image, clip image.Rectangle) {
	if origImg.Bounds().Dx() == width && origImg.Bounds().Dy() == height {
		// do not scale or duplicate image since not needed, draw directly
		drawTex(x, y, width, height, base, origImg, clip)
		return
	}

	scaledBounds := image.Rect(0, 0, width, height)
	scaledImg := image.NewNRGBA(scaledBounds)
	switch mode {
	case canvas.ImageScalePixels:
		draw.NearestNeighbor.Scale(scaledImg, scaledBounds, origImg, origImg.Bounds(), draw.Over, nil)
	case canvas.ImageScaleFastest:
		draw.ApproxBiLinear.Scale(scaledImg, scaledBounds, origImg, origImg.Bounds(), draw.Over, nil)
	default:
		if mode != canvas.ImageScaleSmooth {
			gui.LogError(fmt.Sprintf("Invalid canvas.ImageScale value (%d), using canvas.ImageScaleSmooth as default value", mode), nil)
		}
		draw.CatmullRom.Scale(scaledImg, scaledBounds, origImg, origImg.Bounds(), draw.Over, nil)
	}

	drawTex(x, y, width, height, base, scaledImg, clip)
}

func drawLine(c gui.Canvas, line *canvas.Line, pos gui.Position, base *image.NRGBA, clip image.Rectangle) {
	pad := painter.VectorPad(line)
	scaledWidth := internal.ScaleInt(c, line.Size().Width+pad*2)
	scaledHeight := internal.ScaleInt(c, line.Size().Height+pad*2)
	scaledX, scaledY := internal.ScaleInt(c, pos.X-pad), internal.ScaleInt(c, pos.Y-pad)
	bounds := clip.Intersect(image.Rect(scaledX, scaledY, scaledX+scaledWidth, scaledY+scaledHeight))

	raw := painter.DrawLine(line, pad, func(in float32) float32 {
		return float32(math.Round(float64(in) * float64(c.Scale())))
	})

	// the clip intersect above cannot be negative, so we may need to compensate
	offX, offY := 0, 0
	if scaledX < 0 {
		offX = -scaledX
	}
	if scaledY < 0 {
		offY = -scaledY
	}
	draw.Draw(base, bounds, raw, image.Point{offX, offY}, draw.Over)
}

func drawTex(x, y, width, height int, base *image.NRGBA, tex image.Image, clip image.Rectangle) {
	outBounds := image.Rect(x, y, x+width, y+height)
	clippedBounds := clip.Intersect(outBounds)
	srcPt := image.Point{X: clippedBounds.Min.X - outBounds.Min.X, Y: clippedBounds.Min.Y - outBounds.Min.Y}
	draw.Draw(base, clippedBounds, tex, srcPt, draw.Over)
}

func drawText(c gui.Canvas, text *canvas.Text, pos gui.Position, base *image.NRGBA, clip image.Rectangle) {
	bounds := text.MinSize()
	width := internal.ScaleInt(c, bounds.Width)
	height := internal.ScaleInt(c, bounds.Height)
	txtImg := image.NewRGBA(image.Rect(0, 0, width, height))

	color := text.Color
	if color == nil {
		color = theme.ForegroundColor()
	}

	var opts truetype.Options
	fontSize := text.TextSize * c.Scale()
	opts.Size = float64(fontSize)
	opts.DPI = painter.TextDPI
	face := painter.CachedFontFace(text.TextStyle, &opts)

	painter.DrawString(txtImg, text.Text, color, face, height, text.TextStyle.TabWidth)

	size := text.Size()
	offsetX := float32(0)
	offsetY := float32(0)
	switch text.Alignment {
	case gui.TextAlignTrailing:
		offsetX = size.Width - bounds.Width
	case gui.TextAlignCenter:
		offsetX = (size.Width - bounds.Width) / 2
	}
	if size.Height > bounds.Height {
		offsetY = (size.Height - bounds.Height) / 2
	}
	scaledX := internal.ScaleInt(c, pos.X+offsetX)
	scaledY := internal.ScaleInt(c, pos.Y+offsetY)
	imgBounds := image.Rect(scaledX, scaledY, scaledX+width, scaledY+height)
	clippedBounds := clip.Intersect(imgBounds)
	srcPt := image.Point{X: clippedBounds.Min.X - imgBounds.Min.X, Y: clippedBounds.Min.Y - imgBounds.Min.Y}
	draw.Draw(base, clippedBounds, txtImg, srcPt, draw.Over)
}

func drawRaster(c gui.Canvas, rast *canvas.Raster, pos gui.Position, base *image.NRGBA, clip image.Rectangle) {
	bounds := rast.Size()
	if bounds.IsZero() {
		return
	}
	width := internal.ScaleInt(c, bounds.Width)
	height := internal.ScaleInt(c, bounds.Height)
	scaledX, scaledY := internal.ScaleInt(c, pos.X), internal.ScaleInt(c, pos.Y)

	pix := rast.Generator(width, height)
	if pix.Bounds().Bounds().Dx() != width || pix.Bounds().Dy() != height {
		drawPixels(scaledX, scaledY, width, height, rast.ScaleMode, base, pix, clip)
	} else {
		drawTex(scaledX, scaledY, width, height, base, pix, clip)
	}
}

func drawRectangleStroke(c gui.Canvas, rect *canvas.Rectangle, pos gui.Position, base *image.NRGBA, clip image.Rectangle) {
	pad := painter.VectorPad(rect)
	scaledWidth := internal.ScaleInt(c, rect.Size().Width+pad*2)
	scaledHeight := internal.ScaleInt(c, rect.Size().Height+pad*2)
	scaledX, scaledY := internal.ScaleInt(c, pos.X-pad), internal.ScaleInt(c, pos.Y-pad)
	bounds := clip.Intersect(image.Rect(scaledX, scaledY, scaledX+scaledWidth, scaledY+scaledHeight))

	raw := painter.DrawRectangle(rect, pad, func(in float32) float32 {
		return float32(math.Round(float64(in) * float64(c.Scale())))
	})

	// the clip intersect above cannot be negative, so we may need to compensate
	offX, offY := 0, 0
	if scaledX < 0 {
		offX = -scaledX
	}
	if scaledY < 0 {
		offY = -scaledY
	}
	draw.Draw(base, bounds, raw, image.Point{offX, offY}, draw.Over)
}

func drawRectangle(c gui.Canvas, rect *canvas.Rectangle, pos gui.Position, base *image.NRGBA, clip image.Rectangle) {
	if rect.StrokeColor != nil && rect.StrokeWidth > 0 { // use a rasterizer if there is a stroke
		drawRectangleStroke(c, rect, pos, base, clip)
		return
	}

	scaledWidth := internal.ScaleInt(c, rect.Size().Width)
	scaledHeight := internal.ScaleInt(c, rect.Size().Height)
	scaledX, scaledY := internal.ScaleInt(c, pos.X), internal.ScaleInt(c, pos.Y)
	bounds := clip.Intersect(image.Rect(scaledX, scaledY, scaledX+scaledWidth, scaledY+scaledHeight))
	draw.Draw(base, bounds, image.NewUniform(rect.FillColor), image.Point{}, draw.Over)
}
