package test

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
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/driver/desktop"
	"github.com/bhojpur/gui/pkg/engine/internal/cache"
	"github.com/bhojpur/gui/pkg/engine/internal/driver"
	"github.com/bhojpur/gui/pkg/engine/internal/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// AssertCanvasTappableAt asserts that the canvas is tappable at the given position.
func AssertCanvasTappableAt(t *testing.T, c gui.Canvas, pos gui.Position) bool {
	if o, _ := findTappable(c, pos); o == nil {
		t.Errorf("No tappable found at %#v", pos)
		return false
	}
	return true
}

// AssertImageMatches asserts that the given image is the same as the one stored in the master file.
// The master filename is relative to the `testdata` directory which is relative to the test.
// The test `t` fails if the given image is not equal to the loaded master image.
// In this case the given image is written into a file in `testdata/failed/<masterFilename>` (relative to the test).
// This path is also reported, thus the file can be used as new master.
func AssertImageMatches(t *testing.T, masterFilename string, img image.Image, msgAndArgs ...interface{}) bool {
	return test.AssertImageMatches(t, masterFilename, img, msgAndArgs...)
}

// AssertRendersToMarkup asserts that the given canvas renders the same markup as the one stored in the master file.
// The master filename is relative to the `testdata` directory which is relative to the test.
// The test `t` fails if the rendered markup is not equal to the loaded master markup.
// In this case the rendered markup is written into a file in `testdata/failed/<masterFilename>` (relative to the test).
// This path is also reported, thus the file can be used as new master.
//
// Be aware, that the indentation has to use tab characters ('\t') instead of spaces.
// Every element starts on a new line indented one more than its parent.
// Closing elements stand on their own line, too, using the same indentation as the opening element.
// The only exception to this are text elements which do not contain line breaks unless the text includes them.
//
// Since: 2.0
func AssertRendersToMarkup(t *testing.T, masterFilename string, c gui.Canvas, msgAndArgs ...interface{}) bool {
	wd, err := os.Getwd()
	require.NoError(t, err)

	got := snapshot(c)
	masterPath := filepath.Join(wd, "testdata", masterFilename)
	failedPath := filepath.Join(wd, "testdata/failed", masterFilename)
	_, err = os.Stat(masterPath)
	if os.IsNotExist(err) {
		require.NoError(t, writeMarkup(failedPath, got))
		t.Errorf("Master not found at %s. Markup written to %s might be used as master.", masterPath, failedPath)
		return false
	}

	raw, err := ioutil.ReadFile(masterPath)
	require.NoError(t, err)
	master := strings.ReplaceAll(string(raw), "\r", "")

	var msg string
	if len(msgAndArgs) > 0 {
		msg = fmt.Sprintf(msgAndArgs[0].(string)+"\n", msgAndArgs[1:]...)
	}
	if !assert.Equal(t, master, got, "%sMarkup did not match master. Actual markup written to file://%s.", msg, failedPath) {
		require.NoError(t, writeMarkup(failedPath, got))
		return false
	}
	return true
}

// Drag drags at an absolute position on the canvas.
// deltaX/Y is the dragging distance: <0 for dragging up/left, >0 for dragging down/right.
func Drag(c gui.Canvas, pos gui.Position, deltaX, deltaY float32) {
	matches := func(object gui.CanvasObject) bool {
		if _, ok := object.(gui.Draggable); ok {
			return true
		}
		return false
	}
	o, p, _ := driver.FindObjectAtPositionMatching(pos, matches, c.Overlays().Top(), c.Content())
	if o == nil {
		return
	}
	e := &gui.DragEvent{
		PointEvent: gui.PointEvent{Position: p},
		Dragged:    gui.Delta{DX: deltaX, DY: deltaY},
	}
	o.(gui.Draggable).Dragged(e)
	o.(gui.Draggable).DragEnd()
}

// FocusNext focuses the next focusable on the canvas.
func FocusNext(c gui.Canvas) {
	if tc, ok := c.(*testCanvas); ok {
		tc.focusManager().FocusNext()
	} else {
		gui.LogError("FocusNext can only be called with a test canvas", nil)
	}
}

// FocusPrevious focuses the previous focusable on the canvas.
func FocusPrevious(c gui.Canvas) {
	if tc, ok := c.(*testCanvas); ok {
		tc.focusManager().FocusPrevious()
	} else {
		gui.LogError("FocusPrevious can only be called with a test canvas", nil)
	}
}

// LaidOutObjects returns all gui.CanvasObject starting at the given gui.CanvasObject which is laid out previously.
func LaidOutObjects(o gui.CanvasObject) (objects []gui.CanvasObject) {
	if o != nil {
		objects = layoutAndCollect(objects, o, o.MinSize().Max(o.Size()))
	}
	return objects
}

// MoveMouse simulates a mouse movement to the given position.
func MoveMouse(c gui.Canvas, pos gui.Position) {
	if gui.CurrentDevice().IsMobile() {
		return
	}

	tc, _ := c.(*testCanvas)
	var oldHovered, hovered desktop.Hoverable
	if tc != nil {
		oldHovered = tc.hovered
	}
	matches := func(object gui.CanvasObject) bool {
		if _, ok := object.(desktop.Hoverable); ok {
			return true
		}
		return false
	}
	o, p, _ := driver.FindObjectAtPositionMatching(pos, matches, c.Overlays().Top(), c.Content())
	if o != nil {
		hovered = o.(desktop.Hoverable)
		me := &desktop.MouseEvent{
			PointEvent: gui.PointEvent{
				AbsolutePosition: pos,
				Position:         p,
			},
		}
		if hovered == oldHovered {
			hovered.MouseMoved(me)
		} else {
			if oldHovered != nil {
				oldHovered.MouseOut()
			}
			hovered.MouseIn(me)
		}
	} else if oldHovered != nil {
		oldHovered.MouseOut()
	}
	if tc != nil {
		tc.hovered = hovered
	}
}

// Scroll scrolls at an absolute position on the canvas.
// deltaX/Y is the scrolling distance: <0 for scrolling up/left, >0 for scrolling down/right.
func Scroll(c gui.Canvas, pos gui.Position, deltaX, deltaY float32) {
	matches := func(object gui.CanvasObject) bool {
		if _, ok := object.(gui.Scrollable); ok {
			return true
		}
		return false
	}
	o, _, _ := driver.FindObjectAtPositionMatching(pos, matches, c.Overlays().Top(), c.Content())
	if o == nil {
		return
	}

	e := &gui.ScrollEvent{Scrolled: gui.Delta{DX: deltaX, DY: deltaY}}
	o.(gui.Scrollable).Scrolled(e)
}

// DoubleTap simulates a double left mouse click on the specified object.
func DoubleTap(obj gui.DoubleTappable) {
	ev, c := prepareTap(obj, gui.NewPos(1, 1))
	handleFocusOnTap(c, obj)
	obj.DoubleTapped(ev)
}

// Tap simulates a left mouse click on the specified object.
func Tap(obj gui.Tappable) {
	TapAt(obj, gui.NewPos(1, 1))
}

// TapAt simulates a left mouse click on the passed object at a specified place within it.
func TapAt(obj gui.Tappable, pos gui.Position) {
	ev, c := prepareTap(obj, pos)
	tap(c, obj, ev)
}

// TapCanvas taps at an absolute position on the canvas.
func TapCanvas(c gui.Canvas, pos gui.Position) {
	if o, p := findTappable(c, pos); o != nil {
		tap(c, o.(gui.Tappable), &gui.PointEvent{AbsolutePosition: pos, Position: p})
	}
}

// TapSecondary simulates a right mouse click on the specified object.
func TapSecondary(obj gui.SecondaryTappable) {
	TapSecondaryAt(obj, gui.NewPos(1, 1))
}

// TapSecondaryAt simulates a right mouse click on the passed object at a specified place within it.
func TapSecondaryAt(obj gui.SecondaryTappable, pos gui.Position) {
	ev, c := prepareTap(obj, pos)
	handleFocusOnTap(c, obj)
	obj.TappedSecondary(ev)
}

// Type performs a series of key events to simulate typing of a value into the specified object.
// The focusable object will be focused before typing begins.
// The chars parameter will be input one rune at a time to the focused object.
func Type(obj gui.Focusable, chars string) {
	obj.FocusGained()

	typeChars([]rune(chars), obj.TypedRune)
}

// TypeOnCanvas is like the Type function but it passes the key events to the canvas object
// rather than a focusable widget.
func TypeOnCanvas(c gui.Canvas, chars string) {
	typeChars([]rune(chars), c.OnTypedRune())
}

// ApplyTheme sets the given theme and waits for it to be applied to the current app.
func ApplyTheme(t *testing.T, theme gui.Theme) {
	require.IsType(t, &testApp{}, gui.CurrentApp())
	a := gui.CurrentApp().(*testApp)
	a.Settings().SetTheme(theme)
	for a.lastAppliedTheme() != theme {
		time.Sleep(1 * time.Millisecond)
	}
}

// WidgetRenderer allows test scripts to gain access to the current renderer for a widget.
// This can be used for verifying correctness of rendered components for a widget in unit tests.
func WidgetRenderer(wid gui.Widget) gui.WidgetRenderer {
	return cache.Renderer(wid)
}

// WithTestTheme runs a function with the testTheme temporarily set.
func WithTestTheme(t *testing.T, f func()) {
	settings := gui.CurrentApp().Settings()
	current := settings.Theme()
	ApplyTheme(t, NewTheme())
	defer ApplyTheme(t, current)
	f()
}

func findTappable(c gui.Canvas, pos gui.Position) (o gui.CanvasObject, p gui.Position) {
	matches := func(object gui.CanvasObject) bool {
		_, ok := object.(gui.Tappable)
		return ok
	}
	o, p, _ = driver.FindObjectAtPositionMatching(pos, matches, c.Overlays().Top(), c.Content())
	return
}

func prepareTap(obj interface{}, pos gui.Position) (*gui.PointEvent, gui.Canvas) {
	d := gui.CurrentApp().Driver()
	ev := &gui.PointEvent{Position: pos}
	var c gui.Canvas
	if co, ok := obj.(gui.CanvasObject); ok {
		c = d.CanvasForObject(co)
		ev.AbsolutePosition = d.AbsolutePositionForObject(co).Add(pos)
	}
	return ev, c
}

func tap(c gui.Canvas, obj gui.Tappable, ev *gui.PointEvent) {
	handleFocusOnTap(c, obj)
	obj.Tapped(ev)
}

func handleFocusOnTap(c gui.Canvas, obj interface{}) {
	if c == nil {
		return
	}
	unfocus := true
	if focus, ok := obj.(gui.Focusable); ok {
		if dis, ok := obj.(gui.Disableable); !ok || !dis.Disabled() {
			unfocus = false
			if focus != c.Focused() {
				unfocus = true
			}
		}
	}
	if unfocus {
		c.Unfocus()
	}
}

func typeChars(chars []rune, keyDown func(rune)) {
	for _, char := range chars {
		keyDown(char)
	}
}

func writeMarkup(path string, markup string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return ioutil.WriteFile(path, []byte(markup), 0644)
}
