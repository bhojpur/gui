package dialog

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
	gui "github.com/bhojpur/gui/pkg/engine"
	"github.com/bhojpur/gui/pkg/engine/container"
	"github.com/bhojpur/gui/pkg/engine/layout"
	"github.com/bhojpur/gui/pkg/engine/theme"
	"github.com/bhojpur/gui/pkg/engine/widget"
)

// formDialog is a simple dialog window for displaying FormItems inside a form.
type formDialog struct {
	*dialog
	items   []*widget.FormItem
	confirm *widget.Button
	cancel  *widget.Button
}

// validateItems acts as a validation edge state handler that will respond to an individual widget's validation
// state before checking all others to determine the net validation state. If the error passed is not nil, then the
// confirm button will be disabled. If the error parameter is nil, then all other Validatable widgets in items are
// checked as well to determine whether the confirm button should be disabled.
// This method is passed to each Validatable widget's SetOnValidationChanged method in items by NewFormDialog.
func (d *formDialog) validateItems(err error) {
	if err != nil {
		d.confirm.Disable()
		return
	}
	for _, item := range d.items {
		if validatable, ok := item.Widget.(gui.Validatable); ok {
			if err := validatable.Validate(); err != nil {
				d.confirm.Disable()
				return
			}
		}
	}
	d.confirm.Enable()
}

// NewForm creates and returns a dialog over the specified application using
// the provided FormItems. The cancel button will have the dismiss text set and the confirm button will
// use the confirm text. The response callback is called on user action after validation passes.
// If any Validatable widget reports that validation has failed, then the confirm
// button will be disabled. The initial state of the confirm button will reflect the initial
// validation state of the items added to the form dialog.
//
// Since: 2.0
func NewForm(title, confirm, dismiss string, items []*widget.FormItem, callback func(bool), parent gui.Window) Dialog {
	var itemObjects = make([]gui.CanvasObject, len(items)*2)
	for i, item := range items {
		itemObjects[i*2] = widget.NewLabel(item.Text)
		itemObjects[i*2+1] = item.Widget
	}
	content := gui.NewContainerWithLayout(layout.NewFormLayout(), itemObjects...)

	d := &dialog{content: content, callback: callback, title: title, parent: parent}
	d.layout = &dialogLayout{d: d}
	d.dismiss = &widget.Button{Text: dismiss, Icon: theme.CancelIcon(),
		OnTapped: d.Hide,
	}
	confirmBtn := &widget.Button{Text: confirm, Icon: theme.ConfirmIcon(), Importance: widget.HighImportance,
		OnTapped: func() { d.hideWithResponse(true) },
	}
	formDialog := &formDialog{
		dialog:  d,
		items:   items,
		confirm: confirmBtn,
		cancel:  d.dismiss,
	}

	formDialog.validateItems(nil)

	for _, item := range items {
		if validatable, ok := item.Widget.(gui.Validatable); ok {
			validatable.SetOnValidationChanged(formDialog.validateItems)
		}
	}
	d.setButtons(container.NewHBox(layout.NewSpacer(), d.dismiss, confirmBtn, layout.NewSpacer()))
	return formDialog
}

// ShowForm shows a dialog over the specified application using
// the provided FormItems. The cancel button will have the dismiss text set and the confirm button will
// use the confirm text. The response callback is called on user action after validation passes.
// If any Validatable widget reports that validation has failed, then the confirm
// button will be disabled. The initial state of the confirm button will reflect the initial
// validation state of the items added to the form dialog.
// The MinSize() of the CanvasObject passed will be used to set the size of the window.
//
// Since: 2.0
func ShowForm(title, confirm, dismiss string, content []*widget.FormItem, callback func(bool), parent gui.Window) {
	NewForm(title, confirm, dismiss, content, callback, parent).Show()
}
