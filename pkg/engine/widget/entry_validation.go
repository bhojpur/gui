package widget

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
	"github.com/bhojpur/gui/pkg/engine/canvas"
	"github.com/bhojpur/gui/pkg/engine/theme"
)

var _ gui.Validatable = (*Entry)(nil)

// Validate validates the current text in the widget
func (e *Entry) Validate() error {
	if e.Validator == nil {
		return nil
	}

	err := e.Validator(e.Text)
	e.SetValidationError(err)
	return err
}

// SetOnValidationChanged is intended for parent widgets or containers to hook into the validation.
// The function might be overwritten by a parent that cares about child validation (e.g. widget.Form).
func (e *Entry) SetOnValidationChanged(callback func(error)) {
	if callback != nil {
		e.onValidationChanged = callback
	}
}

// SetValidationError manually updates the validation status until the next input change
func (e *Entry) SetValidationError(err error) {
	if e.Validator == nil {
		return
	}
	if err == nil && e.validationError == nil {
		return
	}

	if (err == nil && e.validationError != nil) || (e.validationError == nil && err != nil) ||
		err.Error() != e.validationError.Error() {
		e.validationError = err

		if e.onValidationChanged != nil {
			e.onValidationChanged(err)
		}

		e.Refresh()
	}
}

var _ gui.Widget = (*validationStatus)(nil)

type validationStatus struct {
	BaseWidget
	entry *Entry
}

func newValidationStatus(e *Entry) *validationStatus {
	rs := &validationStatus{
		entry: e,
	}

	rs.ExtendBaseWidget(rs)
	return rs
}

func (r *validationStatus) CreateRenderer() gui.WidgetRenderer {
	icon := &canvas.Image{}
	icon.Hide()
	return &validationStatusRenderer{
		WidgetRenderer: NewSimpleRenderer(icon),
		icon:           icon,
		entry:          r.entry,
	}
}

var _ gui.WidgetRenderer = (*validationStatusRenderer)(nil)

type validationStatusRenderer struct {
	gui.WidgetRenderer
	entry *Entry
	icon  *canvas.Image
}

func (r *validationStatusRenderer) Layout(size gui.Size) {
	r.icon.Resize(gui.NewSize(theme.IconInlineSize(), theme.IconInlineSize()))
	r.icon.Move(gui.NewPos((size.Width-theme.IconInlineSize())/2, (size.Height-theme.IconInlineSize())/2))
}

func (r *validationStatusRenderer) MinSize() gui.Size {
	return gui.NewSize(theme.IconInlineSize(), theme.IconInlineSize())
}

func (r *validationStatusRenderer) Refresh() {
	r.entry.propertyLock.RLock()
	defer r.entry.propertyLock.RUnlock()
	if r.entry.disabled {
		r.icon.Hide()
		return
	}

	if r.entry.validationError == nil && r.entry.Text != "" {
		r.icon.Resource = theme.ConfirmIcon()
		r.icon.Show()
	} else if r.entry.validationError != nil && !r.entry.focused && r.entry.dirty {
		r.icon.Resource = theme.NewErrorThemedResource(theme.ErrorIcon())
		r.icon.Show()
	} else {
		r.icon.Hide()
	}

	canvas.Refresh(r.icon)
}
