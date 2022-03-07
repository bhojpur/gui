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

import (
	"fmt"
	"image"
	"log"
	"runtime"

	"github.com/bhojpur/gui/pkg/engine/theme"
	"github.com/goki/freetype/truetype"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/canvas"
	"github.com/bhojpur/gui/pkg/engine/internal/cache"
	"github.com/bhojpur/gui/pkg/engine/internal/painter"
)

// Texture represents an uploaded GL texture
type Texture cache.TextureType

var noTexture = Texture(cache.NoTexture)

func logGLError(err uint32) {
	if err == 0 {
		return
	}

	log.Printf("Error %x in GL Renderer", err)
	_, file, line, ok := runtime.Caller(2)
	if ok {
		log.Printf("  At: %s:%d", file, line)
	}
}

func (p *glPainter) getTexture(object gui.CanvasObject, creator func(canvasObject gui.CanvasObject) Texture) (Texture, error) {
	texture, ok := cache.GetTexture(object)

	if !ok {
		texture = cache.TextureType(creator(object))
		cache.SetTexture(object, texture, p.canvas)
	}
	if !cache.IsValid(texture) {
		return noTexture, fmt.Errorf("no texture available")
	}
	return Texture(texture), nil
}

func (p *glPainter) newGlCircleTexture(obj gui.CanvasObject) Texture {
	circle := obj.(*canvas.Circle)
	raw := painter.DrawCircle(circle, painter.VectorPad(circle), p.textureScale)

	return p.imgToTexture(raw, canvas.ImageScaleSmooth)
}

func (p *glPainter) newGlRectTexture(obj gui.CanvasObject) Texture {
	rect := obj.(*canvas.Rectangle)
	if rect.StrokeColor != nil && rect.StrokeWidth > 0 {
		return p.newGlStrokedRectTexture(rect)
	}
	if rect.FillColor == nil {
		return noTexture
	}
	return p.imgToTexture(image.NewUniform(rect.FillColor), canvas.ImageScaleSmooth)
}

func (p *glPainter) newGlStrokedRectTexture(obj gui.CanvasObject) Texture {
	rect := obj.(*canvas.Rectangle)
	raw := painter.DrawRectangle(rect, painter.VectorPad(rect), p.textureScale)

	return p.imgToTexture(raw, canvas.ImageScaleSmooth)
}

func (p *glPainter) newGlTextTexture(obj gui.CanvasObject) Texture {
	text := obj.(*canvas.Text)
	color := text.Color
	if color == nil {
		color = theme.ForegroundColor()
	}

	bounds := text.MinSize()
	width := int(p.textureScale(bounds.Width))
	height := int(p.textureScale(bounds.Height))
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	var opts truetype.Options
	fontSize := float64(text.TextSize) * float64(p.canvas.Scale())
	opts.Size = fontSize
	opts.DPI = float64(painter.TextDPI * p.texScale)
	face := painter.CachedFontFace(text.TextStyle, &opts)

	painter.DrawString(img, text.Text, color, face, height, text.TextStyle.TabWidth)
	return p.imgToTexture(img, canvas.ImageScaleSmooth)
}

func (p *glPainter) newGlImageTexture(obj gui.CanvasObject) Texture {
	img := obj.(*canvas.Image)

	width := p.textureScale(img.Size().Width)
	height := p.textureScale(img.Size().Height)

	tex := painter.PaintImage(img, p.canvas, int(width), int(height))
	if tex == nil {
		return noTexture
	}

	return p.imgToTexture(tex, img.ScaleMode)
}

func (p *glPainter) newGlRasterTexture(obj gui.CanvasObject) Texture {
	rast := obj.(*canvas.Raster)

	width := p.textureScale(rast.Size().Width)
	height := p.textureScale(rast.Size().Height)

	return p.imgToTexture(rast.Generator(int(width), int(height)), rast.ScaleMode)
}

func (p *glPainter) newGlLinearGradientTexture(obj gui.CanvasObject) Texture {
	gradient := obj.(*canvas.LinearGradient)

	width := p.textureScale(gradient.Size().Width)
	height := p.textureScale(gradient.Size().Height)

	return p.imgToTexture(gradient.Generate(int(width), int(height)), canvas.ImageScaleSmooth)
}

func (p *glPainter) newGlRadialGradientTexture(obj gui.CanvasObject) Texture {
	gradient := obj.(*canvas.RadialGradient)

	width := p.textureScale(gradient.Size().Width)
	height := p.textureScale(gradient.Size().Height)

	return p.imgToTexture(gradient.Generate(int(width), int(height)), canvas.ImageScaleSmooth)
}
