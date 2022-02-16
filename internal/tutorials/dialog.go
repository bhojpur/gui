package tutorials

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
	"errors"
	"fmt"
	"image/color"
	"io/ioutil"
	"log"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/canvas"
	"github.com/bhojpur/gui/pkg/engine/container"
	"github.com/bhojpur/gui/pkg/engine/data/validation"
	"github.com/bhojpur/gui/pkg/engine/dialog"
	"github.com/bhojpur/gui/pkg/engine/storage"
	"github.com/bhojpur/gui/pkg/engine/theme"
	"github.com/bhojpur/gui/pkg/engine/widget"
)

func confirmCallback(response bool) {
	fmt.Println("Responded with", response)
}

func colorPicked(c color.Color, w gui.Window) {
	log.Println("Color picked:", c)
	rectangle := canvas.NewRectangle(c)
	size := 2 * theme.IconInlineSize()
	rectangle.SetMinSize(gui.NewSize(size, size))
	dialog.ShowCustom("Color Picked", "Ok", rectangle, w)
}

// dialogScreen loads demos of the dialogs we support
func dialogScreen(win gui.Window) gui.CanvasObject {
	return container.NewVScroll(container.NewVBox(
		widget.NewButton("Info", func() {
			dialog.ShowInformation("Information", "You should know this thing...", win)
		}),
		widget.NewButton("Error", func() {
			err := errors.New("a dummy error message")
			dialog.ShowError(err, win)
		}),
		widget.NewButton("Confirm", func() {
			cnf := dialog.NewConfirm("Confirmation", "Are you enjoying this demo?", confirmCallback, win)
			cnf.SetDismissText("Nah")
			cnf.SetConfirmText("Oh Yes!")
			cnf.Show()
		}),
		widget.NewButton("File Open With Filter (.jpg or .png)", func() {
			fd := dialog.NewFileOpen(func(reader gui.URIReadCloser, err error) {
				if err != nil {
					dialog.ShowError(err, win)
					return
				}
				if reader == nil {
					log.Println("Cancelled")
					return
				}

				imageOpened(reader)
			}, win)
			fd.SetFilter(storage.NewExtensionFileFilter([]string{".png", ".jpg", ".jpeg"}))
			fd.Show()
		}),
		widget.NewButton("File Save", func() {
			dialog.ShowFileSave(func(writer gui.URIWriteCloser, err error) {
				if err != nil {
					dialog.ShowError(err, win)
					return
				}
				if writer == nil {
					log.Println("Cancelled")
					return
				}

				fileSaved(writer, win)
			}, win)
		}),
		widget.NewButton("Folder Open", func() {
			dialog.ShowFolderOpen(func(list gui.ListableURI, err error) {
				if err != nil {
					dialog.ShowError(err, win)
					return
				}
				if list == nil {
					log.Println("Cancelled")
					return
				}

				children, err := list.List()
				if err != nil {
					dialog.ShowError(err, win)
					return
				}
				out := fmt.Sprintf("Folder %s (%d children):\n%s", list.Name(), len(children), list.String())
				dialog.ShowInformation("Folder Open", out, win)
			}, win)
		}),
		widget.NewButton("Color Picker", func() {
			picker := dialog.NewColorPicker("Pick a Color", "What is your favorite color?", func(c color.Color) {
				colorPicked(c, win)
			}, win)
			picker.Show()
		}),
		widget.NewButton("Advanced Color Picker", func() {
			picker := dialog.NewColorPicker("Pick a Color", "What is your favorite color?", func(c color.Color) {
				colorPicked(c, win)
			}, win)
			picker.Advanced = true
			picker.Show()
		}),
		widget.NewButton("Form Dialog (Login Form)", func() {
			username := widget.NewEntry()
			username.Validator = validation.NewRegexp(`^[A-Za-z0-9_-]+$`, "username can only contain letters, numbers, '_', and '-'")
			password := widget.NewPasswordEntry()
			password.Validator = validation.NewRegexp(`^[A-Za-z0-9_-]+$`, "password can only contain letters, numbers, '_', and '-'")
			remember := false
			items := []*widget.FormItem{
				widget.NewFormItem("Username", username),
				widget.NewFormItem("Password", password),
				widget.NewFormItem("Remember me", widget.NewCheck("", func(checked bool) {
					remember = checked
				})),
			}

			dialog.ShowForm("Login...", "Log In", "Cancel", items, func(b bool) {
				if !b {
					return
				}
				var rememberText string
				if remember {
					rememberText = "and remember this login"
				}

				log.Println("Please Authenticate", username.Text, password.Text, rememberText)
			}, win)
		}),
	))
}

func imageOpened(f gui.URIReadCloser) {
	if f == nil {
		log.Println("Cancelled")
		return
	}
	defer f.Close()

	showImage(f)
}

func fileSaved(f gui.URIWriteCloser, w gui.Window) {
	defer f.Close()
	_, err := f.Write([]byte("Written by Bhojpur GUI demo application\n"))
	if err != nil {
		dialog.ShowError(err, w)
	}
	err = f.Close()
	if err != nil {
		dialog.ShowError(err, w)
	}
	log.Println("Saved to...", f.URI())
}

func loadImage(f gui.URIReadCloser) *canvas.Image {
	data, err := ioutil.ReadAll(f)
	if err != nil {
		gui.LogError("Failed to load image data", err)
		return nil
	}
	res := gui.NewStaticResource(f.URI().Name(), data)

	return canvas.NewImageFromResource(res)
}

func showImage(f gui.URIReadCloser) {
	img := loadImage(f)
	if img == nil {
		return
	}
	img.FillMode = canvas.ImageFillOriginal

	w := gui.CurrentApp().NewWindow(f.URI().Name())
	w.SetContent(container.NewScroll(img))
	w.Resize(gui.NewSize(320, 240))
	w.Show()
}
