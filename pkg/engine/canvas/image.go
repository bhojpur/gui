package canvas

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
	"image"
	"io"
	"io/ioutil"
	"path/filepath"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/storage"
)

// ImageFill defines the different type of ways an image can stretch to fill its space.
type ImageFill int

const (
	// ImageFillStretch will scale the image to match the Size() values.
	// This is the default and does not maintain aspect ratio.
	ImageFillStretch ImageFill = iota
	// ImageFillContain makes the image fit within the object Size(),
	// centrally and maintaining aspect ratio.
	// There may be transparent sections top and bottom or left and right.
	ImageFillContain //(Fit)
	// ImageFillOriginal ensures that the container grows to the pixel dimensions
	// required to fit the original image. The aspect of the image will be maintained so,
	// as with ImageFillContain there may be transparent areas around the image.
	// Note that the minSize may be smaller than the image dimensions if scale > 1.
	ImageFillOriginal
)

// ImageScale defines the different scaling filters used to scaling images
type ImageScale int32

const (
	// ImageScaleSmooth will scale the image using ApproxBiLinear filter (or GL equivalent)
	ImageScaleSmooth ImageScale = 0
	// ImageScalePixels will scale the image using NearestNeighbor filter (or GL equivalent)
	ImageScalePixels ImageScale = 1
	// ImageScaleFastest will scale the image using hardware GPU if available
	//
	// Since: 2.0
	ImageScaleFastest ImageScale = 2
)

// Declare conformity with CanvasObject interface
var _ gui.CanvasObject = (*Image)(nil)

// Image describes a drawable image area that can render in a Bhojpur GUI canvas
// The image may be a vector or a bitmap representation and it will fill the area.
// The fill mode can be changed by setting FillMode to a different ImageFill.
type Image struct {
	baseObject

	// one of the following sources will provide our image data
	File     string       // Load the image from a file
	Resource gui.Resource // Load the image from an in-memory resource
	Image    image.Image  // Specify a loaded image to use in this canvas object

	Translucency float64    // Set a translucency value > 0.0 to fade the image
	FillMode     ImageFill  // Specify how the image should expand to fill or fit the available space
	ScaleMode    ImageScale // Specify the type of scaling interpolation applied to the image

}

// Alpha is a convenience function that returns the alpha value for an image
// based on its Translucency value. The result is 1.0 - Translucency.
func (i *Image) Alpha() float64 {
	return 1.0 - i.Translucency
}

// Resize on an image will scale the content or reposition it according to FillMode.
// It will normally cause a Refresh to ensure the pixels are recalculated.
func (i *Image) Resize(s gui.Size) {
	if s == i.Size() {
		return
	}
	if i.FillMode == ImageFillOriginal && i.size.Height > 2 { // don't refresh original scale images after first draw
		return
	}

	i.baseObject.Resize(s)

	Refresh(i)
}

// Refresh causes this object to be redrawn in it's current state
func (i *Image) Refresh() {
	Refresh(i)
}

// NewImageFromFile creates a new image from a local file.
// Images returned from this method will scale to fit the canvas object.
// The method for scaling can be set using the Fill field.
func NewImageFromFile(file string) *Image {
	return &Image{
		File: file,
	}
}

// NewImageFromURI creates a new image from named resource.
// File URIs will read the file path and other schemes will download the data into a resource.
// HTTP and HTTPs URIs will use the GET method by default to request the resource.
// Images returned from this method will scale to fit the canvas object.
// The method for scaling can be set using the Fill field.
//
// Since: 2.0
func NewImageFromURI(uri gui.URI) *Image {
	if uri.Scheme() == "file" && len(uri.String()) > 7 {
		return &Image{
			File: uri.String()[7:],
		}
	}

	var read io.ReadCloser

	read, err := storage.Reader(uri) // attempt unknown / http file type
	if err != nil {
		gui.LogError("Failed to open image URI", err)
		return &Image{}
	}

	defer read.Close()
	return NewImageFromReader(read, filepath.Base(uri.String()))
}

// NewImageFromReader creates a new image from a data stream.
// The name parameter is required to uniquely identify this image (for caching etc).
// If the image in this io.Reader is an SVG, the name should end ".svg".
// Images returned from this method will scale to fit the canvas object.
// The method for scaling can be set using the Fill field.
//
// Since: 2.0
func NewImageFromReader(read io.Reader, name string) *Image {
	data, err := ioutil.ReadAll(read)
	if err != nil {
		gui.LogError("Unable to read image data", err)
		return nil
	}
	res := &gui.StaticResource{
		StaticName:    name,
		StaticContent: data,
	}

	return &Image{
		Resource: res,
	}
}

// NewImageFromResource creates a new image by loading the specified resource.
// Images returned from this method will scale to fit the canvas object.
// The method for scaling can be set using the Fill field.
func NewImageFromResource(res gui.Resource) *Image {
	return &Image{
		Resource: res,
	}
}

// NewImageFromImage returns a new Image instance that is rendered from the Go
// image.Image passed in.
// Images returned from this method will scale to fit the canvas object.
// The method for scaling can be set using the Fill field.
func NewImageFromImage(img image.Image) *Image {
	return &Image{
		Image: img,
	}
}
