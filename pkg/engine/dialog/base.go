// It defines standard dialog windows for application GUIs.
package dialog // import "github.com/bhojpur/gui/pkg/engine/dialog"

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
	"image/color"

	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/canvas"
	"github.com/bhojpur/gui/pkg/engine/container"
	col "github.com/bhojpur/gui/pkg/engine/internal/color"
	"github.com/bhojpur/gui/pkg/engine/layout"
	"github.com/bhojpur/gui/pkg/engine/theme"
	"github.com/bhojpur/gui/pkg/engine/widget"
)

const (
	padWidth  = 32
	padHeight = 16
)

// Dialog is the common API for any dialog window with a single dismiss button
type Dialog interface {
	Show()
	Hide()
	SetDismissText(label string)
	SetOnClosed(closed func())
	Refresh()
	Resize(size gui.Size)

	// Since: 2.1
	MinSize() gui.Size
}

// Declare conformity to Dialog interface
var _ Dialog = (*dialog)(nil)

type dialog struct {
	callback    func(bool)
	title       string
	icon        gui.Resource
	desiredSize gui.Size

	win            *widget.PopUp
	bg             *themedBackground
	content, label gui.CanvasObject
	dismiss        *widget.Button
	parent         gui.Window
	layout         *dialogLayout
}

// NewCustom creates and returns a dialog over the specified application using custom
// content. The button will have the dismiss text set.
// The MinSize() of the CanvasObject passed will be used to set the size of the window.
func NewCustom(title, dismiss string, content gui.CanvasObject, parent gui.Window) Dialog {
	d := &dialog{content: content, title: title, icon: nil, parent: parent}
	d.layout = &dialogLayout{d: d}

	d.dismiss = &widget.Button{Text: dismiss,
		OnTapped: d.Hide,
	}
	d.setButtons(container.NewHBox(layout.NewSpacer(), d.dismiss, layout.NewSpacer()))

	return d
}

// NewCustomConfirm creates and returns a dialog over the specified application using
// custom content. The cancel button will have the dismiss text set and the "OK" will
// use the confirm text. The response callback is called on user action.
// The MinSize() of the CanvasObject passed will be used to set the size of the window.
func NewCustomConfirm(title, confirm, dismiss string, content gui.CanvasObject,
	callback func(bool), parent gui.Window) Dialog {
	d := &dialog{content: content, title: title, icon: nil, parent: parent}
	d.layout = &dialogLayout{d: d}
	d.callback = callback

	d.dismiss = &widget.Button{Text: dismiss, Icon: theme.CancelIcon(),
		OnTapped: d.Hide,
	}
	ok := &widget.Button{Text: confirm, Icon: theme.ConfirmIcon(), Importance: widget.HighImportance,
		OnTapped: func() {
			d.hideWithResponse(true)
		},
	}
	d.setButtons(container.NewHBox(layout.NewSpacer(), d.dismiss, ok, layout.NewSpacer()))

	return d
}

// ShowCustom shows a dialog over the specified application using custom
// content. The button will have the dismiss text set.
// The MinSize() of the CanvasObject passed will be used to set the size of the window.
func ShowCustom(title, dismiss string, content gui.CanvasObject, parent gui.Window) {
	NewCustom(title, dismiss, content, parent).Show()
}

// ShowCustomConfirm shows a dialog over the specified application using custom
// content. The cancel button will have the dismiss text set and the "OK" will use
// the confirm text. The response callback is called on user action.
// The MinSize() of the CanvasObject passed will be used to set the size of the window.
func ShowCustomConfirm(title, confirm, dismiss string, content gui.CanvasObject,
	callback func(bool), parent gui.Window) {
	NewCustomConfirm(title, confirm, dismiss, content, callback, parent).Show()
}

func (d *dialog) Hide() {
	d.hideWithResponse(false)
}

// MinSize returns the size that this dialog should not shrink below
//
// Since: 2.1
func (d *dialog) MinSize() gui.Size {
	return d.win.MinSize()
}

func (d *dialog) Show() {
	if !d.desiredSize.IsZero() {
		d.win.Resize(d.desiredSize)
	}
	d.win.Show()
}

func (d *dialog) Refresh() {
	d.win.Refresh()
}

// Resize dialog, call this function after dialog show
func (d *dialog) Resize(size gui.Size) {
	d.desiredSize = size
	d.win.Resize(size)
}

// SetDismissText allows custom text to be set in the confirmation button
func (d *dialog) SetDismissText(label string) {
	d.dismiss.SetText(label)
	d.win.Refresh()
}

// SetOnClosed allows to set a callback function that is called when
// the dialog is closed
func (d *dialog) SetOnClosed(closed func()) {
	// if there is already a callback set, remember it and call both
	originalCallback := d.callback

	d.callback = func(response bool) {
		closed()
		if originalCallback != nil {
			originalCallback(response)
		}
	}
}

func (d *dialog) hideWithResponse(resp bool) {
	d.win.Hide()
	if d.callback != nil {
		d.callback(resp)
	}
}

func (d *dialog) setButtons(buttons gui.CanvasObject) {
	d.bg = newThemedBackground()
	d.label = widget.NewLabelWithStyle(d.title, gui.TextAlignLeading, gui.TextStyle{Bold: true})

	var content gui.CanvasObject
	if d.icon == nil {
		content = container.New(d.layout,
			&canvas.Image{},
			d.bg,
			d.content,
			buttons,
			d.label,
		)
	} else {
		bgIcon := canvas.NewImageFromResource(d.icon)
		content = container.New(d.layout,
			bgIcon,
			d.bg,
			d.content,
			buttons,
			d.label,
		)
	}

	d.win = widget.NewModalPopUp(content, d.parent.Canvas())
	d.Refresh()
}

func newDialog(title, message string, icon gui.Resource, callback func(bool), parent gui.Window) *dialog {
	d := &dialog{content: newLabel(message), title: title, icon: icon, parent: parent}
	d.layout = &dialogLayout{d: d}

	d.callback = callback

	return d
}

func newLabel(message string) gui.CanvasObject {
	return widget.NewLabelWithStyle(message, gui.TextAlignCenter, gui.TextStyle{})
}

func newButtonList(buttons ...*widget.Button) gui.CanvasObject {
	list := container.New(layout.NewGridLayout(len(buttons)))
	for _, button := range buttons {
		list.Add(button)
	}

	return list
}

// ===============================================================
// ThemedBackground
// ===============================================================

type themedBackground struct {
	widget.BaseWidget
}

func newThemedBackground() *themedBackground {
	t := &themedBackground{}
	t.ExtendBaseWidget(t)
	return t
}

func (t *themedBackground) CreateRenderer() gui.WidgetRenderer {
	t.ExtendBaseWidget(t)
	rect := canvas.NewRectangle(theme.BackgroundColor())
	return &themedBackgroundRenderer{rect, []gui.CanvasObject{rect}}
}

type themedBackgroundRenderer struct {
	rect    *canvas.Rectangle
	objects []gui.CanvasObject
}

func (renderer *themedBackgroundRenderer) Destroy() {
}

func (renderer *themedBackgroundRenderer) Layout(size gui.Size) {
	renderer.rect.Resize(size)
}

func (renderer *themedBackgroundRenderer) MinSize() gui.Size {
	return renderer.rect.MinSize()
}

func (renderer *themedBackgroundRenderer) Objects() []gui.CanvasObject {
	return renderer.objects
}

func (renderer *themedBackgroundRenderer) Refresh() {
	r, g, b, _ := col.ToNRGBA(theme.BackgroundColor())
	bg := &color.NRGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 230}
	renderer.rect.FillColor = bg
}

// ===============================================================
// DialogLayout
// ===============================================================

type dialogLayout struct {
	d *dialog
}

func (l *dialogLayout) Layout(obj []gui.CanvasObject, size gui.Size) {
	l.d.bg.Move(gui.NewPos(0, 0))
	l.d.bg.Resize(size)

	btnMin := obj[3].MinSize()

	// icon
	iconHeight := padHeight*2 + l.d.label.MinSize().Height*2 - theme.Padding()
	obj[0].Resize(gui.NewSize(iconHeight, iconHeight))
	obj[0].Move(gui.NewPos(size.Width-iconHeight+theme.Padding(), -theme.Padding()))

	// buttons
	obj[3].Resize(btnMin)
	obj[3].Move(gui.NewPos(size.Width/2-(btnMin.Width/2), size.Height-padHeight-btnMin.Height))

	// content
	contentStart := l.d.label.Position().Y + l.d.label.MinSize().Height + padHeight
	contentEnd := obj[3].Position().Y - theme.Padding()
	obj[2].Move(gui.NewPos(padWidth/2, l.d.label.MinSize().Height+padHeight))
	obj[2].Resize(gui.NewSize(size.Width-padWidth, contentEnd-contentStart))
}

func (l *dialogLayout) MinSize(obj []gui.CanvasObject) gui.Size {
	contentMin := obj[2].MinSize()
	btnMin := obj[3].MinSize()

	width := gui.Max(gui.Max(contentMin.Width, btnMin.Width), obj[4].MinSize().Width) + padWidth
	height := contentMin.Height + btnMin.Height + l.d.label.MinSize().Height + theme.Padding() + padHeight*2

	return gui.NewSize(width, height)
}
