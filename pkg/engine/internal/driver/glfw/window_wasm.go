//go:build js || wasm || test_web_driver
// +build js wasm test_web_driver

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
	"context"
	_ "image/png" // for the icon
	"runtime"
	"sync"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/driver/desktop"
	"github.com/bhojpur/gui/pkg/engine/internal"
	"github.com/bhojpur/gui/pkg/engine/internal/driver/common"
	"github.com/bhojpur/gui/pkg/engine/internal/painter/gl"
	"github.com/bhojpur/gui/pkg/graphic/glfw"
)

type Cursor struct {
}

// Input modes.
const (
	CursorMode             glfw.InputMode = glfw.CursorMode
	StickyKeysMode         glfw.InputMode = glfw.StickyKeysMode
	StickyMouseButtonsMode glfw.InputMode = glfw.StickyMouseButtonsMode
	LockKeyMods            glfw.InputMode = glfw.LockKeyMods
	RawMouseMotion         glfw.InputMode = glfw.RawMouseMotion
)

// Cursor mode values.
const (
	CursorNormal   int = glfw.CursorNormal
	CursorHidden   int = glfw.CursorHidden
	CursorDisabled int = glfw.CursorDisabled
)

var (
	cursorMap    map[desktop.Cursor]*Cursor
	defaultTitle = "Bhojpur Application"
)

// Declare conformity to Window interface
var _ gui.Window = (*window)(nil)

type window struct {
	common.Window

	viewport   *glfw.Window
	viewLock   sync.RWMutex
	createLock sync.Once
	decorate   bool
	closing    bool
	fixedSize  bool

	cursor   desktop.Cursor
	canvas   *glCanvas
	driver   *gLDriver
	title    string
	icon     gui.Resource
	mainmenu *gui.MainMenu

	clipboard gui.Clipboard

	master     bool
	fullScreen bool
	centered   bool
	visible    bool

	mouseLock            sync.RWMutex
	mousePos             gui.Position
	mouseDragged         gui.Draggable
	mouseDraggedObjStart gui.Position
	mouseDraggedOffset   gui.Position
	mouseDragPos         gui.Position
	mouseDragStarted     bool
	mouseButton          desktop.MouseButton
	mouseOver            desktop.Hoverable
	mouseLastClick       gui.CanvasObject
	mousePressed         gui.CanvasObject
	mouseClickCount      int
	mouseCancelFunc      context.CancelFunc

	onClosed           func()
	onCloseIntercepted func()

	menuTogglePending       gui.KeyName
	menuDeactivationPending gui.KeyName

	xpos, ypos                      int
	width, height                   int
	requestedWidth, requestedHeight int
	shouldWidth, shouldHeight       int
	shouldExpand                    bool

	pending []func()
}

func (w *window) SetFullScreen(full bool) {
	w.fullScreen = true
}

// centerOnScreen handles the logic for centering a window
func (w *window) CenterOnScreen() {
	// FIXME: not supported with WebGL
	w.centered = true
}

func (w *window) doCenterOnScreen() {
	// FIXME: no meaning for defining center on screen in WebGL
}

func (w *window) RequestFocus() {
	// FIXME: no meaning for defining focus in WebGL
}

func (w *window) SetIcon(icon gui.Resource) {
	// FIXME: no support for SetIcon yet
}

func (w *window) SetMaster() {
	// FIXME: there could really only be one window
}

func (w *window) fitContent() {
	w.shouldWidth, w.shouldHeight = w.requestedWidth, w.requestedHeight
}

func (w *window) getMonitorForWindow() *glfw.Monitor {
	return glfw.GetPrimaryMonitor()
}

func scaleForDpi(xdpi int) float32 {
	switch {
	case xdpi > 1000:
		// assume that this is a mistake and bail
		return float32(1.0)
	case xdpi > 192:
		return float32(1.5)
	case xdpi > 144:
		return float32(1.35)
	case xdpi > 120:
		return float32(1.2)
	default:
		return float32(1.0)
	}
}

func (w *window) detectScale() float32 {
	return scaleForDpi(int(96))
}

func (w *window) moved(_ *glfw.Window, x, y int) {
	w.processMoved(x, y)
}

func (w *window) resized(_ *glfw.Window, width, height int) {
	w.canvas.scale = w.calculatedScale()
	w.processResized(width, height)
}

func (w *window) frameSized(_ *glfw.Window, width, height int) {
	w.processFrameSized(width, height)
}

func (w *window) refresh(_ *glfw.Window) {
	w.processRefresh()
}

func (w *window) closed(viewport *glfw.Window) {
	viewport.SetShouldClose(true)

	w.processClosed()
}

func guiToNativeCursor(cursor desktop.Cursor) (*Cursor, bool) {
	return nil, false
}

func (w *window) SetCursor(_ *Cursor) {
}

func (w *window) setCustomCursor(rawCursor *Cursor, isCustomCursor bool) {
}

func (w *window) mouseMoved(_ *glfw.Window, xpos, ypos float64) {
	w.processMouseMoved(xpos, ypos)
}

func (w *window) mouseClicked(viewport *glfw.Window, btn glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	button, modifiers := convertMouseButton(btn, mods)
	mouseAction := convertAction(action)

	w.processMouseClicked(button, mouseAction, modifiers)
}

func (w *window) mouseScrolled(viewport *glfw.Window, xoff, yoff float64) {
	if runtime.GOOS != "darwin" && xoff == 0 &&
		(viewport.GetKey(glfw.KeyLeftShift) == glfw.Press ||
			viewport.GetKey(glfw.KeyRightShift) == glfw.Press) {
		xoff, yoff = yoff, xoff
	}

	w.processMouseScrolled(xoff, yoff)
}

func convertMouseButton(btn glfw.MouseButton, mods glfw.ModifierKey) (desktop.MouseButton, gui.KeyModifier) {
	modifier := desktopModifier(mods)
	var button desktop.MouseButton
	rightClick := false
	if runtime.GOOS == "darwin" {
		if modifier&gui.KeyModifierControl != 0 {
			rightClick = true
			modifier &^= gui.KeyModifierControl
		}
		if modifier&gui.KeyModifierSuper != 0 {
			modifier |= gui.KeyModifierControl
			modifier &^= gui.KeyModifierSuper
		}
	}
	switch btn {
	case glfw.MouseButton1:
		if rightClick {
			button = desktop.RightMouseButton
		} else {
			button = desktop.LeftMouseButton
		}
	case glfw.MouseButton2:
		button = desktop.RightMouseButton
	}
	return button, modifier
}

var keyCodeMap = map[glfw.Key]gui.KeyName{
	// non-printable
	glfw.KeyEscape:    gui.KeyEscape,
	glfw.KeyEnter:     gui.KeyReturn,
	glfw.KeyTab:       gui.KeyTab,
	glfw.KeyBackspace: gui.KeyBackspace,
	glfw.KeyInsert:    gui.KeyInsert,
	glfw.KeyDelete:    gui.KeyDelete,
	glfw.KeyRight:     gui.KeyRight,
	glfw.KeyLeft:      gui.KeyLeft,
	glfw.KeyDown:      gui.KeyDown,
	glfw.KeyUp:        gui.KeyUp,
	glfw.KeyPageUp:    gui.KeyPageUp,
	glfw.KeyPageDown:  gui.KeyPageDown,
	glfw.KeyHome:      gui.KeyHome,
	glfw.KeyEnd:       gui.KeyEnd,

	glfw.KeySpace:   gui.KeySpace,
	glfw.KeyKPEnter: gui.KeyEnter,

	// functions
	glfw.KeyF1:  gui.KeyF1,
	glfw.KeyF2:  gui.KeyF2,
	glfw.KeyF3:  gui.KeyF3,
	glfw.KeyF4:  gui.KeyF4,
	glfw.KeyF5:  gui.KeyF5,
	glfw.KeyF6:  gui.KeyF6,
	glfw.KeyF7:  gui.KeyF7,
	glfw.KeyF8:  gui.KeyF8,
	glfw.KeyF9:  gui.KeyF9,
	glfw.KeyF10: gui.KeyF10,
	glfw.KeyF11: gui.KeyF11,
	glfw.KeyF12: gui.KeyF12,

	// numbers - lookup by code to avoid AZERTY using the symbol name instead of number
	glfw.Key0:   gui.Key0,
	glfw.KeyKP0: gui.Key0,
	glfw.Key1:   gui.Key1,
	glfw.KeyKP1: gui.Key1,
	glfw.Key2:   gui.Key2,
	glfw.KeyKP2: gui.Key2,
	glfw.Key3:   gui.Key3,
	glfw.KeyKP3: gui.Key3,
	glfw.Key4:   gui.Key4,
	glfw.KeyKP4: gui.Key4,
	glfw.Key5:   gui.Key5,
	glfw.KeyKP5: gui.Key5,
	glfw.Key6:   gui.Key6,
	glfw.KeyKP6: gui.Key6,
	glfw.Key7:   gui.Key7,
	glfw.KeyKP7: gui.Key7,
	glfw.Key8:   gui.Key8,
	glfw.KeyKP8: gui.Key8,
	glfw.Key9:   gui.Key9,
	glfw.KeyKP9: gui.Key9,

	// desktop
	glfw.KeyLeftShift:    desktop.KeyShiftLeft,
	glfw.KeyRightShift:   desktop.KeyShiftRight,
	glfw.KeyLeftControl:  desktop.KeyControlLeft,
	glfw.KeyRightControl: desktop.KeyControlRight,
	glfw.KeyLeftAlt:      desktop.KeyAltLeft,
	glfw.KeyRightAlt:     desktop.KeyAltRight,
	glfw.KeyLeftSuper:    desktop.KeySuperLeft,
	glfw.KeyRightSuper:   desktop.KeySuperRight,
	glfw.KeyMenu:         desktop.KeyMenu,
	glfw.KeyPrintScreen:  desktop.KeyPrintScreen,
	glfw.KeyCapsLock:     desktop.KeyCapsLock,
}

var keyCodeMapASCII = map[glfw.Key]gui.KeyName{
	glfw.KeyA: gui.KeyA,
	glfw.KeyB: gui.KeyB,
	glfw.KeyC: gui.KeyC,
	glfw.KeyD: gui.KeyD,
	glfw.KeyE: gui.KeyE,
	glfw.KeyF: gui.KeyF,
	glfw.KeyG: gui.KeyG,
	glfw.KeyH: gui.KeyH,
	glfw.KeyI: gui.KeyI,
	glfw.KeyJ: gui.KeyJ,
	glfw.KeyK: gui.KeyK,
	glfw.KeyL: gui.KeyL,
	glfw.KeyM: gui.KeyM,
	glfw.KeyN: gui.KeyN,
	glfw.KeyO: gui.KeyO,
	glfw.KeyP: gui.KeyP,
	glfw.KeyQ: gui.KeyQ,
	glfw.KeyR: gui.KeyR,
	glfw.KeyS: gui.KeyS,
	glfw.KeyT: gui.KeyT,
	glfw.KeyU: gui.KeyU,
	glfw.KeyV: gui.KeyV,
	glfw.KeyW: gui.KeyW,
	glfw.KeyX: gui.KeyX,
	glfw.KeyY: gui.KeyY,
	glfw.KeyZ: gui.KeyZ,
}

var keyNameMap = map[string]gui.KeyName{
	"'": gui.KeyApostrophe,
	",": gui.KeyComma,
	"-": gui.KeyMinus,
	".": gui.KeyPeriod,
	"/": gui.KeySlash,
	"*": gui.KeyAsterisk,
	"`": gui.KeyBackTick,

	";": gui.KeySemicolon,
	"+": gui.KeyPlus,
	"=": gui.KeyEqual,

	"a": gui.KeyA,
	"b": gui.KeyB,
	"c": gui.KeyC,
	"d": gui.KeyD,
	"e": gui.KeyE,
	"f": gui.KeyF,
	"g": gui.KeyG,
	"h": gui.KeyH,
	"i": gui.KeyI,
	"j": gui.KeyJ,
	"k": gui.KeyK,
	"l": gui.KeyL,
	"m": gui.KeyM,
	"n": gui.KeyN,
	"o": gui.KeyO,
	"p": gui.KeyP,
	"q": gui.KeyQ,
	"r": gui.KeyR,
	"s": gui.KeyS,
	"t": gui.KeyT,
	"u": gui.KeyU,
	"v": gui.KeyV,
	"w": gui.KeyW,
	"x": gui.KeyX,
	"y": gui.KeyY,
	"z": gui.KeyZ,

	"[":  gui.KeyLeftBracket,
	"\\": gui.KeyBackslash,
	"]":  gui.KeyRightBracket,
}

func keyToName(code glfw.Key, scancode int) gui.KeyName {
	if runtime.GOOS == "darwin" && scancode == 0x69 { // TODO remove once fixed upstream glfw/glfw#1786
		code = glfw.KeyPrintScreen
	}

	ret, ok := keyCodeMap[code]
	if ok {
		return ret
	}

	//	keyName := glfw.GetKeyName(code, scancode)
	//	ret, ok = keyNameMap[keyName]
	//	if !ok {
	return gui.KeyUnknown
	//	}

	//	return ret
}

func convertAction(action glfw.Action) action {
	switch action {
	case glfw.Press:
		return press
	case glfw.Release:
		return release
	case glfw.Repeat:
		return repeat
	}
	panic("Could not convert glfw.Action.")
}

func convertASCII(key glfw.Key) gui.KeyName {
	ret, ok := keyCodeMapASCII[key]
	if !ok {
		return gui.KeyUnknown
	}
	return ret
}

func (w *window) keyPressed(viewport *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	keyName := keyToName(key, scancode)
	keyDesktopModifier := desktopModifier(mods)
	keyAction := convertAction(action)
	keyASCII := convertASCII(key)

	w.processKeyPressed(keyName, keyASCII, scancode, keyAction, keyDesktopModifier)
}

func desktopModifier(mods glfw.ModifierKey) gui.KeyModifier {
	var m gui.KeyModifier
	if (mods & glfw.ModShift) != 0 {
		m |= gui.KeyModifierShift
	}
	if (mods & glfw.ModControl) != 0 {
		m |= gui.KeyModifierControl
	}
	if (mods & glfw.ModAlt) != 0 {
		m |= gui.KeyModifierAlt
	}
	if (mods & glfw.ModSuper) != 0 {
		m |= gui.KeyModifierSuper
	}
	return m
}

// charInput defines the character with modifiers callback which is called when a
// Unicode character is input regardless of what modifier keys are used.
//
// Characters do not map 1:1 to physical keys, as a key may produce zero, one or more characters.
func (w *window) charInput(viewport *glfw.Window, char rune) {
	w.processCharInput(char)
}

func (w *window) focused(_ *glfw.Window, focused bool) {
	w.processFocused(focused)
}

func (w *window) DetachCurrentContext() {
	glfw.DetachCurrentContext()
}

func (w *window) rescaleOnMain() {
	if w.viewport == nil {
		return
	}

	//	if w.fullScreen {
	w.width, w.height = w.viewport.GetSize()
	scaledFull := gui.NewSize(
		internal.UnscaleInt(w.canvas, w.width),
		internal.UnscaleInt(w.canvas, w.height))
	w.canvas.Resize(scaledFull)
	return
	//	}

	//	size := w.canvas.size.Union(w.canvas.MinSize())
	//	newWidth, newHeight := w.screenSize(size)
	//	w.viewport.SetSize(newWidth, newHeight)
}

func (w *window) create() {
	runOnMain(func() {
		// we can't hide the window in webgl, so there might be some artifact
		initWindowHints()

		pixWidth, pixHeight := w.screenSize(w.canvas.size)
		pixWidth = int(gui.Max(float32(pixWidth), float32(w.width)))
		if pixWidth == 0 {
			pixWidth = 10
		}
		pixHeight = int(gui.Max(float32(pixHeight), float32(w.height)))
		if pixHeight == 0 {
			pixHeight = 10
		}

		win, err := glfw.CreateWindow(pixWidth, pixHeight, w.title, nil, nil)
		if err != nil {
			w.driver.initFailed("window creation error", err)
			return
		}

		w.viewLock.Lock()
		w.viewport = win
		w.viewLock.Unlock()
	})

	if w.view() == nil { // something went wrong above, it will have been logged
		return
	}

	// run the GL init on the draw thread
	runOnDraw(w, func() {
		w.canvas.SetPainter(gl.NewPainter(w.canvas, w))
		w.canvas.Painter().Init()
	})

	runOnMain(func() {
		w.setDarkMode()

		win := w.view()
		win.SetCloseCallback(w.closed)
		win.SetPosCallback(w.moved)
		win.SetSizeCallback(w.resized)
		win.SetFramebufferSizeCallback(w.frameSized)
		win.SetRefreshCallback(w.refresh)
		win.SetCursorPosCallback(w.mouseMoved)
		win.SetMouseButtonCallback(w.mouseClicked)
		win.SetScrollCallback(w.mouseScrolled)
		win.SetKeyCallback(w.keyPressed)
		win.SetCharCallback(w.charInput)
		win.SetFocusCallback(w.focused)

		w.canvas.detectedScale = w.detectScale()
		w.canvas.scale = w.calculatedScale()
		w.canvas.texScale = w.detectTextureScale()
		// update window size now we have scaled detected
		w.fitContent()

		for _, fn := range w.pending {
			fn()
		}

		w.requestedWidth, w.requestedHeight = w.width, w.height

		width, height := win.GetSize()
		w.processFrameSized(width, height)
		w.processResized(width, height)
	})
}

func (w *window) view() *glfw.Window {
	w.viewLock.RLock()
	defer w.viewLock.RUnlock()

	if w.closing {
		return nil
	}
	return w.viewport
}
