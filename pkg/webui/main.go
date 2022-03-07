package main

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

// It uses the Bhojpur GUI - Wasm Canvas library, and avoids context calls for drawing.

import (
	"encoding/binary"
	"fmt"
	"log"
	"time"

	"github.com/bhojpur/gui/pkg/graphic/gl"
	"github.com/bhojpur/gui/pkg/graphic/glfw"
	"github.com/bhojpur/gui/pkg/graphic/glutil"
	"github.com/bhojpur/gui/pkg/webui/utils"
	"github.com/go-gl/mathgl/mgl32"
	"golang.org/x/mobile/exp/f32"
)

var counter int = -1

// getCounter returns event index.
func getCounter() int {
	counter++
	return counter
}

// Window -> Id.
var windowIds = make(map[*glfw.Window]int)

func getWindowId(w *glfw.Window) int {
	return windowIds[w]
}

var startedProcess = time.Now()

// getTime returns time in seconds since process was started.
func getTime() float64 {
	return time.Since(startedProcess).Seconds()
}

func keyString(key glfw.Key) string {
	switch key {
	// Printable keys.
	case glfw.KeyA:
		return "A"
	case glfw.KeyB:
		return "B"
	case glfw.KeyC:
		return "C"
	case glfw.KeyD:
		return "D"
	case glfw.KeyE:
		return "E"
	case glfw.KeyF:
		return "F"
	case glfw.KeyG:
		return "G"
	case glfw.KeyH:
		return "H"
	case glfw.KeyI:
		return "I"
	case glfw.KeyJ:
		return "J"
	case glfw.KeyK:
		return "K"
	case glfw.KeyL:
		return "L"
	case glfw.KeyM:
		return "M"
	case glfw.KeyN:
		return "N"
	case glfw.KeyO:
		return "O"
	case glfw.KeyP:
		return "P"
	case glfw.KeyQ:
		return "Q"
	case glfw.KeyR:
		return "R"
	case glfw.KeyS:
		return "S"
	case glfw.KeyT:
		return "T"
	case glfw.KeyU:
		return "U"
	case glfw.KeyV:
		return "V"
	case glfw.KeyW:
		return "W"
	case glfw.KeyX:
		return "X"
	case glfw.KeyY:
		return "Y"
	case glfw.KeyZ:
		return "Z"
	case glfw.Key1:
		return "1"
	case glfw.Key2:
		return "2"
	case glfw.Key3:
		return "3"
	case glfw.Key4:
		return "4"
	case glfw.Key5:
		return "5"
	case glfw.Key6:
		return "6"
	case glfw.Key7:
		return "7"
	case glfw.Key8:
		return "8"
	case glfw.Key9:
		return "9"
	case glfw.Key0:
		return "0"
	case glfw.KeySpace:
		return "SPACE"
	case glfw.KeyMinus:
		return "MINUS"
	case glfw.KeyEqual:
		return "EQUAL"
	case glfw.KeyLeftBracket:
		return "LEFT BRACKET"
	case glfw.KeyRightBracket:
		return "RIGHT BRACKET"
	case glfw.KeyBackslash:
		return "BACKSLASH"
	case glfw.KeySemicolon:
		return "SEMICOLON"
	case glfw.KeyApostrophe:
		return "APOSTROPHE"
	case glfw.KeyGraveAccent:
		return "GRAVE ACCENT"
	case glfw.KeyComma:
		return "COMMA"
	case glfw.KeyPeriod:
		return "PERIOD"
	case glfw.KeySlash:
		return "SLASH"
	case glfw.KeyWorld1:
		return "WORLD 1"
	case glfw.KeyWorld2:
		return "WORLD 2"
	// Function keys.
	case glfw.KeyEscape:
		return "ESCAPE"
	case glfw.KeyF1:
		return "F1"
	case glfw.KeyF2:
		return "F2"
	case glfw.KeyF3:
		return "F3"
	case glfw.KeyF4:
		return "F4"
	case glfw.KeyF5:
		return "F5"
	case glfw.KeyF6:
		return "F6"
	case glfw.KeyF7:
		return "F7"
	case glfw.KeyF8:
		return "F8"
	case glfw.KeyF9:
		return "F9"
	case glfw.KeyF10:
		return "F10"
	case glfw.KeyF11:
		return "F11"
	case glfw.KeyF12:
		return "F12"
	case glfw.KeyF13:
		return "F13"
	case glfw.KeyF14:
		return "F14"
	case glfw.KeyF15:
		return "F15"
	case glfw.KeyF16:
		return "F16"
	case glfw.KeyF17:
		return "F17"
	case glfw.KeyF18:
		return "F18"
	case glfw.KeyF19:
		return "F19"
	case glfw.KeyF20:
		return "F20"
	case glfw.KeyF21:
		return "F21"
	case glfw.KeyF22:
		return "F22"
	case glfw.KeyF23:
		return "F23"
	case glfw.KeyF24:
		return "F24"
	case glfw.KeyF25:
		return "F25"
	case glfw.KeyUp:
		return "UP"
	case glfw.KeyDown:
		return "DOWN"
	case glfw.KeyLeft:
		return "LEFT"
	case glfw.KeyRight:
		return "RIGHT"
	case glfw.KeyLeftShift:
		return "LEFT SHIFT"
	case glfw.KeyRightShift:
		return "RIGHT SHIFT"
	case glfw.KeyLeftControl:
		return "LEFT CONTROL"
	case glfw.KeyRightControl:
		return "RIGHT CONTROL"
	case glfw.KeyLeftAlt:
		return "LEFT ALT"
	case glfw.KeyRightAlt:
		return "RIGHT ALT"
	case glfw.KeyTab:
		return "TAB"
	case glfw.KeyEnter:
		return "ENTER"
	case glfw.KeyBackspace:
		return "BACKSPACE"
	case glfw.KeyInsert:
		return "INSERT"
	case glfw.KeyDelete:
		return "DELETE"
	case glfw.KeyPageUp:
		return "PAGE UP"
	case glfw.KeyPageDown:
		return "PAGE DOWN"
	case glfw.KeyHome:
		return "HOME"
	case glfw.KeyEnd:
		return "END"
	case glfw.KeyKP0:
		return "KEYPAD 0"
	case glfw.KeyKP1:
		return "KEYPAD 1"
	case glfw.KeyKP2:
		return "KEYPAD 2"
	case glfw.KeyKP3:
		return "KEYPAD 3"
	case glfw.KeyKP4:
		return "KEYPAD 4"
	case glfw.KeyKP5:
		return "KEYPAD 5"
	case glfw.KeyKP6:
		return "KEYPAD 6"
	case glfw.KeyKP7:
		return "KEYPAD 7"
	case glfw.KeyKP8:
		return "KEYPAD 8"
	case glfw.KeyKP9:
		return "KEYPAD 9"
	case glfw.KeyKPDivide:
		return "KEYPAD DIVIDE"
	case glfw.KeyKPMultiply:
		return "KEYPAD MULTPLY"
	case glfw.KeyKPSubtract:
		return "KEYPAD SUBTRACT"
	case glfw.KeyKPAdd:
		return "KEYPAD ADD"
	case glfw.KeyKPDecimal:
		return "KEYPAD DECIMAL"
	case glfw.KeyKPEqual:
		return "KEYPAD EQUAL"
	case glfw.KeyKPEnter:
		return "KEYPAD ENTER"
	case glfw.KeyPrintScreen:
		return "PRINT SCREEN"
	case glfw.KeyNumLock:
		return "NUM LOCK"
	case glfw.KeyCapsLock:
		return "CAPS LOCK"
	case glfw.KeyScrollLock:
		return "SCROLL LOCK"
	case glfw.KeyPause:
		return "PAUSE"
	case glfw.KeyLeftSuper:
		return "LEFT SUPER"
	case glfw.KeyRightSuper:
		return "RIGHT SUPER"
	case glfw.KeyMenu:
		return "MENU"
	default:
		return "UNKNOWN"
	}
}

func actionString(action glfw.Action) string {
	switch action {
	case glfw.Press:
		return "pressed"
	case glfw.Release:
		return "released"
	case glfw.Repeat:
		return "repeated"
	default:
		return "caused unknown action"
	}
}

func buttonString(button glfw.MouseButton) string {
	switch button {
	case glfw.MouseButtonLeft:
		return "left"
	case glfw.MouseButtonRight:
		return "right"
	case glfw.MouseButtonMiddle:
		return "middle"
	default:
		return fmt.Sprint(button)
	}
}

func modsString(mods glfw.ModifierKey) string {
	if mods == 0 {
		return " no mods"
	}
	var name string
	if mods&glfw.ModShift != 0 {
		name += " shift"
	}
	if mods&glfw.ModControl != 0 {
		name += " control"
	}
	if mods&glfw.ModAlt != 0 {
		name += " alt"
	}
	if mods&glfw.ModSuper != 0 {
		name += " super"
	}
	return name
}

func charString(char rune) string {
	return fmt.Sprintf("%#q", char)
}

func PosCallback(w *glfw.Window, x int, y int) {
	fmt.Printf("%08x to %v at %0.3f: Window position: %v %v\n",
		getCounter(), getWindowId(w), getTime(),
		x, y)
}

func SizeCallback(w *glfw.Window, width int, height int) {
	fmt.Printf("%08x to %v at %0.3f: Window size: %v %v\n",
		getCounter(), getWindowId(w), getTime(),
		width, height)
}

func FramebufferSizeCallback(w *glfw.Window, width int, height int) {
	fmt.Printf("%08x to %v at %0.3f: Framebuffer size: %v %v\n",
		getCounter(), getWindowId(w), getTime(),
		width, height)
}

func CloseCallback(w *glfw.Window) {
	fmt.Printf("%08x to %v at %0.3f: Window close\n",
		getCounter(), getWindowId(w), getTime())
}

func RefreshCallback(w *glfw.Window) {
	fmt.Printf("%08x to %v at %0.3f: Window refresh\n",
		getCounter(), getWindowId(w), getTime())
}

func FocusCallback(w *glfw.Window, focused bool) {
	focusedString := map[bool]string{
		true:  "focused",
		false: "defocused",
	}

	fmt.Printf("%08x to %v at %0.3f: Window %s\n",
		getCounter(), getWindowId(w), getTime(),
		focusedString[focused])
}

func IconifyCallback(w *glfw.Window, iconified bool) {
	iconifiedString := map[bool]string{
		true:  "iconified",
		false: "restored",
	}

	fmt.Printf("%08x to %v at %0.3f: Window was %s\n",
		getCounter(), getWindowId(w), getTime(),
		iconifiedString[iconified])
}

func MouseButtonCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	fmt.Printf("%08x to %v at %0.3f: Mouse button %v (%s) (with%s) was %s\n",
		getCounter(), getWindowId(w), getTime(),
		button, buttonString(button), modsString(mods), actionString(action))
}

func CursorPosCallback(w *glfw.Window, x float64, y float64) {
	fmt.Printf("%08x to %v at %0.3f: Cursor position: %f %f\n",
		getCounter(), getWindowId(w), getTime(),
		x, y)
}

func CursorEnterCallback(w *glfw.Window, entered bool) {
	enteredString := map[bool]string{
		true:  "entered",
		false: "left",
	}

	fmt.Printf("%08x to %v at %0.3f: Cursor %s window\n",
		getCounter(), getWindowId(w), getTime(),
		enteredString[entered])
}

func ScrollCallback(w *glfw.Window, x float64, y float64) {
	fmt.Printf("%08x to %v at %0.3f: Scroll: %0.3f %0.3f\n",
		getCounter(), getWindowId(w), getTime(),
		x, y)
}

func KeyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	fmt.Printf("%08x to %v at %0.3f: Key 0x%04x Scancode 0x%04x (%s) (with%s) was %s\n",
		getCounter(), getWindowId(w), getTime(),
		key, scancode, keyString(key), modsString(mods), actionString(action))
}

func CharCallback(w *glfw.Window, char rune) {
	fmt.Printf("%08x to %v at %0.3f: Character 0x%08x (%s) input\n",
		getCounter(), getWindowId(w), getTime(),
		char, charString(char))
}

func CharModsCallback(w *glfw.Window, char rune, mods glfw.ModifierKey) {
	fmt.Printf("%08x to %v at %0.3f: Character 0x%08x (%s) with modifiers (with%s) input\n",
		getCounter(), getWindowId(w), getTime(),
		char, charString(char), modsString(mods))
}

func DropCallback(w *glfw.Window, names []string) {
	fmt.Printf("%08x to %v at %0.3f: Drop input\n",
		getCounter(), getWindowId(w), getTime())
	for i, name := range names {
		fmt.Printf("  %v: %q\n", i, name)
	}
}

func main() {
	err := run()
	if err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	err := glfw.Init(nil)
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	var windowSize = [2]int{640, 480}
	glfw.WindowHint(glfw.Samples, 8) // Anti-aliasing.
	window, err := glfw.CreateWindow(windowSize[0], windowSize[1], "Bhojpur GUI - Wasm Rendering", nil, nil)
	if err != nil {
		panic(err)
	}
	windowIds[window] = 1 // First (and only) window has id 1.

	// window.MakeContextCurrent()
	window.SetPosCallback(PosCallback)
	window.SetSizeCallback(SizeCallback)
	window.SetFramebufferSizeCallback(FramebufferSizeCallback)
	window.SetCloseCallback(CloseCallback)
	window.SetRefreshCallback(RefreshCallback)
	window.SetFocusCallback(FocusCallback)
	window.SetIconifyCallback(IconifyCallback)
	window.SetMouseButtonCallback(MouseButtonCallback)
	// window.SetCursorPosCallback(CursorPosCallback)
	window.SetCursorEnterCallback(CursorEnterCallback)
	window.SetScrollCallback(ScrollCallback)
	window.SetKeyCallback(KeyCallback)
	window.SetCharCallback(CharCallback)
	window.SetCharModsCallback(CharModsCallback)
	window.SetDropCallback(DropCallback)

	fmt.Printf("GLSL: stage 1")
	var cursorPos = [2]float32{float32(windowSize[0]) / 2, float32(windowSize[1]) / 2}
	var lastCursorPos = cursorPos
	cursorPosCallback := func(_ *glfw.Window, x, y float64) {
		cursorPos[0], cursorPos[1] = float32(x), float32(y)
	}
	window.SetCursorPosCallback(cursorPosCallback)

	fmt.Printf("GLSL: stage 2")
	// Set OpenGL options.
	gl.ClearColor(0, 0, 0, 1)
	gl.Enable(gl.CULL_FACE)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE)
	gl.Enable(gl.BLEND)

	// Init shaders.
	program, err := glutil.CreateProgram(utils.VertexSource, utils.FragmentSource)
	if err != nil {
		return err
	}

	fmt.Printf("GLSL: stage 3")
	gl.ValidateProgram(program)
	if gl.GetProgrami(program, gl.VALIDATE_STATUS) != gl.TRUE {
		return fmt.Errorf("gl validate status: %s", gl.GetProgramInfoLog(program))
	}

	gl.UseProgram(program)

	pMatrixUniform := gl.GetUniformLocation(program, "uPMatrix")
	mvMatrixUniform := gl.GetUniformLocation(program, "uMVMatrix")

	fmt.Printf("GLSL: stage 4")
	tri0v0 := gl.GetUniformLocation(program, "tri0v0")
	tri0v1 := gl.GetUniformLocation(program, "tri0v1")
	tri0v2 := gl.GetUniformLocation(program, "tri0v2")
	tri1v0 := gl.GetUniformLocation(program, "tri1v0")
	tri1v1 := gl.GetUniformLocation(program, "tri1v1")
	tri1v2 := gl.GetUniformLocation(program, "tri1v2")

	vertexPositionAttrib := gl.GetAttribLocation(program, "aVertexPosition")
	gl.EnableVertexAttribArray(vertexPositionAttrib)

	triangleVertexPositionBuffer := gl.CreateBuffer()

	fmt.Printf("GLSL: stage 5")
	// drawTriangle draws a triangle, consisting of 3 vertices, with motion blur corresponding
	// to the provided velocity. The triangle vertices specify its final position (at t = 1.0,
	// the end of frame), and its velocity is used to compute where the triangle is coming from
	// (at t = 0.0, the start of frame).
	drawTriangle := func(triangle [9]float32, velocity mgl32.Vec3) {
		triangle0 := triangle
		for i := 0; i < 3*3; i++ {
			triangle0[i] -= velocity[i%3]
		}
		triangle1 := triangle

		gl.Uniform3f(tri0v0, triangle0[0], triangle0[1], triangle0[2])
		gl.Uniform3f(tri0v1, triangle0[3], triangle0[4], triangle0[5])
		gl.Uniform3f(tri0v2, triangle0[6], triangle0[7], triangle0[8])
		gl.Uniform3f(tri1v0, triangle1[0], triangle1[1], triangle1[2])
		gl.Uniform3f(tri1v1, triangle1[3], triangle1[4], triangle1[5])
		gl.Uniform3f(tri1v2, triangle1[6], triangle1[7], triangle1[8])

		{
			gl.BindBuffer(gl.ARRAY_BUFFER, triangleVertexPositionBuffer)
			vertices := f32.Bytes(binary.LittleEndian,
				triangle0[0], triangle0[1], triangle0[2],
				triangle0[3], triangle0[4], triangle0[5],
				triangle0[6], triangle0[7], triangle0[8],
				triangle1[0], triangle1[1], triangle1[2],
				triangle1[6], triangle1[7], triangle1[8],
				triangle1[3], triangle1[4], triangle1[5],
			)
			gl.BufferData(gl.ARRAY_BUFFER, vertices, gl.DYNAMIC_DRAW)
			itemSize := 3
			itemCount := 6

			gl.VertexAttribPointer(vertexPositionAttrib, itemSize, gl.FLOAT, false, 0, 0)
			gl.DrawArrays(gl.TRIANGLES, 0, itemCount)
		}

		{
			gl.BindBuffer(gl.ARRAY_BUFFER, triangleVertexPositionBuffer)
			vertices := f32.Bytes(binary.LittleEndian,
				triangle0[0], triangle0[1], triangle0[2],
				triangle1[0], triangle1[1], triangle1[2],
				triangle0[3], triangle0[4], triangle0[5],
				triangle1[3], triangle1[4], triangle1[5],
				triangle0[6], triangle0[7], triangle0[8],
				triangle1[6], triangle1[7], triangle1[8],
				triangle0[0], triangle0[1], triangle0[2],
				triangle1[0], triangle1[1], triangle1[2],
			)
			gl.BufferData(gl.ARRAY_BUFFER, vertices, gl.DYNAMIC_DRAW)
			itemSize := 3
			itemCount := 8

			gl.VertexAttribPointer(vertexPositionAttrib, itemSize, gl.FLOAT, false, 0, 0)
			gl.DrawArrays(gl.TRIANGLE_STRIP, 0, itemCount)
		}
	}

	fmt.Printf("GLSL: stage 6")
	if err := gl.GetError(); err != 0 {
		return fmt.Errorf("gl error: %v", err)
	}

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT)

		pMatrix := mgl32.Ortho2D(0, float32(windowSize[0]), float32(windowSize[1]), 0)

		triangle0 := [9]float32{
			-50, -50, 0,
			50, -50, 0,
			-50, 50, 0}
		triangle1 := [9]float32{
			50, 50, 0,
			-50, 50, 0,
			50, -50, 0}

		// Square with motion blur on the left.
		{
			mvMatrix := mgl32.Translate3D(cursorPos[0]-200, cursorPos[1], 0)

			gl.UniformMatrix4fv(pMatrixUniform, pMatrix[:])
			gl.UniformMatrix4fv(mvMatrixUniform, mvMatrix[:])

			velocity := mgl32.Vec3{cursorPos[0] - lastCursorPos[0], cursorPos[1] - lastCursorPos[1], 0}

			drawTriangle(triangle0, velocity)
			drawTriangle(triangle1, velocity)
		}

		// Square without motion blur on the right.
		{
			mvMatrix := mgl32.Translate3D(cursorPos[0]+200, cursorPos[1], 0)

			gl.UniformMatrix4fv(pMatrixUniform, pMatrix[:])
			gl.UniformMatrix4fv(mvMatrixUniform, mvMatrix[:])

			drawTriangle(triangle0, mgl32.Vec3{})
			drawTriangle(triangle1, mgl32.Vec3{})
		}

		lastCursorPos = cursorPos

		window.SwapBuffers()
		glfw.PollEvents()
	}

	return nil
}
