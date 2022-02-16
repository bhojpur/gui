//go:build !ci && !mobile
// +build !ci,!mobile

package glfw

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
	"sync"
	"testing"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/container"
	"github.com/bhojpur/gui/pkg/engine/internal/cache"
	"github.com/bhojpur/gui/pkg/engine/widget"

	"github.com/stretchr/testify/assert"
)

func Test_gLDriver_AbsolutePositionForObject(t *testing.T) {
	w := createWindow("Test").(*window)

	cr1c1 := widget.NewLabel("row 1 col 1")
	cr1c2 := widget.NewLabel("row 1 col 2")
	cr1c3 := widget.NewLabel("row 1 col 3")
	cr2c1 := widget.NewLabel("row 2 col 1")
	cr2c2 := widget.NewLabel("row 2 col 2")
	cr2c3 := widget.NewLabel("row 2 col 3")
	cr3c1 := widget.NewLabel("row 3 col 1")
	cr3c2 := widget.NewLabel("row 3 col 2")
	cr3c3 := widget.NewLabel("row 3 col 3")
	cr1 := container.NewHBox(cr1c1, cr1c2, cr1c3)
	cr2 := container.NewHBox(cr2c1, cr2c2, cr2c3)
	cr3 := container.NewHBox(cr3c1, cr3c2, cr3c3)
	content := container.NewVBox(cr1, cr2, cr3)
	cr2c2.Hide()

	mm := gui.NewMainMenu(
		gui.NewMenu("Menu 1", gui.NewMenuItem("Menu 1 Item", nil)),
		gui.NewMenu("Menu 2", gui.NewMenuItem("Menu 2 Item", nil)),
	)
	// We want to test the handling of the canvas' Bhojpur GUI menu here.
	// We work around w.SetMainMenu because on MacOS the main menu is a native menu.
	c := w.Canvas().(*glCanvas)
	movl := buildMenuOverlay(mm, w)
	c.Lock()
	c.setMenuOverlay(movl)
	c.Unlock()
	w.SetContent(content)
	w.Resize(gui.NewSize(200, 199))

	ovli1 := widget.NewLabel("Overlay Item 1")
	ovli2 := widget.NewLabel("Overlay Item 2")
	ovli3 := widget.NewLabel("Overlay Item 3")
	ovlContent := container.NewVBox(ovli1, ovli2, ovli3)
	ovl := widget.NewModalPopUp(ovlContent, c)
	ovl.Show()

	repaintWindow(w)
	// accessing the menu bar's actual CanvasObjects isn't straight forward
	// 0 is the shadow
	// 1 is the menu bar’s underlay
	// 2 is the menu bar's background
	// 3 is the container holding the items
	mbarCont := cache.Renderer(movl.(gui.Widget)).Objects()[3].(*gui.Container)
	m2 := mbarCont.Objects[1]

	tests := map[string]struct {
		object       gui.CanvasObject
		wantX, wantY int
	}{
		"a cell": {
			object: cr1c3,
			wantX:  197,
			wantY:  32,
		},
		"a row": {
			object: cr2,
			wantX:  4,
			wantY:  73,
		},
		"the window content": {
			object: content,
			wantX:  4,
			wantY:  32,
		},
		"a hidden element": {
			object: cr2c2,
			wantX:  0,
			wantY:  0,
		},

		"a menu": {
			object: m2,
			wantX:  77,
			wantY:  0,
		},

		"an overlay item": {
			object: ovli2,
			wantX:  87,
			wantY:  81,
		},
		"the overlay content": {
			object: ovlContent,
			wantX:  87,
			wantY:  40,
		},
		"the overlay": {
			object: ovl,
			wantX:  0,
			wantY:  0,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			pos := d.AbsolutePositionForObject(tt.object)
			assert.Equal(t, tt.wantX, int(pos.X))
			assert.Equal(t, tt.wantY, int(pos.Y))
		})
	}
}

var mainRoutineID int

func init() {
	mainRoutineID = goroutineID()
}

func TestGoroutineID(t *testing.T) {
	assert.Equal(t, 1, mainRoutineID)

	var childID1, childID2 int
	testID1 := goroutineID()
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		childID1 = goroutineID()
		wg.Done()
	}()
	go func() {
		childID2 = goroutineID()
		wg.Done()
	}()
	wg.Wait()
	testID2 := goroutineID()

	assert.Equal(t, testID1, testID2)
	assert.Greater(t, childID1, 0)
	assert.NotEqual(t, testID1, childID1)
	assert.Greater(t, childID2, 0)
	assert.NotEqual(t, childID1, childID2)
}
