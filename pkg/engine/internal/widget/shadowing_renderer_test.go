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
	w "github.com/bhojpur/gui/pkg/engine/internal/widget"
	"github.com/bhojpur/gui/pkg/engine/widget"
)

func TestShadowingRenderer_Objects(t *testing.T) {
	tests := map[string]struct {
		level                w.ElevationLevel
		wantPrependedObjects []gui.CanvasObject
	}{
		"with shadow": {
			12,
			[]gui.CanvasObject{w.NewShadow(w.ShadowAround, 12)},
		},
		"without shadow": {
			0,
			[]gui.CanvasObject{},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			objects := []gui.CanvasObject{widget.NewLabel("A"), widget.NewLabel("B")}
			r := w.NewShadowingRenderer(objects, tt.level)
			assert.Equal(t, append(tt.wantPrependedObjects, objects...), r.Objects())

			otherObjects := []gui.CanvasObject{widget.NewLabel("X"), widget.NewLabel("Y")}
			r.SetObjects(otherObjects)
			assert.Equal(t, append(tt.wantPrependedObjects, otherObjects...), r.Objects())
		})
	}
}
